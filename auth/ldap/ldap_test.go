package ldap

import (
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	libkvmock "github.com/docker/libkv/store/mock"
	"github.com/docker/orca/auth"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var (
	datastoreAccounts = datastoreVersion + "/accounts"
	datastoreCerts    = datastoreVersion + "/controller_pub_certs"
	idString          = "orca"
)

func TestChangePassword(t *testing.T) {
	libkv.AddStore("mock", libkvmock.New)
	kvStore, err := libkv.NewStore(
		"mock",
		[]string{},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	s := kvStore.(*libkvmock.Mock)
	s.On("List", datastoreAccounts).Return([]*store.KVPair{}, nil)
	s.On("Put", datastoreCerts+"/"+idString, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("*store.WriteOptions")).Return(nil)

	a := NewAuthenticator(&kvStore, idString, nil, nil, &auth.LDAPSettings{})
	err = a.ChangePassword(nil, "foobar", "baz", "hunter2")
	if err == nil {
		t.Error("Change password not supported")
	}
}
