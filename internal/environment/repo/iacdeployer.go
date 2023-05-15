package repo

import (
	"github.com/swizzleio/swiz/internal/environment/model"
)

type State int

const (
	StateUnknown State = iota
	StateCreating
	StateUpdating
	StateDeleting
	StateRollingBack
	StateFailed
	StateComplete
	StateDryRun
)

type StackInfo struct {
	Name         string
	DeployStatus DeployStatus
}

type EnvironmentInfo struct {
	EnvironmentName   string
	DeployStatus      DeployStatus
	StackDeployStatus []DeployStatus
}

type DeployStatus struct {
	Name    string
	State   State
	Reason  string
	Details string
}

type IacDeployer interface {
	CreateStack(enclave model.Enclave, name string, template string, params map[string]string) error
	DeleteStack(enclave model.Enclave, name string) error
	UpdateStack(enclave model.Enclave, name string, template string, params map[string]string) error
	GetStackInfo(enclave model.Enclave, name string) (*StackInfo, error)
	GetStackOutputs(enclave model.Enclave, name string) (map[string]string, error)
	ListStacks(enclave model.Enclave, envName string) ([]string, error)
	ListEnvironments(enclave model.Enclave) ([]string, error)
	GetEnvironment(enclave model.Enclave, envName string) (*EnvironmentInfo, error)
}
