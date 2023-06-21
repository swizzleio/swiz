package repo

import (
	"context"
	"fmt"
	appcli "github.com/swizzleio/swiz/pkg/cli"
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
	envs    map[string]*model.EnvironmentInfo
	stacks  map[string]*DummyStack
	enclave model.Enclave
	cl      appcli.SwizClier
}

func NewDummyDeployRepo(config appconfig.AppConfig, enclave model.Enclave, provider *model.EncProvider) IacDeployer {
	return &DummyDeployRepo{
		envs:    map[string]*model.EnvironmentInfo{},
		stacks:  map[string]*DummyStack{},
		enclave: enclave,
		cl:      appcli.NewCli(nil, nil),
	}
}

func (r *DummyDeployRepo) outputParams(params map[string]string) string {
	output := ""
	for k, v := range params {
		output += fmt.Sprintf("  %s : %s\n", k, v)
	}
	return output
}

func (r *DummyDeployRepo) CreateStack(ctx context.Context, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {
	r.cl.Info("CreateStack: %v with template %v in enclave %v. Params:\n", name, template, r.enclave.Name)
	r.cl.Info("Metadata:\n%v\n", r.outputParams(metadata))

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
		StackInfo: []model.StackInfo{
			{
				Name: "swiz-boot",
				DeployStatus: model.DeployStatus{
					Name:    "swiz-boot",
					State:   model.StateComplete,
					Reason:  "It's done",
					Details: "An awesome stack has been created",
				},
				NextAction: model.NextActionNone,
				Resources:  []string{},
			},
			{
				Name: "swiz-sleep",
				DeployStatus: model.DeployStatus{
					Name:    "swiz-sleep",
					State:   model.StateComplete,
					Reason:  "It's done",
					Details: "An awesome stack has been created",
				},
				NextAction: model.NextActionNone,
				Resources:  []string{},
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
		Resources:  []string{},
	}, nil
}

func (r *DummyDeployRepo) DeleteStack(ctx context.Context, name string, dryRun bool) (*model.StackInfo, error) {
	r.cl.Info("DeleteStack: %v in enclave %v\n", name, r.enclave.Name)

	return &model.StackInfo{
		Name: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
		NextAction: model.NextActionDelete,
		Resources:  []string{},
	}, nil
}

func (r *DummyDeployRepo) UpdateStack(ctx context.Context, name string, template string,
	params map[string]string, metadata map[string]string, dryRun bool) (*model.StackInfo, error) {
	r.cl.Info("UpdateStack: %v with template %v in enclave %v. Params: \n", name, template, r.enclave.Name)
	r.cl.Info("Metadata:\n%v\n", r.outputParams(metadata))

	return &model.StackInfo{
		Name: name,
		DeployStatus: model.DeployStatus{
			Name:    name,
			State:   model.StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
		NextAction: model.NextActionUpdate,
		Resources:  []string{},
	}, nil
}

func (r *DummyDeployRepo) GetStackInfo(ctx context.Context, name string) (*model.StackInfo, error) {
	r.cl.Info("GetStackInfo: %v in enclave %v\n", name, r.enclave.Name)

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
		NextAction: model.NextActionNone,
		Resources:  []string{},
	}

	return stackInfo, nil
}

func (r *DummyDeployRepo) GetStackOutputs(ctx context.Context, name string) (map[string]string, error) {
	r.cl.Info("GetStackOutputs: %v in enclave %v\n", name, r.enclave.Name)

	if r.stacks[name] == nil {
		return nil, apperr.NewNotFoundError("stack", name)
	}

	outputs := map[string]string{
		"SleepTestFunctionArn": "arn:aws:lambda:us-east-1:123456789:function:SleepTestFunction",
	}
	return outputs, nil
}

func (r *DummyDeployRepo) ListStacks(ctx context.Context, envName string) ([]model.StackInfo, error) {
	r.cl.Info("ListStacks: %v in enclave %v\n", envName, r.enclave.Name)

	stacks := []model.StackInfo{
		{
			Name: "swizboot",
			DeployStatus: model.DeployStatus{
				Name:    "swizboot",
				State:   model.StateComplete,
				Reason:  "It's done",
				Details: "An awesome stack has been created",
			},
			NextAction: model.NextActionNone,
			Resources:  []string{},
		},
		{
			Name: "swizsleep",
			DeployStatus: model.DeployStatus{
				Name:    "swizsleep",
				State:   model.StateComplete,
				Reason:  "It's done",
				Details: "An awesome stack has been created",
			},
			NextAction: model.NextActionNone,
			Resources:  []string{},
		},
		{
			Name: "swizrogue",
			DeployStatus: model.DeployStatus{
				Name:    "swizrogue",
				State:   model.StateComplete,
				Reason:  "It's done",
				Details: "An awesome stack has been created",
			},
			NextAction: model.NextActionNone,
			Resources:  []string{},
		},
	}

	return stacks, nil
}

func (r *DummyDeployRepo) ListEnvironments(ctx context.Context) ([]string, error) {
	r.cl.Info("ListEnvironments in enclave %v\n", r.enclave.Name)

	envList := []string{}
	for k := range r.envs {
		envList = append(envList, k)
	}

	if len(envList) == 0 {
		envList = append(envList, "SomeEnvironment")
		envList = append(envList, "AnotherEnvironment")
	}

	return envList, nil
}

func (r *DummyDeployRepo) GetEnvironment(ctx context.Context, envName string) (*model.EnvironmentInfo, error) {
	r.cl.Info("GetEnvironment: %v in enclave %v\n", envName, r.enclave.Name)

	env := r.envs[envName]

	if envName == "AnotherEnvironment" {
		env = &model.EnvironmentInfo{
			EnvironmentName: envName,
			DeployStatus: model.DeployStatus{
				Name:    envName,
				State:   model.StateComplete,
				Reason:  "It's done",
				Details: "An awesome environment has been created",
			},
			StackInfo: []model.StackInfo{
				{
					Name: "swiz-boot",
					DeployStatus: model.DeployStatus{
						Name:    "swiz-boot",
						State:   model.StateComplete,
						Reason:  "It's done",
						Details: "An awesome stack has been created",
					},
					NextAction: model.NextActionNone,
					Resources:  []string{},
				},
				{
					Name: "swiz-sleep",
					DeployStatus: model.DeployStatus{
						Name:    "swiz-sleep",
						State:   model.StateComplete,
						Reason:  "It's done",
						Details: "An awesome stack has been created",
					},
					NextAction: model.NextActionNone,
					Resources:  []string{},
				},
			},
		}
	}

	if env == nil {
		return nil, apperr.NewNotFoundError("environment", envName)
	}

	return env, nil
}

func (r *DummyDeployRepo) IsEnvironmentInState(ctx context.Context, envName string, stacks []string, states []model.State) (bool, []string, error) {
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

	stackCompleteList := []string{}

	if createState {
		for _, stack := range stacks {
			stackInfo := r.stacks[stack]
			if nil == stackInfo {
				return false, stackCompleteList, apperr.NewNotFoundError("stack", stack)
			}
			if stackInfo.DeployTime.Before(time.Now()) {
				stackCompleteList = append(stackCompleteList, stack)
			}
		}
	} else {
		stackCompleteList = stacks
	}

	return len(stackCompleteList) == len(stacks), stackCompleteList, nil
}
