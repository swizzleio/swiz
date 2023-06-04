package repo

import (
	"context"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
)

const defaultIacType = model.IacTypeCf

type IacDeployer interface {
	CreateStack(ctx context.Context, name string, template string, params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error)
	DeleteStack(ctx context.Context, name string, dryRun bool) (*model.StackInfo, error)
	UpdateStack(ctx context.Context, name string, template string, params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error)
	GetStackInfo(ctx context.Context, name string) (*model.StackInfo, error)
	GetStackOutputs(ctx context.Context, name string) (map[string]string, error)
	ListStacks(ctx context.Context, envName string) ([]model.StackInfo, error)
	ListEnvironments(ctx context.Context) ([]string, error)
	GetEnvironment(ctx context.Context, envName string) (*model.EnvironmentInfo, error)
	IsEnvironmentInState(ctx context.Context, envName string, stacks []string, states []model.State) (bool, []string, error)
}

type iacRepoMapping struct {
	provider string
	iacType  string
}

type IacRepoFactory struct {
	config appconfig.AppConfig
	iacMap map[iacRepoMapping]IacDeployer
}

func NewIacRepoFactory(config appconfig.AppConfig) *IacRepoFactory {
	return &IacRepoFactory{
		config: config,
		iacMap: map[iacRepoMapping]IacDeployer{},
	}
}

func (f IacRepoFactory) GetDeployer(enclave model.Enclave, providerName string, iacType string) (IacDeployer, error) {

	provider := enclave.GetProvider(providerName)
	if provider == nil {
		return nil, apperr.NewNotFoundError("provider", providerName)
	}

	if iacType == "" {
		iacType = enclave.DefaultIac
		if iacType == "" {
			iacType = defaultIacType
		}
	}

	mapping := iacRepoMapping{
		provider: providerName,
		iacType:  iacType,
	}

	if f.iacMap[mapping] == nil {
		switch iacType {
		case model.IacTypeCf:
			f.iacMap[mapping] = NewCloudFormationRepo(f.config, enclave, provider)
		case model.IacTypeDummy:
			f.iacMap[mapping] = NewDummyDeployRepo(f.config, enclave, provider)
		default:
			return nil, apperr.NewNotFoundError("iac type", iacType)
		}
	}

	return f.iacMap[mapping], nil
}
