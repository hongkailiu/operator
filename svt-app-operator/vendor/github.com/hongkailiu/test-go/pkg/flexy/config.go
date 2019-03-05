package flexy

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"gopkg.in/yaml.v2"
)

// Config is the configuration for flexy
type Config struct {
	MasterGroup    []Host
	ETCDGroup      []Host
	NodeGroup      []Host
	LBGroup        []Host
	GlusterFSGroup []Host
	OCVars         map[string]string
}

//
const (
	OCPRoleMaster    = "master"
	OCPRoleInfra     = "infra"
	OCPRoleCompute   = "compute"
	OCPRoleGlusterFS = "glusterfs"

	CloudProviderAWS       = "aws"
	CloudProviderGCE       = "gce"
	CloudProviderDryRunner = "dry-runner"
)

var (
	// OCPRoles is the list of OCP roles
	OCPRoles = []string{OCPRoleMaster, OCPRoleInfra, OCPRoleCompute, OCPRoleGlusterFS}
)

// ValidateOCPRole validates an OCP role
func ValidateOCPRole(role string) error {
	for _, r := range OCPRoles {
		if role == r {
			return nil
		}
	}
	return fmt.Errorf("invalid OCP role: %s", role)
}

// Host represent a host
type Host struct {
	ID                  string
	PublicDNS           string
	OCNodeGroupName     string
	OCMasterSchedulable bool
	IPv4PublicIP        string
}

// OCPClusterConfig represents OCP cluster configuration
type OCPClusterConfig struct {
	AllInOne               bool              `yaml:"allInOne"`
	CloudProvider          string            `yaml:"cloudProvider"`
	OCPRoles               []OCPRole         `yaml:"ocpRoles"`
	KubernetesClusterValue string            `yaml:"kubernetesClusterValue"`
	ImageID                string            `yaml:"imageID"`
	InstancePrefix         string            `yaml:"instancePrefix"`
	OpenshiftAnsibleVar    map[string]string `yaml:"openshiftAnsibleVar"`
	PublicKeyFile          string            `yaml:"publicKeyFile"`
	GCSKeyfile             string            `yaml:"gcsKeyfile"`
}

// OCPRole represents an OCP role
type OCPRole struct {
	Name      string
	Size      int
	EC2Params EC2Params `yaml:"ec2Params"`
	GCEParams GCEParams `yaml:"gceParams"`
}

// EC2Params represents parameters for EC2 instances
type EC2Params struct {
	InstanceType        ec2.InstanceType         `yaml:"instanceType"`
	BlockDeviceMappings []ec2.BlockDeviceMapping `yaml:"blockDeviceMappings"`
}

// GCEParams represents parameters for GCE instances
type GCEParams struct {
	MachineType string       `yaml:"machineType"`
	Disks       []DiskParams `yaml:"disks"`
}

// DiskParams represents disk parameters
type DiskParams struct {
	DiskType   string `yaml:"diskType"`
	DiskSizeGb int    `yaml:"diskSizeGb"`
}

// LoadOCPClusterConfig loads OCP cluster configuration
func LoadOCPClusterConfig(file string, config *OCPClusterConfig) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	return nil
}
