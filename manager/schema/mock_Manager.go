package schema

import "github.com/stretchr/testify/mock"

import "github.com/docker/dhe-deploy/manager/schema/interfaces"

type MockManager struct {
	mock.Mock
}

func (m *MockManager) Initialize() error {
	ret := m.Called()

	r0 := ret.Error(0)

	return r0
}
func (m *MockManager) Migrate() error {
	ret := m.Called()

	r0 := ret.Error(0)

	return r0
}
func (m *MockManager) GetPropertyManager() interfaces.PropertyManager {
	ret := m.Called()

	r0 := ret.Get(0).(interfaces.PropertyManager)

	return r0
}
