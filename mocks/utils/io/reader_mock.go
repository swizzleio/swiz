package mockio

import "github.com/stretchr/testify/mock"

// Read is a mock implementation of io.Reader using testify/mock
type Read struct {
	mock.Mock
}

func (m *Read) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}
