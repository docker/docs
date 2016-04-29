package store

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDir = "/tmp/testFilesystemStore/"

func TestNewFilesystemStore(t *testing.T) {
	_, err := NewFilesystemStore(testDir, "metadata", "json")
	require.Nil(t, err, "Initializing FilesystemStore returned unexpected error: %v", err)
	defer os.RemoveAll(testDir)

	info, err := os.Stat(path.Join(testDir, "metadata"))
	require.Nil(t, err, "Error attempting to stat metadata dir: %v", err)
	require.NotNil(t, info, "Nil FileInfo from stat on metadata dir")
	require.True(t, 0700&info.Mode() != 0, "Metadata directory is not writable")
}

func TestSetMeta(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.Nil(t, err, "Initializing FilesystemStore returned unexpected error: %v", err)
	defer os.RemoveAll(testDir)

	testContent := []byte("test data")

	err = s.SetMeta("testMeta", testContent)
	require.Nil(t, err, "SetMeta returned unexpected error: %v", err)

	content, err := ioutil.ReadFile(path.Join(testDir, "metadata", "testMeta.json"))
	require.Nil(t, err, "Error reading file: %v", err)
	require.Equal(t, testContent, content, "Content written to file was corrupted.")
}

func TestSetMetaWithNoParentDirectory(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.Nil(t, err, "Initializing FilesystemStore returned unexpected error: %v", err)
	defer os.RemoveAll(testDir)

	testContent := []byte("test data")

	err = s.SetMeta("noexist/"+"testMeta", testContent)
	require.Nil(t, err, "SetMeta returned unexpected error: %v", err)

	content, err := ioutil.ReadFile(path.Join(testDir, "metadata", "noexist/testMeta.json"))
	require.Nil(t, err, "Error reading file: %v", err)
	require.Equal(t, testContent, content, "Content written to file was corrupted.")
}

// if something already existed there, remove it first and write a new file
func TestSetMetaRemovesExistingFileBeforeWriting(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.Nil(t, err, "Initializing FilesystemStore returned unexpected error: %v", err)
	defer os.RemoveAll(testDir)

	// make a directory where we want metadata to go
	os.Mkdir(filepath.Join(testDir, "metadata", "root.json"), 0700)

	testContent := []byte("test data")
	err = s.SetMeta("root", testContent)
	require.NoError(t, err, "SetMeta returned unexpected error: %v", err)

	content, err := ioutil.ReadFile(path.Join(testDir, "metadata", "root.json"))
	require.NoError(t, err, "Error reading file: %v", err)
	require.Equal(t, testContent, content, "Content written to file was corrupted.")
}

func TestGetMeta(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.Nil(t, err, "Initializing FilesystemStore returned unexpected error: %v", err)
	defer os.RemoveAll(testDir)

	testContent := []byte("test data")

	ioutil.WriteFile(path.Join(testDir, "metadata", "testMeta.json"), testContent, 0600)

	content, err := s.GetMeta("testMeta", int64(len(testContent)))
	require.Nil(t, err, "GetMeta returned unexpected error: %v", err)

	require.Equal(t, testContent, content, "Content read from file was corrupted.")

	// Check that NoSizeLimit size reads everything
	content, err = s.GetMeta("testMeta", NoSizeLimit)
	require.Nil(t, err, "GetMeta returned unexpected error: %v", err)

	require.Equal(t, testContent, content, "Content read from file was corrupted.")

	// Check that we return only up to size bytes
	content, err = s.GetMeta("testMeta", 4)
	require.Nil(t, err, "GetMeta returned unexpected error: %v", err)

	require.Equal(t, []byte("test"), content, "Content read from file was corrupted.")
}

func TestGetSetMetadata(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.NoError(t, err, "Initializing FilesystemStore returned unexpected error", err)
	defer os.RemoveAll(testDir)

	testGetSetMeta(t, func() MetadataStore { return s })
}

func TestRemoveMetadata(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.NoError(t, err, "Initializing FilesystemStore returned unexpected error", err)
	defer os.RemoveAll(testDir)

	testRemoveMeta(t, func() MetadataStore { return s })
}

func TestRemoveAll(t *testing.T) {
	s, err := NewFilesystemStore(testDir, "metadata", "json")
	require.Nil(t, err, "Initializing FilesystemStore returned unexpected error: %v", err)
	defer os.RemoveAll(testDir)

	testContent := []byte("test data")

	// Write some files in metadata and targets dirs
	metaPath := path.Join(testDir, "metadata", "testMeta.json")
	ioutil.WriteFile(metaPath, testContent, 0600)

	// Remove all
	err = s.RemoveAll()
	require.Nil(t, err, "Removing all from FilesystemStore returned unexpected error: %v", err)

	// Test that files no longer exist
	_, err = ioutil.ReadFile(metaPath)
	require.True(t, os.IsNotExist(err))

	// Removing the empty filestore returns nil
	require.Nil(t, s.RemoveAll())
}
