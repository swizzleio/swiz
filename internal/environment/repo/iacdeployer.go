package repo

import (
	"context"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type IacDeployer interface {
	CreateStack(ctx context.Context, name string, template string, params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error)
	DeleteStack(ctx context.Context, name string, dryRun bool) (*model.StackInfo, error)
	UpdateStack(ctx context.Context, name string, template string, params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error)
	GetStackInfo(ctx context.Context, name string) (*model.StackInfo, error)
	GetStackOutputs(ctx context.Context, name string) (map[string]string, error)
	ListStacks(ctx context.Context, envName string) ([]string, error)
	ListEnvironments(ctx context.Context) ([]string, error)
	GetEnvironment(ctx context.Context, envName string) (*model.EnvironmentInfo, error)
	IsEnvironmentInState(ctx context.Context, envName string, stacks []string, states []model.State) (bool, error)
}

type IacRepoFactory struct {
	config appconfig.AppConfig
	iacMap map[string]IacDeployer
}

func NewIacRepoFactory(config appconfig.AppConfig) *IacRepoFactory {
	return &IacRepoFactory{
		config: config,
		iacMap: map[string]IacDeployer{},
	}
}

func (f IacRepoFactory) GetDeployer(enclave model.Enclave, providerName string) (IacDeployer, error) {

	provider := enclave.GetProvider(providerName)
	if provider == nil {
		return nil, apperr.NewNotFoundError("enclave", providerName)
	}

	if f.iacMap["Dummy"] == nil {
		f.iacMap["Dummy"] = NewDummyDeployRepo(f.config, enclave)
	}

	return f.iacMap["Dummy"], nil
}
