package environment

import (
	"context"
	"errors"
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/internal/environment/repo"
	"github.com/swizzleio/swiz/pkg/configutil"
	"github.com/swizzleio/swiz/pkg/preprocessor"
	"os"
	"time"
)

type EnvService struct {
	envRepo    *repo.EnvironmentRepo
	iacFactory *repo.IacRepoFactory
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
		envRepo:    envRepo,
		iacFactory: repo.NewIacRepoFactory(*config),
	}, nil
}

func (s EnvService) DeployEnvironment(ctx context.Context, enclaveName string, envDef string, envName string, dryRun bool,
	noUpdate bool) ([]*model.StackInfo, error) {
	// Get environment definition
	env, enclave, err := s.getEnvEnclave(enclaveName, envDef)
	if err != nil {
		return nil, err
	}

	// Determine no update behavior
	noUpdate = configutil.FlagOrConfig(noUpdate, enclave.EnvBehavior.NoUpdateDeploy)

	// Init param store
	ps := preprocessor.NewParamStore(enclave.Parameters)

	// Determine dependency order
	stackDeps := s.buildDependencyOrder(env.Stacks, false)

	// Create stacks
	stackInfoList := []*model.StackInfo{}
	for _, stackDep := range stackDeps {
		waitList := make([]string, len(stackDep))
		for i, stack := range stackDep {
			params := ps.GetParams(stack.Parameters)

			// Upsert stack
			stackInfo, createUpErr := s.upsertStack(ctx, env, enclave, envName, stack, params, noUpdate, dryRun)
			if createUpErr != nil {
				return nil, createUpErr
			}

			stackInfoList = append(stackInfoList, stackInfo)
			waitList[i] = stackInfo.Name
		}

		// Wait for completion
		err = s.waitForStacksComplete(ctx, enclave, envName, waitList, model.StateComplete)
		if err != nil {
			return nil, err
		}

		// Get outputs
		for _, stackName := range waitList {
			iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "")
			if iacErr != nil {
				return nil, iacErr
			}

			out, oerr := iacDeploy.GetStackOutputs(ctx, stackName)
			if oerr != nil {
				return nil, oerr
			}

			ps.SetParams(stackName, out)
		}
	}

	return stackInfoList, nil
}

func (s EnvService) DeleteEnvironment(ctx context.Context, enclaveName string, envDef string, envName string, dryRun bool,
	noOrphanDelete bool, fastDelete bool) ([]model.StackInfo, error) {

	// Get environment definition
	env, enclave, err := s.getEnvEnclave(enclaveName, envDef)
	if err != nil {
		return nil, err
	}

	// Determine flag behavior
	noOrphanDelete = configutil.FlagOrConfig(noOrphanDelete, enclave.EnvBehavior.NoOrphanDelete)
	fastDelete = configutil.FlagOrConfig(noOrphanDelete, enclave.EnvBehavior.FastDelete)

	// Determine dependency order
	stackDeps := s.buildDependencyOrder(env.Stacks, true)

	// Delete stacks
	stackInfoList := []*model.StackInfo{}
	stackDeleted := map[string]bool{}
	for _, stackDep := range stackDeps {
		waitList := make([]string, len(stackDep))
		for i, stack := range stackDep {
			iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "")
			if iacErr != nil {
				return nil, iacErr
			}

			// Generate stack name
			stack.Name = s.generateStackName(env, envName, stack.Name)
			stackInfo, deleteErr := iacDeploy.DeleteStack(ctx, stack.Name, dryRun)
			if deleteErr != nil {
				return nil, deleteErr
			}

			stackInfoList = append(stackInfoList, stackInfo)
			waitList[i] = stackInfo.Name
			stackDeleted[stack.Name] = true
		}
		if !fastDelete {
			// Wait for completion
			err = s.waitForStacksComplete(ctx, enclave, envName, waitList, model.StateDeleted)
			if err != nil {
				return nil, err
			}
		}
	}

	// Find orphaned stacks
	if !noOrphanDelete {
		// Get list of stacks
		iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "")
		if iacErr != nil {
			return nil, iacErr
		}

		stackList, listErr := iacDeploy.ListStacks(ctx, envName)
		if listErr != nil {
			return nil, listErr
		}

		waitList := []string{}
		for _, stack := range stackList {
			stackName := s.generateStackName(env, envName, stack)
			if _, ok := stackDeleted[stackName]; !ok {
				iacDeploy, iacErr = s.iacFactory.GetDeployer(*enclave, "")
				if iacErr != nil {
					return nil, iacErr
				}

				// Delete stack
				stackInfo, deleteErr := iacDeploy.DeleteStack(ctx, stackName, dryRun)
				if deleteErr != nil {
					return nil, deleteErr
				}

				stackInfoList = append(stackInfoList, stackInfo)
				waitList = append(waitList, stackInfo.Name)
			}
		}

		if !fastDelete {
			// Wait for completion
			err = s.waitForStacksComplete(ctx, enclave, envName, waitList, model.StateDeleted)
			if err != nil {
				return nil, err
			}
		}
	}

	return []model.StackInfo{}, nil
}

func (s EnvService) ListEnvironments(ctx context.Context, enclaveName string, envDef string) ([]string, error) {
	// Get environment definition
	_, enclave, err := s.getEnvEnclave(enclaveName, envDef)
	if err != nil {
		return nil, err
	}

	iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "")
	if iacErr != nil {
		return nil, iacErr
	}

	return iacDeploy.ListEnvironments(ctx)
}

func (s EnvService) GetEnvironmentInfo(ctx context.Context, enclaveName string, envDef string, envName string) (*model.EnvironmentInfo, error) {
	// Get environment definition
	_, enclave, err := s.getEnvEnclave(enclaveName, envDef)
	if err != nil {
		return nil, err
	}

	iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "")
	if iacErr != nil {
		return nil, iacErr
	}

	return iacDeploy.GetEnvironment(ctx, envName)
}

func (s EnvService) upsertStack(ctx context.Context, env *model.EnvironmentConfig, enclave *model.Enclave, envName string, stack *model.StackConfig, params map[string]string, noUpdate bool, dryRun bool) (*model.StackInfo, error) {
	var err error
	var stackInfo *model.StackInfo

	iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "")
	if iacErr != nil {
		return nil, iacErr
	}

	// Generate stack name
	stack.Name = s.generateStackName(env, envName, stack.Name)

	// Check to see if stack exists
	_, getErr := iacDeploy.GetStackInfo(ctx, stack.Name)
	if getErr != nil {
		if errors.Is(getErr, apperr.GenNotFoundError) {
			// No new stack, create one
			stackInfo, err = iacDeploy.CreateStack(ctx, stack.Name, stack.TemplateFile, params, s.generateMetadata(envName, env.EnvDefName, enclave.Name, true), dryRun)
		} else {
			return nil, getErr
		}
	} else if !noUpdate {
		// Update stack
		stackInfo, err = iacDeploy.UpdateStack(ctx, stack.Name, stack.TemplateFile, params, s.generateMetadata(envName, env.EnvDefName, enclave.Name, false), dryRun)
	} else {
		// Stacks exists and no update requested
		return nil, apperr.NewExistsError("stack", stack.Name)
	}

	return stackInfo, err
}

func (s EnvService) getEnvEnclave(enclaveName string, envDef string) (*model.EnvironmentConfig, *model.Enclave, error) {
	// Get environment definition
	env, err := s.envRepo.GetEnvironmentByDef(envDef)
	if err != nil {
		return nil, nil, err
	}

	// Get enclave
	enclaveRepo := repo.NewEnclaveRepo(*env)
	enclave, err := enclaveRepo.GetEnclave(enclaveName)
	if err != nil {
		return nil, nil, err
	}
	if enclave == nil {
		return nil, nil, apperr.NewNotFoundError("enclave", enclaveName)
	}
	return env, enclave, nil
}

func (s EnvService) generateStackName(env *model.EnvironmentConfig, envName string, stackName string) string {
	template := env.NamingScheme
	if env.NamingScheme == "" {
		template = "{{env_name:32}}-{{stack_name:32}}"
	}

	return preprocessor.ParseTemplateTokens(template, map[string]string{
		"env_name":   envName,
		"stack_name": stackName,
	})
}

func (s EnvService) generateMetadata(envName string, envDef string, enclaveName string, isCreate bool) map[string]string {
	retVal := map[string]string{
		"SwzEnv": envName,
	}

	if isCreate {
		retVal["SwzCreateDate"] = time.Now().Format(time.RFC3339)
		retVal["SwzCreateUser"] = os.Getenv("USER")
		retVal["SwzEnvDef"] = envDef
		retVal["SwzEnclave"] = enclaveName
	}

	return retVal
}

func (s EnvService) waitForStacksComplete(ctx context.Context, enclave *model.Enclave, envName string, stackList []string, state model.State) error {
	iacDeploy, iacErr := s.iacFactory.GetDeployer(*enclave, "") // This will need refactoring when we support multiple IACs
	if iacErr != nil {
		return iacErr
	}

	stopPoll := false
	for !stopPoll {
		var envErr error
		stopPoll, envErr = iacDeploy.IsEnvironmentInState(ctx, envName, stackList, []model.State{state})
		if envErr != nil {
			return envErr
		}

		if !stopPoll {
			time.Sleep(POLL_INTERVAL_SEC * time.Second)
		}
	}
	return nil
}

func (s EnvService) buildDependencyOrder(stacks map[string]*model.StackConfig, reverseOrder bool) [][]*model.StackConfig {
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

	if reverseOrder {
		// Reverse order
		for i, j := 0, len(retVal)-1; i < j; i, j = i+1, j-1 {
			retVal[i], retVal[j] = retVal[j], retVal[i]
		}
	}

	return retVal
}
