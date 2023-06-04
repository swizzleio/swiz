package awswrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var DefaultAccountName = "dev"
var DefaultRegion = "us-east-1"

type AwsConfig struct {
	Name      string
	Profile   string
	AccountId string
	Region    string
	Endpoint  string
}

func GetDefaultConfig() (*AwsConfig, error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// Create STS client
	svc := sts.NewFromConfig(cfg)

	// Get account ID
	resp, err := svc.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	// Create IAM client and get account alias
	svcIam := iam.NewFromConfig(cfg)
	iamResp, err := svcIam.ListAccountAliases(context.Background(), &iam.ListAccountAliasesInput{})
	accountName := DefaultAccountName
	if err != nil && len(iamResp.AccountAliases) != 0 {
		accountName = iamResp.AccountAliases[0]
	}

	return &AwsConfig{
		Name:      accountName,
		AccountId: *resp.Account,
		Region:    cfg.Region,
	}, nil
}

func GetAllOrgAccounts() ([]AwsConfig, error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// Get accounts
	orgSvc := organizations.NewFromConfig(cfg)
	resp, err := orgSvc.ListAccounts(context.Background(), &organizations.ListAccountsInput{})
	if err != nil {
		return nil, err
	}

	retVal := make([]AwsConfig, len(resp.Accounts))
	for _, acct := range resp.Accounts {

		retVal = append(retVal, AwsConfig{
			Name:      *acct.Name,
			AccountId: *acct.Id,
			Region:    cfg.Region,
		})
	}

	return retVal, nil
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
				PartitionID:   "awswrap",
				URL:           a.Endpoint,
				SigningRegion: a.Region,
			}, nil
		})

		cfgOpt = config.WithEndpointResolverWithOptions(customResolver)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpt,
		config.WithSharedConfigProfile(a.Profile))
	if err != nil {
		// handle error
		log.Fatalf("creating awswrap session %v", err)
	}

	return cfg
}
