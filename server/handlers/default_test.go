package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	ctxu "github.com/docker/distribution/context"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"

	"github.com/docker/notary/tuf/testutils"
	"github.com/docker/notary/utils"
	"github.com/stretchr/testify/assert"
)

type handlerState struct {
	// interface{} so we can test invalid values
	store   interface{}
	crypto  interface{}
	keyAlgo interface{}
}

func defaultState() handlerState {
	return handlerState{
		store:   storage.NewMemStorage(),
		crypto:  signed.NewEd25519(),
		keyAlgo: data.ED25519Key,
	}
}

func getContext(h handlerState) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "metaStore", h.store)
	ctx = context.WithValue(ctx, "keyAlgorithm", h.keyAlgo)
	ctx = context.WithValue(ctx, "cryptoService", h.crypto)
	return ctxu.WithLogger(ctx, ctxu.GetRequestLogger(ctx))
}

func TestMainHandlerGet(t *testing.T) {
	hand := utils.RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
	handler := hand(MainHandler)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Received error on GET /: %s", err.Error())
	}
}

func TestMainHandlerNotGet(t *testing.T) {
	hand := utils.RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
	handler := hand(MainHandler)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Head(ts.URL)
	if err != nil {
		t.Fatalf("Received error on GET /: %s", err.Error())
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404, received %d", res.StatusCode)
	}
}

// GetKeyHandler needs to have access to a metadata store and cryptoservice,
// a key algorithm
func TestGetKeyHandlerInvalidConfiguration(t *testing.T) {
	noStore := defaultState()
	noStore.store = nil

	invalidStore := defaultState()
	invalidStore.store = "not a store"

	noCrypto := defaultState()
	noCrypto.crypto = nil

	invalidCrypto := defaultState()
	invalidCrypto.crypto = "not a cryptoservice"

	noKeyAlgo := defaultState()
	noKeyAlgo.keyAlgo = ""

	invalidKeyAlgo := defaultState()
	invalidKeyAlgo.keyAlgo = 1

	invalidStates := map[string][]handlerState{
		"no storage":       {noStore, invalidStore},
		"no cryptoservice": {noCrypto, invalidCrypto},
		"no keyalgorithm":  {noKeyAlgo, invalidKeyAlgo},
	}

	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   data.CanonicalTimestampRole,
	}
	req := &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(nil))}
	for errString, states := range invalidStates {
		for _, s := range states {
			err := getKeyHandler(getContext(s), httptest.NewRecorder(), req, vars)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), errString)
		}
	}
}

// GetKeyHandler needs to be set up such that an imageName and tufRole are both
// provided and non-empty.
func TestGetKeyHandlerNoRoleOrRepo(t *testing.T) {
	state := defaultState()
	req := &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(nil))}

	for _, key := range []string{"imageName", "tufRole"} {
		vars := map[string]string{
			"imageName": "gun",
			"tufRole":   data.CanonicalTimestampRole,
		}

		// not provided
		delete(vars, key)
		err := getKeyHandler(getContext(state), httptest.NewRecorder(), req, vars)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown")

		// empty
		vars[key] = ""
		err = getKeyHandler(getContext(state), httptest.NewRecorder(), req, vars)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown")
	}
}

// Getting a key for a non-supported role results in a 400.
func TestGetKeyHandlerInvalidRole(t *testing.T) {
	state := defaultState()
	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   data.CanonicalRootRole,
	}
	req := &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(nil))}

	err := getKeyHandler(getContext(state), httptest.NewRecorder(), req, vars)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid role")
}

// Getting the key for a valid role and gun succeeds
func TestGetKeyHandlerCreatesOnce(t *testing.T) {
	state := defaultState()
	roles := []string{data.CanonicalTimestampRole, data.CanonicalSnapshotRole}
	req := &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(nil))}

	for _, role := range roles {
		vars := map[string]string{"imageName": "gun", "tufRole": role}
		recorder := httptest.NewRecorder()
		err := getKeyHandler(getContext(state), recorder, req, vars)
		assert.NoError(t, err)
		assert.True(t, len(recorder.Body.String()) > 0)
	}
}

func TestGetHandlerRoot(t *testing.T) {
	store := storage.NewMemStorage()
	_, repo, _ := testutils.EmptyRepo()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "metaStore", store)

	root, err := repo.SignRoot(data.DefaultExpires("root"))
	rootJSON, err := json.Marshal(root)
	assert.NoError(t, err)
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "root", Version: 1, Data: rootJSON})

	req := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   "root",
	}

	rw := httptest.NewRecorder()

	err = getHandler(ctx, rw, req, vars)
	assert.NoError(t, err)
}

func TestGetHandlerTimestamp(t *testing.T) {
	store := storage.NewMemStorage()
	_, repo, crypto := testutils.EmptyRepo()

	ctx := getContext(handlerState{store: store, crypto: crypto})

	sn, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	snJSON, err := json.Marshal(sn)
	assert.NoError(t, err)
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "snapshot", Version: 1, Data: snJSON})

	ts, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	tsJSON, err := json.Marshal(ts)
	assert.NoError(t, err)
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "timestamp", Version: 1, Data: tsJSON})

	req := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   "timestamp",
	}

	rw := httptest.NewRecorder()

	err = getHandler(ctx, rw, req, vars)
	assert.NoError(t, err)
}

func TestGetHandlerSnapshot(t *testing.T) {
	store := storage.NewMemStorage()
	_, repo, crypto := testutils.EmptyRepo()

	ctx := getContext(handlerState{store: store, crypto: crypto})

	sn, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	snJSON, err := json.Marshal(sn)
	assert.NoError(t, err)
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "snapshot", Version: 1, Data: snJSON})

	req := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   "snapshot",
	}

	rw := httptest.NewRecorder()

	err = getHandler(ctx, rw, req, vars)
	assert.NoError(t, err)
}

func TestGetHandler404(t *testing.T) {
	store := storage.NewMemStorage()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "metaStore", store)

	req := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   "root",
	}

	rw := httptest.NewRecorder()

	err := getHandler(ctx, rw, req, vars)
	assert.Error(t, err)
}

func TestGetHandlerNilData(t *testing.T) {
	store := storage.NewMemStorage()
	store.UpdateCurrent("gun", storage.MetaUpdate{Role: "root", Version: 1, Data: nil})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "metaStore", store)

	req := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	vars := map[string]string{
		"imageName": "gun",
		"tufRole":   "root",
	}

	rw := httptest.NewRecorder()

	err := getHandler(ctx, rw, req, vars)
	assert.Error(t, err)
}

func TestGetHandlerNoStorage(t *testing.T) {
	ctx := context.Background()

	req := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	err := GetHandler(ctx, nil, req)
	assert.Error(t, err)
}
