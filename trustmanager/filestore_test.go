package trustmanager

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	// Check to see if file exists
	b, err := ioutil.ReadFile(expectedFilePath)
	assert.NoError(t, err)
	assert.Equal(t, testData, b, "unexpected content in the file: %s", expectedFilePath)
}

func TestRemoveFile(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+testExt)

	_, err = generateRandomFile(expectedFilePath, perms)
	assert.NoError(t, err)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the Remove function
	err = store.Remove(testName)
	assert.NoError(t, err)

	// Check to see if file exists
	_, err = os.Stat(expectedFilePath)
	assert.Error(t, err)
}

func TestRemoveDir(t *testing.T) {
	testName := "docker.com/diogomonica/"
	testExt := ".key"
	perms := os.FileMode(0700)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+testExt)

	_, err = generateRandomFile(expectedFilePath, perms)
	assert.NoError(t, err)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the RemoveDir function
	err = store.RemoveDir(testName)
	assert.NoError(t, err)

	expectedDirectory := filepath.Dir(expectedFilePath)
	// Check to see if file exists
	_, err = os.Stat(expectedDirectory)
	assert.Error(t, err)
}

func TestListFiles(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	var expectedFilePath string
	// Create 10 randomfiles
	for i := 1; i <= 10; i++ {
		// Since we're generating this manually we need to add the extension '.'
		expectedFilename := testName + strconv.Itoa(i) + "." + testExt
		expectedFilePath = filepath.Join(tempBaseDir, expectedFilename)
		_, err = generateRandomFile(expectedFilePath, perms)
		assert.NoError(t, err)
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the List function. Expect 10 files
	files := store.ListFiles()
	assert.Len(t, files, 10)
}

func TestListDir(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	var expectedFilePath string
	// Create 10 randomfiles
	for i := 1; i <= 10; i++ {
		// Since we're generating this manually we need to add the extension '.'
		fileName := fmt.Sprintf("%s-%s.%s", testName, strconv.Itoa(i), testExt)
		expectedFilePath = filepath.Join(tempBaseDir, fileName)
		_, err = generateRandomFile(expectedFilePath, perms)
		assert.NoError(t, err)
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the ListDir function
	files := store.ListDir("docker.com/")
	assert.Len(t, files, 10)
	files = store.ListDir("docker.com/notary")
	assert.Len(t, files, 10)
	files = store.ListDir("fakedocker.com/")
	assert.Len(t, files, 0)
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
	assert.Equal(t, firstPath, result, "unexpected error from GetPath: %v", err)

	result, err = store.GetPath("/docker.io/testing-dashes/@#$%^&()")
	assert.Equal(t, secondPath, result, "unexpected error from GetPath: %v", err)
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
	assert.Error(t, err)
	assert.Equal(t, ErrPathOutsideStore, err)

	_, err = store.GetPath("private/../../../etc/passwd")
	assert.Error(t, err)
	assert.Equal(t, ErrPathOutsideStore, err)

	// Convoluted paths should work as long as they end up inside the store
	expected := "/path/to/filestore/filename.crt"
	result, err := store.GetPath("private/../../filestore/./filename")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Repeat tests with a relative baseDir
	relStore := &SimpleFileStore{
		baseDir: "relative/file/path",
		fileExt: testExt,
		perms:   perms,
	}

	// Should deny requests for paths outside the filestore
	_, err = relStore.GetPath("../../etc/passwd")
	assert.Error(t, err)
	assert.Equal(t, ErrPathOutsideStore, err)
	_, err = relStore.GetPath("private/../../../etc/passwd")
	assert.Error(t, err)
	assert.Equal(t, ErrPathOutsideStore, err)

	// Convoluted paths should work as long as they end up inside the store
	expected = "relative/file/path/filename.crt"
	result, err = relStore.GetPath("private/../../path/./filename")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetData(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := ".crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+testExt)

	expectedData, err := generateRandomFile(expectedFilePath, perms)
	assert.NoError(t, err)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}
	testData, err := store.Get(testName)
	assert.NoError(t, err, "failed to get data from: %s", testName)
	assert.Equal(t, expectedData, testData, "unexpected content for the file: %s", expectedFilePath)
}

func TestCreateDirectory(t *testing.T) {
	testDir := "fake/path/to/directory"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	dirPath := filepath.Join(tempBaseDir, testDir)

	// Call createDirectory
	CreateDirectory(dirPath)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	assert.NoError(t, err)

	// Check to see if it is a directory
	assert.True(t, fi.IsDir(), "expected to be directory: %s", dirPath)

	// Check to see if the permissions match
	assert.Equal(t, "drwxr-xr-x", fi.Mode().String(), "permissions are wrong for: %s. Got: %s", dirPath, fi.Mode().String())
}

func TestCreatePrivateDirectory(t *testing.T) {
	testDir := "fake/path/to/private/directory"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	dirPath := filepath.Join(tempBaseDir, testDir)

	// Call createDirectory
	CreatePrivateDirectory(dirPath)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	assert.NoError(t, err)

	// Check to see if it is a directory
	assert.True(t, fi.IsDir(), "expected to be directory: %s", dirPath)

	// Check to see if the permissions match
	assert.Equal(t, "drwx------", fi.Mode().String(), "permissions are wrong for: %s. Got: %s", dirPath, fi.Mode().String())
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
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	s, err := NewPrivateSimpleFileStore(tempBaseDir, "txt")
	assert.NoError(t, err)

	file1Data := make([]byte, 20)
	_, err = rand.Read(file1Data)
	assert.NoError(t, err)

	file2Data := make([]byte, 20)
	_, err = rand.Read(file2Data)
	assert.NoError(t, err)

	file3Data := make([]byte, 20)
	_, err = rand.Read(file3Data)
	assert.NoError(t, err)

	file1Path := "file1"
	file2Path := "path/file2"
	file3Path := "long/path/file3"

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
		assert.True(t, ok, fmt.Sprintf("returned path not found: %s", p))
		d, err := s.Get(p)
		assert.NoError(t, err)
		assert.Equal(t, paths[p], d)
	}

}
