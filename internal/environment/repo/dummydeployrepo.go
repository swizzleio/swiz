package repo

import (
	"fmt"

	"github.com/swizzleio/swiz/internal/appconfig"
	"github.com/swizzleio/swiz/internal/environment/model"
)

type DummyDeloyRepo struct {
	envs map[string]*EnvironmentInfo
}

func NewDummyDeloyRepo(config appconfig.AppConfig) IacDeployer {
	return &DummyDeloyRepo{
		envs: map[string]*EnvironmentInfo{},
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

	r.envs[name] = &EnvironmentInfo{
		EnvironmentName: name,
		DeployStatus: DeployStatus{
			Name:    name,
			State:   StateComplete,
			Reason:  "It's done",
			Details: "An awesome environment has been created",
		},
		StackDeployStatus: []DeployStatus{
			{
				Name:    "swiz-boot",
				State:   StateComplete,
				Reason:  "It's done",
				Details: "An awesome stack has been created",
			},
			{
				Name:    "swiz-sleep",
				State:   StateComplete,
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

func (r *DummyDeloyRepo) GetStackInfo(enclave model.Enclave, name string) (*StackInfo, error) {
	fmt.Printf("GetStackInfo: %v in enclave %v\n", name, enclave.Name)

	stackInfo := &StackInfo{
		Name: name,
		DeployStatus: DeployStatus{
			Name:    name,
			State:   StateComplete,
			Reason:  "It's done",
			Details: "An awesome stack has been created",
		},
	}

	return stackInfo, nil
}

func (r *DummyDeloyRepo) GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error) {
	fmt.Printf("GetStackOutputs: %v in enclave %v\n", name, enclave.Name)

	outputs := map[string]string{
		"swiz-boot.SleepTestFunctionArn": "arn:aws:lambda:us-east-1:123456789:function:SleepTestFunction",
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

func (r *DummyDeloyRepo) GetEnvironment(enclave model.Enclave, envName string) (*EnvironmentInfo, error) {
	fmt.Printf("GetEnvironment: %v in enclave %v\n", envName, enclave.Name)

	return r.envs[envName], nil
}
