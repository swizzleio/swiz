package errtype

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrList_Error(t *testing.T) {
	errList := &ErrList{
		Errors: []error{
			errors.New("error1"),
			errors.New("error2"),
			errors.New("error3"),
		},
	}

	expectedMessage := "multiple errors: error1, error2, error3"
	assert.Equal(t, expectedMessage, errList.Error(), "Expected error message to match")

	emptyErrList := &ErrList{}
	assert.Equal(t, "multiple errors: ", emptyErrList.Error(), "Expected error message to match")
}

func TestErrList_ErrOrNil(t *testing.T) {
	errList := &ErrList{
		Errors: []error{
			errors.New("error1"),
		},
	}

	assert.Equal(t, errList, errList.ErrOrNil(), "Expected ErrOrNil to return the error list")

	emptyErrList := &ErrList{}
	assert.Nil(t, emptyErrList.ErrOrNil(), "Expected ErrOrNil to return nil")
}

func TestErrList_Add(t *testing.T) {
	errList := &ErrList{}
	err := errors.New("error1")
	errList.Add(err)

	assert.Equal(t, 1, len(errList.Errors), "Expected the length of Errors to be 1")
	assert.Equal(t, err, errList.Errors[0], "Expected the first error to match")
}
