package gc

import (
	"fmt"
)

// SweepError is a type used for recording errors during the sweep phase
// of GC. If any given blob fails to delete we record the error in this
// type's map and continue deleting. After we've iterated all blobs we
// then return this error if .Blobs has values.
type SweepError struct {
	Blobs map[string]string
}

func (s *SweepError) AddBlobError(dgst string, err error) {
	// we don't want a nil map error
	if s.Blobs == nil {
		s.Blobs = map[string]string{}
	}
	s.Blobs[dgst] = err.Error()
}

func (s *SweepError) Error() string {
	return fmt.Sprintf("error deleting the following blobs: %s", s.Blobs)
}
