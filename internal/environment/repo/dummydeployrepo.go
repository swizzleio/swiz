package repo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type DummyStack struct {
	Name       string
	DeployTime time.Time
}

type DummyDeployRepo struct {
	envs   map[string]*model.EnvironmentInfo
	stacks map[string]*DummyStack
}

func NewDummyDeployRepo(config appconfig.AppConfig) IacDeployer {
	return &DummyDeployRepo{
		envs:   map[string]*model.EnvironmentInfo{},
		stacks: map[string]*DummyStack{},
	}
}

func (r *DummyDeployRepo) outputParams(params map[string]string) string {
	output := ""
	for k, v := range params {
		output += fmt.Sprintf("  %s : %s\n", k, v)
	}
	return output
}

func (r *DummyDeployRepo) CreateStack(enclave model.Enclave, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {
	fmt.Printf("CreateStack: %v with template %v in enclave %v. Params:\n", name, template, enclave.Name)
	fmt.Printf(r.outputParams(params))
	fmt.Printf("Metadata:\n%v\n", r.outputParams(metadata))

	/*
		cfYaml, err := fileutil.YamlFromLocation[map[string]interface{}](template)
		if err != nil {
			return err
		}
		fmt.Println("%v\n", cfYaml)
	*/

	// Set a dummy deploy time
	timeLen, err := strconv.Atoi(params["SleepTestTime"])
	if err != nil {
		timeLen = 2
	}
	r.stacks[name] = &DummyStack{
		Name:       name,
		DeployTime: time.Now().Add(time.Duration(timeLen+4) * time.Second),
	}

	// Print the environment info
	r.envs[name] = &model.EnvironmentInfo{
		EnvironmentName: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome environment has been created",
		},
		StackDeployStatus: []model.DeployStatus{
			{
				Name:    "swiz-boot",
				State:   model.StateComplete,
				Reason:  "It's done",
				Details: "An awesome stack has been created",
			},
			{
				Name:    "swiz-sleep",
				State:   model.StateComplete,
				Reason:  "It's done",
				Details: "An awesome stack has been created",
			},
		},
	}

	return &model.StackInfo{
		Name: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
		NextAction: model.NextActionCreate,
	}, nil
}

func (r *DummyDeployRepo) DeleteStack(enclave model.Enclave, name string, dryRun bool) (*model.StackInfo, error) {
	fmt.Printf("DeleteStack: %v in enclave %v\n", name, enclave.Name)

	return &model.StackInfo{
		Name: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
		NextAction: model.NextActionDelete,
	}, nil
}

func (r *DummyDeployRepo) UpdateStack(enclave model.Enclave, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {
	fmt.Printf("UpdateStack: %v with template %v in enclave %v. Params: \n", name, template, enclave.Name)
	fmt.Printf(r.outputParams(params))
	fmt.Printf("Metadata:\n%v\n", r.outputParams(metadata))

	return &model.StackInfo{
		Name: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
		NextAction: model.NextActionUpdate,
	}, nil
}

func (r *DummyDeployRepo) GetStackInfo(enclave model.Enclave, name string) (*model.StackInfo, error) {
	fmt.Printf("GetStackInfo: %v in enclave %v\n", name, enclave.Name)

	if r.stacks[name] == nil {
		return nil, apperr.NewNotFoundError("stack", name)
	}

	stackInfo := &model.StackInfo{
		Name: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
	}

	return stackInfo, nil
}

func (r *DummyDeployRepo) GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error) {
	fmt.Printf("GetStackOutputs: %v in enclave %v\n", name, enclave.Name)

	if r.stacks[name] == nil {
		return nil, apperr.NewNotFoundError("stack", name)
	}

	outputs := map[string]string{
		"SleepTestFunctionArn": "arn:aws:lambda:us-east-1:123456789:function:SleepTestFunction",
	}
	return outputs, nil
}

func (r *DummyDeployRepo) ListStacks(enclave model.Enclave, envName string) ([]string, error) {
	fmt.Printf("ListStacks: %v in enclave %v\n", envName, enclave.Name)

	stacks := []string{
		"swizboot",
		"swizsleep",
		"swizrogue",
	}

	return stacks, nil
}

func (r *DummyDeployRepo) ListEnvironments(enclave model.Enclave) ([]string, error) {
	fmt.Printf("ListEnvironments in enclave %v\n", enclave.Name)

	envList := []string{}
	for k, _ := range r.envs {
		envList = append(envList, k)
	}

	return envList, nil
}

func (r *DummyDeployRepo) GetEnvironment(enclave model.Enclave, envName string) (*model.EnvironmentInfo, error) {
	fmt.Printf("GetEnvironment: %v in enclave %v\n", envName, enclave.Name)

	env := r.envs[envName]

	if env == nil {
		return nil, apperr.NewNotFoundError("environment", envName)
	}

	return env, nil
}

func (r *DummyDeployRepo) IsEnvironmentInState(enclave model.Enclave, envName string, stacks []string, states []model.State) (bool, error) {
	// check to see if r.deployTime[name] is past the current time
	// if it is, then set the state to complete
	// if it isn't, then set the state to in progress

	createState := false
	for _, state := range states {
		if state == model.StateComplete {
			createState = true
			break
		}
	}

	stacksComplete := 0

	if createState {
		for _, stack := range stacks {
			stackInfo := r.stacks[stack]
			if nil == stackInfo {
				return false, apperr.NewNotFoundError("stack", stack)
			}
			if stackInfo.DeployTime.Before(time.Now()) {
				stacksComplete++
			}
		}
	} else {
		stacksComplete = len(stacks)
	}

	return stacksComplete == len(stacks), nil
}
