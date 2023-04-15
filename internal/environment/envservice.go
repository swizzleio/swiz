package environment

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/repo"
)

type EnvService struct {
	enclaveRepo *repo.EnclaveRepo
	envRepo     *repo.EnvironmentRepo
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
	}, nil
}

func (s *EnvService) CreateEnvironment(enclaveName string, envDef string) error {

	enclave, err := s.enclaveRepo.GetEnclave(enclaveName)
	if err != nil {
		return err
	}

	env, err := s.envRepo.GetEnvironmentByDef(envDef)
	if err != nil {
		return err
	}

	for _, stack := range env.Stacks {
		fmt.Printf("Creating %v in %v with %v\n", stack.Name, enclave.Name, stack.TemplateFile)
	}

	return nil
}
