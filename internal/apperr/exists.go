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

func (m ExistsErr) Error() string {

	return fmt.Sprintf("%v %v already exists", m.Subject, m.Noun)
}

func (e *ExistsErr) Is(tgt error) bool {
	_, ok := tgt.(*NotFoundErr)
	if !ok {
		return false
	}
	return true
}

var GenExistsError *ExistsErr = &ExistsErr{}
