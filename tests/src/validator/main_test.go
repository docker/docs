package main

import (
	"os"
	"testing"
)

// TestMain is used to add extra setup or
// teardown before or after testing
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
