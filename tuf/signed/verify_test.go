package signed

import (
	"errors"
	"testing"
	"time"

	"github.com/docker/go/canonical/json"
	"github.com/stretchr/testify/require"

	"github.com/docker/notary/tuf/data"
)

func TestRoleNoKeys(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	require.NoError(t, err)
	roleWithKeys := data.BaseRole{Name: "root", Keys: data.Keys{}, Threshold: 1}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	require.NoError(t, err)
	s := &data.Signed{Signed: (*json.RawMessage)(&b)}
	Sign(cs, s, []data.PublicKey{k}, 1, nil)
	err = Verify(s, roleWithKeys, 1)
	require.IsType(t, ErrRoleThreshold{}, err)
}

func TestNotEnoughSigs(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	require.NoError(t, err)
	roleWithKeys := data.BaseRole{Name: "root", Keys: data.Keys{k.ID(): k}, Threshold: 2}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	require.NoError(t, err)
	s := &data.Signed{Signed: (*json.RawMessage)(&b)}
	Sign(cs, s, []data.PublicKey{k}, 1, nil)
	err = Verify(s, roleWithKeys, 1)
	require.IsType(t, ErrRoleThreshold{}, err)
}

func TestMoreThanEnoughSigs(t *testing.T) {
	cs := NewEd25519()
	k1, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	k2, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	roleWithKeys := data.BaseRole{Name: "root", Keys: data.Keys{k1.ID(): k1, k2.ID(): k2}, Threshold: 1}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	require.NoError(t, err)
	s := &data.Signed{Signed: (*json.RawMessage)(&b)}
	Sign(cs, s, []data.PublicKey{k1, k2}, 2, nil)
	require.Equal(t, 2, len(s.Signatures))

	err = Verify(s, roleWithKeys, 1)
	require.NoError(t, err)
}

func TestValidSigWithIncorrectKeyID(t *testing.T) {
	cs := NewEd25519()
	k1, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	roleWithKeys := data.BaseRole{Name: "root", Keys: data.Keys{"invalidIDA": k1}, Threshold: 1}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	require.NoError(t, err)
	s := &data.Signed{Signed: (*json.RawMessage)(&b)}
	Sign(cs, s, []data.PublicKey{k1}, 1, nil)
	require.Equal(t, 1, len(s.Signatures))
	s.Signatures[0].KeyID = "invalidIDA"
	err = Verify(s, roleWithKeys, 1)
	require.Error(t, err)
	require.IsType(t, ErrInvalidKeyID{}, err)
}

func TestDuplicateSigs(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	roleWithKeys := data.BaseRole{Name: "root", Keys: data.Keys{k.ID(): k}, Threshold: 2}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	require.NoError(t, err)
	s := &data.Signed{Signed: (*json.RawMessage)(&b)}
	Sign(cs, s, []data.PublicKey{k}, 1, nil)
	s.Signatures = append(s.Signatures, s.Signatures[0])
	err = Verify(s, roleWithKeys, 1)
	require.IsType(t, ErrRoleThreshold{}, err)
}

func TestUnknownKeyBelowThreshold(t *testing.T) {
	cs := NewEd25519()
	k, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	unknown, err := cs.Create("root", "", data.ED25519Key)
	require.NoError(t, err)
	roleWithKeys := data.BaseRole{Name: "root", Keys: data.Keys{k.ID(): k}, Threshold: 2}

	meta := &data.SignedCommon{Type: "Root", Version: 1, Expires: data.DefaultExpires("root")}

	b, err := json.MarshalCanonical(meta)
	require.NoError(t, err)
	s := &data.Signed{Signed: (*json.RawMessage)(&b)}
	Sign(cs, s, []data.PublicKey{k, unknown}, 2, nil)
	s.Signatures = append(s.Signatures)
	err = Verify(s, roleWithKeys, 1)
	require.IsType(t, ErrRoleThreshold{}, err)
}

func Test(t *testing.T) {
	cryptoService := NewEd25519()
	type test struct {
		name     string
		roleData data.BaseRole
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
			k, _ := cryptoService.Create("root", "", data.ED25519Key)
			run.roleData = data.BaseRole{Name: "root", Keys: data.Keys{k.ID(): k}, Threshold: 1}
			meta := &data.SignedCommon{Type: run.typ, Version: run.ver, Expires: *run.exp}

			b, err := json.MarshalCanonical(meta)
			require.NoError(t, err)
			s := &data.Signed{Signed: (*json.RawMessage)(&b)}
			Sign(cryptoService, s, []data.PublicKey{k}, 1, nil)
			run.s = s
		}
		if run.mut != nil {
			run.mut(&run)
		}

		err := Verify(run.s, run.roleData, minVer)
		if e, ok := run.err.(ErrExpired); ok {
			requireErrExpired(t, err, e)
		} else {
			require.Equal(t, run.err, err)
		}
	}
}

func requireErrExpired(t *testing.T, err error, expected ErrExpired) {
	actual, ok := err.(ErrExpired)
	if !ok {
		t.Fatalf("expected err to have type ErrExpired, got %T", err)
	}
	require.Equal(t, actual.Expired, expected.Expired)
}
