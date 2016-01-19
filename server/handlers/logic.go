package handlers

import (
	"io"

	"github.com/Sirupsen/logrus"
	ctxu "github.com/docker/distribution/context"
	"golang.org/x/net/context"

	"github.com/docker/notary/server/errors"
	"github.com/docker/notary/server/snapshot"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/server/timestamp"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

func getRole(ctx context.Context, logger ctxu.Logger, w io.Writer, store storage.MetaStore, gun, role, checksum string) error {
	var (
		out []byte
		err error
	)
	if checksum == "" {
		// the timestamp and snapshot might be server signed so are
		// handled specially
		switch role {
		case data.CanonicalTimestampRole, data.CanonicalSnapshotRole:
			return getMaybeServerSigned(ctx, w, logger, store, gun, role)
		}
		out, err = store.GetCurrent(gun, role)
	} else {
		out, err = store.GetChecksum(gun, role, checksum)
	}

	if err != nil {
		if _, ok := err.(storage.ErrNotFound); ok {
			logrus.Errorf("404 GET %s:%s@%s, error: %v", gun, role, checksum, err)
			return errors.ErrMetadataNotFound.WithDetail(nil)
		}
		logger.Error("500 GET")
		return errors.ErrUnknown.WithDetail(err)
	}
	if out == nil {
		logger.Error("404 GET")
		return errors.ErrMetadataNotFound.WithDetail(nil)
	}
	w.Write(out)
	logger.Debug("200 GET")

	return nil
}

// getMaybeServerSigned writes the current snapshot or timestamp (based on the
// role passed) to the provided writer or returns an error. In retrieving
// the timestamp and snapshot, based on the keys held by the server, a new one
// might be generated and signed due to expiry of the previoud one or updates
// to other roles.
func getMaybeServerSigned(ctx context.Context, w io.Writer, logger ctxu.Logger, store storage.MetaStore, gun, role string) error {
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
			logger.Errorf("404 GET %s", role)
			return errors.ErrMetadataNotFound.WithDetail(nil)
		default:
			logger.Errorf("500 GET %s", role)
			return errors.ErrUnknown.WithDetail(err)
		}
	}

	logger.Debugf("200 GET %s", role)
	w.Write(out)
	return nil
}
