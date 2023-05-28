package repo

import (
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type CloudFormationRepo struct {
}

func NewCloudFormationRepo(config appconfig.AppConfig) (IacDeployer, error) {
	return &CloudFormationRepo{}, nil
}

func (r *CloudFormationRepo) CreateStack(enclave model.Enclave, name string, template string,
	params map[string]string, dryRun bool) (*model.StackInfo, error) {
	// For dry run: aws cloudformation get-template-summary --template-body file://bootstrap.yaml --profile myprofile
	return &model.StackInfo{}, nil
}

func (r *CloudFormationRepo) DeleteStack(enclave model.Enclave, name string, dryRun bool) (*model.StackInfo, error) {
	return &model.StackInfo{}, nil
}

func (r *CloudFormationRepo) UpdateStack(enclave model.Enclave, name string, template string,
	params map[string]string, dryRun bool) (*model.StackInfo, error) {

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

func (r *CloudFormationRepo) GetStackInfo(enclave model.Enclave, name string) (*model.StackInfo, error) {
	return nil, nil
}

func (r *CloudFormationRepo) GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) ListStacks(enclave model.Enclave, envName string) ([]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) ListEnvironments(enclave model.Enclave) ([]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) GetEnvironment(enclave model.Enclave, envName string) (*model.EnvironmentInfo, error) {
	return nil, nil
}

func (r *CloudFormationRepo) IsEnvironmentReady(enclave model.Enclave, envName string, stacks []string) (bool, error) {
	return true, nil
}
