package model

type State int

const (
	StateUnknown State = iota
	StateDryRun
	StateComplete
	StateDeleted
	StateCreating
	StateUpdating
	StateDeleting
	StateRollingBack
	StateFailed
)

func (e State) GetPriority(newState State) State {
	if newState > e {
		return newState
	}
	return e
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

type NextAction int

const (
	NextActionUnknown NextAction = iota
	NextActionCreate
	NextActionUpdate
	NextActionDelete
	NextActionNone
)

func (e NextAction) String() string {
	switch e {
	case NextActionUnknown:
		return "Unknown"
	case NextActionCreate:
		return "Create"
	case NextActionUpdate:
		return "Update"
	case NextActionDelete:
		return "Delete"
	case NextActionNone:
		return "None"
	default:
		return "Unknown"
	}
}

type DeployStatus struct {
	Name    string
	State   State
	Reason  string
	Details string
}
