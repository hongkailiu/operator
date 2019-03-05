package flexy

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	nameKey              = "Name"
	securityGroupID      = "sg-5c5ace38"
	keyNameValue         = "id_rsa_perf"
	subnetID             = "subnet-4879292d"
	kubernetesClusterKey = "KubernetesCluster"
)

// CreateInstancesOnEC2 creates instances on EC2
func CreateInstancesOnEC2(svc *ec2.EC2, name string, imageID string, count int64, instanceType ec2.InstanceType, kubernetesClusterValue string, blockDeviceMappings []ec2.BlockDeviceMapping) ([]ec2.Instance, error) {
	log.WithFields(log.Fields{"name": name}).Info("instance creating")
	req := svc.RunInstancesRequest(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String(imageID),
		InstanceType: instanceType,
		MinCount:     aws.Int64(count),
		MaxCount:     aws.Int64(count),
		SecurityGroupIds: []string{
			securityGroupID,
		},
		KeyName:             aws.String(keyNameValue),
		SubnetId:            aws.String(subnetID),
		BlockDeviceMappings: blockDeviceMappings,
		TagSpecifications: []ec2.TagSpecification{

			{
				ResourceType: ec2.ResourceTypeInstance,
				Tags: []ec2.Tag{
					{
						Key:   &nameKey,
						Value: &name,
					},
					{
						Key:   &kubernetesClusterKey,
						Value: &kubernetesClusterValue,
					},
				},
			},
		},
	})

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	return resp.Instances, nil
}

// WaitUntilRunningOnEC2 waits until the instance is running on EC2
func WaitUntilRunningOnEC2(svc *ec2.EC2, instanceID string, timeout time.Duration, host *Host) error {
	return wait.Poll(10*time.Second, timeout,
		func() (bool, error) {
			log.Debug(fmt.Sprintf("checking if the instance %s is running ...", instanceID))
			instance, err := DescribeAnInstance(svc, instanceID)
			if err != nil {
				return false, nil
			}
			log.WithFields(log.Fields{"instance.State.Code": strconv.FormatInt(*instance.State.Code, 10)}).Debug("instance code found")
			if *(instance.State.Code) == int64(16) {
				host.PublicDNS = *instance.PublicDnsName
				host.IPv4PublicIP = *instance.PublicIpAddress
				return true, nil
			}
			return false, nil
		})
}

// DescribeAnInstance describe an instance
func DescribeAnInstance(svc *ec2.EC2, instanceID string) (*ec2.Instance, error) {
	req := svc.DescribeInstancesRequest(&ec2.DescribeInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	})

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	if resp.Reservations == nil || resp.Reservations[0].Instances == nil {
		return nil, fmt.Errorf("no found instance by id: %s", instanceID)
	}

	return &resp.Reservations[0].Instances[0], nil
}
