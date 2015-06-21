package trustmanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const visible os.FileMode = 0755
const private os.FileMode = 0700

// FileStore is the interface for all FileStores
type FileStore interface {
	Add(fileName string, data []byte) error
	Remove(fileName string) error
	RemoveGUN(gun string) error
	GetData(fileName string) ([]byte, error)
	GetPath(fileName string) string
	List() []string
	ListGUN(gun string) []string
}

type fileStore struct {
	baseDir string
	fileExt string
	perms   os.FileMode
}

func NewFileStore(baseDir string, fileExt string) (FileStore, error) {
	if err := CreateDirectory(baseDir); err != nil {
		return nil, err
	}

	return &fileStore{
		baseDir: baseDir,
		fileExt: fileExt,
		perms:   visible,
	}, nil
}

func NewPrivateFileStore(baseDir string, fileExt string) (FileStore, error) {
	if err := CreatePrivateDirectory(baseDir); err != nil {
		return nil, err
	}

	return &fileStore{
		baseDir: baseDir,
		fileExt: fileExt,
		perms:   private,
	}, nil
}

func (f *fileStore) Add(name string, data []byte) error {
	filePath := f.genFilePath(name)
	createDirectory(filepath.Dir(filePath), f.perms)
	if err := ioutil.WriteFile(filePath, data, f.perms); err != nil {
		return err
	}

	return nil
}

func (f *fileStore) Remove(name string) error {
	filePath := f.genFilePath(name)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (f *fileStore) RemoveGUN(gun string) error {
	dirPath := filepath.Join(f.baseDir, gun)

	// Check to see if file exists
	fi, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	// Check to see if it is a directory
	if !fi.IsDir() {
		return fmt.Errorf("GUN not found: %s", gun)
	}

	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}

	return nil
}

func (f *fileStore) GetData(name string) ([]byte, error) {
	filePath := f.genFilePath(name)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *fileStore) GetPath(name string) string {
	return f.genFilePath(name)
}

func (f *fileStore) List() []string {
	files := make([]string, 0, 0)
	filepath.Walk(f.baseDir, func(fp string, fi os.FileInfo, err error) error {
		// If there are errors, ignore this particular file
		if err != nil {
			return nil
		}
		// Ignore if it is a directory
		if fi.IsDir() {
			return nil
		}
		// Only allow matches that end with our certificate extension (e.g. *.crt)
		matched, _ := filepath.Match("*"+f.fileExt, fi.Name())

		if matched {
			files = append(files, fp)
		}
		return nil
	})
	return files
}

func (f *fileStore) ListGUN(gun string) []string {
	files := make([]string, 0, 0)
	// This prevents someone passing /path/to/dir and 'dir' not being included
	// If two '//' exist, Walk deals it with correctly
	gunPath := filepath.Join(f.baseDir, gun) + "/"
	filepath.Walk(gunPath, func(fp string, fi os.FileInfo, err error) error {
		// If there are errors, ignore this particular file
		if err != nil {
			return nil
		}
		// Ignore if it is a directory
		if fi.IsDir() {
			return nil
		}
		// Only allow matches that end with our certificate extension (e.g. *.crt)
		matched, _ := filepath.Match("*"+f.fileExt, fi.Name())

		if matched {
			files = append(files, fp)
		}
		return nil
	})
	return files
}

func (f *fileStore) genFilePath(name string) string {
	fileName := fmt.Sprintf("%s.%s", name, f.fileExt)
	filePath := filepath.Join(f.baseDir, fileName)
	return filePath
}

func CreateDirectory(dir string) error {
	return createDirectory(dir, visible)
}

func CreatePrivateDirectory(dir string) error {
	return createDirectory(dir, private)
}

// createDirectory receives a string of the path to a directory.
// It does not support passing files, so the caller has to remove
// the filename by doing filepath.Dir(full_path_to_file)
func createDirectory(dir string, perms os.FileMode) error {
	// This prevents someone passing /path/to/dir and 'dir' not being created
	// If two '//' exist, MkdirAll deals it with correctly
	dir = dir + "/"
	if err := os.MkdirAll(dir, perms); err != nil {
		return err
	}
	return nil
}
