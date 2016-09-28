package adminserver

import "github.com/stretchr/testify/mock"

type MockSyslogWriter struct {
	mock.Mock
}

// Info provides a mock function with given fields: _a0
func (_m *MockSyslogWriter) Info(_a0 ...interface{}) {
	_m.Called(_a0)
}

// Debug provides a mock function with given fields: _a0
func (_m *MockSyslogWriter) Debug(_a0 ...interface{}) {
	_m.Called(_a0)
}

// Warning provides a mock function with given fields: _a0
func (_m *MockSyslogWriter) Warning(_a0 ...interface{}) {
	_m.Called(_a0)
}

// Error provides a mock function with given fields: _a0
func (_m *MockSyslogWriter) Error(_a0 ...interface{}) {
	_m.Called(_a0)
}
