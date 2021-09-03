package common

type NotFoundError struct {
	Subject string
}

func (e NotFoundError) Error() string {
	return e.Subject + " not found"
}
