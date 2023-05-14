package repo

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
	CreateStack(name string, template string) error
	DeleteStack(name string) error
	UpdateStack(name string, template string) error
	GetStackInfo(name string) (*StackInfo, error)
	GetStackOutputs(name string) (map[string]string, error)
	ListStacks(envName string) ([]string, error)
	ListEnvironments() ([]string, error)
	GetEnvironment(envName string) (*EnvironmentInfo, error)
}
