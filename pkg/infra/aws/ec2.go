package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

// ListEc2 lists all the EC2 instances
func ListEc2() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon EC2 service client
	client := ec2.NewFromConfig(cfg)

	// Describe EC2 instances with paginator
	params := &ec2.DescribeInstancesInput{}

	paginator := ec2.NewDescribeInstancesPaginator(client, params)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		for _, instanceRes := range output.Reservations {
			for _, instance := range instanceRes.Instances {
				name := getTagValue("Name", instance.Tags)
				plat := getPlatformType("Os", instance)
				log.Printf("%v (%v): %v %v %v\n", name, instance.InstanceId, plat, strOrEmpty(instance.PrivateIpAddress),
					strOrEmpty(instance.PublicIpAddress))
			}
		}
	}
}

// getPlatformType returns the platform type based on an instance value or a tag value
func getPlatformType(platformTypeTag string, instance types.Instance) string {
	platType := string(instance.Platform)
	if platType == "" &&
		platformTypeTag != "" {
		return getTagValue(platformTypeTag, instance.Tags)
	}

	return platType
}
