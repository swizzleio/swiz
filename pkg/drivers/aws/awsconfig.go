package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AwsConfig struct {
	AccountId string
	Region    string
	Endpoint  string
}

// GenerateConfig generates an AWS config
func (a AwsConfig) GenerateConfig() aws.Config {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.

	cfgOpt := config.WithRegion(a.Region)
	if "" != a.Endpoint {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           a.Endpoint,
				SigningRegion: a.Region,
			}, nil
		})

		cfgOpt = config.WithEndpointResolverWithOptions(customResolver)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpt)
	if err != nil {
		// handle error
		log.Fatalf("creating aws session %v", err)
	}

	return cfg
}
