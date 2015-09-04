package signed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/jfrazelle/go/canonical/json"
)

func Test(t *testing.T) {
	cryptoService := NewEd25519()
	type test struct {
		name  string
		keys  []data.PublicKey
		roles map[string]*data.Role
		s     *data.Signed
		ver   int
		exp   *time.Time
		typ   string
		role  string
		err   error
		mut   func(*test)
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
			err:  ErrUnknownRole,
		},
		//{
		//	name: "wrong signature method",
		//	mut:  func(t *test) { t.s.Signatures[0].Method = "foo" },
		//	err:  ErrWrongMethod,
		//},
		//	{
		//		name: "signature wrong length",
		//		mut:  func(t *test) { t.s.Signatures[0].Signature = []byte{0} },
		//		err:  ErrInvalid,
		//	},
		{
			name: "key missing from role",
			mut:  func(t *test) { t.roles["root"].KeyIDs = nil },
			err:  ErrRoleThreshold{},
		},
		//	{
		//		name: "invalid signature",
		//		mut:  func(t *test) { t.s.Signatures[0].Signature = make([]byte, ed25519.SignatureSize) },
		//		err:  ErrInvalid,
		//	},
		{
			name: "not enough signatures",
			mut:  func(t *test) { t.roles["root"].Threshold = 2 },
			err:  ErrRoleThreshold{},
		},
		{
			name: "exactly enough signatures",
		},
		{
			name: "more than enough signatures",
			mut: func(t *test) {
				k, _ := cryptoService.Create("root", data.ED25519Key)
				Sign(cryptoService, t.s, k)
				t.keys = append(t.keys, k)
				t.roles["root"].KeyIDs = append(t.roles["root"].KeyIDs, k.ID())
			},
		},
		{
			name: "duplicate key id",
			mut: func(t *test) {
				t.roles["root"].Threshold = 2
				t.s.Signatures = append(t.s.Signatures, t.s.Signatures[0])
			},
			err: ErrRoleThreshold{},
		},
		{
			name: "unknown key",
			mut: func(t *test) {
				k, _ := cryptoService.Create("root", data.ED25519Key)
				Sign(cryptoService, t.s, k)
			},
		},
		{
			name: "unknown key below threshold",
			mut: func(t *test) {
				k, _ := cryptoService.Create("root", data.ED25519Key)
				Sign(cryptoService, t.s, k)
				t.roles["root"].Threshold = 2
			},
			err: ErrRoleThreshold{},
		},
		{
			name: "unknown keys in db",
			mut: func(t *test) {
				k, _ := cryptoService.Create("root", data.ED25519Key)
				Sign(cryptoService, t.s, k)
				t.keys = append(t.keys, k)
			},
		},
		{
			name: "unknown keys in db below threshold",
			mut: func(t *test) {
				k, _ := cryptoService.Create("root", data.ED25519Key)
				Sign(cryptoService, t.s, k)
				t.keys = append(t.keys, k)
				t.roles["root"].Threshold = 2
			},
			err: ErrRoleThreshold{},
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
		if run.keys == nil && run.s == nil {
			k, _ := cryptoService.Create("root", data.ED25519Key)
			meta := &data.SignedCommon{Type: run.typ, Version: run.ver, Expires: *run.exp}

			b, err := json.MarshalCanonical(meta)
			assert.NoError(t, err)
			s := &data.Signed{Signed: b}
			Sign(cryptoService, s, k)
			run.s = s
			run.keys = []data.PublicKey{k}
		}
		if run.roles == nil {
			run.roles = map[string]*data.Role{
				"root": &data.Role{
					RootRole: data.RootRole{
						KeyIDs:    []string{run.keys[0].ID()},
						Threshold: 1,
					},
					Name: "root",
				},
			}
		}
		if run.mut != nil {
			run.mut(&run)
		}

		db := keys.NewDB()
		for _, k := range run.keys {
			db.AddKey(k)
		}
		for _, r := range run.roles {
			err := db.AddRole(r)
			assert.NoError(t, err)
		}

		err := Verify(run.s, run.role, minVer, db)
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
