package store

import (
	"fmt"
)

type ErrMetaNotFound struct {
	role string
}

func (err ErrMetaNotFound) Error() string {
	return fmt.Sprintf("no metadata for %s", err.role)
}
