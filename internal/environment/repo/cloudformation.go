package repo

import "github.com/swizzleio/swiz/internal/appconfig"

type CloudFormationRepo struct {
}

func NewCloudFormationRepo(config appconfig.AppConfig) IacDeployer {
	return &CloudFormationRepo{}
}

func (r *CloudFormationRepo) CreateStack(name string, template string) error {
	return nil
}

func (r *CloudFormationRepo) DeleteStack(name string) error {
	return nil
}

func (r *CloudFormationRepo) UpdateStack(name string, template string) error {
	return nil
}

func (r *CloudFormationRepo) GetStackInfo(name string) (*StackInfo, error) {
	return nil, nil
}

func (r *CloudFormationRepo) GetStackOutputs(name string) (map[string]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) ListStacks(envName string) ([]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) ListEnvironments() ([]string, error) {
	return nil, nil
}

func (r *CloudFormationRepo) GetEnvironment(envName string) (*EnvironmentInfo, error) {
	return nil, nil
}
