package changelist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemChangelist(t *testing.T) {
	cl := memChangelist{}

	c := NewTufChange(ActionCreate, "targets", "target", "test/targ", []byte{1})

	err := cl.Add(c)
	assert.Nil(t, err, "Non-nil error while adding change")

	cs := cl.List()

	assert.Equal(t, 1, len(cs), "List should have returned exactly one item")
	assert.Equal(t, c.Action(), cs[0].Action(), "Action mismatch")
	assert.Equal(t, c.Scope(), cs[0].Scope(), "Scope mismatch")
	assert.Equal(t, c.Type(), cs[0].Type(), "Type mismatch")
	assert.Equal(t, c.Path(), cs[0].Path(), "Path mismatch")
	assert.Equal(t, c.Content(), cs[0].Content(), "Content mismatch")

	err = cl.Clear("")
	assert.Nil(t, err, "Non-nil error while clearing")

	cs = cl.List()
	assert.Equal(t, 0, len(cs), "List should be empty")
}
