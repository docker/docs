// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package proto

/*
 * Types and routines for supporting protocol buffer extensions.
 */

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

// ErrMissingExtension is the error returned by GetExtension if the named extension is not in the message.
var ErrMissingExtension = errors.New("proto: missing extension")

// ExtensionRange represents a range of message extensions for a protocol buffer.
// Used in code generated by the protocol compiler.
type ExtensionRange struct {
	Start, End int32 // both inclusive
}

// extendableProto is an interface implemented by any protocol buffer that may be extended.
type extendableProto interface {
	Message
	ExtensionRangeArray() []ExtensionRange
}

type extensionsMap interface {
	extendableProto
	ExtensionMap() map[int32]Extension
}

type extensionsBytes interface {
	extendableProto
	GetExtensions() *[]byte
}

var extendableProtoType = reflect.TypeOf((*extendableProto)(nil)).Elem()

// ExtensionDesc represents an extension specification.
// Used in generated code from the protocol compiler.
type ExtensionDesc struct {
	ExtendedType  Message     // nil pointer to the type that is being extended
	ExtensionType interface{} // nil pointer to the extension type
	Field         int32       // field number
	Name          string      // fully-qualified name of extension, for text formatting
	Tag           string      // protobuf tag style
}

func (ed *ExtensionDesc) repeated() bool {
	t := reflect.TypeOf(ed.ExtensionType)
	return t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
}

// Extension represents an extension in a message.
type Extension struct {
	// When an extension is stored in a message using SetExtension
	// only desc and value are set. When the message is marshaled
	// enc will be set to the encoded form of the message.
	//
	// When a message is unmarshaled and contains extensions, each
	// extension will have only enc set. When such an extension is
	// accessed using GetExtension (or GetExtensions) desc and value
	// will be set.
	desc  *ExtensionDesc
	value interface{}
	enc   []byte
}

// SetRawExtension is for testing only.
func SetRawExtension(base extendableProto, id int32, b []byte) {
	if ebase, ok := base.(extensionsMap); ok {
		ebase.ExtensionMap()[id] = Extension{enc: b}
	} else if ebase, ok := base.(extensionsBytes); ok {
		clearExtension(base, id)
		ext := ebase.GetExtensions()
		*ext = append(*ext, b...)
	} else {
		panic("unreachable")
	}
}

// isExtensionField returns true iff the given field number is in an extension range.
func isExtensionField(pb extendableProto, field int32) bool {
	for _, er := range pb.ExtensionRangeArray() {
		if er.Start <= field && field <= er.End {
			return true
		}
	}
	return false
}

// checkExtensionTypes checks that the given extension is valid for pb.
func checkExtensionTypes(pb extendableProto, extension *ExtensionDesc) error {
	// Check the extended type.
	if a, b := reflect.TypeOf(pb), reflect.TypeOf(extension.ExtendedType); a != b {
		return errors.New("proto: bad extended type; " + b.String() + " does not extend " + a.String())
	}
	// Check the range.
	if !isExtensionField(pb, extension.Field) {
		return errors.New("proto: bad extension number; not in declared ranges")
	}
	return nil
}

// extPropKey is sufficient to uniquely identify an extension.
type extPropKey struct {
	base  reflect.Type
	field int32
}

var extProp = struct {
	sync.RWMutex
	m map[extPropKey]*Properties
}{
	m: make(map[extPropKey]*Properties),
}

func extensionProperties(ed *ExtensionDesc) *Properties {
	key := extPropKey{base: reflect.TypeOf(ed.ExtendedType), field: ed.Field}

	extProp.RLock()
	if prop, ok := extProp.m[key]; ok {
		extProp.RUnlock()
		return prop
	}
	extProp.RUnlock()

	extProp.Lock()
	defer extProp.Unlock()
	// Check again.
	if prop, ok := extProp.m[key]; ok {
		return prop
	}

	prop := new(Properties)
	prop.Init(reflect.TypeOf(ed.ExtensionType), "unknown_name", ed.Tag, nil)
	extProp.m[key] = prop
	return prop
}

// encodeExtensionMap encodes any unmarshaled (unencoded) extensions in m.
func encodeExtensionMap(m map[int32]Extension) error {
	for k, e := range m {
		err := encodeExtension(&e)
		if err != nil {
			return err
		}
		m[k] = e
	}
	return nil
}

func encodeExtension(e *Extension) error {
	if e.value == nil || e.desc == nil {
		// Extension is only in its encoded form.
		return nil
	}
	// We don't skip extensions that have an encoded form set,
	// because the extension value may have been mutated after
	// the last time this function was called.

	et := reflect.TypeOf(e.desc.ExtensionType)
	props := extensionProperties(e.desc)

	p := NewBuffer(nil)
	// If e.value has type T, the encoder expects a *struct{ X T }.
	// Pass a *T with a zero field and hope it all works out.
	x := reflect.New(et)
	x.Elem().Set(reflect.ValueOf(e.value))
	if err := props.enc(p, props, toStructPointer(x)); err != nil {
		return err
	}
	e.enc = p.buf
	return nil
}

func sizeExtensionMap(m map[int32]Extension) (n int) {
	for _, e := range m {
		if e.value == nil || e.desc == nil {
			// Extension is only in its encoded form.
			n += len(e.enc)
			continue
		}

		// We don't skip extensions that have an encoded form set,
		// because the extension value may have been mutated after
		// the last time this function was called.

		et := reflect.TypeOf(e.desc.ExtensionType)
		props := extensionProperties(e.desc)

		// If e.value has type T, the encoder expects a *struct{ X T }.
		// Pass a *T with a zero field and hope it all works out.
		x := reflect.New(et)
		x.Elem().Set(reflect.ValueOf(e.value))
		n += props.size(props, toStructPointer(x))
	}
	return
}

// HasExtension returns whether the given extension is present in pb.
func HasExtension(pb extendableProto, extension *ExtensionDesc) bool {
	// TODO: Check types, field numbers, etc.?
	if epb, doki := pb.(extensionsMap); doki {
		_, ok := epb.ExtensionMap()[extension.Field]
		return ok
	} else if epb, doki := pb.(extensionsBytes); doki {
		ext := epb.GetExtensions()
		buf := *ext
		o := 0
		for o < len(buf) {
			tag, n := DecodeVarint(buf[o:])
			fieldNum := int32(tag >> 3)
			if int32(fieldNum) == extension.Field {
				return true
			}
			wireType := int(tag & 0x7)
			o += n
			l, err := size(buf[o:], wireType)
			if err != nil {
				return false
			}
			o += l
		}
		return false
	}
	panic("unreachable")
}

func deleteExtension(pb extensionsBytes, theFieldNum int32, offset int) int {
	ext := pb.GetExtensions()
	for offset < len(*ext) {
		tag, n1 := DecodeVarint((*ext)[offset:])
		fieldNum := int32(tag >> 3)
		wireType := int(tag & 0x7)
		n2, err := size((*ext)[offset+n1:], wireType)
		if err != nil {
			panic(err)
		}
		newOffset := offset + n1 + n2
		if fieldNum == theFieldNum {
			*ext = append((*ext)[:offset], (*ext)[newOffset:]...)
			return offset
		}
		offset = newOffset
	}
	return -1
}

func clearExtension(pb extendableProto, fieldNum int32) {
	if epb, doki := pb.(extensionsMap); doki {
		delete(epb.ExtensionMap(), fieldNum)
	} else if epb, doki := pb.(extensionsBytes); doki {
		offset := 0
		for offset != -1 {
			offset = deleteExtension(epb, fieldNum, offset)
		}
	} else {
		panic("unreachable")
	}
}

// ClearExtension removes the given extension from pb.
func ClearExtension(pb extendableProto, extension *ExtensionDesc) {
	// TODO: Check types, field numbers, etc.?
	clearExtension(pb, extension.Field)
}

// GetExtension parses and returns the given extension of pb.
// If the extension is not present it returns ErrMissingExtension.
func GetExtension(pb extendableProto, extension *ExtensionDesc) (interface{}, error) {
	if err := checkExtensionTypes(pb, extension); err != nil {
		return nil, err
	}

	if epb, doki := pb.(extensionsMap); doki {
		emap := epb.ExtensionMap()
		e, ok := emap[extension.Field]
		if !ok {
			// defaultExtensionValue returns the default value or
			// ErrMissingExtension if there is no default.
			return defaultExtensionValue(extension)
		}
		if e.value != nil {
			// Already decoded. Check the descriptor, though.
			if e.desc != extension {
				// This shouldn't happen. If it does, it means that
				// GetExtension was called twice with two different
				// descriptors with the same field number.
				return nil, errors.New("proto: descriptor conflict")
			}
			return e.value, nil
		}

		v, err := decodeExtension(e.enc, extension)
		if err != nil {
			return nil, err
		}

		// Remember the decoded version and drop the encoded version.
		// That way it is safe to mutate what we return.
		e.value = v
		e.desc = extension
		e.enc = nil
		emap[extension.Field] = e
		return e.value, nil
	} else if epb, doki := pb.(extensionsBytes); doki {
		ext := epb.GetExtensions()
		o := 0
		for o < len(*ext) {
			tag, n := DecodeVarint((*ext)[o:])
			fieldNum := int32(tag >> 3)
			wireType := int(tag & 0x7)
			l, err := size((*ext)[o+n:], wireType)
			if err != nil {
				return nil, err
			}
			if int32(fieldNum) == extension.Field {
				v, err := decodeExtension((*ext)[o:o+n+l], extension)
				if err != nil {
					return nil, err
				}
				return v, nil
			}
			o += n + l
		}
		return defaultExtensionValue(extension)
	}
	panic("unreachable")
}

// defaultExtensionValue returns the default value for extension.
// If no default for an extension is defined ErrMissingExtension is returned.
func defaultExtensionValue(extension *ExtensionDesc) (interface{}, error) {
	t := reflect.TypeOf(extension.ExtensionType)
	props := extensionProperties(extension)

	sf, _, err := fieldDefault(t, props)
	if err != nil {
		return nil, err
	}

	if sf == nil || sf.value == nil {
		// There is no default value.
		return nil, ErrMissingExtension
	}

	if t.Kind() != reflect.Ptr {
		// We do not need to return a Ptr, we can directly return sf.value.
		return sf.value, nil
	}

	// We need to return an interface{} that is a pointer to sf.value.
	value := reflect.New(t).Elem()
	value.Set(reflect.New(value.Type().Elem()))
	if sf.kind == reflect.Int32 {
		// We may have an int32 or an enum, but the underlying data is int32.
		// Since we can't set an int32 into a non int32 reflect.value directly
		// set it as a int32.
		value.Elem().SetInt(int64(sf.value.(int32)))
	} else {
		value.Elem().Set(reflect.ValueOf(sf.value))
	}
	return value.Interface(), nil
}

// decodeExtension decodes an extension encoded in b.
func decodeExtension(b []byte, extension *ExtensionDesc) (interface{}, error) {
	o := NewBuffer(b)

	t := reflect.TypeOf(extension.ExtensionType)
	rep := extension.repeated()

	props := extensionProperties(extension)

	// t is a pointer to a struct, pointer to basic type or a slice.
	// Allocate a "field" to store the pointer/slice itself; the
	// pointer/slice will be stored here. We pass
	// the address of this field to props.dec.
	// This passes a zero field and a *t and lets props.dec
	// interpret it as a *struct{ x t }.
	value := reflect.New(t).Elem()

	for {
		// Discard wire type and field number varint. It isn't needed.
		if _, err := o.DecodeVarint(); err != nil {
			return nil, err
		}

		if err := props.dec(o, props, toStructPointer(value.Addr())); err != nil {
			return nil, err
		}

		if !rep || o.index >= len(o.buf) {
			break
		}
	}
	return value.Interface(), nil
}

// GetExtensions returns a slice of the extensions present in pb that are also listed in es.
// The returned slice has the same length as es; missing extensions will appear as nil elements.
func GetExtensions(pb Message, es []*ExtensionDesc) (extensions []interface{}, err error) {
	epb, ok := pb.(extendableProto)
	if !ok {
		err = errors.New("proto: not an extendable proto")
		return
	}
	extensions = make([]interface{}, len(es))
	for i, e := range es {
		extensions[i], err = GetExtension(epb, e)
		if err == ErrMissingExtension {
			err = nil
		}
		if err != nil {
			return
		}
	}
	return
}

// SetExtension sets the specified extension of pb to the specified value.
func SetExtension(pb extendableProto, extension *ExtensionDesc, value interface{}) error {
	if err := checkExtensionTypes(pb, extension); err != nil {
		return err
	}
	typ := reflect.TypeOf(extension.ExtensionType)
	if typ != reflect.TypeOf(value) {
		return errors.New("proto: bad extension value type")
	}
	// nil extension values need to be caught early, because the
	// encoder can't distinguish an ErrNil due to a nil extension
	// from an ErrNil due to a missing field. Extensions are
	// always optional, so the encoder would just swallow the error
	// and drop all the extensions from the encoded message.
	if reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("proto: SetExtension called with nil value of type %T", value)
	}
	return setExtension(pb, extension, value)
}

func setExtension(pb extendableProto, extension *ExtensionDesc, value interface{}) error {
	if epb, doki := pb.(extensionsMap); doki {
		epb.ExtensionMap()[extension.Field] = Extension{desc: extension, value: value}
	} else if epb, doki := pb.(extensionsBytes); doki {
		ClearExtension(pb, extension)
		ext := epb.GetExtensions()
		et := reflect.TypeOf(extension.ExtensionType)
		props := extensionProperties(extension)
		p := NewBuffer(nil)
		x := reflect.New(et)
		x.Elem().Set(reflect.ValueOf(value))
		if err := props.enc(p, props, toStructPointer(x)); err != nil {
			return err
		}
		*ext = append(*ext, p.buf...)
	}
	return nil
}

// A global registry of extensions.
// The generated code will register the generated descriptors by calling RegisterExtension.

var extensionMaps = make(map[reflect.Type]map[int32]*ExtensionDesc)

// RegisterExtension is called from the generated code.
func RegisterExtension(desc *ExtensionDesc) {
	st := reflect.TypeOf(desc.ExtendedType).Elem()
	m := extensionMaps[st]
	if m == nil {
		m = make(map[int32]*ExtensionDesc)
		extensionMaps[st] = m
	}
	if _, ok := m[desc.Field]; ok {
		panic("proto: duplicate extension registered: " + st.String() + " " + strconv.Itoa(int(desc.Field)))
	}
	m[desc.Field] = desc
}

// RegisteredExtensions returns a map of the registered extensions of a
// protocol buffer struct, indexed by the extension number.
// The argument pb should be a nil pointer to the struct type.
func RegisteredExtensions(pb Message) map[int32]*ExtensionDesc {
	return extensionMaps[reflect.TypeOf(pb).Elem()]
}
