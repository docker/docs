package strings

import (
	"bytes"
)

// MergeBytes merges an array of []byte strings into one []byte.
//
// Example
//
//     one := []byte("Hello")
//     two := []byte(" ")
//     three := []byte("Stretchr!")
//
//     // merge the bytes
//     mergedBytes := strings.MergeBytes(one, two, three)
//
//     log.Print(mergedBytes)
//     // prints "Hello Stretchr!"
func MergeBytes(byteArray ...[]byte) []byte {

	var buffer bytes.Buffer
	for _, v := range byteArray {
		buffer.Write(v)
	}
	return buffer.Bytes()

}

// JoinStrings joins many []byte strings together separated by the specified separator.
func JoinBytes(separator []byte, byteArray ...[]byte) []byte {

	var buffer bytes.Buffer
	var max int = len(byteArray) - 1
	for vi, v := range byteArray {
		buffer.Write(v)
		if vi < max {
			buffer.Write(separator)
		}
	}
	return buffer.Bytes()

}
