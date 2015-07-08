package changelist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Sirupsen/logrus"
)

type fileChangelist struct {
	dir string
}

func NewFileChangelist(dir string) (*fileChangelist, error) {
	logrus.Debug("Making dir path: ", dir)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return nil, err
	}
	return &fileChangelist{dir: dir}, nil
}

func (cl fileChangelist) List() []Change {
	changes := make([]Change, 0)
	dir, err := os.Open(cl.dir)
	if err != nil {
		return changes
	}
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		return changes
	}
	sort.Sort(fileChanges(fileInfos))
	for _, f := range fileInfos {
		if f.IsDir() {
			continue
		}
		raw, err := ioutil.ReadFile(path.Join(cl.dir, f.Name()))
		if err != nil {
			// TODO(david): How should we handle this?
			fmt.Println(err.Error())
			continue
		}
		c := &tufChange{}
		err = json.Unmarshal(raw, c)
		if err != nil {
			// TODO(david): How should we handle this?
			fmt.Println(err.Error())
			continue
		}
		changes = append(changes, c)
	}
	return changes
}

func (cl fileChangelist) Add(c Change) error {
	cJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%020d_%s.change", time.Now().UnixNano(), uuid.New())
	return ioutil.WriteFile(path.Join(cl.dir, filename), cJson, 0644)
}

func (cl fileChangelist) Clear(archive string) error {
	dir, err := os.Open(cl.dir)
	if err != nil {
		return err
	}
	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, f := range files {
		os.Remove(path.Join(cl.dir, f.Name()))
	}
	return nil
}

func (cl fileChangelist) Close() error {
	// Nothing to do here
	return nil
}

type fileChanges []os.FileInfo

func (cs fileChanges) Len() int {
	return len(cs)
}

func (cs fileChanges) Less(i, j int) bool {
	return cs[i].Name() < cs[j].Name()
}

func (cs fileChanges) Swap(i, j int) {
	tmp := cs[i]
	cs[i] = cs[j]
	cs[j] = tmp
}
