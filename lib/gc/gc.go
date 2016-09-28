package gc

import (
	"github.com/docker/dhe-deploy/events"
	"github.com/docker/dhe-deploy/manager/schema"

	// register all storage and auth drivers
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
	_ "github.com/docker/distribution/registry/auth/token"
	_ "github.com/docker/distribution/registry/proxy"
	_ "github.com/docker/distribution/registry/storage/driver/azure"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/cloudfront"
	_ "github.com/docker/distribution/registry/storage/driver/oss"
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	_ "github.com/docker/distribution/registry/storage/driver/swift"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution"
	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/registry/storage"
	"github.com/docker/distribution/registry/storage/driver"
	"github.com/docker/distribution/registry/storage/driver/factory"
	"github.com/docker/libtrust"
	"github.com/palantir/stacktrace"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

const (
	ModeByTag      = "tag"
	ModeByManifest = "manifest"
	ModeLegacy     = "legacy"
)

// NewGC returns a new GC struct for running GC with the desired options
func NewGC(opts Opts) *GC {
	return &GC{
		mode:      opts.Mode,
		Context:   context.Background(),
		config:    opts.Config,
		tagMgr:    schema.NewTagManager(opts.Session),
		mfstMgr:   schema.NewManifestRegistryManager(opts.Session),
		eventMgr:  schema.NewEventManager(opts.Session),
		markSet:   map[string]struct{}{},
		deleteSet: map[string]struct{}{},
	}
}

// Opts is a configuration struct passed to NewGC to configure a GC run
type Opts struct {
	Session *rethink.Session
	Config  configuration.Configuration
	Mode    string
}

type GC struct {
	context.Context

	mode string

	tagMgr   *schema.TagManager
	mfstMgr  *schema.ManifestRegistryManager
	eventMgr schema.EventManager

	config   configuration.Configuration
	registry distribution.Namespace
	driver   driver.StorageDriver

	tags      schema.Tags
	markSet   map[string]struct{}
	deleteSet map[string]struct{}
}

// Run kicks off a new garbage collection. GC may be configured to run in one
// of three ways:
//
//  - Mark by tag: untagged but undeleted manifests are collected
//  - Mark by manifest: only deleted manifests are collected
//  - Legacy: GC runs using distribution's GC
func (gc *GC) Run() error {
	var chain []func() error

	log.WithField("mode", gc.mode).Info("starting gc")

	switch gc.mode {
	case ModeByTag:
		chain = []func() error{
			gc.newRegistry,
			gc.markByTag,
			gc.sweep,
			gc.pruneManifests,
		}
	case ModeByManifest:
		chain = []func() error{
			gc.newRegistry,
			gc.markByManifest,
			gc.sweep,
		}
	case ModeLegacy:
		chain = []func() error{
			gc.newRegistry,
			gc.legacyMarkAndSweep,
		}
	}

	for _, flow := range chain {
		if err := flow(); err != nil {
			return stacktrace.Propagate(err, "")
		}
	}

	return nil
}

// markByTag runs GC by marking all manifests referenced by undeleted, valid
// tags only. This means that any manifest not referenced by a tag will be
// marked for deletion and will be deleted.
func (gc *GC) markByTag() error {
	var err error
	// store all tags for future use when pruning dangling manifests
	if gc.tags, err = gc.tagMgr.AllTags(); err != nil {
		return stacktrace.Propagate(err, "")
	}
	log.WithField("tagLen", len(gc.tags)).Info("starting mark on tags")
	for _, t := range gc.tags {
		if err = gc.markManifest(&t.Manifest); err != nil {
			return stacktrace.Propagate(err, "")
		}
	}
	return nil
}

// markByManifest marks all blobs by iterating through undeleted manifests.
// This means that any dangling manifests will also have their blobs marked
// and will not be cleaned up.
func (gc *GC) markByManifest() error {
	manifests, err := gc.mfstMgr.GetAllManifests(gc)
	if err != nil {
		return stacktrace.Propagate(err, "")
	}
	log.WithField("manifestLen", len(manifests)).Info("starting mark on manifests")
	for _, m := range manifests {
		if err = gc.markManifest(m); err != nil {
			return stacktrace.Propagate(err, "")
		}
	}
	return nil
}

// markManifest takes a distribution.Manifest and adds all blobs
// from the manifest into the mark set.
//
// This works for both v1 and v2 manifests.
func (gc *GC) markManifest(m *schema.Manifest) error {
	// TODO: unit test this with v1 and v2 manifests stored in metadata store
	// mark the manifest blob so the manifest itself is kept in the blobstore
	gc.markSet[m.Digest] = struct{}{}
	log.Debugf("marking manifest %s", m.Digest)

	// mark the config blob for v2 manifests; this should never be deleted
	gc.markSet[m.ConfigDigest] = struct{}{}
	log.Debugf("marking manifest config blob %s", m.ConfigDigest)

	for _, layer := range m.Layers {
		gc.markSet[layer] = struct{}{}
		log.Debugf("marking blob %s", layer)
	}

	return nil
}

// sweep iterates through the blobstore and deletes any blobs which haven't
// been referenced during the mark phase of GC
func (gc *GC) sweep() error {
	blobService := gc.registry.Blobs()
	blobService.Enumerate(gc, func(dgst digest.Digest) error {
		str := dgst.String()
		if _, ok := gc.markSet[str]; !ok {
			gc.deleteSet[str] = struct{}{}
		}
		return nil
	})

	log.WithFields(log.Fields{
		"markSetLen":   len(gc.markSet),
		"deleteSetLen": len(gc.deleteSet),
	}).Info("marked blobs as eligible for deletion")

	// create our error type which records blob errors
	errs := new(SweepError)
	// begin iteration over all of the delete set and start deleting blobs
	vacuum := storage.NewVacuum(gc, gc.driver)
	for dgst := range gc.deleteSet {
		log.WithField("blob", dgst).Debugf("deleting blob")
		if err := vacuum.RemoveBlob(dgst); err != nil {
			errs.AddBlobError(dgst, err)
			log.WithFields(log.Fields{
				"error": err,
				"blob":  dgst,
			}).Errorf("failed to delete blob")
			continue
		}
		if err := events.GCDeleteBlobEvent(gc.eventMgr, dgst); err != nil {
			log.WithField("error", err).Errorf("failed to create GC blob deletion event")
		}
	}

	// only return an error if any of the blobs failed to delete. During
	// MarkByTag mode this will ensure we fail here - we won't prune manifests
	// potentially incorrectly, leaving orphaned blobs.
	if len(errs.Blobs) > 0 {
		return errs
	}

	return nil
}

// pruneManifests marks dangling manifests which have been sweeped by GC as
// deleted in the database. This is a complement to the `markByTag` phase of
// GC; this marks only referenced manifests.
//
// The strategy is to:
//
//  - fetch all manifest hashes
//  - create a `referencedManifest` set
//    (via iterating over all manifests fetched during `markByTag`)
//  - create a `danglingManifest` set
//    (via iterating over all manifests and checking  the existence of the
//    manifest's hash in `referencedManifest`)
//  - Mark all manifests in `danglingManifest` as deleted within mfstMgr
//
// Because we store all manifests in the tagstore specific to repositories (by
// concatenating the repository and diget together ie `namespace/repo@sha256:...`),
// we need to operate over the DigestPK value to delete manifests.
//
// This is necessary if we ever show manifests in the UI for consistency with
// the ACTUAL blobs stored, but is not necessary for registry to work.
func (gc *GC) pruneManifests() error {
	log.Info("pruning manifests")

	manifests, err := gc.mfstMgr.GetAllManifests(gc)
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	log.WithField("totalManifests", len(manifests)).Info("loaded manifests")

	var refMnfsts = map[string]struct{}{}
	var delMnfsts = map[string]struct{}{}

	for _, t := range gc.tags {
		refMnfsts[t.DigestPK] = struct{}{}
	}

	// if the manifest is not referenced within our tag set we can safely
	// mark for deletion
	for _, m := range manifests {
		if _, ok := refMnfsts[m.PK]; !ok {
			delMnfsts[m.PK] = struct{}{}
		}
	}

	// delete all manifests. note that the rethink driver expects a slice
	// of interfaces, so instead of creating a slice of strings then converting
	// by iterating through we'll just go ahead and make a slice of interfaces
	// now
	delSlice := make([]interface{}, len(delMnfsts))
	var i int
	for dgst := range delMnfsts {
		delSlice[i] = dgst
		i++
	}

	log.WithFields(log.Fields{
		"manifestLen": len(delSlice),
	}).Info("marked manifests eligble for deletion")

	if err = gc.mfstMgr.DeleteManifests(gc, delSlice); err != nil {
		return stacktrace.Propagate(err, "unable to delete all manifests")
	}
	return nil
}

// legacyMarkAndSweep runs the garbage collector within docker/distribution.
// It is included for backwards compatibility and as a fail-safe only.
func (gc *GC) legacyMarkAndSweep() error {
	if err := storage.MarkAndSweep(gc, gc.driver, gc.registry, false); err != nil {
		return stacktrace.Propagate(err, "")
	}
	return nil
}

// newRegistry instatiates a new driver and registry based off of
// the previously set `config configuration.Configuration` field
func (gc *GC) newRegistry() error {
	var err error

	gc.driver, err = factory.Create(gc.config.Storage.Type(), gc.config.Storage.Parameters())
	if err != nil {
		return stacktrace.Propagate(err, "failed to construct %s driver", gc.config.Storage.Type())
	}

	k, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return stacktrace.Propagate(err, "unable to generate schema1 signature key")
	}

	if gc.registry, err = storage.NewRegistry(gc, gc.driver, storage.DisableSchema1Signatures, storage.Schema1SigningKey(k)); err != nil {
		return stacktrace.Propagate(err, "failed to construct registry")
	}

	return nil
}
