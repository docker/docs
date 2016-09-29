package nsuserdefaults

import (
	// "github.com/Sirupsen/logrus"
	// "strconv"
	"unsafe"
)

/*
#include <stdlib.h>     // for free()
#import "preferences.h"
#cgo CFLAGS: -fobjc-arc -x objective-c -I .
#cgo LDFLAGS: -framework Foundation
*/
import "C"

// KeyExists checks if a key exists
func KeyExists(key string) bool {
	var cKey = C.CString(key)
	var keyExists = int(C.keyExists(cKey))
	C.free(unsafe.Pointer(cKey))
	return keyExists == 1
}

// BoolForKey returns false if key doesn't exist
func BoolForKey(key string) bool {
	var cKey = C.CString(key)
	var result = int(C.boolForKey(cKey))
	C.free(unsafe.Pointer(cKey))
	return result == 1
}
