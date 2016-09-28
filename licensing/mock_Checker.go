package licensing

import "github.com/stretchr/testify/mock"

import "time"

import "github.com/docker/dhe-deploy/hubconfig"

type MockChecker struct {
	mock.Mock
}

func (m *MockChecker) Initialize() error {
	ret := m.Called()

	r0 := ret.Error(0)

	return r0
}
func (m *MockChecker) IsValid() bool {
	ret := m.Called()

	r0 := ret.Get(0).(bool)

	return r0
}
func (m *MockChecker) BeginLicenseSyncing() {
	m.Called()
}
func (m *MockChecker) IsExpired() bool {
	ret := m.Called()

	r0 := ret.Get(0).(bool)

	return r0
}
func (m *MockChecker) LicensingEnforced() bool {
	ret := m.Called()

	r0 := ret.Get(0).(bool)

	return r0
}
func (m *MockChecker) LicenseTier() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *MockChecker) LicenseType() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *MockChecker) GetLicenseID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *MockChecker) Expiration() time.Time {
	ret := m.Called()

	r0 := ret.Get(0).(time.Time)

	return r0
}
func (m *MockChecker) ToggleAutoRefresh(autoRefresh bool) error {
	ret := m.Called(autoRefresh)

	r0 := ret.Error(0)

	return r0
}
func (m *MockChecker) ChangeLicenseFromId(keyID string, privateKey string) error {
	ret := m.Called(keyID, privateKey)

	r0 := ret.Error(0)

	return r0
}
func (m *MockChecker) LoadLicenseFromConfig(config *hubconfig.LicenseConfig, newLicense bool) error {
	ret := m.Called(config, newLicense)

	r0 := ret.Error(0)

	return r0
}
