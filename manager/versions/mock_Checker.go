package versions

import "github.com/stretchr/testify/mock"

import "github.com/samalba/dockerclient"

type MockChecker struct {
	mock.Mock
}

func (m *MockChecker) NewestVersion(_a0 *dockerclient.AuthConfig) (string, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(string)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *MockChecker) VersionList(_a0 *dockerclient.AuthConfig) (ManagerVersionList, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(ManagerVersionList)
	r1 := ret.Error(1)

	return r0, r1
}
