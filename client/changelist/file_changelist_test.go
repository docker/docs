package changelist

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "test")
	defer os.RemoveAll(tmpDir)

	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(tmpDir)

	cl, err := NewFileChangelist(tmpDir)
	assert.Nil(t, err, "Error initializing fileChangelist")

	c := NewTufChange(ActionCreate, "targets", "target", "test/targ", []byte{1})
	err = cl.Add(c)
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

	err = os.Remove(tmpDir) // will error if anything left in dir
	assert.Nil(t, err, "Clear should have left the tmpDir empty")
}

func TestListOrder(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "test")
	defer os.RemoveAll(tmpDir)

	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(tmpDir)

	cl, err := NewFileChangelist(tmpDir)
	assert.Nil(t, err, "Error initializing fileChangelist")

	c1 := NewTufChange(ActionCreate, "targets", "target", "test/targ1", []byte{1})
	err = cl.Add(c1)
	assert.Nil(t, err, "Non-nil error while adding change")

	c2 := NewTufChange(ActionCreate, "targets", "target", "test/targ2", []byte{1})
	err = cl.Add(c2)
	assert.Nil(t, err, "Non-nil error while adding change")

	cs := cl.List()

	assert.Equal(t, 2, len(cs), "List should have returned exactly one item")
	assert.Equal(t, c1.Action(), cs[0].Action(), "Action mismatch")
	assert.Equal(t, c1.Scope(), cs[0].Scope(), "Scope mismatch")
	assert.Equal(t, c1.Type(), cs[0].Type(), "Type mismatch")
	assert.Equal(t, c1.Path(), cs[0].Path(), "Path mismatch")
	assert.Equal(t, c1.Content(), cs[0].Content(), "Content mismatch")

	assert.Equal(t, c2.Action(), cs[1].Action(), "Action 2 mismatch")
	assert.Equal(t, c2.Scope(), cs[1].Scope(), "Scope 2 mismatch")
	assert.Equal(t, c2.Type(), cs[1].Type(), "Type 2 mismatch")
	assert.Equal(t, c2.Path(), cs[1].Path(), "Path 2 mismatch")
	assert.Equal(t, c2.Content(), cs[1].Content(), "Content 2 mismatch")
}
