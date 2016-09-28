package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/docker/dhe-deploy/hubconfig"
)

type keyValueStore struct {
	rootDirectory string
}

func NewKeyValueStore(rootDirectory string) hubconfig.KeyValueStore {
	return keyValueStore{
		rootDirectory: rootDirectory,
	}
}

func (s keyValueStore) Get(filePath string) ([]byte, error) {
	if content, err := ioutil.ReadFile(path.Join(s.rootDirectory, filePath)); os.IsNotExist(err) {
		return nil, nil
	} else {
		return content, err
	}
}

// TODO: implement
func (s keyValueStore) List(filePath string) ([]string, error) {
	return []string{}, nil
}

// TODO: implement
func (s keyValueStore) Lock(filePath string, value []byte, ttl time.Duration) error {
	return nil
}

func (s keyValueStore) Delete(filePath string) error {
	panic("Not implemented")
	return nil
}

func (s keyValueStore) SemWait(filePath string, stopCh <-chan struct{}) (bool, error) {
	panic("Not implemented")
	return false, nil
}

func (s keyValueStore) SemSignal(filePath string) error {
	panic("Not implemented")
	return nil
}

func (s keyValueStore) Put(filePath string, content []byte) error {
	fullFilePath := path.Join(s.rootDirectory, filePath)
	dirPath := path.Dir(fullFilePath)
	if err := os.MkdirAll(dirPath, 02750); err != nil {
		return err
	}

	tmpF, err := ioutil.TempFile(dirPath, "")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %s", err)
	}

	if _, err = tmpF.Write(content); err != nil {
		tmpF.Close()
		return err
	}
	if err := tmpF.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpF.Name(), fullFilePath); err != nil {
		return err
	}
	return nil
}

func (s keyValueStore) Ping() error {
	_, err := os.Stat(s.rootDirectory)
	return err
}
