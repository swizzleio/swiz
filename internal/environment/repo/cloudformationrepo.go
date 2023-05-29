package repo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/pkg/drivers/awswrap"
)

type CloudFormationRepo struct {
	client *cloudformation.Client
}

func NewCloudFormationRepo(config appconfig.AppConfig, cfg awswrap.AwsConfig, enclave model.Enclave) (IacDeployer, error) {
	return &CloudFormationRepo{
		client: cloudformation.NewFromConfig(cfg.GenerateConfig()),
	}, nil
}

func (r *CloudFormationRepo) CreateStack(ctx context.Context, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {
	// For dry run: aws cloudformation get-template-summary --template-body file://bootstrap.yaml --profile myprofile

	cfParams := []types.Parameter{}
	for k, v := range params {
		cfParams = append(cfParams, types.Parameter{
			ParameterKey:   &k,
			ParameterValue: &v,
		})
	}

	tags := []types.Tag{}
	for k, v := range metadata {
		tags = append(tags, types.Tag{
			Key:   &k,
			Value: &v,
		})
	}

	_, err := r.client.CreateStack(context.TODO(), &cloudformation.CreateStackInput{
		StackName:   &name,
		TemplateURL: &template,
		//TemplateBody: &templateBody,
		Parameters: cfParams,
		Tags:       tags,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to create stack: %w", err)
	}

	return &model.StackInfo{}, nil
}

func (r *CloudFormationRepo) DeleteStack(ctx context.Context, name string, dryRun bool) (*model.StackInfo, error) {
	return &model.StackInfo{}, nil
}

func (r *CloudFormationRepo) UpdateStack(ctx context.Context, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {

	// TODO: Mimic the cloudformation deploy command:
	// https://stackoverflow.com/questions/49945531/aws-cloudformation-create-stack-vs-deploy
	// https://www.quora.com/How-does-AWS-CloudFormation-determine-whether-to-create-new-resources-or-updating-existing-ones-when-doing-a-deploy
	// https://blog.boltops.com/2017/04/07/a-simple-introduction-to-aws-cloudformation-part-4-change-sets-dry-run-mode/

	// This will use a create-change set and then execute it.
	// aws cloudformation create-change-set --stack-name Foobar --change-set-name cs-1  --template-body file://bootstrap.yaml --capabilities CAPABILITY_NAMED_IAM --profile myprofile
	// aws cloudformation describe-change-set --stack-name Foobar --change-set-name cs-1 --profile myprofile
	// aws cloudformation list-change-sets --stack-name Foobar --profile myprofile
	// aws cloudformation execute-change-set --stack-name Foobar --change-set-name cs-1 --profile myprofile

	return &model.StackInfo{}, nil
}

func (r *CloudFormationRepo) GetStackInfo(ctx context.Context, name string) (*model.StackInfo, error) {
	return nil, nil
}

func (r *CloudFormationRepo) GetStackOutputs(ctx context.Context, name string) (map[string]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) ListStacks(ctx context.Context, envName string) ([]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) ListEnvironments(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) GetEnvironment(ctx context.Context, envName string) (*model.EnvironmentInfo, error) {
	return nil, nil
}

func (r *CloudFormationRepo) IsEnvironmentInState(ctx context.Context, envName string, stacks []string, states []model.State) (bool, error) {
	return true, nil
}
