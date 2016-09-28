package schema

import (
	"fmt"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type TagRegistryManager struct {
	session *rethink.Session
}

func (mgr *TagRegistryManager) PutTag(ctx context.Context, repo distribution.Repository, key string, val distribution.Descriptor) error {
	tag := &Tag{
		PK:         repo.Named().Name() + ":" + key,
		Name:       key,
		Digest:     val.Digest.String(),
		DigestPK:   repo.Named().String() + "@" + val.Digest.String(),
		Repository: repo.Named().Name(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}
	tag.Author, _ = ctx.Value(auth.UserNameKey).(string)

	// TODO (tonyhb): on conflict use a custom function to update all fields but `CreatedAt`
	if _, err := tagsTable.Term().Insert(tag, rethink.InsertOpts{Conflict: "replace"}).RunWrite(mgr.session); err != nil {
		return fmt.Errorf("unable to create tag in database: %s", err)
	}

	return nil
}

func (mgr *TagRegistryManager) GetTag(ctx context.Context, repo distribution.Repository, key string) (distribution.Descriptor, error) {
	tag, err := fetchTag(mgr.session, repo.Named().Name(), key)
	if err != nil {
		return distribution.Descriptor{}, err
	}
	return tag.Descriptor(), nil
}

// DeleteTag removes a tag from the metadata store.
func (mgr *TagRegistryManager) DeleteTag(ctx context.Context, repo distribution.Repository, key string) error {
	return deleteTag(mgr.session, repo.Named().Name(), key)
}

// AllTags fetches all tags for a given repository
func (mgr *TagRegistryManager) AllTags(ctx context.Context, repo distribution.Repository) ([]string, error) {
	tags, err := fetchTags(mgr.session, repo.Named().Name())
	if err != nil {
		return nil, err
	}
	return tags.Names(), nil
}

func (mgr *TagRegistryManager) LookupTags(ctx context.Context, repo distribution.Repository, desc distribution.Descriptor) ([]string, error) {
	cursor, err := tagsTable.
		Term().
		Between(
			[]interface{}{repo.Named().Name(), desc.Digest.String(), rethink.MinVal},
			[]interface{}{repo.Named().Name(), desc.Digest.String(), rethink.MaxVal},
			rethink.BetweenOpts{Index: "repository_digest_updatedAt"},
		).
		OrderBy(rethink.OrderByOpts{Index: "repository_digest_updatedAt"}).
		Run(mgr.session)

	result := Tags{}
	if err = cursor.All(&result); err != nil {
		return nil, fmt.Errorf("unable to scan tag results: %s", err)
	}

	return result.FilterValid().Names(), nil
}
