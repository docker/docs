package version

// Package is the overall, canonical project import path under which the
// package was built.
var Package = "github.com/docker/pinata/v1"

// Version indicates which version of the binary is running. This is set to
// the latest release tag by hand, always suffixed by "+unknown". During
// build, it will be replaced by the actual version.
var Version = "v0.0.0+unknown"
