package data

import (
	"bytes"
	rjson "encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	cjson "github.com/docker/go/canonical/json"
	"github.com/stretchr/testify/require"
)

type errorSerializer struct {
	canonicalJSON
}

func (e errorSerializer) MarshalCanonical(from interface{}) ([]byte, error) {
	return nil, fmt.Errorf("bad")
}

func TestToSignedMarshalsSignedPortionWithCanonicalJSON(t *testing.T) {
	r := SignedRoot{Signed: Root{Type: "root", Version: 2, Expires: time.Now()}}
	signedCanonical, err := r.ToSigned()
	require.NoError(t, err)

	canonicalSignedPortion, err := cjson.MarshalCanonical(r.Signed)
	require.NoError(t, err)

	castedCanonical := rjson.RawMessage(canonicalSignedPortion)

	// don't bother testing regular JSON because it might not be different

	require.True(t, bytes.Equal(signedCanonical.Signed, castedCanonical),
		"expected %v == %v", signedCanonical.Signed, castedCanonical)
}

func TestToSignCopiesSignatures(t *testing.T) {
	r := SignedRoot{
		Signed: Root{Type: "root", Version: 2, Expires: time.Now()},
		Signatures: []Signature{
			{KeyID: "key1", Method: "method1", Signature: []byte("hello")},
		},
	}
	signed, err := r.ToSigned()
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(r.Signatures, signed.Signatures),
		"expected %v == %v", r.Signatures, signed.Signatures)

	r.Signatures[0].KeyID = "changed"
	require.False(t, reflect.DeepEqual(r.Signatures, signed.Signatures),
		"expected %v != %v", r.Signatures, signed.Signatures)
}

func TestToSignedMarshallingErrorsPropagated(t *testing.T) {
	setDefaultSerializer(errorSerializer{})
	defer setDefaultSerializer(canonicalJSON{})
	r := SignedRoot{
		Signed: Root{Type: "root", Version: 2, Expires: time.Now()},
	}
	_, err := r.ToSigned()
	require.EqualError(t, err, "bad")
}

func TestMarshalJSONMarshalsSignedWithRegularJSON(t *testing.T) {
	r := SignedRoot{
		Signed: Root{Type: "root", Version: 2, Expires: time.Now()},
		Signatures: []Signature{
			{KeyID: "key1", Method: "method1", Signature: []byte("hello")},
			{KeyID: "key2", Method: "method2", Signature: []byte("there")},
		},
	}
	serialized, err := r.MarshalJSON()
	require.NoError(t, err)

	signed, err := r.ToSigned()
	require.NoError(t, err)

	// don't bother testing canonical JSON because it might not be different

	regular, err := rjson.Marshal(signed)
	require.NoError(t, err)

	require.True(t, bytes.Equal(serialized, regular),
		"expected %v != %v", serialized, regular)
}

func TestMarshalJSONMarshallingErrorsPropagated(t *testing.T) {
	setDefaultSerializer(errorSerializer{})
	defer setDefaultSerializer(canonicalJSON{})
	r := SignedRoot{
		Signed: Root{Type: "root", Version: 2, Expires: time.Now()},
	}
	_, err := r.MarshalJSON()
	require.EqualError(t, err, "bad")
}
