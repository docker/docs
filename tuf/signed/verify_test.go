package signed

import (
	"errors"
	"testing"
	"time"

	"github.com/docker/go/canonical/json"
	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/tuf/data"
)

func TestRoleNoKeys(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	r, err := data.NewRole(
		"root",
		1,
		[]string{},
		nil,
		nil,
	)
	assert.NoError(t, err)
	roleWithKeys := &data.RoleWithKeys{Role: *r, Keys: map[string]data.PublicKey{k.ID(): k}}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	assert.NoError(t, err)
	s := &data.Signed{Signed: b}
	Sign(cs, s, k)
	err = Verify(s, roleWithKeys, 1)
	assert.IsType(t, ErrRoleThreshold{}, err)
}

func TestNotEnoughSigs(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	r, err := data.NewRole(
		"root",
		2,
		[]string{k.ID()},
		nil,
		nil,
	)
	assert.NoError(t, err)
	roleWithKeys := &data.RoleWithKeys{Role: *r, Keys: map[string]data.PublicKey{k.ID(): k}}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	assert.NoError(t, err)
	s := &data.Signed{Signed: b}
	Sign(cs, s, k)
	err = Verify(s, roleWithKeys, 1)
	assert.IsType(t, ErrRoleThreshold{}, err)
}

func TestMoreThanEnoughSigs(t *testing.T) {
	cs := NewEd25519()
	k1, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	k2, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	r, err := data.NewRole(
		"root",
		1,
		[]string{k1.ID(), k2.ID()},
		nil,
		nil,
	)
	assert.NoError(t, err)
	roleWithKeys := &data.RoleWithKeys{Role: *r, Keys: map[string]data.PublicKey{k1.ID(): k1, k2.ID(): k2}}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	assert.NoError(t, err)
	s := &data.Signed{Signed: b}
	Sign(cs, s, k1, k2)
	assert.Equal(t, 2, len(s.Signatures))
	err = Verify(s, roleWithKeys, 1)
	assert.NoError(t, err)
}

func TestDuplicateSigs(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	r, err := data.NewRole(
		"root",
		2,
		[]string{k.ID()},
		nil,
		nil,
	)
	assert.NoError(t, err)
	roleWithKeys := &data.RoleWithKeys{Role: *r, Keys: map[string]data.PublicKey{k.ID(): k}}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	assert.NoError(t, err)
	s := &data.Signed{Signed: b}
	Sign(cs, s, k)
	s.Signatures = append(s.Signatures, s.Signatures[0])
	err = Verify(s, roleWithKeys, 1)
	assert.IsType(t, ErrRoleThreshold{}, err)
}

func TestUnknownKeyBelowThreshold(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	unknown, err := cs.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	r, err := data.NewRole(
		"root",
		2,
		[]string{k.ID()},
		nil,
		nil,
	)
	assert.NoError(t, err)
	roleWithKeys := &data.RoleWithKeys{Role: *r, Keys: map[string]data.PublicKey{k.ID(): k, unknown.ID(): unknown}}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	assert.NoError(t, err)
	s := &data.Signed{Signed: b}
	Sign(cs, s, k, unknown)
	s.Signatures = append(s.Signatures)
	err = Verify(s, roleWithKeys, 1)
	assert.IsType(t, ErrRoleThreshold{}, err)
}

func Test(t *testing.T) {
	cryptoService := NewEd25519()
	type test struct {
		name     string
		roleData *data.RoleWithKeys
		s        *data.Signed
		ver      int
		exp      *time.Time
		typ      string
		role     string
		err      error
		mut      func(*test)
	}

	expiredTime := time.Now().Add(-time.Hour)
	minVer := 10
	tests := []test{
		{
			name: "no signatures",
			mut:  func(t *test) { t.s.Signatures = []data.Signature{} },
			err:  ErrNoSignatures,
		},
		{
			name: "unknown role",
			role: "foo",
			err:  errors.New("tuf: meta file has wrong type"),
		},
		{
			name: "exactly enough signatures",
		},
		{
			name: "wrong type",
			typ:  "bar",
			err:  ErrWrongType,
		},
		{
			name: "low version",
			ver:  minVer - 1,
			err:  ErrLowVersion{minVer - 1, minVer},
		},
		{
			role: "root",
			name: "expired",
			exp:  &expiredTime,
			err:  ErrExpired{"root", expiredTime.Format("Mon Jan 2 15:04:05 MST 2006")},
		},
	}
	for _, run := range tests {
		if run.role == "" {
			run.role = "root"
		}
		if run.ver == 0 {
			run.ver = minVer
		}
		if run.exp == nil {
			expires := time.Now().Add(time.Hour)
			run.exp = &expires
		}
		if run.typ == "" {
			run.typ = data.TUFTypes[run.role]
		}
		if run.s == nil {
			k, _ := cryptoService.Create("root", data.ED25519Key)
			r, err := data.NewRole(
				"root",
				1,
				[]string{k.ID()},
				nil,
				nil,
			)
			assert.NoError(t, err)
			run.roleData = &data.RoleWithKeys{Role: *r, Keys: map[string]data.PublicKey{k.ID(): k}}
			meta := &data.SignedCommon{Type: run.typ, Version: run.ver, Expires: *run.exp}

			b, err := json.MarshalCanonical(meta)
			assert.NoError(t, err)
			s := &data.Signed{Signed: b}
			Sign(cryptoService, s, k)
			run.s = s
		}
		if run.mut != nil {
			run.mut(&run)
		}

		err := Verify(run.s, run.roleData, minVer)
		if e, ok := run.err.(ErrExpired); ok {
			assertErrExpired(t, err, e)
		} else {
			assert.Equal(t, run.err, err)
		}
	}
}

func assertErrExpired(t *testing.T, err error, expected ErrExpired) {
	actual, ok := err.(ErrExpired)
	if !ok {
		t.Fatalf("expected err to have type ErrExpired, got %T", err)
	}
	assert.Equal(t, actual.Expired, expected.Expired)
}
