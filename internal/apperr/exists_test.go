package apperr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExistsError(t *testing.T) {
	err := NewExistsError("subject", "noun")

	assert.NotNil(t, err, "NewExistsError should not return nil")
	assert.Equal(t, "subject", err.Subject, "Expected subject 'subject'")
	assert.Equal(t, "noun", err.Noun, "Expected noun 'noun'")
}

func TestExistsErr_Error(t *testing.T) {
	err := &ExistsErr{
		Subject: "subject",
		Noun:    "noun",
	}

	expectedMessage := "subject noun already exists"
	assert.Equal(t, expectedMessage, err.Error(), "Expected error message to match")
}

func TestExistsErr_Is(t *testing.T) {
	err := &ExistsErr{
		Subject: "subject",
		Noun:    "noun",
	}

	assert.True(t, err.Is(GenExistsError), "Expected Is method to return true")
	assert.False(t, err.Is(nil), "Expected Is method to return false")
}
