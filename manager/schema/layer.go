package schema

type Layer struct {
	// PK is the digest of the layer's contents
	PK string `gorethink:"pk"`
	// Size is the size of the layer in bytes
	Size uint64 `gorethink:"size"`
	// MediaType is the media type of the layer, typically
	// "application/vnd.docker.image.rootfs.diff.tar.gzip"
	MediaType string `gorethink:"mediaType"`
}
