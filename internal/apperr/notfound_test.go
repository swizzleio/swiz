package apperr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("subject", "noun")

	assert.NotNil(t, err, "NewNotFoundError should not return nil")
	assert.Equal(t, "subject", err.Subject, "Expected subject 'subject'")
	assert.Equal(t, "noun", err.Noun, "Expected noun 'noun'")
}

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundErr{
		Subject: "subject",
		Noun:    "noun",
	}

	expectedMessage := "subject noun not found"
	assert.Equal(t, expectedMessage, err.Error(), "Expected error message to match")
}

func TestNotFoundError_Is(t *testing.T) {
	err := &NotFoundErr{
		Subject: "subject",
		Noun:    "noun",
	}

	assert.True(t, err.Is(GenNotFoundError), "Expected Is method to return true")
	assert.False(t, err.Is(nil), "Expected Is method to return false")
}
