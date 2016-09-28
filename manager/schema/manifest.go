package schema

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/dhe-deploy"
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	containertypes "github.com/docker/engine-api/types/container"
)

const digestSHA256GzippedEmptyTar = digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")

var manifestsTable = table{
	db:         deploy.DTRDBName,
	name:       "manifests",
	primaryKey: "pk",
	secondaryIndexes: map[string][]string{
		"deleted": nil,
		// For ordering manifests which are not deleted by created at
		"deleted_repository_createdAt": {"deleted", "repository", "createdAt"},
	},
}

// Manifest represents an image manifest stored within the registry.
//
// Within DTR the rethink database is the source of truth for image metadata;
// we wrap the registry with middleware which saves all tags and manifests in
// the database.
//
// This allows for much faster GC, easier tag deletion and a richer API for
// listing tags.
//
// XXX: In the future we may store ManifestLists, which have many architectures
//      and OS types.
type Manifest struct {
	// PK is the namespace/repo concatenated with the digest of this manifest
	PK string `gorethink:"pk"`
	// Digest is the digest of the actual manifest
	Digest string `gorethink:"digest"`
	// Repository is the namespace/repository string for this manifest. It's
	// stored for manifest lookups by repository.
	Repository string `gorethink:"repository"`
	// OS
	OS string `gorethink:"os"`
	// Architecture lists what arch the manifest runs on
	Architecture string `gorethink:"architecture"`
	// MediaType stores the media type of the manifest (which encompasses
	// v1, v2, and lists)
	MediaType string `gorethink:"mediaType"`
	// Payload stores the actual manifest
	Payload []byte `gorethink:"payload"`
	// Layers stores a list of layer hashes used in this manifest, if this is
	// not a manifest list. This is stored in V2 order (index 0 == base image)
	Layers []string `gorethink:"layers"`
	// ConfigDigest stores the checksum of the `config` file for V2 manifests.
	// We need to ensure we store this so we can mark it during GC phases.
	ConfigDigest string `gorethink:"configDigest"`
	// Config stores the config blob contents for V2 manifests. This is
	// necessary for some nautilus operations.
	Config []byte `gorethink:"config"`
	// Size of the image in bytes. This is a sum of each layer's size.
	Size int64 `gorethink:"size"`
	// OriginalAuthor stores the name of the account that first pushed this
	// manifest.
	OriginalAuthor string    `gorethink:"originalAuthor"`
	CreatedAt      time.Time `gorethink:"createdAt"`
	// IsDeleted is a tombstone marker which depicts whether the manifest is
	// marked for deletion.
	IsDeleted bool `gorethink:"deleted"`
}

// LoadMetadata takes a v1 or v2 distribtion manifest and returns a
// filled *Manifest type by loading the os, arch, media type, layers
// and size from the concrete distribution manifest.
func (m *Manifest) LoadMetadata(ctx context.Context, manifest distribution.Manifest) (err error) {
	mediaType, payload, err := manifest.Payload()
	if err != nil {
		return err
	}
	// Only set mediatype and payload if there
	m.MediaType = mediaType
	m.Payload = payload

	// Tackle manifest metadata such as layers, arch and OS
	switch manifest.(type) {
	case *schema1.SignedManifest:
		return m.loadV1Metadata(ctx, manifest.(*schema1.SignedManifest))
	case *schema2.DeserializedManifest:
		return m.loadV2Metadata(ctx, manifest.(*schema2.DeserializedManifest))
	case *manifestlist.DeserializedManifestList:
		// No metadata to retrieve
		return nil
	default:
		return fmt.Errorf("unknown manifest type %T", manifest)
	}

	return nil
}

func (m *Manifest) loadV1Metadata(ctx context.Context, manifest *schema1.SignedManifest) (err error) {
	// The OS and arch are hidden within the first entry of "history"
	if len(manifest.History) > 0 {
		if err = m.parseConfig([]byte(manifest.History[0].V1Compatibility)); err != nil {
			return err
		}
	}

	m.Layers = []string{}
	// V1 layers are stored such that the base image is the last entry in the
	// array. We'll reverse that order to keep things standard with V2
	for i := len(manifest.FSLayers) - 1; i >= 0; i-- {
		m.Layers = append(m.Layers, string(manifest.FSLayers[i].BlobSum))
		// Note that V1 manifests did not include layer size within the manifest
		// itself so we can't get the manifest size. Upgrade, yo.
	}

	return nil
}

func (m *Manifest) loadV2Metadata(ctx context.Context, manifest *schema2.DeserializedManifest) (err error) {
	// The config for V2 manifests is stored within a separate blob. Our
	// middleware verifies that this blob exists by retrieving the blob (vs Stat
	// for layers). This is the stored in "target".
	config, ok := ctx.Value("target").([]byte)
	if !ok {
		return fmt.Errorf("unable to retrieve metadata target info")
	}

	if err = m.parseConfig(config); err != nil {
		return err
	}

	m.Layers = make([]string, len(manifest.Layers))
	m.ConfigDigest = manifest.Target().Digest.String()
	for i, l := range manifest.Layers {
		m.Size += l.Size
		m.Layers[i] = string(l.Digest)
	}

	return nil
}

// parseConfig adds OS and Architecture given a marshaled config struct
func (m *Manifest) parseConfig(config []byte) error {
	var c ManifestConfig
	if err := json.Unmarshal(config, &c); err != nil {
		return err
	}
	m.OS = c.OS
	m.Architecture = c.Architecture
	m.Config = config
	return nil
}

// Config represents docker configuration for a manifest. This information is
// stored in the first History[] entry in V1 manifests and a separate Config
// blob in V2 manifests. Only V2 Config types are fully unmarshalled;
// marshalling and unmarshalling Config for V1 manifests **will** result in
// data loss.
//
// NOTE: This config is not part of distribution and is treated as opaque;
//       it comes directly from the docker client. We only rely on this for OS
//       and architecture.
type ManifestConfig struct {
	Architecture    string                `json:"architecture"`
	Author          string                `json:"author"`
	Container       string                `json:"container"`
	Created         time.Time             `json:"created"`
	DockerVersion   string                `json:"docker_version"`
	OS              string                `json:"os"`
	Config          containertypes.Config `json:"config"`
	ContainerConfig containertypes.Config `json:"container_config"`
	RootFS          FS                    // V2 only
	// Note that other V1 fields exist ("throwaway", "parent") but are ignored
}

type FS struct {
	Type  string   `json:"type"`
	Diffs []string `json:"diff_ids"` // slice of "sha256:xyz" strings
}

type History struct {
	Created    time.Time `json:"created"`
	Author     string    `json:"author,omitempty"`
	CreatedBy  string    `json:"created_by,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	EmptyLayer bool      `json:"empty_layer,omitempty"`
}
