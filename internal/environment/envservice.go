package environment

import (
	"errors"
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/internal/environment/repo"
	"time"
)

type EnvService struct {
	envRepo   *repo.EnvironmentRepo
	iacDeploy repo.IacDeployer
}

const (
	POLL_INTERVAL_SEC = 5
)

func NewEnvService(config *appconfig.AppConfig) (*EnvService, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	envRepo, err := repo.NewEnvironmentRepo(*config)
	if err != nil {
		return nil, err
	}
	return &EnvService{
		envRepo:   envRepo,
		iacDeploy: repo.NewDummyDeloyRepo(*config),
	}, nil
}

func (s EnvService) DeployEnvironment(enclaveName string, envDef string, envName string, dryRun bool,
	noUpdate bool) error {

	// TODO: Mimic the cloudformation deploy command:
	// https://stackoverflow.com/questions/49945531/aws-cloudformation-create-stack-vs-deploy
	// https://www.quora.com/How-does-AWS-CloudFormation-determine-whether-to-create-new-resources-or-updating-existing-ones-when-doing-a-deploy

	// Get environment definition
	env, err := s.envRepo.GetEnvironmentByDef(envDef)
	if err != nil {
		return err
	}

	enclaveRepo := repo.NewEnclaveRepo(*env)
	enclave, err := enclaveRepo.GetEnclave(enclaveName)
	if err != nil {
		return err
	}
	if enclave == nil {
		return apperr.NewNotFoundError("enclave", enclaveName)
	}

	// Init param store
	ps := NewParamStore(enclave.Parameters)

	// Check if environment already exists
	envInfo, err := s.iacDeploy.GetEnvironment(*enclave, envName)
	if err != nil && !errors.Is(err, apperr.GenNotFoundError) {
		return err
	}
	if envInfo != nil {
		return fmt.Errorf("environment %s already exists", envName)
		// TODO: Handle update if the env already exists
	}

	// Determine dependency order
	stackDeps := s.buildDependencyOrder(env.Stacks)

	// Create stacks
	for _, stackDep := range stackDeps {
		stackList := make([]string, len(stackDep))
		for i, stack := range stackDep {
			params := ps.getParams(stack.Parameters)

			createErr := s.iacDeploy.CreateStack(*enclave, stack.Name, stack.TemplateFile, params)
			if createErr != nil {
				return err
			}

			out, oerr := s.iacDeploy.GetStackOutputs(*enclave, stack.Name)
			if oerr != nil {
				return oerr
			}

			ps.setParams(stack.Name, out)
			stackList[i] = stack.Name
		}

		// Wait for completion
		stopPoll := false
		for !stopPoll {
			var envErr error
			stopPoll, envErr = s.iacDeploy.IsEnvironmentReady(*enclave, envName, stackList)
			if envErr != nil {
				return envErr
			}

			if !stopPoll {
				time.Sleep(POLL_INTERVAL_SEC * time.Second)
			}
		}
	}

	return nil
}

func (s EnvService) DeleteEnvironment(enclaveName string, envDef string, envName string, dryRun bool,
	noOrphanDelete bool, fastDelete bool) error {
	// TODO: Implement

	return nil
}

func (s EnvService) ListEnvironments(enclaveName string) ([]string, error) {
	// TODO: Implement

	return []string{}, nil
}

func (s EnvService) GetEnvironmentInfo(enclaveName string, envName string) (*model.EnvironmentInfo, error) {
	// TODO: Implement

	return &model.EnvironmentInfo{
		StackDeployStatus: []model.DeployStatus{},
	}, nil
}

func (s EnvService) buildDependencyOrder(stacks map[string]*model.StackConfig) [][]*model.StackConfig {
	// Figure out how many stack order buckets
	maxSize := 0
	for _, stack := range stacks {
		if stack.Order < 0 {
			continue
		}

		if stack.Order > maxSize {
			maxSize = stack.Order
		}
	}

	// Iterate through the stacks and determine the order. Note, 0 is a valid number.
	// Yes this could be more efficent but it's not worth it
	retVal := make([][]*model.StackConfig, maxSize+1)
	for _, stack := range stacks {
		if retVal[stack.Order] == nil {
			retVal[stack.Order] = []*model.StackConfig{}
		}

		retVal[stack.Order] = append(retVal[stack.Order], stack)
	}

	return retVal
}
