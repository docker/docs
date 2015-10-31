// +build pkcs11,darwin

package api

var (
	pkcs11Lib = "/usr/local/lib/libykcs11.dylib"
)

func init() {
	// TODO(diogomonica): all the crap for darwin to configure
	// the variable pkcs11 to find the right one in the right dir
}
