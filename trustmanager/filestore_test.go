package trustmanager

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestAddFile(t *testing.T) {
	testData := []byte("This test data should be part of the file.")
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the Add function
	err = store.Add(testName, testData)
	if err != nil {
		t.Fatalf("failed to add file to store: %v", err)
	}

	// Check to see if file exists
	b, err := ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("expected file not found: %v", err)
	}

	if !bytes.Equal(b, testData) {
		t.Fatalf("unexpected content in the file: %s", expectedFilePath)
	}
}

func TestRemoveFile(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	_, err = generateRandomFile(expectedFilePath, perms)
	if err != nil {
		t.Fatalf("failed to generate random file: %v", err)
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the Remove function
	err = store.Remove(testName)
	if err != nil {
		t.Fatalf("failed to remove file from store: %v", err)
	}

	// Check to see if file exists
	_, err = os.Stat(expectedFilePath)
	if err == nil {
		t.Fatalf("expected not to find file: %s", expectedFilePath)
	}
}

func TestRemoveDir(t *testing.T) {
	testName := "docker.com/diogomonica/"
	testExt := "key"
	perms := os.FileMode(0700)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	_, err = generateRandomFile(expectedFilePath, perms)
	if err != nil {
		t.Fatalf("failed to generate random file: %v", err)
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the RemoveDir function
	err = store.RemoveDir(testName)
	if err != nil {
		t.Fatalf("failed to remove directory: %v", err)
	}

	expectedDirectory := filepath.Dir(expectedFilePath)
	// Check to see if file exists
	_, err = os.Stat(expectedDirectory)
	if err == nil {
		t.Fatalf("expected not to find directory: %s", expectedDirectory)
	}
}

func TestListFiles(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	var expectedFilePath string
	// Create 10 randomfiles
	for i := 1; i <= 10; i++ {
		// Since we're generating this manually we need to add the extension '.'
		expectedFilename := testName + strconv.Itoa(i) + "." + testExt
		expectedFilePath = filepath.Join(tempBaseDir, expectedFilename)
		_, err = generateRandomFile(expectedFilePath, perms)
		if err != nil {
			t.Fatalf("failed to generate random file: %v", err)
		}
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the List function. Expect 10 files
	files := store.ListFiles()
	if len(files) != 10 {
		t.Fatalf("expected 10 files in listing, got: %d", len(files))
	}
}

func TestListDir(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	var expectedFilePath string
	// Create 10 randomfiles
	for i := 1; i <= 10; i++ {
		// Since we're generating this manually we need to add the extension '.'
		fileName := fmt.Sprintf("%s-%s.%s", testName, strconv.Itoa(i), testExt)
		expectedFilePath = filepath.Join(tempBaseDir, fileName)
		_, err = generateRandomFile(expectedFilePath, perms)
		if err != nil {
			t.Fatalf("failed to generate random file: %v", err)
		}
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}

	// Call the ListDir function
	files := store.ListDir("docker.com/")
	if len(files) != 10 {
		t.Fatalf("expected 10 files in listing, got: %d", len(files))
	}
	files = store.ListDir("docker.com/notary")
	if len(files) != 10 {
		t.Fatalf("expected 10 files in listing, got: %d", len(files))
	}
	files = store.ListDir("fakedocker.com/")
	if len(files) != 0 {
		t.Fatalf("expected 0 files in listing, got: %d", len(files))
	}
}

func TestGetPath(t *testing.T) {
	testExt := "crt"
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
	if err != nil {
		t.Fatalf("unexpected error from GetPath: %v", err)
	}
	if result != firstPath {
		t.Fatalf("Expecting: %s", firstPath)
	}

	result, err = store.GetPath("/docker.io/testing-dashes/@#$%^&()")
	if err != nil {
		t.Fatalf("unexpected error from GetPath: %v", err)
	}
	if result != secondPath {
		t.Fatalf("Expecting: %s", secondPath)
	}
}

func TestGetPathProtection(t *testing.T) {
	testExt := "crt"
	perms := os.FileMode(0755)

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: "/path/to/filestore/",
		fileExt: testExt,
		perms:   perms,
	}

	// Should deny requests for paths outside the filestore
	if _, err := store.GetPath("../../etc/passwd"); err != ErrPathOutsideStore {
		t.Fatalf("expected ErrPathOutsideStore error from GetPath")
	}
	if _, err := store.GetPath("private/../../../etc/passwd"); err != ErrPathOutsideStore {
		t.Fatalf("expected ErrPathOutsideStore error from GetPath")
	}

	// Convoluted paths should work as long as they end up inside the store
	expected := "/path/to/filestore/filename.crt"
	result, err := store.GetPath("private/../../filestore/./filename")
	if err != nil {
		t.Fatalf("unexpected error from GetPath: %v", err)
	}
	if result != expected {
		t.Fatalf("Expecting: %s (got: %s)", expected, result)
	}

	// Repeat tests with a relative baseDir
	relStore := &SimpleFileStore{
		baseDir: "relative/file/path",
		fileExt: testExt,
		perms:   perms,
	}

	// Should deny requests for paths outside the filestore
	if _, err := relStore.GetPath("../../etc/passwd"); err != ErrPathOutsideStore {
		t.Fatalf("expected ErrPathOutsideStore error from GetPath")
	}
	if _, err := relStore.GetPath("private/../../../etc/passwd"); err != ErrPathOutsideStore {
		t.Fatalf("expected ErrPathOutsideStore error from GetPath")
	}

	// Convoluted paths should work as long as they end up inside the store
	expected = "relative/file/path/filename.crt"
	result, err = relStore.GetPath("private/../../path/./filename")
	if err != nil {
		t.Fatalf("unexpected error from GetPath: %v", err)
	}
	if result != expected {
		t.Fatalf("Expecting: %s (got: %s)", expected, result)
	}
}

func TestGetData(t *testing.T) {
	testName := "docker.com/notary/certificate"
	testExt := "crt"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	expectedData, err := generateRandomFile(expectedFilePath, perms)
	if err != nil {
		t.Fatalf("failed to generate random file: %v", err)
	}

	// Create our SimpleFileStore
	store := &SimpleFileStore{
		baseDir: tempBaseDir,
		fileExt: testExt,
		perms:   perms,
	}
	testData, err := store.Get(testName)
	if err != nil {
		t.Fatalf("failed to get data from: %s", testName)

	}
	if !bytes.Equal(testData, expectedData) {
		t.Fatalf("unexpected content for the file: %s", expectedFilePath)
	}
}

func TestCreateDirectory(t *testing.T) {
	testDir := "fake/path/to/directory"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	dirPath := filepath.Join(tempBaseDir, testDir)

	// Call createDirectory
	CreateDirectory(dirPath)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("expected find directory: %s", dirPath)
	}

	// Check to see if it is a directory
	if !fi.IsDir() {
		t.Fatalf("expected to be directory: %s", dirPath)
	}

	// Check to see if the permissions match
	if fi.Mode().String() != "drwxr-xr-x" {
		t.Fatalf("permissions are wrong for: %s. Got: %s", dirPath, fi.Mode().String())
	}
}

func TestCreatePrivateDirectory(t *testing.T) {
	testDir := "fake/path/to/private/directory"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	dirPath := filepath.Join(tempBaseDir, testDir)

	// Call createDirectory
	CreatePrivateDirectory(dirPath)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("expected find directory: %s", dirPath)
	}

	// Check to see if it is a directory
	if !fi.IsDir() {
		t.Fatalf("expected to be directory: %s", dirPath)
	}

	// Check to see if the permissions match
	if fi.Mode().String() != "drwx------" {
		t.Fatalf("permissions are wrong for: %s. Got: %s", dirPath, fi.Mode().String())
	}
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
