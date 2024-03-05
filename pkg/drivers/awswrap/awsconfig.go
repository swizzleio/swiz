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

type AwsConfig struct {
	//Name string
	Profile   string
	AccountId string
	Region    string
	Endpoint  string
}

type AwsConfigManager interface {
	GetDefaultConfig() (*AwsConfig, error)
	GetAllOrgAccounts() ([]AwsConfig, error)
}

type AwsConfiger interface {
	GenerateConfig() aws.Config
}

type AwsConfigManage struct {
	cfg aws.Config
	iam Iamer
	sts Stser
	org Orger
}

func NewAwsConfigManage() (AwsConfigManager, error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &AwsConfigManage{
		cfg: cfg,
		iam: iam.NewFromConfig(cfg),
		sts: sts.NewFromConfig(cfg),
		org: organizations.NewFromConfig(cfg),
	}, nil
}

func (c AwsConfigManage) GetDefaultConfig() (*AwsConfig, error) {
	// Get account ID
	resp, err := c.sts.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	// Get account alias
	iamResp, err := c.iam.ListAccountAliases(context.Background(), &iam.ListAccountAliasesInput{})
	accountName := DefaultAccountName
	if err == nil && len(iamResp.AccountAliases) != 0 {
		accountName = iamResp.AccountAliases[0]
	}

	return &AwsConfig{
		Profile:   accountName,
		AccountId: *resp.Account,
		Region:    c.cfg.Region,
	}, nil
}

func (c AwsConfigManage) GetAllOrgAccounts() ([]AwsConfig, error) {
	// Get accounts
	resp, err := c.org.ListAccounts(context.Background(), &organizations.ListAccountsInput{})
	if err != nil {
		return nil, err
	}

	retVal := make([]AwsConfig, len(resp.Accounts))
	for i, acct := range resp.Accounts {

		retVal[i] = AwsConfig{
			Profile:   *acct.Name,
			AccountId: *acct.Id,
			Region:    c.cfg.Region,
		}
	}

	return retVal, nil
}

func NewAwsConfig(name, accountId, region string) AwsConfiger {
	return &AwsConfig{
		Profile:   name,
		AccountId: accountId,
		Region:    region,
	}
}

// GenerateConfig generates an AWS specific config
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
