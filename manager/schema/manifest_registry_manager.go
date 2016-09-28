package schema

import (
	"fmt"
	"time"

	middlewareErrors "github.com/docker/dhe-deploy/registry/middleware/errors"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth"
	rethink "gopkg.in/dancannon/gorethink.v2"

	"github.com/palantir/stacktrace"
)

// ManifestRegistryManager handles all database related actions for a manifest.
// This fulfils our middleware's ManifestStore interface; this interface's
// functions are called directly from the registry when saving and loading
// manifests.
//
// ManifestRegistryManager should be used only within our registry middleware
type ManifestRegistryManager struct {
	session *rethink.Session
}

func NewManifestRegistryManager(s *rethink.Session) *ManifestRegistryManager {
	return &ManifestRegistryManager{s}
}

// PutManifest is called from our middleware in the registry when saving a
// manifest. It should not be called outside any other context.
func (mgr *ManifestRegistryManager) PutManifest(ctx context.Context, repo, digest string, val distribution.Manifest) error {
	pk := fmt.Sprintf("%s@%s", repo, digest)
	m := &Manifest{
		PK:         pk,
		Repository: repo,
		Digest:     digest,
		CreatedAt:  time.Now(),
		IsDeleted:  false,
	}
	// Registry authorizes a user for each action and stores info in the context
	m.OriginalAuthor, _ = ctx.Value(auth.UserNameKey).(string)
	m.LoadMetadata(ctx, val)

	if _, err := manifestsTable.Term().Insert(m, rethink.InsertOpts{Conflict: "replace"}).RunWrite(mgr.session); err != nil {
		return stacktrace.Propagate(err, "error writing manifest")
	}
	return nil
}

func (mgr *ManifestRegistryManager) GetManifest(ctx context.Context, pk string) ([]byte, error) {
	var m Manifest
	if err := manifestsTable.getRowByIndexVal(mgr.session, "pk", pk, &m, middlewareErrors.ErrNotFound); err != nil {
		return nil, err
	}
	if m.IsDeleted {
		return nil, middlewareErrors.ErrNotFound
	}
	return m.Payload, nil
}

// DeleteManifest removes a manifest from the metadata store by marking the tombstone column as deleted.
// Note that when we mark a manifest as deleted all tags referencing this manifest must also be updated.
func (mgr *ManifestRegistryManager) DeleteManifest(ctx context.Context, key string) error {
	if _, err := mgr.GetManifest(ctx, key); err == middlewareErrors.ErrNotFound {
		return nil
	}

	resp, err := manifestsTable.Term().Insert(
		map[string]interface{}{
			"pk":      key,
			"deleted": true,
		},
		rethink.InsertOpts{Conflict: "update"},
	).RunWrite(mgr.session)
	if err != nil {
		return stacktrace.Propagate(err, "unable to delete manifest from database")
	}
	if resp.Replaced != 1 {
		return stacktrace.NewError("unexpected rows updated in database: %d", resp.Replaced)
	}
	return nil
}

// GetRepoManifests returns all manifests for a repository
func (mgr *ManifestRegistryManager) GetRepoManifests(ctx context.Context, named string) ([]*Manifest, error) {
	cursor, err := manifestsTable.Term().
		Between(
			[]interface{}{false, named, rethink.MinVal},
			[]interface{}{false, named, rethink.MaxVal},
			rethink.BetweenOpts{Index: "deleted_repository_createdAt"},
		).
		OrderBy(rethink.OrderByOpts{
			Index: "deleted_repository_createdAt",
		}).
		Run(mgr.session)
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to query db for repo manifests")
	}
	result := []*Manifest{}
	if err = cursor.All(&result); err != nil {
		return nil, stacktrace.Propagate(err, "unable to scan manifests")
	}
	return result, nil
}

// GetAllManifests returns all manifests from the metadata store that
// are not marked for deletion
func (mgr *ManifestRegistryManager) GetAllManifests(ctx context.Context) ([]*Manifest, error) {
	cursor, err := manifestsTable.Term().GetAllByIndex("deleted", false).Run(mgr.session)
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to query db for all manifests")
	}
	result := []*Manifest{}
	if err = cursor.All(&result); err != nil {
		return nil, stacktrace.Propagate(err, "unable to scan manifests")
	}
	return result, nil
}

func (mgr *ManifestRegistryManager) DeleteManifests(ctx context.Context, pks []interface{}) error {
	resp, err := manifestsTable.Term().
		GetAllByIndex("pk", pks...).
		Update(map[string]interface{}{
			"deleted": true,
		}).
		RunWrite(mgr.session)

	if err != nil {
		return stacktrace.Propagate(err, "error deleting all manifests")
	}
	if resp.Replaced != len(pks) {
		return stacktrace.NewError("unexpected number of updates; expected %d, got %d", len(pks), resp.Replaced)
	}
	return nil
}
