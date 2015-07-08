package changelist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type appendChangelist struct {
	path   string
	file   *os.File
	closed bool
}

func NewAppendChangelist(path string) (*appendChangelist, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return &appendChangelist{
		path: path,
		file: file,
	}, nil
}

func (cl *appendChangelist) List() []Change {
	cl.file.Seek(0, 0) // seek to start of file
	changes := make([]Change, 0)
	scnr := bufio.NewScanner(cl.file)
	for scnr.Scan() {
		line := scnr.Bytes()
		c := &tufChange{}
		err := json.Unmarshal(line, c)
		if err != nil {
			// TODO(david): How should we handle this?
			fmt.Println(err.Error())
			continue
		}
		changes = append(changes, c)
	}
	return changes
}

func (cl *appendChangelist) Add(c Change) error {
	cl.file.Seek(0, 2) // seek to end of file
	entry, err := json.Marshal(c)
	if err != nil {
		return err
	}
	n, err := cl.file.Write(entry)
	if err != nil {
		if n > 0 {
			// trim partial write if necessary
			size, _ := cl.file.Seek(-int64(n), 2)
			cl.file.Truncate(size)
		}
		return err
	}
	cl.file.Write([]byte("\n"))
	cl.file.Sync()
	return nil
}

// Clear empties the changelist file. It does not currently
// support archiving
func (cl *appendChangelist) Clear(archive string) error {
	cl.file.Seek(0, 0)  // seek to start
	cl.file.Truncate(0) // truncate
	cl.file.Sync()
	return nil
}

func (cl *appendChangelist) Close() error {
	cl.file.Sync()
	cl.closed = true
	return cl.file.Close()
}

// memChangeList implements a simple in memory change list.
type memChangelist struct {
	changes []Change
}

func (cl memChangelist) List() []Change {
	return cl.changes
}

func (cl *memChangelist) Add(c Change) error {
	cl.changes = append(cl.changes, c)
	return nil
}

func (cl *memChangelist) Clear(archive string) error {
	// appending to a nil list initializes it.
	cl.changes = nil
	return nil
}

func (cl *memChangelist) Close() error {
	return nil
}
