package handlers

import (
	"io"

	"golang.org/x/net/context"

	"github.com/docker/notary/server/errors"
	"github.com/docker/notary/server/snapshot"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/server/timestamp"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

func getRole(ctx context.Context, w io.Writer, store storage.MetaStore, gun, role, checksum string) error {
	var (
		out []byte
		err error
	)
	if checksum == "" {
		// the timestamp and snapshot might be server signed so are
		// handled specially
		switch role {
		case data.CanonicalTimestampRole, data.CanonicalSnapshotRole:
			return getMaybeServerSigned(ctx, w, store, gun, role)
		}
		out, err = store.GetCurrent(gun, role)
	} else {
		out, err = store.GetChecksum(gun, role, checksum)
	}

	if err != nil {
		if _, ok := err.(storage.ErrNotFound); ok {
			return errors.ErrMetadataNotFound.WithDetail(err)
		}
		return errors.ErrUnknown.WithDetail(err)
	}
	if out == nil {
		return errors.ErrMetadataNotFound.WithDetail(nil)
	}
	w.Write(out)

	return nil
}

// getMaybeServerSigned writes the current snapshot or timestamp (based on the
// role passed) to the provided writer or returns an error. In retrieving
// the timestamp and snapshot, based on the keys held by the server, a new one
// might be generated and signed due to expiry of the previous one or updates
// to other roles.
func getMaybeServerSigned(ctx context.Context, w io.Writer, store storage.MetaStore, gun, role string) error {
	cryptoServiceVal := ctx.Value("cryptoService")
	cryptoService, ok := cryptoServiceVal.(signed.CryptoService)
	if !ok {
		return errors.ErrNoCryptoService.WithDetail(nil)
	}

	var (
		out []byte
		err error
	)
	switch role {
	case data.CanonicalSnapshotRole:
		out, err = snapshot.GetOrCreateSnapshot(gun, store, cryptoService)
	case data.CanonicalTimestampRole:
		out, err = timestamp.GetOrCreateTimestamp(gun, store, cryptoService)
	}
	if err != nil {
		switch err.(type) {
		case *storage.ErrNoKey, storage.ErrNotFound:
			return errors.ErrMetadataNotFound.WithDetail(err)
		default:
			return errors.ErrUnknown.WithDetail(err)
		}
	}

	w.Write(out)
	return nil
}
