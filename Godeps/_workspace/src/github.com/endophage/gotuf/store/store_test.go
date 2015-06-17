package store

import "bytes"
import "reflect"
import "testing"

func testMeta(t *testing.T, s *MetadataStore) {
	storeName := reflect.TypeOf(s).Name()
	testData := []byte{}
	err := s.SetMeta(testName, testData)
	if err != nil {
		t.Fatal(err)
	}
	out, err := s.GetMeta(testName, len(testData))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(testData, out) != 0 {
		t.Fatalf("%s mangled data", storeName)
	}
}
