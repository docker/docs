package changelist

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "test")
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
func TestErrorConditions(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "test")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(tmpDir)

	cl, err := NewFileChangelist(tmpDir)
	// Attempt to unmarshall a bad JSON file. Note: causes a WARN on the console.
	ioutil.WriteFile(path.Join(tmpDir, "broken_file.change"), []byte{5}, 0644)
	noItems := cl.List()
	assert.Len(t, noItems, 0, "List returns zero items on bad JSON file error")

	os.RemoveAll(tmpDir)
	err = cl.Clear("")
	assert.Error(t, err, "Clear on missing change list should return err")

	noItems = cl.List()
	assert.Len(t, noItems, 0, "List returns zero items on directory read error")
}

func TestListOrder(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "test")
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

func TestFileChangeIterator(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "test")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(tmpDir)

	cl, err := NewFileChangelist(tmpDir)
	assert.Nil(t, err, "Error initializing fileChangelist")

	it, err := cl.NewIterator()
	assert.Nil(t, err, "Error initializing iterator")
	assert.False(t, it.HasNext(), "HasNext returns false for empty ChangeList")

	c1 := NewTufChange(ActionCreate, "t1", "target1", "test/targ1", []byte{1})
	cl.Add(c1)

	c2 := NewTufChange(ActionUpdate, "t2", "target2", "test/targ2", []byte{2})
	cl.Add(c2)

	c3 := NewTufChange(ActionUpdate, "t3", "target3", "test/targ3", []byte{3})
	cl.Add(c3)

	cs := cl.List()
	index := 0
	it, err = cl.NewIterator()
	assert.Nil(t, err, "Error initializing iterator")
	for it.HasNext() {
		c, err := it.Next()
		assert.Nil(t, err, "Next err should be false")
		assert.Equal(t, c.Action(), cs[index].Action(), "Action mismatch")
		assert.Equal(t, c.Scope(), cs[index].Scope(), "Scope mismatch")
		assert.Equal(t, c.Type(), cs[index].Type(), "Type mismatch")
		assert.Equal(t, c.Path(), cs[index].Path(), "Path mismatch")
		assert.Equal(t, c.Content(), cs[index].Content(), "Content mismatch")
		index++
	}
	assert.Equal(t, index, len(cs), "Iterator produced all data in ChangeList")

	// negative test case: index out of range
	_, err = it.Next()
	assert.Error(t, err, "Next errors gracefully when exhausted")
	var iterError IteratorBoundsError
	assert.IsType(t, iterError, err, "IteratorBoundsError type")
	assert.Regexp(t, "out of bounds", err, "Message for iterator bounds error")

	// negative test case: changelist files missing
	it, err = cl.NewIterator()
	assert.Nil(t, err, "Error initializing iterator")
	for it.HasNext() {
		cl.Clear("")
		_, err := it.Next()
		assert.Error(t, err, "Next() error for missing changelist files")
	}

	// negative test case: bad JSON file to unmarshall via Next()
	cl.Clear("")
	ioutil.WriteFile(path.Join(tmpDir, "broken_file.change"), []byte{5}, 0644)
	it, err = cl.NewIterator()
	assert.Nil(t, err, "Error initializing iterator")
	for it.HasNext() {
		_, err := it.Next()
		assert.Error(t, err, "Next should indicate error for bad JSON file")
	}

	// negative test case: changelist directory does not exist
	os.RemoveAll(tmpDir)
	it, err = cl.NewIterator()
	assert.Error(t, err, "Initializing iterator without underlying file store")
}
