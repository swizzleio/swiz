package repo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/apperr"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type DummyDeloyRepo struct {
	envs       map[string]*model.EnvironmentInfo
	deployTime map[string]time.Time
}

func NewDummyDeloyRepo(config appconfig.AppConfig) IacDeployer {
	return &DummyDeloyRepo{
		envs:       map[string]*model.EnvironmentInfo{},
		deployTime: map[string]time.Time{},
	}
}

func (r *DummyDeloyRepo) outputParams(params map[string]string) string {
	output := ""
	for k, v := range params {
		output += fmt.Sprintf("  %s : %s\n", k, v)
	}
	return output
}

func (r *DummyDeloyRepo) CreateStack(enclave model.Enclave, name string, template string,
	params map[string]string) error {
	fmt.Printf("CreateStack: %v with template %v in enclave %v. Params:\n", name, template, enclave.Name)
	fmt.Println(r.outputParams(params))

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
	r.deployTime[name] = time.Now().Add(time.Duration(timeLen+4) * time.Second)

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

	return nil
}

func (r *DummyDeloyRepo) DeleteStack(enclave model.Enclave, name string) error {
	fmt.Printf("DeleteStack: %v in enclave %v\n", name, enclave.Name)

	return nil
}

func (r *DummyDeloyRepo) UpdateStack(enclave model.Enclave, name string, template string,
	params map[string]string) error {
	fmt.Printf("UpdateStack: %v with template %v in enclave %v. Params: \n", name, template, enclave.Name)
	fmt.Println(r.outputParams(params))

	return nil
}

func (r *DummyDeloyRepo) GetStackInfo(enclave model.Enclave, name string) (*model.StackInfo, error) {
	fmt.Printf("GetStackInfo: %v in enclave %v\n", name, enclave.Name)

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

func (r *DummyDeloyRepo) GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error) {
	fmt.Printf("GetStackOutputs: %v in enclave %v\n", name, enclave.Name)

	outputs := map[string]string{
		"SleepTestFunctionArn": "arn:aws:lambda:us-east-1:123456789:function:SleepTestFunction",
	}
	return outputs, nil
}

func (r *DummyDeloyRepo) ListStacks(enclave model.Enclave, envName string) ([]string, error) {
	fmt.Printf("ListStacks: %v in enclave %v\n", envName, enclave.Name)

	stacks := []string{
		"swiz-boot",
		"swiz-sleep",
	}

	return stacks, nil
}

func (r *DummyDeloyRepo) ListEnvironments(enclave model.Enclave) ([]string, error) {
	fmt.Printf("ListEnvironments in enclave %v\n", enclave.Name)

	envList := []string{}
	for k, _ := range r.envs {
		envList = append(envList, k)
	}

	return envList, nil
}

func (r *DummyDeloyRepo) GetEnvironment(enclave model.Enclave, envName string) (*model.EnvironmentInfo, error) {
	fmt.Printf("GetEnvironment: %v in enclave %v\n", envName, enclave.Name)

	env := r.envs[envName]

	if env == nil {
		return nil, apperr.NewNotFoundError("environment", envName)
	}

	return env, nil
}

func (r *DummyDeloyRepo) IsEnvironmentReady(enclave model.Enclave, envName string, stacks []string) (bool, error) {
	// check to see if r.deployTime[name] is past the current time
	// if it is, then set the state to complete
	// if it isn't, then set the state to in progress
	stacksComplete := 0

	for _, stack := range stacks {
		if r.deployTime[stack].Before(time.Now()) {
			stacksComplete++
		}
	}

	return stacksComplete == len(stacks), nil
}
