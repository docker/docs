package data

import (
	"bytes"
	rjson "encoding/json"
	"reflect"
	"testing"
	"time"

	cjson "github.com/docker/go/canonical/json"
	"github.com/stretchr/testify/require"
)

func validTimestampTemplate() *SignedTimestamp {
	return &SignedTimestamp{
		Signed: Timestamp{
			Type: "Timestamp", Version: 1, Expires: time.Now(), Meta: Files{
				CanonicalSnapshotRole: FileMeta{},
			}},
		Signatures: []Signature{
			{KeyID: "key1", Method: "method1", Signature: []byte("hello")},
		},
	}
}

func TestTimestampToSignedMarshalsSignedPortionWithCanonicalJSON(t *testing.T) {
	ts := SignedTimestamp{Signed: Timestamp{Type: "Timestamp", Version: 1, Expires: time.Now()}}
	signedCanonical, err := ts.ToSigned()
	require.NoError(t, err)

	canonicalSignedPortion, err := cjson.MarshalCanonical(ts.Signed)
	require.NoError(t, err)

	castedCanonical := rjson.RawMessage(canonicalSignedPortion)

	// don't bother testing regular JSON because it might not be different

	require.True(t, bytes.Equal(signedCanonical.Signed, castedCanonical),
		"expected %v == %v", signedCanonical.Signed, castedCanonical)
}

func TestTimestampToSignCopiesSignatures(t *testing.T) {
	ts := SignedTimestamp{
		Signed: Timestamp{Type: "Timestamp", Version: 2, Expires: time.Now()},
		Signatures: []Signature{
			{KeyID: "key1", Method: "method1", Signature: []byte("hello")},
		},
	}
	signed, err := ts.ToSigned()
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(ts.Signatures, signed.Signatures),
		"expected %v == %v", ts.Signatures, signed.Signatures)

	ts.Signatures[0].KeyID = "changed"
	require.False(t, reflect.DeepEqual(ts.Signatures, signed.Signatures),
		"expected %v != %v", ts.Signatures, signed.Signatures)
}

func TestTimestampToSignedMarshallingErrorsPropagated(t *testing.T) {
	setDefaultSerializer(errorSerializer{})
	defer setDefaultSerializer(canonicalJSON{})
	ts := SignedTimestamp{
		Signed: Timestamp{Type: "Timestamp", Version: 2, Expires: time.Now()},
	}
	_, err := ts.ToSigned()
	require.EqualError(t, err, "bad")
}

func TestTimestampMarshalJSONMarshalsSignedWithRegularJSON(t *testing.T) {
	ts := SignedTimestamp{
		Signed: Timestamp{Type: "Timestamp", Version: 1, Expires: time.Now()},
		Signatures: []Signature{
			{KeyID: "key1", Method: "method1", Signature: []byte("hello")},
			{KeyID: "key2", Method: "method2", Signature: []byte("there")},
		},
	}
	serialized, err := ts.MarshalJSON()
	require.NoError(t, err)

	signed, err := ts.ToSigned()
	require.NoError(t, err)

	// don't bother testing canonical JSON because it might not be different

	regular, err := rjson.Marshal(signed)
	require.NoError(t, err)

	require.True(t, bytes.Equal(serialized, regular),
		"expected %v != %v", serialized, regular)
}

func TestTimestampMarshalJSONMarshallingErrorsPropagated(t *testing.T) {
	setDefaultSerializer(errorSerializer{})
	defer setDefaultSerializer(canonicalJSON{})
	ts := SignedTimestamp{
		Signed: Timestamp{Type: "Timestamp", Version: 2, Expires: time.Now()},
	}
	_, err := ts.MarshalJSON()
	require.EqualError(t, err, "bad")
}

func TestTimestampFromSignedUnmarshallingErrorsPropagated(t *testing.T) {
	signed, err := validTimestampTemplate().ToSigned()
	require.NoError(t, err)

	setDefaultSerializer(errorSerializer{})
	defer setDefaultSerializer(canonicalJSON{})

	_, err = TimestampFromSigned(signed)
	require.EqualError(t, err, "bad")
}

// TimestampFromSigned succeeds if the timestamp is valid, and copies the signatures
// rather than assigns them
func TestTimestampFromSignedCopiesSignatures(t *testing.T) {
	signed, err := validTimestampTemplate().ToSigned()
	require.NoError(t, err)

	signedTimestamp, err := TimestampFromSigned(signed)
	require.NoError(t, err)

	signed.Signatures[0] = Signature{KeyID: "key3", Method: "method3", Signature: []byte("world")}

	require.Equal(t, "key3", signed.Signatures[0].KeyID)
	require.Equal(t, "key1", signedTimestamp.Signatures[0].KeyID)
}

// If the snapshot metadata is missing, the timestamp metadata fails to validate
// and thus fails to convert into a SignedTimestamp
func TestTimestampFromSignedValidatesMeta(t *testing.T) {
	ts := validTimestampTemplate()

	// no timestamp meta
	delete(ts.Signed.Meta, CanonicalSnapshotRole)
	s, err := ts.ToSigned()
	require.NoError(t, err)
	_, err = TimestampFromSigned(s)
	require.IsType(t, ErrInvalidMeta{}, err)

	// add some extra metadata to make sure it's not failing because the metadata
	// is empty
	ts.Signed.Meta[CanonicalTimestampRole] = FileMeta{}
	s, err = ts.ToSigned()
	require.NoError(t, err)
	_, err = TimestampFromSigned(s)
	require.IsType(t, ErrInvalidMeta{}, err)
}

// Type must be "Timestamp"
func TestTimestampFromSignedValidatesRoleType(t *testing.T) {
	ts := validTimestampTemplate()

	for _, invalid := range []string{" Timestamp", CanonicalSnapshotRole, "TIMESTAMP"} {
		ts.Signed.Type = invalid
		s, err := ts.ToSigned()
		require.NoError(t, err)
		_, err = TimestampFromSigned(s)
		require.IsType(t, ErrInvalidMeta{}, err)
	}

	ts = validTimestampTemplate()
	ts.Signed.Type = "Timestamp"
	s, err := ts.ToSigned()
	require.NoError(t, err)
	sTimestamp, err := TimestampFromSigned(s)
	require.NoError(t, err)
	require.Equal(t, "Timestamp", sTimestamp.Signed.Type)
}

// GetSnapshot returns the snapshot checksum, or an error if it is missing.
func TestTimestampGetSnapshot(t *testing.T) {
	ts := validTimestampTemplate()
	f, err := ts.GetSnapshot()
	require.NoError(t, err)
	require.IsType(t, &FileMeta{}, f)

	// no timestamp meta
	delete(ts.Signed.Meta, CanonicalSnapshotRole)
	f, err = ts.GetSnapshot()
	require.Error(t, err)
	require.IsType(t, ErrMissingMeta{}, err)
	require.Nil(t, f)
}
