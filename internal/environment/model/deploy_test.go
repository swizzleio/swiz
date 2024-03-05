package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_GetPriority(t *testing.T) {
	st1 := StateDryRun
	st2 := StateUpdating

	assert.Equal(t, StateUpdating, st1.GetPriority(st2))
	assert.Equal(t, StateUpdating, st2.GetPriority(st1))
}

func TestState_String(t *testing.T) {
	st := StateUnknown
	assert.Equal(t, "Unknown", st.String())
}

func TestNextAction_String(t *testing.T) {
	act := NextActionCreate
	assert.Equal(t, "Create", act.String())
}
