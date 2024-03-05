package apperr

import (
	"fmt"
)

type ExistsErr struct {
	Subject string
	Noun    string
}

func NewExistsError(subject string, noun string) *ExistsErr {
	return &ExistsErr{
		Subject: subject,
		Noun:    noun,
	}
}

func (e *ExistsErr) Error() string {

	return fmt.Sprintf("%v %v already exists", e.Subject, e.Noun)
}

func (e *ExistsErr) Is(tgt error) bool {
	_, ok := tgt.(*ExistsErr)
	if !ok {
		return false
	}
	return true
}

var GenExistsError = &ExistsErr{}
