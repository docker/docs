// +build pkcs11,linux

package api

var (
	pkcs11Lib = "/usr/local/lib/libykcs11.so"
)

func init() {
	// TODO(diogomonica): all the crap for linux to configure
	// the variable pkcs11 to find the right one in the right dir
}
