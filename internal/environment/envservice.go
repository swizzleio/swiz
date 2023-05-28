package environment

import (
	"errors"
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/internal/environment/repo"
	"github.com/swizzleio/swiz/pkg/preprocessor"
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
		iacDeploy: repo.NewDummyDeployRepo(*config),
	}, nil
}

func (s EnvService) getStackName(envName string, stackName string) string {

	return fmt.Sprintf("%s-%s", envName, stackName)
}

func (s EnvService) DeployEnvironment(enclaveName string, envDef string, envName string, dryRun bool,
	noUpdate bool) ([]*model.StackInfo, error) {

	// Get environment definition
	env, err := s.envRepo.GetEnvironmentByDef(envDef)
	if err != nil {
		return nil, err
	}

	enclaveRepo := repo.NewEnclaveRepo(*env)
	enclave, err := enclaveRepo.GetEnclave(enclaveName)
	if err != nil {
		return nil, err
	}
	if enclave == nil {
		return nil, apperr.NewNotFoundError("enclave", enclaveName)
	}

	// Init param store
	ps := preprocessor.NewParamStore(enclave.Parameters)

	// Determine dependency order
	stackDeps := s.buildDependencyOrder(env.Stacks)

	// Create stacks
	stackInfoList := []*model.StackInfo{}
	for _, stackDep := range stackDeps {
		stackList := make([]string, len(stackDep))
		for i, stack := range stackDep {
			params := ps.GetParams(stack.Parameters)

			// Upsert stack
			stackInfo, createUpErr := s.upsertStack(enclave, stack, params, noUpdate, dryRun)
			if createUpErr != nil {
				return nil, createUpErr
			}

			stackInfoList = append(stackInfoList, stackInfo)
			stackList[i] = stack.Name
		}

		// Wait for completion
		err = s.waitForStacksComplete(enclave, envName, stackList)
		if err != nil {
			return nil, err
		}

		// Get outputs
		for _, stackName := range stackList {
			out, oerr := s.iacDeploy.GetStackOutputs(*enclave, stackName)
			if oerr != nil {
				return nil, oerr
			}

			ps.SetParams(stackName, out)
		}
	}

	return stackInfoList, nil
}

func (s EnvService) DeleteEnvironment(enclaveName string, envDef string, envName string, dryRun bool,
	noOrphanDelete bool, fastDelete bool) ([]model.StackInfo, error) {
	// TODO: Implement

	return []model.StackInfo{}, nil
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

func (s EnvService) upsertStack(enclave *model.Enclave, stack *model.StackConfig, params map[string]string, noUpdate bool, dryRun bool) (*model.StackInfo, error) {
	var err error
	var stackInfo *model.StackInfo
	// Check to see if stack exists
	_, getErr := s.iacDeploy.GetStackInfo(*enclave, stack.Name)
	if getErr != nil {
		if errors.Is(getErr, apperr.GenNotFoundError) {
			// No new stack, create one
			stackInfo, err = s.iacDeploy.CreateStack(*enclave, stack.Name, stack.TemplateFile, params, dryRun)
		} else {
			return nil, getErr
		}
	} else if !noUpdate {
		// Update stack
		stackInfo, err = s.iacDeploy.UpdateStack(*enclave, stack.Name, stack.TemplateFile, params, dryRun)
	} else {
		// Stacks exists and no update requested
		return nil, apperr.NewExistsError("stack", stack.Name)
	}

	return stackInfo, err
}

func (s EnvService) waitForStacksComplete(enclave *model.Enclave, envName string, stackList []string) error {
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
	return nil
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
