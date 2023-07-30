package mockio

import "github.com/stretchr/testify/mock"

// WriteRaw is a mock implementation of io.Writer using testify/mock where the length can be overridden
type WriteRaw struct {
	mock.Mock
}

func (m *WriteRaw) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// Write is a mock implementation of io.Writer using testify/mock
type Write struct {
	mock.Mock
}

func (m *Write) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return len(p), args.Error(1)
}
