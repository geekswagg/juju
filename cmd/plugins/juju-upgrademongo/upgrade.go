// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	goyaml "gopkg.in/yaml.v1"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/juju/agent"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/juju"
	"github.com/juju/juju/utils/ssh"
	"github.com/juju/loggo"
	"github.com/juju/utils"
)

const noauth = `
NoAuth() {
    sed -i "s/--auth/--noauth/" /etc/systemd/system/juju-db*.service;
    sed -i "s/--keyFile '\/var\/lib\/juju\/shared-secret'//" /etc/systemd/system/juju-db*.service;
}

Auth() {
    sed -i "s/--noauth/--auth --keyFile '\/var\/lib\/juju\/shared-secret'/" /etc/systemd/system/juju-db*.service;
}
`

const replset = `
NoReplSet() {
    sed -i "s/--replSet juju//" /etc/systemd/system/juju-db*.service;
}

ReplSet() {
    sed -i "s/mongod /mongod --replSet juju /" /etc/systemd/system/juju-db*.service;
}
`

const upgradeMongoInAgentConfig = `
ReplaceVersion() {
    	sed -i "s/mongoversion:.*/mongoversion: \"${1}\"/" /var/lib/juju/agents/machine-0/agent.conf
}

AddVersion () {
    echo "mongoversion: \"${1}\"" >> /var/lib/juju/agents/machine-0/agent.conf
}

UpgradeMongoVersion() {
    VERSIONKEY=$(grep mongoversion /var/lib/juju/agents/machine-0/agent.conf)
    if [ -n "$VERSIONKEY" ]; then
	ReplaceVersion $1;
    else
	AddVersion $1;
    fi
}
`

const upgradeMongoBinary = `
UpgradeMongoBinary() {
    sed -i "s/juju\/bin/juju\/mongo${1}\/bin/" /etc/systemd/system/juju-db.service;
    sed -i "s/juju\/mongo.*\/bin/juju\/mongo${1}\/bin/" /etc/systemd/system/juju-db.service;
    if [ "$2" == "storage" ]; then
	sed -i "s/--smallfiles//" /etc/systemd/system/juju-db.service;
	sed -i "s/--noprealloc/--storageEngine wiredTiger/" /etc/systemd/system/juju-db.service;
    fi
    cat  /etc/systemd/system/juju-db.service
}
`

const mongoEval = `
mongoAdminEval() {
        attempts=0
        until [ $attempts -ge 60 ]
        do
    	    mongo --ssl -u admin -p {{.OldPassword | shquote}} localhost:{{.StatePort}}/admin --eval "printjson($1)" && break
            echo "attempt $attempts"
            attempts=$[$attempts+1]
            sleep 10
        done
}

mongoAnonEval() {
        attempts=0
        until [ $attempts -ge 60 ]
        do
    	    mongo --ssl  localhost:{{.StatePort}}/admin --eval "printjson($1)" && break
            echo "attempt $attempts"
            attempts=$[$attempts+1]
            sleep 10
        done
}
`

func main() {
	Main(os.Args)
}

// Main is the entry point for this plugins.
func Main(args []string) {
	ctx, err := cmd.DefaultContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	if err := juju.InitJujuHome(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
	os.Exit(cmd.Main(envcmd.Wrap(&upgradeMongoCommand{}), ctx, args[1:]))
}

const upgradeDoc = ``

var logger = loggo.GetLogger("juju.plugins.upgrademongo")

type upgradeMongoCommand struct {
	envcmd.EnvCommandBase
}

func (c *upgradeMongoCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "juju-upgrademongo",
		Purpose: "Upgrade from mongo 2.4 to 3.1",
		Args:    "",
		Doc:     upgradeDoc,
	}
}

func progress(f string, a ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(f, a...))
}

func mustParseTemplate(templ string) *template.Template {
	t := template.New("").Funcs(template.FuncMap{
		"shquote": utils.ShQuote,
	})
	return template.Must(t.Parse(templ))
}

type sshParams struct {
	OldPassword string
	StatePort   int
}

// runViaSSH will run arbitrary code in the remote machine.
func runViaSSH(addr, script string, params sshParams, stderr, stdout *bytes.Buffer, verbose bool) error {
	// This is taken from cmd/juju/ssh.go there is no other clear way to set user
	userAddr := "ubuntu@" + addr
	functions := upgradeMongoInAgentConfig + upgradeMongoBinary + mongoEval + replset + noauth
	var callable string
	if verbose {
		callable = `set -xu
`
	}
	callable += functions + script
	tmpl := mustParseTemplate(callable)
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, params)
	if err != nil {
		panic(errors.Annotate(err, "template error"))
	}

	userCmd := ssh.Command(userAddr, []string{"-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", "sudo", "-n", "bash", "-c " + utils.ShQuote(buf.String())}, nil)
	userCmd.Stderr = stderr
	userCmd.Stdout = stdout
	err = userCmd.Run()
	if err != nil {
		fmt.Println("\nErr: ")
		fmt.Printf("%s", fmt.Sprintf("%s", stderr.String()))
		fmt.Println("\nOut: ")
		fmt.Printf("%s", fmt.Sprintf("%s", stdout.String()))
		return errors.Annotatef(err, "ssh command failed: See above")
	}
	fmt.Printf("%s\n\n", fmt.Sprintf("%s", stderr.String()))
	fmt.Printf("%s", fmt.Sprintf("%s", stdout.String()))
	progress("ssh command succedded: %s", "see above")
	return nil
}

// runViaJujuSSH will run arbitrary code in the remote machine.
func runViaJujuSSH(machine, script string, params sshParams, stdout, stderr *bytes.Buffer) error {
	// This is taken from cmd/juju/ssh.go there is no other clear way to set user
	functions := upgradeMongoInAgentConfig + upgradeMongoBinary + mongoEval
	script = functions + script
	tmpl := mustParseTemplate(script)
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, params)
	if err != nil {
		panic(errors.Annotate(err, "template error"))
	}
	cmd := exec.Command("juju", []string{"ssh", machine, "sudo -n bash -c " + utils.ShQuote(buf.String())}...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("\nErr:")
		fmt.Println(fmt.Sprintf("%s", stderr.String()))
		fmt.Println("\nOut:")
		fmt.Println(fmt.Sprintf("%s", stdout.String()))
		return errors.Annotatef(err, "ssh command failed: (%q)", stderr.String())
	}
	progress("ssh command succedded: %q", fmt.Sprintf("%s", stdout.String()))
	return nil
}

func addPPA(addr string) error {
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	// beware, juju-mongodb3 only works in vivid.
	addPPACommand := `apt-add-repository -y ppa:hduran-8/juju-mongodb2.6
apt-add-repository -y ppa:hduran-8/juju-mongodb3
apt-get update
apt-get install juju-mongodb2.6 juju-mongodb3
apt-get --option=Dpkg::Options::=--force-confold --option=Dpkg::options::=--force-unsafe-io --assume-yes --quiet install mongodb-clients`
	return runViaSSH(addr, addPPACommand, sshParams{}, &stdoutBuf, &stderrBuf, true)
}

func upgradeTo26(addr, password string, port int) error {
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	upgradeTo26Command := `/usr/lib/juju/bin/mongodump --ssl -u admin -p {{.OldPassword | shquote}} --port {{.StatePort}} --out ~/migrateTo26dump
echo "dumped mongo"

systemctl stop jujud-machine-0.service
echo "stoped juju"

systemctl stop juju-db.service
echo "stoped mongo"

UpgradeMongoVersion 2.6
echo "upgraded to 2.6 in conf"

UpgradeMongoBinary 2.6 only
echo "upgraded to 2.6 in systemd"

systemctl daemon-reload
echo "realoaded systemd"

/usr/lib/juju/mongo2.6/bin/mongod --dbpath /var/lib/juju/db --replSet juju --upgrade
echo "upgraded mongo 2.6"

systemctl start juju-db.service
echo "starting mongodb 2.6"

sleep 120
mongoAdminEval 'db.getSiblingDB("admin").runCommand({authSchemaUpgrade: 1 })'
echo "upgraded auth schema."

systemctl restart juju-db.service
sleep 60
`
	return runViaSSH(addr, upgradeTo26Command, sshParams{OldPassword: password, StatePort: port}, &stdoutBuf, &stderrBuf, true)
}

func upgradeTo3(addr, password string, port int) error {
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	upgradeTo3Command := `attempts=0
until [ $attempts -ge 60 ]
do
/usr/lib/juju/mongo2.6/bin/mongodump --ssl -u admin -p {{.OldPassword | shquote}} --port {{.StatePort}} --out ~/migrateTo3dump&& break
    echo "attempt $attempts"
    attempts=$[$attempts+1]
    sleep 10
done
echo "dumped for migration to 3"

systemctl stop jujud-machine-0.service
echo "stopped juju"

systemctl stop juju-db.service
echo "stopped mongo"

UpgradeMongoVersion 3.1
echo "upgrade version in agent.conf"

UpgradeMongoBinary 3 only
echo "upgrade systemctl call"

systemctl daemon-reload
echo "reload systemctl"

systemctl start juju-db.service
echo "start mongo 3 without wt"

echo "will wait"
attempts=0
until [ $attempts -ge 60 ]
do
    /usr/lib/juju/mongo3/bin/mongodump --ssl -u admin -p {{.OldPassword | shquote}} --port {{.StatePort}} --out ~/migrateToTigerDump  && break
    echo "attempt $attempts"
    attempts=$[$attempts+1]
    sleep 10
done
echo "perform migration dump"

systemctl stop juju-db.service
echo "stopped mongo"

UpgradeMongoBinary 3 storage
NoAuth
NoReplSet
cat /etc/systemd/system/juju-db.service
echo "upgrade mongo including storage"

systemctl daemon-reload
echo "reload systemctl"

mv /var/lib/juju/db /var/lib/juju/db.old
echo "move db"

mkdir /var/lib/juju/db
echo "create new db"

systemctl start juju-db.service
echo "start mongo"

sleep 60
echo "initiated peergrouper"

/usr/lib/juju/mongo3/bin/mongorestore -vvvvv --ssl --sslAllowInvalidCertificates --port {{.StatePort}} --host localhost  ~/migrateToTigerDump
echo "restored backup to wt"

systemctl stop juju-db.service
echo "stop mongo"

Auth
ReplSet
cat /etc/systemd/system/juju-db.service
systemctl daemon-reload
echo "reload systemctl"

systemctl start juju-db.service
echo "start mongo"

sleep 60
mongoAdminEval 'rs.initiate()'

systemctl start jujud-machine-0.service
echo "start juju"

systemctl start juju-db.service
echo "start mongo"
`
	return runViaSSH(addr, upgradeTo3Command, sshParams{OldPassword: password, StatePort: port}, &stdoutBuf, &stderrBuf, true)
}

func ensureRunningServices(addr string) error {
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	command := `
sleep 60
systemctl start juju-db.service
echo "start mongo"
`
	return runViaSSH(addr, command, sshParams{}, &stdoutBuf, &stderrBuf, true)
}

func (c *upgradeMongoCommand) agentConfig(addr string) (agent.Config, error) {
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	err := runViaJujuSSH(addr, "cat /var/lib/juju/agents/machine-0/agent.conf", sshParams{}, &stdoutBuf, &stderrBuf)
	if err != nil {
		return nil, errors.Annotate(err, "cannot obtain agent config")
	}
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, errors.Annotate(err, "cannot write temporary file for agent config")
	}
	defer os.Remove(f.Name())
	rawAgent := stdoutBuf.String()
	index := strings.Index(rawAgent, "# format")
	if index > 0 {
		rawAgent = rawAgent[index:]
	}

	_, err = f.Write([]byte(rawAgent))
	if err != nil {
		return nil, errors.Annotate(err, "cannot write config in temporary file")
	}
	return agent.ReadConfig(f.Name())
}

func externalIPFromStatus() (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command("juju", "status")
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", errors.Annotate(err, "cannot get juju status")
	}

	var status map[string]interface{}
	err = goyaml.Unmarshal(stdout.Bytes(), &status)
	if err != nil {
		return "", errors.Annotate(err, "cannot unmarshall status")
	}
	machines := status["machines"].(map[interface{}]interface{})
	machine0 := machines["0"].(map[interface{}]interface{})
	dnsname := machine0["dns-name"].(string)
	return dnsname, nil
}

func (c *upgradeMongoCommand) Run(ctx *cmd.Context) error {

	dnsname, err := externalIPFromStatus()
	if err != nil {
		return errors.Annotate(err, "cannot determine api addresses")
	}
	addr := dnsname
	config, err := c.agentConfig("0")
	if err != nil {
		return errors.Annotate(err, "cannot determine agent config")
	}

	err = addPPA(addr)
	if err != nil {
		return errors.Annotate(err, "cannot add mongo 2.6 and 3 ppas")
	}
	info, _ := config.StateServingInfo()
	err = upgradeTo26(addr, config.OldPassword(), info.StatePort)
	if err != nil {
		return errors.Annotate(err, "cannot upgrade to 2.6")
	}

	err = upgradeTo3(addr, config.OldPassword(), info.StatePort)
	if err != nil {
		return errors.Annotate(err, "cannot upgrade to 3")
	}
	err = ensureRunningServices(addr)
	if err != nil {
		return errors.Annotate(err, "cannot ensure services")
	}

	return nil
}
