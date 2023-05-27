package repo

import (
	"github.com/swizzleio/swiz/internal/environment/model"
)

type IacDeployer interface {
	CreateStack(enclave model.Enclave, name string, template string, params map[string]string) error
	DeleteStack(enclave model.Enclave, name string) error
	UpdateStack(enclave model.Enclave, name string, template string, params map[string]string) error
	GetStackInfo(enclave model.Enclave, name string) (*model.StackInfo, error)
	GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error)
	ListStacks(enclave model.Enclave, envName string) ([]string, error)
	ListEnvironments(enclave model.Enclave) ([]string, error)
	GetEnvironment(enclave model.Enclave, envName string) (*model.EnvironmentInfo, error)
	IsEnvironmentReady(enclave model.Enclave, envName string, stacks []string) (bool, error)
}
