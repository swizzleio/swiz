package repo

import (
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type IacDeployer interface {
	CreateStack(enclave model.Enclave, name string, template string, params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error)
	DeleteStack(enclave model.Enclave, name string, dryRun bool) (*model.StackInfo, error)
	UpdateStack(enclave model.Enclave, name string, template string, params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error)
	GetStackInfo(enclave model.Enclave, name string) (*model.StackInfo, error)
	GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error)
	ListStacks(enclave model.Enclave, envName string) ([]string, error)
	ListEnvironments(enclave model.Enclave) ([]string, error)
	GetEnvironment(enclave model.Enclave, envName string) (*model.EnvironmentInfo, error)
	IsEnvironmentInState(enclave model.Enclave, envName string, stacks []string, states []model.State) (bool, error)
}

type IacRepoFactory struct {
	config appconfig.AppConfig
	iacMap map[string]IacDeployer
}

func NewIacRepoFactory(config appconfig.AppConfig) *IacRepoFactory {
	return &IacRepoFactory{
		config: config,
		iacMap: map[string]IacDeployer{
			"Dummy": NewDummyDeployRepo(config),
		},
	}
}

func (f IacRepoFactory) GetDeployer(enclave model.Enclave, providerName string) (IacDeployer, error) {

	provider := enclave.GetProvider(providerName)
	if provider == nil {
		return nil, apperr.NewNotFoundError("enclave", providerName)
	}

	return f.iacMap["Dummy"], nil
}
