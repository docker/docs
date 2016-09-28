package strings

import (
	"strconv"
	"strings"
)

const (
	literalTrue        string = "true"
	literalFalse       string = "false"
	literalNull        string = "null"
	literalDoubleQuote string = `"`
	literalSingleQuote string = `'`
)

// Parse tries to create a native object from the given
// string, or just returns the string if nothing takes.
//
// Values wrapped in "quotes" or 'single quotes' will always
// be treated as a string, but the quotes will be removed.
//
// This method knows about all number types, and will always
// look for the smallest type to fit the number.  It also handles
// the boolean literals 'true' and 'false'.
//
// An empty string ("") will return nil.
func Parse(s string) interface{} {

	// nothing is nil
	if len(s) == 0 {
		return nil
	}

	/*
	   Is it forced to be a string with quotes?
	*/
	if strings.HasPrefix(s, literalDoubleQuote) && strings.HasSuffix(s, literalDoubleQuote) {
		return strings.Trim(s, literalDoubleQuote)
	}
	if strings.HasPrefix(s, literalSingleQuote) && strings.HasSuffix(s, literalSingleQuote) {
		return strings.Trim(s, literalSingleQuote)
	}

	/*
	   Check literals
	*/
	switch strings.ToLower(s) {
	/*
	   Booleans
	*/
	case literalTrue:
		return true
	case literalFalse:
		return false
	/*
		Other
	*/
	case literalNull:
		return nil
	}

	/*
	   Numbers
	*/

	// try int (most common type)
	if val, err := strconv.ParseInt(s, 10, 0); err == nil {
		return int(val)
	}

	/*
	   ints
	*/

	// try int8
	if val, err := strconv.ParseInt(s, 10, 8); err == nil {
		return val
	}

	// try int16
	if val, err := strconv.ParseInt(s, 10, 16); err == nil {
		return val
	}

	// try int32
	if val, err := strconv.ParseInt(s, 10, 32); err == nil {
		return val
	}

	// try int64
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}

	/*
	   uints
	*/
	// try uint8
	if val, err := strconv.ParseUint(s, 10, 8); err == nil {
		return val
	}

	// try uint16
	if val, err := strconv.ParseUint(s, 10, 16); err == nil {
		return val
	}

	// try uint32
	if val, err := strconv.ParseUint(s, 10, 32); err == nil {
		return val
	}

	// try uint64
	if val, err := strconv.ParseUint(s, 10, 64); err == nil {
		return val
	}

	/*
	   floats
	*/

	// try float32
	if val, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(val)
	}

	// try float64
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}

	/*
	   Nothing - just return the string
	*/
	return s
}
