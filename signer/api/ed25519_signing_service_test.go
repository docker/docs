package api_test

import (
	"testing"

	"github.com/docker/rufus/api"
	"github.com/docker/rufus/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/docker/rufus/proto"
)

type FakeKeyDB struct {
	mock.Mock
}

func (m *FakeKeyDB) CreateKey() (*pb.PublicKey, error) {
	args := m.Mock.Called()
	return args.Get(0).(*pb.PublicKey), args.Error(1)
}

func (m *FakeKeyDB) AddKey(key *keys.Key) error {
	args := m.Mock.Called()
	return args.Error(0)
}

func (m *FakeKeyDB) DeleteKey(keyInfo *pb.KeyInfo) (*pb.Void, error) {
	args := m.Mock.Called(keyInfo.ID)
	return nil, args.Error(0)
}

func (m *FakeKeyDB) KeyInfo(keyInfo *pb.KeyInfo) (*pb.PublicKey, error) {
	args := m.Mock.Called(keyInfo.ID)
	return args.Get(0).(*pb.PublicKey), args.Error(1)
}

func (m *FakeKeyDB) GetKey(keyInfo *pb.KeyInfo) (*keys.Key, error) {
	args := m.Mock.Called(keyInfo.ID)
	return args.Get(0).(*keys.Key), args.Error(1)
}

func TestDeleteKey(t *testing.T) {
	fakeKeyID := "830158bb5a4af00a3f689a8f29120f0fa7f8ae57cf00ce1fede8ae8652b5181a"

	m := FakeKeyDB{}
	sigService := api.NewEdDSASigningService(&m)

	m.On("DeleteKey", fakeKeyID).Return(nil).Once()
	_, err := sigService.DeleteKey(&pb.KeyInfo{ID: fakeKeyID})

	m.Mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestKeyInfo(t *testing.T) {
	fakeKeyID := "830158bb5a4af00a3f689a8f29120f0fa7f8ae57cf00ce1fede8ae8652b5181a"

	m := FakeKeyDB{}
	sigService := api.NewEdDSASigningService(&m)

	m.On("KeyInfo", fakeKeyID).Return(&pb.PublicKey{}, nil).Once()
	_, err := sigService.KeyInfo(&pb.KeyInfo{ID: fakeKeyID})

	m.Mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestCreateKey(t *testing.T) {
	m := FakeKeyDB{}
	sigService := api.NewEdDSASigningService(&m)

	m.On("AddKey").Return(nil).Once()

	_, err := sigService.CreateKey()

	m.Mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestSigner(t *testing.T) {
	fakeKeyID := "830158bb5a4af00a3f689a8f29120f0fa7f8ae57cf00ce1fede8ae8652b5181a"
	m := FakeKeyDB{}
	sigService := api.NewEdDSASigningService(&m)

	m.On("GetKey", fakeKeyID).Return(&keys.Key{}, nil).Once()
	_, err := sigService.Signer(&pb.KeyInfo{ID: fakeKeyID})

	m.Mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func BenchmarkCreateKey(b *testing.B) {
	sigService := api.NewEdDSASigningService(keys.NewKeyDB())
	for n := 0; n < b.N; n++ {
		_, _ = sigService.CreateKey()
	}
}
