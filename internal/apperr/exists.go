package apperr

import (
	"fmt"
)

type ExistsErr struct {
	Subject string
	Noun    string
}

func NewExistsError(subject string, noun string) *NotFoundErr {
	return &NotFoundErr{
		Subject: subject,
		Noun:    noun,
	}
}

func (e *ExistsErr) Error() string {

	return fmt.Sprintf("%v %v already exists", e.Subject, e.Noun)
}

func (e *ExistsErr) Is(tgt error) bool {
	_, ok := tgt.(*NotFoundErr)
	return ok
}

var GenExistsError *ExistsErr = &ExistsErr{}
