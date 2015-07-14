package api_test

import (
	"testing"

	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/signer/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/docker/notary/proto"
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

func (m *FakeKeyDB) DeleteKey(keyID *pb.KeyID) (*pb.Void, error) {
	args := m.Mock.Called(keyID.ID)
	return nil, args.Error(0)
}

func (m *FakeKeyDB) KeyInfo(keyID *pb.KeyID) (*pb.PublicKey, error) {
	args := m.Mock.Called(keyID.ID)
	return args.Get(0).(*pb.PublicKey), args.Error(1)
}

func (m *FakeKeyDB) GetKey(keyID *pb.KeyID) (*keys.Key, error) {
	args := m.Mock.Called(keyID.ID)
	return args.Get(0).(*keys.Key), args.Error(1)
}

func TestDeleteKey(t *testing.T) {
	fakeKeyID := "830158bb5a4af00a3f689a8f29120f0fa7f8ae57cf00ce1fede8ae8652b5181a"

	m := FakeKeyDB{}
	sigService := api.NewEdDSASigningService(&m)

	m.On("DeleteKey", fakeKeyID).Return(nil).Once()
	_, err := sigService.DeleteKey(&pb.KeyID{ID: fakeKeyID})

	m.Mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestKeyInfo(t *testing.T) {
	fakeKeyID := "830158bb5a4af00a3f689a8f29120f0fa7f8ae57cf00ce1fede8ae8652b5181a"

	m := FakeKeyDB{}
	sigService := api.NewEdDSASigningService(&m)

	m.On("KeyInfo", fakeKeyID).Return(&pb.PublicKey{}, nil).Once()
	_, err := sigService.KeyInfo(&pb.KeyID{ID: fakeKeyID})

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
	_, err := sigService.Signer(&pb.KeyInfo{KeyID: &pb.KeyID{ID: fakeKeyID}})

	m.Mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func BenchmarkCreateKey(b *testing.B) {
	sigService := api.NewEdDSASigningService(keys.NewKeyDB())
	for n := 0; n < b.N; n++ {
		_, _ = sigService.CreateKey()
	}
}
