package model

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

type DeployStatus struct {
	Name    string
	State   State
	Reason  string
	Details string
}

