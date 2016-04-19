package trustmanager

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestAddFile(t *testing.T) {
	testData := []byte("This test data should be part of the file.")
	testName := "docker.com/notary/certificate"
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+testExt)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the Add function
	err = store.Add(testName, testData)
	require.NoError(t, err)

	// Check to see if file exists
	b, err := ioutil.ReadFile(expectedFilePath)
	require.NoError(t, err)
	require.Equal(t, testData, b, "unexpected content in the file: %s", expectedFilePath)
}

func TestRemoveFile(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+testExt)

	_, err = generateRandomFile(expectedFilePath, perms)
	require.NoError(t, err)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the Remove function
	err = store.Remove(testName)
	require.NoError(t, err)

	// Check to see if file exists
	_, err = os.Stat(expectedFilePath)
	require.Error(t, err)
}

func TestListFiles(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	var expectedFilePath string
	// Create 10 randomfiles
	for i := 1; i <= 10; i++ {
		// Since we're generating this manually we need to add the extension '.'
		expectedFilename := testName + strconv.Itoa(i) + "." + testExt
		expectedFilePath = filepath.Join(tempBaseDir, expectedFilename)
		_, err = generateRandomFile(expectedFilePath, perms)
		require.NoError(t, err)
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the List function. Expect 10 files
	files := store.ListFiles()
	require.Len(t, files, 10)
}

func TestGetPath(t *testing.T) {
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: "",
		fileExt: testExt,
		perms:   perms,
	}

	firstPath := "diogomonica.com/openvpn/0xdeadbeef.crt"
	secondPath := "/docker.io/testing-dashes/@#$%^&().crt"

	result, err := store.GetPath("diogomonica.com/openvpn/0xdeadbeef")
	require.Equal(t, firstPath, result, "unexpected error from GetPath: %v", err)

	result, err = store.GetPath("/docker.io/testing-dashes/@#$%^&()")
	require.Equal(t, secondPath, result, "unexpected error from GetPath: %v", err)
}

func TestGetPathProtection(t *testing.T) {
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: "/path/to/filestore/",
		fileExt: testExt,
		perms:   perms,
	}

	// Should deny requests for paths outside the filestore
	_, err := store.GetPath("../../etc/passwd")
	require.Error(t, err)
	require.Equal(t, ErrPathOutsideStore, err)

	_, err = store.GetPath("private/../../../etc/passwd")
	require.Error(t, err)
	require.Equal(t, ErrPathOutsideStore, err)

	// Convoluted paths should work as long as they end up inside the store
	expected := "/path/to/filestore/filename.crt"
	result, err := store.GetPath("private/../../filestore/./filename")
	require.NoError(t, err)
	require.Equal(t, expected, result)

	// Repeat tests with a relative baseDir
	relStore := &SimpleFileStore{
		baseDir: "relative/file/path",
		fileExt: testExt,
		perms:   perms,
	}

	// Should deny requests for paths outside the filestore
	_, err = relStore.GetPath("../../etc/passwd")
	require.Error(t, err)
	require.Equal(t, ErrPathOutsideStore, err)
	_, err = relStore.GetPath("private/../../../etc/passwd")
	require.Error(t, err)
	require.Equal(t, ErrPathOutsideStore, err)

	// Convoluted paths should work as long as they end up inside the store
	expected = "relative/file/path/filename.crt"
	result, err = relStore.GetPath("private/../../path/./filename")
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestGetData(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+testExt)

	expectedData, err := generateRandomFile(expectedFilePath, perms)
	require.NoError(t, err)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}
	testData, err := store.Get(testName)
	require.NoError(t, err, "failed to get data from: %s", testName)
	require.Equal(t, expectedData, testData, "unexpected content for the file: %s", expectedFilePath)
}

func TestCreateDirectory(t *testing.T) {
	testDir := "fake/path/to/directory"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	dirPath := filepath.Join(tempBaseDir, testDir)

	// Call createDirectory
	createDirectory(dirPath, visible)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	require.NoError(t, err)

	// Check to see if it is a directory
	require.True(t, fi.IsDir(), "expected to be directory: %s", dirPath)

	// Check to see if the permissions match
	require.Equal(t, "drwxr-xr-x", fi.Mode().String(), "permissions are wrong for: %s. Got: %s", dirPath, fi.Mode().String())
}

func TestCreatePrivateDirectory(t *testing.T) {
	testDir := "fake/path/to/private/directory"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	dirPath := filepath.Join(tempBaseDir, testDir)

	// Call createDirectory
	createDirectory(dirPath, private)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	require.NoError(t, err)

	// Check to see if it is a directory
	require.True(t, fi.IsDir(), "expected to be directory: %s", dirPath)

	// Check to see if the permissions match
	require.Equal(t, "drwx------", fi.Mode().String(), "permissions are wrong for: %s. Got: %s", dirPath, fi.Mode().String())
}

func generateRandomFile(filePath string, perms os.FileMode) ([]byte, error) {
	rndBytes := make([]byte, 10)
	_, err := rand.Read(rndBytes)
	if err != nil {
		return nil, err
	}

	os.MkdirAll(filepath.Dir(filePath), perms)
	if err = ioutil.WriteFile(filePath, rndBytes, perms); err != nil {
		return nil, err
	}

	return rndBytes, nil
}

func TestFileStoreConsistency(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir2)

	s, err := NewPrivateSimpleFileStore(tempBaseDir, "txt")
	require.NoError(t, err)

	s2, err := NewPrivateSimpleFileStore(tempBaseDir2, ".txt")
	require.NoError(t, err)

	file1Data := make([]byte, 20)
	_, err = rand.Read(file1Data)
	require.NoError(t, err)

	file2Data := make([]byte, 20)
	_, err = rand.Read(file2Data)
	require.NoError(t, err)

	file3Data := make([]byte, 20)
	_, err = rand.Read(file3Data)
	require.NoError(t, err)

	file1Path := "file1"
	file2Path := "path/file2"
	file3Path := "long/path/file3"

	for _, s := range []Storage{s, s2} {
		s.Add(file1Path, file1Data)
		s.Add(file2Path, file2Data)
		s.Add(file3Path, file3Data)

		paths := map[string][]byte{
			file1Path: file1Data,
			file2Path: file2Data,
			file3Path: file3Data,
		}
		for _, p := range s.ListFiles() {
			_, ok := paths[p]
			require.True(t, ok, fmt.Sprintf("returned path not found: %s", p))
			d, err := s.Get(p)
			require.NoError(t, err)
			require.Equal(t, paths[p], d)
		}
	}

}
