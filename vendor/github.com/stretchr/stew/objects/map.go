package objects

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/signature"
	"io/ioutil"
	"net/url"
	"strings"
)

var (
	// PathSeparator is the character used to separate the elements
	// of the keypath.
	//
	// For example, `location.address.city`
	PathSeparator string = "."

	// SignatureSeparator is the character that is used to
	// separate the Base64 string from the security signature.
	SignatureSeparator = "_"
)

// Map is a map[string]interface{} with additional helpful functionality.
//
// You can use Map functionality on any map[string]interface{} using the following
// format:
//
//     data := map[string]interface{}{"name": "Stew"}
//     objects.Map(data).Get("name")
//     // returns "Stew"
type Map map[string]interface{}

// NewMap creates a new map.  You may also use the M shortcut method.
//
// The arguments follow a key, value pattern.
//
// Panics
//
// Panics if any key arugment is non-string or if there are an odd number of arguments.
//
// Example
//
// To easily create Maps:
//
//     m := objects.M("name", "Mat", "age", 29, "subobj", objects.M("active", true))
//
//     // creates a Map equivalent to
//     m := map[string]interface{}{"name": "Mat", "age": 29, "subobj": map[string]interface{}{"active": true}}
func NewMap(keyAndValuePairs ...interface{}) Map {

	newMap := make(Map)
	keyAndValuePairsLen := len(keyAndValuePairs)

	if keyAndValuePairsLen%2 != 0 {
		panic("NewMap must have an even number of arguments following the 'key, value' pattern.")
	}

	for i := 0; i < keyAndValuePairsLen; i = i + 2 {

		key := keyAndValuePairs[i]
		value := keyAndValuePairs[i+1]

		// make sure the key is a string
		keyString, keyStringOK := key.(string)
		if !keyStringOK {
			panic(fmt.Sprintf("NewMap must follow 'string, interface{}' pattern.  %s is not a valid key.", keyString))
		}

		newMap[keyString] = value

	}

	return newMap
}

// M is a shortcut method for NewMap.
func M(keyAndValuePairs ...interface{}) Map {
	return NewMap(keyAndValuePairs...)
}

// NewMapFromJSON creates a new map from a JSON string representation
func NewMapFromJSON(data string) (Map, error) {

	var unmarshalled map[string]interface{}

	err := json.Unmarshal([]byte(data), &unmarshalled)

	if err != nil {
		return nil, errors.New("Map: JSON decode failed with: " + err.Error())
	}

	return Map(unmarshalled), nil

}

// NewMapFromBase64String creates a new map from a Base64 string representation
func NewMapFromBase64String(data string) (Map, error) {

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))

	decoded, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}

	return NewMapFromJSON(string(decoded))

}

// NewMapFromSignedBase64String creates a new map from a signed Base64 string representation
func NewMapFromSignedBase64String(data, key string) (Map, error) {

	parts := strings.Split(data, SignatureSeparator)
	if len(parts) != 2 {
		return nil, errors.New("Map: Signed base64 string is malformed.")
	}

	sig := signature.HashWithKey([]byte(parts[0]), []byte(key))
	if parts[1] != sig {
		return nil, errors.New("Map: Signature for Base64 data does not match.")
	}

	return NewMapFromBase64String(parts[0])

}

// Get gets the value from the map.  Supports deep nesting of other maps,
// For example:
//
//     m = Map{"name":Map{"First": "Mat", "Last": "Ryer"}}
//
//     m.Get("name.Last")
//     // returns "Ryer"
func (d Map) Get(keypath string) interface{} {

	var segs []string = strings.Split(keypath, PathSeparator)

	obj := d

	for fieldIndex, field := range segs {

		if fieldIndex == len(segs)-1 {
			return obj[field]
		}

		switch obj[field].(type) {
		case Map:
			obj = obj[field].(Map)
		case map[string]interface{}:
			obj = Map(obj[field].(map[string]interface{}))
		}

	}

	return obj

}

// GetMap gets another Map from this one, or panics if the object is missing or not a Map.
func (d Map) GetMap(keypath string) Map {
	return d.Get(keypath).(Map)
}

// GetString gets a string value from the map at the given keypath, or panics if one
// is not available, or is of the wrong type.
func (d Map) GetString(keypath string) string {
	return d.Get(keypath).(string)
}

// GetWithDefault gets the value at the specified keypath, or returns the defaultValue if
// none could be found.
func (d Map) GetOrDefault(keypath string, defaultValue interface{}) interface{} {
	obj := d.Get(keypath)
	if obj == nil {
		return defaultValue
	}
	return obj
}

// GetStringOrDefault gets the string value at the specified keypath,
// or returns the defaultValue if none could be found.  Will panic if the
// object is there but of the wrong type.
func (d Map) GetStringOrDefault(keypath, defaultValue string) string {
	obj := d.Get(keypath)
	if obj == nil {
		return defaultValue
	}
	return obj.(string)
}

// GetStringOrEmpty gets the string value at the specified keypath or returns
// an empty string if none could be fo und. Will panic if the object is there
// but of the wrong type.
func (d Map) GetStringOrEmpty(keypath string) string {
	return d.GetStringOrDefault(keypath, "")
}

// Set sets a value in the map.  Supports dot syntax to set deep values.
//
// For example,
//
//     m.Set("name.first", "Mat")
//
// The above code sets the 'first' field on the 'name' object in the m Map.
//
// If objects are nil along the way, Set creates new Map objects as needed.
func (d Map) Set(keypath string, value interface{}) Map {

	var segs []string
	segs = strings.Split(keypath, PathSeparator)

	obj := d

	for fieldIndex, field := range segs {

		if fieldIndex == len(segs)-1 {
			obj[field] = value
		}

		if _, exists := obj[field]; !exists {
			obj[field] = make(Map)
			obj = obj[field].(Map)
		} else {
			switch obj[field].(type) {
			case Map:
				obj = obj[field].(Map)
			case map[string]interface{}:
				obj = Map(obj[field].(map[string]interface{}))
			}
		}

	}

	// chain
	return d
}

// Exclude returns a new Map with the keys in the specified []string
// excluded.
func (d Map) Exclude(exclude []string) Map {

	excluded := make(Map)
	for k, v := range d {
		var shouldInclude bool = true
		for _, toExclude := range exclude {
			if k == toExclude {
				shouldInclude = false
				break
			}
		}
		if shouldInclude {
			excluded[k] = v
		}
	}

	return excluded
}

// Copy creates a shallow copy of the Map.
func (d Map) Copy() Map {
	copied := make(Map)
	for k, v := range d {
		copied[k] = v
	}
	return copied
}

// Merge blends the specified map with a copy of this map and returns the result.
//
// Keys that appear in both will be selected from the specified map.
func (d Map) Merge(merge Map) Map {
	return d.Copy().MergeHere(merge)
}

// Merge blends the specified map with this map and returns the current map.
//
// Keys that appear in both will be selected from the specified map.  The original map
// will be modified.
func (d Map) MergeHere(merge Map) Map {

	for k, v := range merge {
		d[k] = v
	}

	return d

}

// Has gets whether the Map has the specified field or not. Supports deep nesting of other maps.
//
// For example:
//     m := map[string]interface{}{"parent": map[string]interface{}{"childname": "Luke"}}
//     m.Has("parent.childname")
//     // return true
func (d Map) Has(path string) bool {
	return d.Get(path) != nil
}

// MSI is a shortcut method to get the current map as a
// normal map[string]interface{}.
func (d Map) MSI() map[string]interface{} {
	return map[string]interface{}(d)
}

// JSON converts the map to a JSON string
func (d Map) JSON() (string, error) {

	result, err := json.Marshal(d)

	if err != nil {
		err = errors.New("Map: JSON encode failed with: " + err.Error())
	}

	return string(result), err

}

// Base64 converts the map to a base64 string
func (d Map) Base64() (string, error) {

	var buf bytes.Buffer

	jsonData, err := d.JSON()
	if err != nil {
		return "", err
	}

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(jsonData))
	encoder.Close()

	return buf.String(), nil

}

// SignedBase64 converts the map to a base64 string and signs it using the
// provided key. The returned data is the base64 string plus an appended signature.
//
// Will return an error if Base64ing the map fails.
func (d Map) SignedBase64(key string) (string, error) {

	base64, err := d.Base64()
	if err != nil {
		return "", err
	}

	sig := signature.HashWithKey([]byte(base64), []byte(key))

	return base64 + SignatureSeparator + sig, nil

}

// Hash gets the hash of the map with no security key.
//
// Will return an error if Base64ing the map fails.
func (d Map) Hash() (string, error) {
	return d.HashWithKey("")
}

// HashWithKey gets the a hash of the map, signed by the
// specified security key.
//
// Will return an error if Base64ing the map fails.
func (d Map) HashWithKey(key string) (string, error) {

	base64, err := d.Base64()
	if err != nil {
		return "", err
	}

	sig := signature.HashWithKey([]byte(base64), []byte(key))

	return sig, nil

}

/*
	URL Query
	------------------------------------------------
*/

// NewMapFromURLQuery generates a new map by parsing the specified
// query.
//
// For queries with multiple values, the first value is selected.
func NewMapFromURLQuery(query string) (Map, error) {

	vals, err := url.ParseQuery(query)

	if err != nil {
		return nil, err
	}

	return NewMapFromURLValues(vals)
}

func (d Map) URLValues() url.Values {

	vals := make(url.Values)

	for k, v := range d {
		vals.Set(k, fmt.Sprintf("%v", v))
	}

	return vals
}

func NewMapFromURLValues(vals url.Values) (Map, error) {
	m := make(Map)
	for k, vals := range vals {
		m.Set(k, vals[0])
	}

	return m, nil
}

// URLQuery gets an encoded URL query representing the given
// map.
func (d Map) URLQuery() (string, error) {
	return d.URLValues().Encode(), nil
}

// Transform builds a new map giving the transformer a chance
// to change the keys and values as it goes.
func (d Map) Transform(transformer func(key string, value interface{}) (string, interface{})) Map {
	m := make(Map)
	for k, v := range d {
		modifiedKey, modifiedVal := transformer(k, v)
		m[modifiedKey] = modifiedVal
	}
	return m
}

// TransformKeys builds a new map using the specified key mapping.
//
// Unspecified keys will be unaltered.
func (d Map) TransformKeys(mapping map[string]string) Map {
	return d.Transform(func(key string, value interface{}) (string, interface{}) {

		if newKey, ok := mapping[key]; ok {
			return newKey, value
		}

		return key, value
	})
}
