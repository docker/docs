package enzi

import (
	"crypto"

	libkv "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
)

// LegacyAdaptor allows for setting user roles and adding public keys. It is
// very useful for migrating data from the builtin authenticator to the eNZi
// authenticator.
type LegacyAdaptor interface {
	AddUserPublicKey(user *auth.Account, label string, publicKey crypto.PublicKey) error
	SetUserRole(username string, role auth.Role) error
}

// NewLegacyAdaptor creates a limited version of the eNZi authenticator
// which is able to handle Adding public keys for users and setting roles for
// users.
func NewLegacyAdaptor(kvStore libkv.Store) LegacyAdaptor {
	return &Authenticator{
		kvStore: kvStore,
	}
}
