package validation

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewSerializableError errors if some random error is not returned
func TestNewSerializableErrorNonValidationError(t *testing.T) {
	_, err := NewSerializableError(fmt.Errorf("not validation error"))
	assert.Error(t, err)
}

// NewSerializableError succeeds if a validation error is passed to it
func TestNewSerializableErrorValidationError(t *testing.T) {
	vError := ErrValidation{"validation error"}
	s, err := NewSerializableError(vError)
	assert.NoError(t, err)
	assert.Equal(t, "ErrValidation", s.Name)
	assert.Equal(t, vError, s.Error)
}

// We can unmarshal a marshalled SerializableError
func TestUnmarshalSerialiableErrorSuccessfully(t *testing.T) {
	origS, err := NewSerializableError(ErrBadHierarchy{Missing: "root", Msg: "badness"})
	assert.NoError(t, err)

	b, err := json.Marshal(origS)
	assert.NoError(t, err)

	jsonBytes := [][]byte{
		b,
		[]byte(`{"Name":"ErrBadHierarchy","Error":{"Missing":"root","Msg":"badness"}}`),
	}

	for _, toUnmarshal := range jsonBytes {
		var newS SerializableError
		err = json.Unmarshal(toUnmarshal, &newS)
		assert.NoError(t, err)

		assert.Equal(t, "ErrBadHierarchy", newS.Name)
		e, ok := newS.Error.(ErrBadHierarchy)
		assert.True(t, ok)
		assert.Equal(t, "root", e.Missing)
		assert.Equal(t, "badness", e.Msg)
	}
}

// If the name is unrecognized, unmarshalling will error
func TestUnmarshalUnknownErrorName(t *testing.T) {
	origS := SerializableError{Name: "boop", Error: ErrBadRoot{"bad"}}
	b, err := json.Marshal(origS)
	assert.NoError(t, err)

	var newS SerializableError
	err = json.Unmarshal(b, &newS)
	assert.Error(t, err)
}

// If the error is unmarshallable, unmarshalling will error even if the name
// is valid
func TestUnmarshalInvalidError(t *testing.T) {
	var newS SerializableError
	err := json.Unmarshal([]byte(`{"Name": "ErrBadRoot", "Error": "meh"}`), &newS)
	assert.Error(t, err)
}

// If there is no name, unmarshalling will error even if the error is valid
func TestUnmarshalNoName(t *testing.T) {
	origS := SerializableError{Error: ErrBadRoot{"bad"}}
	b, err := json.Marshal(origS)
	assert.NoError(t, err)

	var newS SerializableError
	err = json.Unmarshal(b, &newS)
	assert.Error(t, err)
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	var newS SerializableError
	err := json.Unmarshal([]byte("{"), &newS)
	assert.Error(t, err)
}
