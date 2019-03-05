package flexy

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	// GCEProjectName is GCE project name
	GCEProjectName = "openshift-gce-devel"
	// GCEZone is GCE zone
	GCEZone = "us-central1-a"
	// GCEPrefix is GCE prefix
	GCEPrefix = "https://www.googleapis.com/compute/v1/projects/" + GCEProjectName
)

var (
	// DryRunnerCounter is the instance counter for dry-runner
	DryRunnerCounter = 0
	// value2 is defined for gce instances created by cucushift flexy
	value2 = `#cloud-config

# see
# http://www.marcoberube.com/archives/236
# https://cloudinit.readthedocs.org/en/latest/topics/examples.html
# https://help.ubuntu.com/community/CloudInit

disable_root: false
# preserve_hostname: true 
# Remove update_hostname as we want preserve the hostname after rebooted
cloud_init_modules:
 - migrator
 - bootcmd
 - write-files
 - growpart
 - resizefs
 - set_hostname
 - update_etc_hosts
 - rsyslog
 - users-groups
 - ssh

# sudo without tty
runcmd:
- "sed -i -e 's/^.*requiretty$/# \\0/' /etc/sudoers"
- "sed -i -e 's/^.*visiblepw$/# \\0/' /etc/sudoers"
# for GCE until https://bugzilla.redhat.com/show_bug.cgi?id=1310649 is fixed
- "curl 'http://metadata.google.internal/computeMetadata/v1/instance/attributes/sshKeys' -H 'Metadata-Flavor: Google' | sed -r -e 's/(^|,)[^\\S]*:/\\1/g' -e 's/,/\\n/g' >> /root/.ssh/authorized_keys"`
)

// CloudProvider defines functions for cloud providers
type CloudProvider interface {
	CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error
	WaitUntilRunning(host *Host, timeout time.Duration) error
}

// AWS represents cloud provider AWS
type AWS struct {
	SVC *ec2.EC2
}

// GCE represents cloud provider GCE
type GCE struct {
	SVC *compute.Service
}

// CreateAnInstance creates an instance on AWS
func (aws AWS) CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error {
	name := configParams["name"]
	imageID := configParams["imageID"]
	kubernetesClusterValue := configParams["kubernetesClusterValue"]

	instances, err := CreateInstancesOnEC2(aws.SVC, name, imageID, int64(1), role.EC2Params.InstanceType, kubernetesClusterValue, role.EC2Params.BlockDeviceMappings)
	if err != nil {
		return err
	}
	if len(instances) != 1 {
		return fmt.Errorf("NOT 1 instance: %d", len(instances))
	}
	instance := instances[0]
	log.WithFields(log.Fields{"instance.InstanceId": *instance.InstanceId}).Info("instance created")
	(*host).ID = *instance.InstanceId
	return nil
}

// WaitUntilRunning waits until the instance is running on AWS
func (aws AWS) WaitUntilRunning(host *Host, timeout time.Duration) error {
	return WaitUntilRunningOnEC2(aws.SVC, (*host).ID, timeout, host)
}

// DryRunner represents a fake cloud provider
type DryRunner struct {
}

// CreateAnInstance creates an instance for dry-runner
func (dr DryRunner) CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error {
	host.ID = fmt.Sprintf("dry-run-%s", uuid.New().String())
	host.IPv4PublicIP = fmt.Sprintf("23.23.23.%d", DryRunnerCounter+23)
	host.PublicDNS = fmt.Sprintf("dry-run-23-23-23-%d.example.com", DryRunnerCounter+23)
	DryRunnerCounter++
	return nil
}

// WaitUntilRunning waits until the instance is running for dry-runner
func (dr DryRunner) WaitUntilRunning(host *Host, timeout time.Duration) error {
	return nil
}

// CreateAnInstance creates an instance on GCE
func (g GCE) CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error {
	name := configParams["name"]
	imageID := configParams["imageID"]
	publicKeyFile := configParams["publicKeyFile"]

	log.WithFields(log.Fields{"publicKeyFile": publicKeyFile}).Info("file:")
	bytes, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return err
	}

	value1 := "root:" + string(bytes)

	var disks []*compute.AttachedDisk

	for index, d := range role.GCEParams.Disks {
		disk := &compute.AttachedDisk{
			AutoDelete: true,
			Boot:       index == 0,
			Type:       "PERSISTENT",
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskName:    fmt.Sprintf("%s-%d", name, index),
				SourceImage: GCEPrefix + "/global/images/" + imageID,
				//SourceImage: "https://www.googleapis.com/compute/v1/projects/rhel-cloud" + "/global/images/family/" + "rhel-7",
				DiskType:   GCEPrefix + "/zones/" + GCEZone + "/diskTypes/" + d.DiskType,
				DiskSizeGb: int64(d.DiskSizeGb),
			},
		}
		if index > 0 {
			disk.InitializeParams = &compute.AttachedDiskInitializeParams{
				DiskName: fmt.Sprintf("%s-%d", name, index),
			}
		}
		disks = append(disks, disk)
	}

	instance := &compute.Instance{
		Name:        name,
		Description: "compute sample instance",
		MachineType: GCEPrefix + "/zones/" + GCEZone + "/machineTypes/" + role.GCEParams.MachineType,
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "sshKeys",
					Value: &value1,
				},
				{
					Key:   "user-data",
					Value: &value2,
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: GCEPrefix + "/global/networks/default",
			},
		},
		Disks: disks,
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "1043659492591-0e9slvv63c1s4tqgsbsqlv8imruusjr9@developer.gserviceaccount.com",
				Scopes: []string{
					compute.CloudPlatformScope,
				},
			},
		},
	}
	op, err := g.SVC.Instances.Insert(GCEProjectName, GCEZone, instance).Do()
	//log.Printf("Got compute.Operation, err: %#v, %v", op, err)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"name": name, "op.TargetLink": op.TargetLink}).Info("instance created on gce")
	host.ID = name
	return nil
}

// WaitUntilRunning waits until the instance is running on GCE
func (g GCE) WaitUntilRunning(host *Host, timeout time.Duration) error {
	return wait.Poll(10*time.Second, timeout,
		func() (bool, error) {
			log.Debug(fmt.Sprintf("checking if the instance %s is running ...", host.ID))
			instance, err := g.SVC.Instances.Get(GCEProjectName, GCEZone, host.ID).Do()
			//inst.Metadata
			//log.Printf("Got compute.Instance, err: %#v, %v", instance, err)
			if err != nil {
				return false, nil
			}
			log.WithFields(log.Fields{"instance.Status": instance.Status}).Debug("instance.Status found")
			if instance.Status == "RUNNING" {
				if len(instance.NetworkInterfaces) != 1 {
					return false, fmt.Errorf("length of instance.NetworkInterfaces is %d", len(instance.NetworkInterfaces))
				}
				if len(instance.NetworkInterfaces[0].AccessConfigs) != 1 {
					return false, fmt.Errorf("length of instance.NetworkInterfaces[0].AccessConfigs is %d", len(instance.NetworkInterfaces[0].AccessConfigs))
				}
				ip := instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
				host.PublicDNS = "ocp." + ip + ".xip.io"

				host.IPv4PublicIP = ip
				return true, nil
			}
			return false, nil
		})

}
