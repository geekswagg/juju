// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package instances

import (
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/environs/context"
	"github.com/juju/juju/network"
)

// Instance represents the the realization of a machine in state.
type Instance interface {
	// Id returns a provider-generated identifier for the Instance.
	Id() instance.Id

	// Status returns the provider-specific status for the instance.
	Status(context.ProviderCallContext) instance.Status

	// Addresses returns a list of hostnames or ip addresses
	// associated with the instance.
	Addresses(context.ProviderCallContext) ([]network.Address, error)
}

// InstanceFirewaller provides instance-level firewall functionality
type InstanceFirewaller interface {
	// OpenPorts opens the given port ranges on the instance, which
	// should have been started with the given machine id.
	OpenPorts(ctx context.ProviderCallContext, machineId string, rules []network.IngressRule) error

	// ClosePorts closes the given port ranges on the instance, which
	// should have been started with the given machine id.
	ClosePorts(ctx context.ProviderCallContext, machineId string, rules []network.IngressRule) error

	// IngressRules returns the set of ingress rules for the instance,
	// which should have been applied to the given machine id. The
	// rules are returned as sorted by network.SortIngressRules().
	// It is expected that there be only one ingress rule result for a given
	// port range - the rule's SourceCIDRs will contain all applicable source
	// address rules for that port range.
	IngressRules(ctx context.ProviderCallContext, machineId string) ([]network.IngressRule, error)
}
