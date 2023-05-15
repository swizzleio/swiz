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
	params map[string]string) error {
	return nil
}

func (r *CloudFormationRepo) DeleteStack(enclave model.Enclave, name string) error {
	return nil
}

func (r *CloudFormationRepo) UpdateStack(enclave model.Enclave, name string, template string,
	params map[string]string) error {
	return nil
}

func (r *CloudFormationRepo) GetStackInfo(enclave model.Enclave, name string) (*StackInfo, error) {
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

func (r *CloudFormationRepo) GetEnvironment(enclave model.Enclave, envName string) (*EnvironmentInfo, error) {
	return nil, nil
}
