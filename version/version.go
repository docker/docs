package version

var (
	Version    = "1.2.0"
	PreRelease = ""
	GitCommit  = "HEAD"
	BuildTime  = "<unknown>"
)

func FullVersion() string {
	if PreRelease != "" {
		return Version + "-" + PreRelease + " (" + GitCommit + ")"
	} else {
		return Version + " (" + GitCommit + ")"
	}
}

func TagVersion() string {
	if PreRelease != "" {
		return Version + "-" + PreRelease
	} else {
		return Version
	}
}
