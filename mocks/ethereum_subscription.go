package mocks

import "github.com/stretchr/testify/mock"

type MockEthereumSubscription struct {
	mock.Mock
}

func (m *MockEthereumSubscription) Unsubscribe() {
	m.Called()
}

func (m *MockEthereumSubscription) Err() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}
