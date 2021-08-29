package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

func InitService() {
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
				log.Printf("%v\n", instance.InstanceId)
			}
		}
	}
}
