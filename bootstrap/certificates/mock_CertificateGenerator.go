package certificates

import "github.com/stretchr/testify/mock"

type MockCertificateGenerator struct {
	mock.Mock
}

func (m *MockCertificateGenerator) Generate() error {
	ret := m.Called()

	r0 := ret.Error(0)

	return r0
}
