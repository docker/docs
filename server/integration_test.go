// This makes sure that the server is compatible with the tuf httpstore.

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/store"
	"github.com/docker/notary/tuf/testutils"
	"github.com/docker/notary/tuf/validation"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

// Ensures that the httpstore can interpret the errors returned from the server
func TestValidationErrorFormat(t *testing.T) {
	ctx := context.WithValue(
		context.Background(), "metaStore", storage.NewMemStorage())
	ctx = context.WithValue(ctx, "keyAlgorithm", data.ED25519Key)

	handler := RootHandler(nil, ctx, signed.NewEd25519(), nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := store.NewHTTPStore(
		fmt.Sprintf("%s/v2/docker.com/notary/_trust/tuf/", server.URL),
		"",
		"json",
		"key",
		http.DefaultTransport,
	)

	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	r, tg, sn, ts, err := testutils.Sign(repo)
	require.NoError(t, err)
	rs, rt, _, _, err := testutils.Serialize(r, tg, sn, ts)
	require.NoError(t, err)

	// No snapshot is passed, and the server doesn't have the snapshot key,
	// so ErrBadHierarchy
	err = client.SetMultiMeta(map[string][]byte{
		data.CanonicalRootRole:    rs,
		data.CanonicalTargetsRole: rt,
	})
	require.Error(t, err)
	require.IsType(t, validation.ErrBadHierarchy{}, err)
}
