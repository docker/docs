package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	ctxu "github.com/docker/distribution/context"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/server/errors"
	"github.com/docker/notary/server/snapshot"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/server/timestamp"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

// MainHandler is the default handler for the server
func MainHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		_, err := w.Write([]byte("{}"))
		if err != nil {
			return errors.ErrUnknown.WithDetail(err)
		}
	} else {
		return errors.ErrGenericNotFound.WithDetail(nil)
	}
	return nil
}

// AtomicUpdateHandler will accept multiple TUF files and ensure that the storage
// backend is atomically updated with all the new records.
func AtomicUpdateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	vars := mux.Vars(r)
	return atomicUpdateHandler(ctx, w, r, vars)
}

func atomicUpdateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	gun := vars["imageName"]
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
	if !ok {
		return errors.ErrNoStorage.WithDetail(nil)
	}
	cryptoServiceVal := ctx.Value("cryptoService")
	cryptoService, ok := cryptoServiceVal.(signed.CryptoService)
	if !ok {
		return errors.ErrNoCryptoService.WithDetail(nil)
	}

	reader, err := r.MultipartReader()
	if err != nil {
		return errors.ErrMalformedUpload.WithDetail(nil)
	}
	var updates []storage.MetaUpdate
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		role := strings.TrimSuffix(part.FileName(), ".json")
		if role == "" {
			return errors.ErrNoFilename.WithDetail(nil)
		} else if !data.ValidRole(role) {
			return errors.ErrInvalidRole.WithDetail(role)
		}
		meta := &data.SignedMeta{}
		var input []byte
		inBuf := bytes.NewBuffer(input)
		dec := json.NewDecoder(io.TeeReader(part, inBuf))
		err = dec.Decode(meta)
		if err != nil {
			return errors.ErrMalformedJSON.WithDetail(nil)
		}
		version := meta.Signed.Version
		updates = append(updates, storage.MetaUpdate{
			Role:    role,
			Version: version,
			Data:    inBuf.Bytes(),
		})
	}
	updates, err = validateUpdate(cryptoService, gun, updates, store)
	if err != nil {
		return errors.ErrInvalidUpdate.WithDetail(err)
	}
	err = store.UpdateMany(gun, updates)
	if err != nil {
		return errors.ErrUpdating.WithDetail(err)
	}
	return nil
}

// GetHandler returns the json for a specified role and GUN.
func GetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	vars := mux.Vars(r)
	return getHandler(ctx, w, r, vars)
}

func getHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	gun := vars["imageName"]
	tufRole := vars["tufRole"]
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
	if !ok {
		return errors.ErrNoStorage.WithDetail(nil)
	}

	logger := ctxu.GetLoggerWithFields(ctx, map[string]interface{}{"gun": gun, "tufRole": tufRole})

	switch data.CanonicalRole(tufRole) {
	case data.CanonicalTimestampRole:
		return getTimestamp(ctx, w, logger, store, gun)
	case data.CanonicalSnapshotRole:
		return getSnapshot(ctx, w, logger, store, gun)
	}

	out, err := store.GetCurrent(gun, tufRole)
	if err != nil {
		if _, ok := err.(*storage.ErrNotFound); ok {
			logrus.Error(gun + ":" + tufRole)
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

// DeleteHandler deletes all data for a GUN. A 200 responses indicates success.
func DeleteHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
	if !ok {
		return errors.ErrNoStorage.WithDetail(nil)
	}
	vars := mux.Vars(r)
	gun := vars["imageName"]
	logger := ctxu.GetLoggerWithField(ctx, gun, "gun")
	err := store.Delete(gun)
	if err != nil {
		logger.Error("500 DELETE repository")
		return errors.ErrUnknown.WithDetail(err)
	}
	return nil
}

// getTimestampHandler returns a timestamp.json given a GUN
func getTimestamp(ctx context.Context, w http.ResponseWriter, logger ctxu.Logger, store storage.MetaStore, gun string) error {
	cryptoServiceVal := ctx.Value("cryptoService")
	cryptoService, ok := cryptoServiceVal.(signed.CryptoService)
	if !ok {
		return errors.ErrNoCryptoService.WithDetail(nil)
	}

	out, err := timestamp.GetOrCreateTimestamp(gun, store, cryptoService)
	if err != nil {
		switch err.(type) {
		case *storage.ErrNoKey, *storage.ErrNotFound:
			logger.Error("404 GET timestamp")
			return errors.ErrMetadataNotFound.WithDetail(nil)
		default:
			logger.Error("500 GET timestamp")
			return errors.ErrUnknown.WithDetail(err)
		}
	}

	logger.Debug("200 GET timestamp")
	w.Write(out)
	return nil
}

// getTimestampHandler returns a timestamp.json given a GUN
func getSnapshot(ctx context.Context, w http.ResponseWriter, logger ctxu.Logger, store storage.MetaStore, gun string) error {
	cryptoServiceVal := ctx.Value("cryptoService")
	cryptoService, ok := cryptoServiceVal.(signed.CryptoService)
	if !ok {
		return errors.ErrNoCryptoService.WithDetail(nil)
	}

	out, err := snapshot.GetOrCreateSnapshot(gun, store, cryptoService)
	if err != nil {
		switch err.(type) {
		case *storage.ErrNoKey, *storage.ErrNotFound:
			logger.Error("404 GET snapshot")
			return errors.ErrMetadataNotFound.WithDetail(nil)
		default:
			logger.Error("500 GET snapshot")
			return errors.ErrUnknown.WithDetail(err)
		}
	}

	logger.Debug("200 GET snapshot")
	w.Write(out)
	return nil
}

// GetKeyHandler returns a public key for the specified role, creating a new key-pair
// it if it doesn't yet exist
func GetKeyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	vars := mux.Vars(r)
	return getKeyHandler(ctx, w, r, vars)
}

func getKeyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	gun, ok := vars["imageName"]
	if !ok || gun == "" {
		return errors.ErrUnknown.WithDetail("no gun")
	}
	role, ok := vars["tufRole"]
	if !ok || role == "" {
		return errors.ErrUnknown.WithDetail("no role")
	}

	logger := ctxu.GetLoggerWithField(ctx, gun, "gun")

	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
	if !ok || store == nil {
		logger.Error("500 GET storage not configured")
		return errors.ErrNoStorage.WithDetail(nil)
	}
	c := ctx.Value("cryptoService")
	crypto, ok := c.(signed.CryptoService)
	if !ok || crypto == nil {
		logger.Error("500 GET crypto service not configured")
		return errors.ErrNoCryptoService.WithDetail(nil)
	}
	algo := ctx.Value("keyAlgorithm")
	keyAlgo, ok := algo.(string)
	if !ok || keyAlgo == "" {
		logger.Error("500 GET key algorithm not configured")
		return errors.ErrNoKeyAlgorithm.WithDetail(nil)
	}
	keyAlgorithm := keyAlgo

	var (
		key data.PublicKey
		err error
	)
	switch role {
	case data.CanonicalTimestampRole:
		key, err = timestamp.GetOrCreateTimestampKey(gun, store, crypto, keyAlgorithm)
	case data.CanonicalSnapshotRole:
		key, err = snapshot.GetOrCreateSnapshotKey(gun, store, crypto, keyAlgorithm)
	default:
		logger.Errorf("400 GET %s key: %v", role, err)
		return errors.ErrInvalidRole.WithDetail(role)
	}
	if err != nil {
		logger.Errorf("500 GET %s key: %v", role, err)
		return errors.ErrUnknown.WithDetail(err)
	}

	out, err := json.Marshal(key)
	if err != nil {
		logger.Errorf("500 GET %s key", role)
		return errors.ErrUnknown.WithDetail(err)
	}
	logger.Debugf("200 GET %s key", role)
	w.Write(out)
	return nil
}

// NotFoundHandler is used as a generic catch all handler to return the ErrMetadataNotFound
// 404 response
func NotFoundHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return errors.ErrMetadataNotFound.WithDetail(nil)
}
