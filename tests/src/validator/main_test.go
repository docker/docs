package main

import (
	"os"
	"testing"
)

var docsHtmlWithoutRedirects = "/docs-without-redirects"
var docsHtmlWithRedirects = "/docs-with-redirects"
var docsSource = "/docs-source"

// TestMain is used to add extra setup or
// teardown before or after testing
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
