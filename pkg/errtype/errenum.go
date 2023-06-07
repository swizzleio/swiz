package errtype

import (
	"fmt"
)

type ErrEnum[T fmt.Stringer] struct {
	Enum    T
	Subject string
	Noun    string
	details string
}

func NewErrEnum[T fmt.Stringer](enum T, subject string, noun string, details string) *ErrEnum[T] {
	return &ErrEnum[T]{
		Enum:    enum,
		Subject: subject,
		Noun:    noun,
		details: details,
	}
}

func NewSimpleErrEnum[T fmt.Stringer](enum T, details string) *ErrEnum[T] {
	return &ErrEnum[T]{
		Enum:    enum,
		details: details,
	}
}

func (e *ErrEnum[T]) Error() string {
	prepend := ""
	if e.Subject != "" {
		prepend = fmt.Sprintf(" in %v %v", e.Subject, e.Noun)
	}
	retVal := fmt.Sprintf("error%v of type %v", prepend, e.Enum.String())
	if e.details != "" {
		retVal += fmt.Sprintf(": %v", e.details)
	}

	return retVal
}

func (e *ErrEnum[T]) Is(err error) bool {
	_, ok := err.(*ErrEnum[T])
	return ok
}
