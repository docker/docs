package schema

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/docker/dhe-deploy"
	merr "github.com/docker/dhe-deploy/registry/middleware/errors"

	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
	"github.com/palantir/stacktrace"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var tagsTable = table{
	db:         deploy.DTRDBName,
	name:       "tags",
	primaryKey: "pk",
	secondaryIndexes: map[string][]string{
		// Used within GC to look up all undeleted tags
		"deleted": nil,
		// For undeleted tag lookup by repository, ordered by updated at
		"repository_deleted_updatedAt": {"repository", "deleted", "updatedAt"},
		// For finding all tags pointing to a manifest for a repository, ordered by updated at.
		"repository_digest_updatedAt": {"repository", "digest", "updatedAt"},
	},
}

var (
	ErrNoSuchTag = fmt.Errorf("tag not found")
)

// Tag represents a tag pointing to a manifest for a given repsoitory,
// storing additional metadata about the tag such as author and date
// pushed.
//
// A tag is created via the metadata store middleware applied to the registry
type Tag struct {
	// PK is the string "namespace/repo:tagname"
	PK string `gorethink:"pk"`
	// Name represents the name of the tag
	Name string `gorethink:"name"`
	// Repository is a concatenation of the namespace and repo name that the
	// tag refers to (ie. reference.Named interface within distribution)
	Repository string `gorethink:"repository"`
	// Digest stores the hash of the manifest that this tag points to
	Digest string `gorethink:"digest"`
	// DigestPK stores the concatenated repository and digest which is the PK
	// of the tag's manifest. We need to store this to properly use rethink's
	// merge functionality, as string concatenation within a merge fails
	DigestPK string `gorethink:"digestPK"`
	// Author is the user that last pushed this tag
	Author string `gorethink:"authorNamespace"`
	// UpdatedAt is when this tag was last pushed
	UpdatedAt time.Time `gorethink:"updatedAt"`
	// CreatedAt refers to when this tag was first created
	CreatedAt time.Time `gorethink:"createdAt"`
	// IsDeleted is a tombstone marker which depicts whether a tag is deleted.
	IsDeleted bool `gorethink:"deleted"`
	// Manifest is used to pull the associated manifest from the tag,
	// when reading tags from the database
	Manifest Manifest `gorethink:"manifest,reference,omitempty" gorethink_ref:"pk"`
}

func (t Tag) Descriptor() distribution.Descriptor {
	dgst, _ := digest.ParseDigest(t.Digest)
	return distribution.Descriptor{
		Digest: dgst,
	}
}

// DigestMatches returns whether the tag's digest is equal to the given
// *decoded* checksum
func (t Tag) DigestMatches(hash []byte) (matches bool, err error) {
	defer func() {
		// digest.Digest(t.Digest) can cause a panic with invalid data
		if r := recover(); r != nil {
			matches = false
			err = fmt.Errorf("invalid tag hash")
		}
	}()

	dgst := digest.Digest(t.Digest)
	decoded, err := hex.DecodeString(dgst.Hex())
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash, decoded), nil
}

// DigestMismatches returns true if the tag's digest is different to the
// given decoded checksum. This is a helper method inverting DigestMatches
// for easy assignment to Tag.HashMismatch attributes.
func (t Tag) DigestMismatches(hash []byte) (bool, error) {
	match, err := t.DigestMatches(hash)
	if err != nil {
		// by default if there's an error return a mismatch
		return true, err
	}
	return !match, nil
}

// HasManifest is used to ernsure that Tag.Manifest is valid. We need this
// because:
//  - A user can delete manifests via the registry or DTR API
//  - Deleting manifests doesn't mark tags as deleted
//  - Therefore we can load undeleted tags which are invalid as their manifest
//    is deleted.
// This lets us filter out invalid tags due to deleted manifests.
func (t Tag) HasManifest() bool {
	// Ensure that the tag's digest matches the manifest PK
	if t.DigestPK == "" || t.DigestPK != t.Manifest.PK {
		return false
	}
	// assert that the manifest is all around valid
	return len(t.Manifest.Payload) != 0 && !t.Manifest.IsDeleted
}

type Tags []Tag

// Names returns a slice of tag names
func (t Tags) Names() []string {
	r := make([]string, len(t))
	for i, v := range t {
		r[i] = v.Name
	}
	return r
}

// Undeleted returns a subset of Tags representing all valid tags which
// have not been marked for deletion AND still have existing manifests.
//
// This is because we can delete manifests without deleting tags - see
// Tag.HasManifest for more information.
func (t Tags) FilterValid() Tags {
	set := Tags{}
	for _, v := range t {
		if !v.IsDeleted && v.HasManifest() {
			set = append(set, v)
		}
	}
	return set
}

type TagManager struct {
	session *rethink.Session
}

func NewTagManager(s *rethink.Session) *TagManager {
	return &TagManager{s}
}

// RepositoryTags finds all tags for a given "namespace/repo" string.
// This returns a Tags collection which includes metadata regarding each tag;
// because of this it does not conform to any registry interface and is to be
// used in DTR only.
func (mgr *TagManager) RepositoryTags(named string) (Tags, error) {
	return fetchTags(mgr.session, named)
}

// RepostoryTag finds and returns a tag for a given "namespace/repo" repository
func (mgr *TagManager) RepositoryTag(named, tagName string) (Tag, error) {
	tag, err := fetchTag(mgr.session, named, tagName)

	// fetchTag returns `merr.ErrNotFound` to comply with our middleware
	// interface; the rest of our app code expects schema.ErrNoSuchTag
	if err != nil && err == merr.ErrNotFound {
		return tag, ErrNoSuchTag
	}

	return tag, err
}

// DeleteTag delegates to deleteTag in order to remove a tag from registry
func (mgr *TagManager) DeleteTag(named, tagName string) error {
	return deleteTag(mgr.session, named, tagName)
}

// AllTags returns all undeleted tags within the database for every repository
func (mgr *TagManager) AllTags() (Tags, error) {
	cursor, err := tagsTable.Term().
		GetAllByIndex("deleted", false).
		Merge(func(result rethink.Term) interface{} {
			return map[string]interface{}{
				"manifest": manifestsTable.Term().Get(result.Field("digestPK")),
			}
		}).
		Run(mgr.session)

	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to query db for tags")
	}

	result := Tags{}
	if err = cursor.All(&result); err != nil {
		return nil, stacktrace.Propagate(err, "unable to scan tag results")
	}

	return result.FilterValid(), nil
}

// fetchTag returns a single tag given a rethink session via its primary key
// (namespace/repo:tag).
func fetchTag(session *rethink.Session, named, tagName string) (Tag, error) {
	var tag Tag

	cursor, err := tagsTable.
		Term().
		Get(named + ":" + tagName).
		// The "Do" function allows us to perform the merge ONLY if a tag was
		// found; without this rethink will throw an error when merging into
		// NULL which isn't caught as the error is not of type
		// rethink.ErrEmptyResult
		Do(func(row rethink.Term) interface{} {
			return rethink.Branch(
				row,
				row.Merge(func(result rethink.Term) interface{} {
					return map[string]interface{}{
						"manifest": manifestsTable.Term().Get(
							result.Field("digestPK"),
						),
					}
				}),
				nil,
			)
		}).
		Run(session)

	if err != nil {
		return tag, stacktrace.Propagate(err, "unable to query db for tag")
	}

	if err = cursor.One(&tag); err != nil {
		if err == rethink.ErrEmptyResult {
			return tag, merr.ErrNotFound
		}
		return tag, stacktrace.Propagate(err, "unable to scan tag result")
	}

	// We still use the PK for looking up a single tag; check for deletions in memory.
	// TODO (tonyhb): potentially reduce network traffic by querying for tag by
	// PK and deleted flag. Need to test index performance and think about HasManifest
	// checks within the database potentially.
	if tag.IsDeleted || !tag.HasManifest() {
		return Tag{}, merr.ErrNotFound
	}

	return tag, nil
}

// fetchTags is a utility function for fetching tags from a given repo.
// This is used across both TagManager and TagRegistryManager.
func fetchTags(session *rethink.Session, named string) (Tags, error) {
	cursor, err := tagsTable.Term().
		// Select all tags which aren't deleted for the given repository
		Between(
			[]interface{}{named, false, rethink.MinVal},
			[]interface{}{named, false, rethink.MaxVal},
			rethink.BetweenOpts{Index: "repository_deleted_updatedAt"},
		).
		OrderBy(rethink.OrderByOpts{
			Index: "repository_deleted_updatedAt",
		}).
		// Merge the tag's manifest into each tag, adding things such as OS and
		// architecture.
		Merge(func(result rethink.Term) interface{} {
			return map[string]interface{}{
				"manifest": manifestsTable.Term().Get(result.Field("digestPK")),
			}
		}).
		Run(session)

	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to query db for tags")
	}

	result := Tags{}
	if err = cursor.All(&result); err != nil {
		return nil, stacktrace.Propagate(err, "unable to scan tag results")
	}

	return result.FilterValid(), nil
}

// deleteTag removes a tag from the metadata store.
// Note that if this is the last tag referencing a manifest the manifest will
// not be deleted. Dangling manifests must be cleaned up via the garbage
// collection job.
func deleteTag(session *rethink.Session, named, tagName string) (err error) {
	// See if the tag exists. We don't use fetchTag as fetchTag verifies that
	// the manifest is valid and undeleted; we still want to be able to delete
	// tags that have their manifest deleted (ie push tag A, delete manifest
	// for tag A, then delete tag A separately). This is especially a need for
	// our integration tests.
	//
	// Note we can't use `getRowByIndexVal` as rethink returns an empty string
	// for the manifest; go tries marshalling this into type Manifest which
	// fails. For this reason we must also merge in an empty manifest.
	var tag Tag
	cursor, err := tagsTable.
		Term().
		Get(named + ":" + tagName).
		// The "Do" function allows us to perform the merge ONLY if a tag was
		// found; without this rethink will throw an error when merging into
		// NULL which isn't caught as the error is not of type
		// rethink.ErrEmptyResult
		Do(func(row rethink.Term) interface{} {
			return rethink.Branch(
				row,
				row.Merge(func(result rethink.Term) interface{} {
					return map[string]interface{}{
						"manifest": map[string]interface{}{},
					}
				}),
				nil,
			)
		}).
		Run(session)
	if err != nil {
		return stacktrace.Propagate(err, "unable to query db for tag")
	}
	if err = cursor.One(&tag); err != nil {
		if err == rethink.ErrEmptyResult {
			return merr.ErrNotFound
		}
		return stacktrace.Propagate(err, "unable to scan tag result")
	}

	// Update the tag to deleted
	resp, err := tagsTable.Term().Insert(
		map[string]interface{}{"pk": named + ":" + tagName, "deleted": true},
		rethink.InsertOpts{Conflict: "update"},
	).RunWrite(session)
	if err != nil {
		return stacktrace.Propagate(err, "unable to delete tag from database")
	}
	if resp.Replaced != 1 {
		return stacktrace.NewError("unexpected rows updated in database: %d", resp.Replaced)
	}

	return nil
}
