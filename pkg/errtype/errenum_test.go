package errtype

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockStringer is a mock type for fmt.Stringer to be used in ErrEnum.
type MockStringer int

func (ms MockStringer) String() string {
	return fmt.Sprintf("MockStringer %d", ms)
}

func TestNewErrEnum(t *testing.T) {
	err := NewErrEnum(MockStringer(1), "subject", "noun", "details")
	assert.NotNil(t, err, "NewErrEnum should not return nil")
	assert.Equal(t, MockStringer(1), err.Enum, "Expected Enum to match")
	assert.Equal(t, "subject", err.Subject, "Expected subject 'subject'")
	assert.Equal(t, "noun", err.Noun, "Expected noun 'noun'")
	assert.Equal(t, "details", err.details, "Expected details 'details'")
}

func TestNewSimpleErrEnum(t *testing.T) {
	err := NewSimpleErrEnum(MockStringer(1), "details")
	assert.NotNil(t, err, "NewSimpleErrEnum should not return nil")
	assert.Equal(t, MockStringer(1), err.Enum, "Expected Enum to match")
	assert.Equal(t, "details", err.details, "Expected details 'details'")
}

func TestErrEnum_Error(t *testing.T) {
	err := &ErrEnum[MockStringer]{
		Enum:    MockStringer(1),
		Subject: "subject",
		Noun:    "noun",
		details: "details",
	}
	expectedMessage := "error in subject noun of type MockStringer 1: details"
	assert.Equal(t, expectedMessage, err.Error(), "Expected error message to match")

	err = &ErrEnum[MockStringer]{
		Enum:    MockStringer(1),
		details: "details",
	}
	expectedMessage = "error of type MockStringer 1: details"
	assert.Equal(t, expectedMessage, err.Error(), "Expected error message to match")
}

func TestErrEnum_Is(t *testing.T) {
	err := &ErrEnum[MockStringer]{
		Enum:    MockStringer(1),
		Subject: "subject",
		Noun:    "noun",
		details: "details",
	}
	anotherErr := &ErrEnum[MockStringer]{
		Enum:    MockStringer(2),
		Subject: "anotherSubject",
		Noun:    "anotherNoun",
		details: "anotherDetails",
	}
	assert.True(t, err.Is(anotherErr), "Expected Is method to return true")
	assert.False(t, err.Is(nil), "Expected Is method to return false")
}
