package errtype

import (
	"fmt"
	"strings"
)

type ErrList struct {
	Errors []error
}

func (e *ErrList) Error() string {
	errStr := []string{}
	for _, err := range e.Errors {
		errStr = append(errStr, err.Error())
	}

	return fmt.Sprintf("multiple errors: %v", strings.Join(errStr, ", "))
}

func (e *ErrList) ErrOrNil() error {
	if (e.Errors != nil) &&
		(len(e.Errors) > 0) {
		return e
	}

	return nil
}

func (e *ErrList) Add(err error) {
	if e.Errors == nil {
		e.Errors = []error{}
	}

	e.Errors = append(e.Errors, err)
}
