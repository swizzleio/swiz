package errtype

import (
	"fmt"
)

type ErrEnum[T fmt.Stringer] struct {
	Enum    T
	details string
}

func (e *ErrEnum[T]) Error() string {
	retVal := fmt.Sprintf("error of type %v", e.Enum.String())
	if e.details != "" {
		retVal += fmt.Sprintf(": %v", e.details)
	}

	return retVal
}

func (e *ErrEnum[T]) Is(err error) bool {
	_, ok := err.(*ErrEnum[T])
	return ok
}
