package apperr

import (
	"fmt"
)

type NotFoundErr struct {
	Subject string
	Noun    string
}

func NewNotFoundError(subject string, noun string) *NotFoundErr {
	return &NotFoundErr{
		Subject: subject,
		Noun:    noun,
	}
}

func (m NotFoundErr) Error() string {

	return fmt.Sprintf("%v %v not found", m.Subject, m.Noun)
}

func (e *NotFoundErr) Is(tgt error) bool {
	_, ok := tgt.(*NotFoundErr)
	if !ok {
		return false
	}
	return true
}

var GenNotFoundError *NotFoundErr = &NotFoundErr{}
