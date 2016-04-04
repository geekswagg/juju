// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package maas

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/gomaasapi"

	"github.com/juju/juju/instance"
	"github.com/juju/juju/network"
	"github.com/juju/juju/status"
)

type maasInstance struct {
	maasObject   *gomaasapi.MAASObject
	environ      *maasEnviron
	statusGetter func(instance.Id) (string, string)
}

var _ instance.Instance = (*maasInstance)(nil)

// Override for testing.
var resolveHostnames = func(addrs []network.Address) []network.Address {
	return network.ResolvableHostnames(addrs)
}

func (mi *maasInstance) String() string {
	hostname, err := mi.hostname()
	if err != nil {
		// This is meant to be impossible, but be paranoid.
		hostname = fmt.Sprintf("<DNSName failed: %q>", err)
	}
	return fmt.Sprintf("%s:%s", hostname, mi.Id())
}

func (mi *maasInstance) Id() instance.Id {
	return maasObjectId(mi.maasObject)
}

func maasObjectId(maasObject *gomaasapi.MAASObject) instance.Id {
	// Use the node's 'resource_uri' value.
	return instance.Id(maasObject.URI().String())
}

// Status returns a juju status based on the maas instance returned
// status message.
func (mi *maasInstance) Status() instance.InstanceStatus {
	maasInstanceStatus := status.StatusEmpty
	statusMsg, substatus := mi.statusGetter(mi.Id())
	switch statusMsg {
	case "":
		logger.Debugf("unable to obtain status of instance %s", mi.Id())
		statusMsg = "error in getting status"
	case "Deployed":
		maasInstanceStatus = status.StatusRunning
	case "Deploying":
		maasInstanceStatus = status.StatusAllocating
		if substatus != "" {
			statusMsg = fmt.Sprintf("%s: %s", statusMsg, substatus)
		}
	case "Failed Deployment":
		maasInstanceStatus = status.StatusProvisioningError
		if substatus != "" {
			statusMsg = fmt.Sprintf("%s: %s", statusMsg, substatus)
		}
	default:
		maasInstanceStatus = status.StatusEmpty
		statusMsg = fmt.Sprintf("%s: %s", statusMsg, substatus)
	}

	return instance.InstanceStatus{
		Status:  maasInstanceStatus,
		Message: statusMsg,
	}
}

func (mi *maasInstance) Addresses() ([]network.Address, error) {
	interfaceAddresses, err := mi.interfaceAddresses()
	if errors.IsNotSupported(err) {
		logger.Warningf(
			"using legacy approach to get instance addresses as %q API capability is not supported: %v",
			capNetworkDeploymentUbuntu, err,
		)
		return mi.legacyAddresses()
	} else if err != nil {
		return nil, errors.Annotate(err, "getting node interfaces")
	}

	logger.Debugf("instance %q has interface addresses: %+v", mi.Id(), interfaceAddresses)
	return interfaceAddresses, nil
}

// legacyAddresses is used to extract all IP addresses of the node when not
// using MAAS 1.9+ API.
func (mi *maasInstance) legacyAddresses() ([]network.Address, error) {
	name, err := mi.hostname()
	if err != nil {
		return nil, err
	}
	// MAAS prefers to use the dns name for intra-node communication.
	// When Juju looks up the address to use for communicating between
	// nodes, it looks up the address by scope. So we add a cloud
	// local address for that purpose.
	addrs := network.NewAddresses(name, name)
	addrs[0].Scope = network.ScopePublic
	addrs[1].Scope = network.ScopeCloudLocal

	// Append any remaining IP addresses after the preferred ones. We have to do
	// this the hard way, since maasObject doesn't have this built-in yet
	addressArray := mi.maasObject.GetMap()["ip_addresses"]
	if !addressArray.IsNil() {
		// Older MAAS versions do not return ip_addresses.
		objs, err := addressArray.GetArray()
		if err != nil {
			return nil, err
		}
		for _, obj := range objs {
			s, err := obj.GetString()
			if err != nil {
				return nil, err
			}
			addrs = append(addrs, network.NewAddress(s))
		}
	}

	// Although we would prefer a DNS name there's no point
	// returning unresolvable names because activities like 'juju
	// ssh 0' will instantly fail.
	return resolveHostnames(addrs), nil
}

var refreshMAASObject = func(maasObject *gomaasapi.MAASObject) (gomaasapi.MAASObject, error) {
	// Defined like this to allow patching in tests to overcome limitations of
	// gomaasapi's test server.
	return maasObject.Get()
}

// interfaceAddresses fetches a fresh copy of the node details from MAAS and
// extracts all addresses from the node's interfaces. Returns an error
// satisfying errors.IsNotSupported() if MAAS API does not report interfaces
// information.
func (mi *maasInstance) interfaceAddresses() ([]network.Address, error) {
	// Fetch a fresh copy of the instance JSON first.
	obj, err := refreshMAASObject(mi.maasObject)
	if err != nil {
		return nil, errors.Annotate(err, "getting instance details")
	}

	subnetsMap, err := mi.environ.subnetToSpaceIds()
	if err != nil {
		return nil, errors.Trace(err)
	}
	// Get all the interface details and extract the addresses.
	interfaces, err := maasObjectNetworkInterfaces(&obj, subnetsMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var addresses []network.Address
	for _, iface := range interfaces {
		if iface.Address.Value != "" {
			logger.Debugf("found address %q on interface %q", iface.Address, iface.InterfaceName)
			addresses = append(addresses, iface.Address)
		} else {
			logger.Infof("no address found on interface %q", iface.InterfaceName)
		}
	}
	return addresses, nil
}

func (mi *maasInstance) architecture() (arch, subarch string, err error) {
	// MAAS may return an architecture of the form, for example,
	// "amd64/generic"; we only care about the major part.
	arch, err = mi.maasObject.GetField("architecture")
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(arch, "/", 2)
	arch = parts[0]
	if len(parts) == 2 {
		subarch = parts[1]
	}
	return arch, subarch, nil
}

func (mi *maasInstance) zone() (string, error) {
	// TODO (anastasiamac 2016-03-31)
	// This code is needed until gomaasapi testing code is
	// updated to align with MAAS.
	// Currently, "zone" property is still treated as field
	// by gomaasi infrastructure and is searched for
	// using matchField(node, "zone", zoneName) instead of
	// getMap.
	// @see gomaasapi/testservice.go#findFreeNode
	// bug https://bugs.launchpad.net/gomaasapi/+bug/1563631
	zone, fieldErr := mi.maasObject.GetField("zone")
	if fieldErr == nil && zone != "" {
		return zone, nil
	}

	obj := mi.maasObject.GetMap()["zone"]
	if obj.IsNil() {
		return "", errors.New("zone property not set on maas")
	}
	zoneMap, err := obj.GetMap()
	if err != nil {
		return "", errors.New("zone property is not an expected type")
	}
	nameObj, ok := zoneMap["name"]
	if !ok {
		return "", errors.New("zone property is not set correctly: name is missing")
	}
	str, err := nameObj.GetString()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (mi *maasInstance) cpuCount() (uint64, error) {
	count, err := mi.maasObject.GetMap()["cpu_count"].GetFloat64()
	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}

func (mi *maasInstance) memory() (uint64, error) {
	mem, err := mi.maasObject.GetMap()["memory"].GetFloat64()
	if err != nil {
		return 0, err
	}
	return uint64(mem), nil
}

func (mi *maasInstance) tagNames() ([]string, error) {
	obj := mi.maasObject.GetMap()["tag_names"]
	if obj.IsNil() {
		return nil, errors.NotFoundf("tag_names")
	}
	array, err := obj.GetArray()
	if err != nil {
		return nil, err
	}
	tags := make([]string, len(array))
	for i, obj := range array {
		tag, err := obj.GetString()
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}
	return tags, nil
}

func (mi *maasInstance) hardwareCharacteristics() (*instance.HardwareCharacteristics, error) {
	nodeArch, _, err := mi.architecture()
	if err != nil {
		return nil, errors.Annotate(err, "error determining architecture")
	}
	nodeCpuCount, err := mi.cpuCount()
	if err != nil {
		return nil, errors.Annotate(err, "error determining cpu count")
	}
	nodeMemoryMB, err := mi.memory()
	if err != nil {
		return nil, errors.Annotate(err, "error determining available memory")
	}
	zone, err := mi.zone()
	if err != nil {
		return nil, errors.Annotate(err, "error determining availability zone")
	}
	hc := &instance.HardwareCharacteristics{
		Arch:             &nodeArch,
		CpuCores:         &nodeCpuCount,
		Mem:              &nodeMemoryMB,
		AvailabilityZone: &zone,
	}
	nodeTags, err := mi.tagNames()
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Annotate(err, "error determining tag names")
	}
	if len(nodeTags) > 0 {
		hc.Tags = &nodeTags
	}
	return hc, nil
}

func (mi *maasInstance) hostname() (string, error) {
	// A MAAS instance has its DNS name immediately.
	return mi.maasObject.GetField("hostname")
}

// MAAS does not do firewalling so these port methods do nothing.
func (mi *maasInstance) OpenPorts(machineId string, ports []network.PortRange) error {
	logger.Debugf("unimplemented OpenPorts() called")
	return nil
}

func (mi *maasInstance) ClosePorts(machineId string, ports []network.PortRange) error {
	logger.Debugf("unimplemented ClosePorts() called")
	return nil
}

func (mi *maasInstance) Ports(machineId string) ([]network.PortRange, error) {
	logger.Debugf("unimplemented Ports() called")
	return nil, nil
}
