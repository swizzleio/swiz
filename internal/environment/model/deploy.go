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

func (e State) String() string {
	switch e {
	case StateUnknown:
		return "Unknown"
	case StateCreating:
		return "Creating"
	case StateUpdating:
		return "Updating"
	case StateDeleting:
		return "Deleting"
	case StateRollingBack:
		return "RollingBack"
	case StateFailed:
		return "Failed"
	case StateComplete:
		return "Complete"
	case StateDryRun:
		return "DryRun"
	default:
		return "Unknown"
	}
}
