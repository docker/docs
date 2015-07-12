package changelist

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
)

// AppendChangelist represents a list of TUF changes
type AppendChangelist struct {
	path   string
	file   *os.File
	closed bool
}

// NewAppendChangelist is a convinience method that returns an append only TUF
// change list
func NewAppendChangelist(path string) (*AppendChangelist, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return &AppendChangelist{
		path: path,
		file: file,
	}, nil
}

// List returns a list of Changes
func (cl *AppendChangelist) List() []Change {
	cl.file.Seek(0, 0) // seek to start of file
	var changes []Change
	scnr := bufio.NewScanner(cl.file)
	for scnr.Scan() {
		line := scnr.Bytes()
		c := &TufChange{}
		err := json.Unmarshal(line, c)
		if err != nil {
			// TODO(david): How should we handle this?
			logrus.Warn(err.Error())
			continue
		}
		changes = append(changes, c)
	}
	return changes
}

// Add adds a change to the append only changelist
func (cl *AppendChangelist) Add(c Change) error {
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
func (cl *AppendChangelist) Clear(archive string) error {
	cl.file.Seek(0, 0)  // seek to start
	cl.file.Truncate(0) // truncate
	cl.file.Sync()
	return nil
}

// Close marks the change list as closed
func (cl *AppendChangelist) Close() error {
	cl.file.Sync()
	cl.closed = true
	return cl.file.Close()
}

// memChangeList implements a simple in memory change list.
type memChangelist struct {
	changes []Change
}

// List returns a list of Changes
func (cl memChangelist) List() []Change {
	return cl.changes
}

// Add adds a change to the in-memory change list
func (cl *memChangelist) Add(c Change) error {
	cl.changes = append(cl.changes, c)
	return nil
}

// Clear empties the changelist file.
func (cl *memChangelist) Clear(archive string) error {
	// appending to a nil list initializes it.
	cl.changes = nil
	return nil
}

// Close is a no-op in this in-memory change-list
func (cl *memChangelist) Close() error {
	return nil
}
