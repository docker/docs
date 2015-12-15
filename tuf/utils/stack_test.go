package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStack(t *testing.T) {
	s := NewStack()
	assert.NotNil(t, s)
}

func TestPush(t *testing.T) {
	s := NewStack()
	s.Push("foo")
	assert.Len(t, *s, 1)
	assert.Equal(t, "foo", (*s)[0])
}

func TestPop(t *testing.T) {
	s := NewStack()
	s.Push("foo")
	i, err := s.Pop()
	assert.NoError(t, err)
	assert.Len(t, *s, 0)
	assert.IsType(t, "", i)
	assert.Equal(t, "foo", i)
}

func TestPopEmpty(t *testing.T) {
	s := NewStack()
	_, err := s.Pop()
	assert.Error(t, err)
	assert.IsType(t, ErrEmptyStack{}, err)
}

func TestPopString(t *testing.T) {
	s := NewStack()
	s.Push("foo")
	i, err := s.PopString()
	assert.NoError(t, err)
	assert.Len(t, *s, 0)
	assert.Equal(t, "foo", i)
}

func TestPopStringWrongType(t *testing.T) {
	s := NewStack()
	s.Push(123)
	_, err := s.PopString()
	assert.Error(t, err)
	assert.IsType(t, ErrBadTypeCast{}, err)
}

func TestPopStringEmpty(t *testing.T) {
	s := NewStack()
	_, err := s.PopString()
	assert.Error(t, err)
	assert.IsType(t, ErrEmptyStack{}, err)
}

func TestEmpty(t *testing.T) {
	s := NewStack()
	assert.True(t, s.Empty())
	s.Push("foo")
	assert.False(t, s.Empty())
	s.Pop()
	assert.True(t, s.Empty())
}
