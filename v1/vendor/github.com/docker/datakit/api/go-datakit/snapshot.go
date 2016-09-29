package datakit

import (
	"bytes"
	"io"
	"log"
	"strings"

	p9p "github.com/docker/go-p9p"
	"golang.org/x/net/context"
)

type SnapshotKind uint8

const (
	COMMIT SnapshotKind = 0x01 // from a commit hash
	OBJECT SnapshotKind = 0x02 // from an object hash
)

type snapshot struct {
	client *Client
	kind   SnapshotKind
	thing  string
}

type Snapshot struct {
	snapshot
}

// NewSnapshot opens a new snapshot referencing the given object.
func NewSnapshot(ctx context.Context, client *Client, kind SnapshotKind, thing string) *Snapshot {
	return &Snapshot{snapshot{client: client, kind: kind, thing: thing}}
}

// Head retrieves the commit sha of the given branch
func Head(ctx context.Context, client *Client, fromBranch string) (string, error) {
	// SHA=$(cat branch/<fromBranch>/head)
	file, err := client.Open(ctx, p9p.ORDWR, "branch", fromBranch, "head")
	if err != nil {
		log.Println("Failed to open branch/", fromBranch, "/head")
		return "", err
	}
	defer file.Close(ctx)
	buf := make([]byte, 512)
	n, err := file.Read(ctx, buf, 0)
	if err != nil {
		log.Println("Failed to Read branch", fromBranch, "head", err)
		return "", err
	}
	return strings.TrimSpace(string(buf[0:n])), nil
}

// Read reads a value from the snapshot
func (s *Snapshot) Read(ctx context.Context, path []string) (string, error) {
	var p []string

	switch s.kind {
	case COMMIT:
		p = []string{"snapshots", s.thing, "ro"}
	case OBJECT:
		p = []string{"trees", s.thing}
	}

	for _, element := range path {
		p = append(p, element)
	}
	file, err := s.client.Open(ctx, p9p.OREAD, p...)
	if err != nil {
		if err == enoent {
			return "", nil
		}
		return "", err
	}
	defer file.Close(ctx)
	reader := file.NewFileReader(ctx)
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, reader)
	return string(buf.Bytes()), nil
}
