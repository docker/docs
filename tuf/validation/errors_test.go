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

// We can unmarshal a marshalled SerializableError for all validation errors
func TestUnmarshalSerialiableErrorSuccessfully(t *testing.T) {
	validationErrors := []error{
		ErrValidation{"bad validation"},
		ErrBadHierarchy{Missing: "root", Msg: "badness"},
		ErrBadRoot{"bad root"},
		ErrBadTargets{"bad targets"},
		ErrBadSnapshot{"bad snapshot"},
	}

	for _, validError := range validationErrors {
		origS, err := NewSerializableError(validError)
		assert.NoError(t, err)
		jsonBytes, err := json.Marshal(origS)
		assert.NoError(t, err)

		var newS SerializableError
		err = json.Unmarshal(jsonBytes, &newS)
		assert.NoError(t, err)

		assert.Equal(t, validError, newS.Error)
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
