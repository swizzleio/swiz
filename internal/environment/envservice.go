package environment

import (
	"fmt"

	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/repo"
)

type EnvService struct {
	enclaveRepo *repo.EnclaveRepo
	envRepo     *repo.EnvironmentRepo
	iacDeploy   repo.IacDeployer
}

func NewEnvService(config *appconfig.AppConfig) (*EnvService, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	envRepo, err := repo.NewEnvironmentRepo(*config)
	if err != nil {
		return nil, err
	}
	enclaveRepo := repo.NewEnclaveRepo(*config)
	return &EnvService{
		envRepo:     envRepo,
		enclaveRepo: enclaveRepo,
		iacDeploy:   repo.NewDummyDeloyRepo(*config),
	}, nil
}

func (s *EnvService) CreateEnvironment(enclaveName string, envDef string, envName string) error {

	enclave, err := s.enclaveRepo.GetEnclave(enclaveName)
	if err != nil {
		return err
	}
	if enclave == nil {
		return fmt.Errorf("enclave %s not found", enclaveName)
	}

	// Get environment definition
	env, err := s.envRepo.GetEnvironmentByDef(envDef)
	if err != nil {
		return err
	}

	// TODO: Determine dependency order
	// TODO: Check if environment already exists

	// Create stacks
	for _, stack := range env.Stacks {
		s.iacDeploy.CreateStack(*enclave, stack.Name, stack.TemplateFile)
	}

	return nil
}
