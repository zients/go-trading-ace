package mocks

import "github.com/stretchr/testify/mock"

type MockABI struct {
	mock.Mock
}

func (m *MockABI) UnpackIntoInterface(i interface{}, name string, data []byte) error {
	args := m.Called(i, name, data)
	return args.Error(0)
}
