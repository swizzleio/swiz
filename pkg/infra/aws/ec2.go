package aws

import (
	"context"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/infra/model"
	"getswizzle.io/swiz/pkg/network"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const DefaultPlatform = common.OsLinux
const DefaultPlatformTagName = "Os"

type Ec2 struct {
	client *ec2.Client
}

func NewEc2(cfg aws.Config) Ec2 {
	return Ec2{
		client: ec2.NewFromConfig(cfg),
	}
}

// ListInstances lists all the EC2 instances and returns a mapping by unique id
func (e Ec2) ListInstances() (map[string]model.TargetInstance, error) {

	// Describe EC2 instances with paginator
	instances := map[string]model.TargetInstance{}

	params := &ec2.DescribeInstancesInput{}
	paginator := ec2.NewDescribeInstancesPaginator(e.client, params)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		// Iterate across all of the returned reservations and instances
		for _, instanceRes := range output.Reservations {
			for _, instance := range instanceRes.Instances {
				name := getTagValue("Name", instance.Tags)
				plat := e.getPlatformType(DefaultPlatformTagName, DefaultPlatform, instance)
				//log.Printf("%v (%v): %v %v %v\n", name, instance.InstanceId, plat, strOrEmpty(instance.PrivateIpAddress),
				//	strOrEmpty(instance.PublicIpAddress))

				vm := model.TargetInstance{
					Id:        strOrEmpty(instance.InstanceId),
					Name:      name,
					Os:        plat,
					Endpoints: []network.Endpoint{},
				}

				private := e.getEndpoint(instance.PrivateIpAddress)
				public := e.getEndpoint(instance.PublicIpAddress)
				if private != nil {
					vm.Endpoints = append(vm.Endpoints, *private)
				}
				if public != nil {
					vm.Endpoints = append(vm.Endpoints, *public)
				}

				instances[vm.Id] = vm
			}
		}
	}

	return instances, nil
}

// getEndpoint returns an endpoint on a valid ip address
func (e Ec2) getEndpoint(ip *string) *network.Endpoint {
	if ip == nil {
		return nil
	}

	endpoint := network.NewEndpointFromHostString(*ip)

	return &endpoint
}

// getPlatformType returns the platform type based on an instance value or a tag value
func (e Ec2) getPlatformType(platformTypeTag string, defaultPlatform string, instance types.Instance) string {
	platType := string(instance.Platform)
	if platType == "" &&
		platformTypeTag != "" {
		platType = getTagValue(platformTypeTag, instance.Tags)
	}

	if platType == "" {
		platType = defaultPlatform
	}

	return platType
}
