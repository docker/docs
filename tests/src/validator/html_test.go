package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var countLinks = 0
var countImages = 0

// TestURLs tests if we're not using absolute paths for URLs
// when pointing to local pages.
func TestURLs(t *testing.T) {
	count := 0

	filepath.Walk("/usr/src/app/allvbuild", func(path string, info os.FileInfo, err error) error {

		relPath := strings.TrimPrefix(path, "/usr/src/app/allvbuild")

		if err != nil {
			t.Error(err.Error(), "-", relPath)
		}
		b, htmlBytes, err := isHTML(path)
		if err != nil {
			t.Error(err.Error(), "-", relPath)
		}
		// don't check non-html files
		if b == false {
			return nil
		}

		count++

		err = testURLs(htmlBytes)
		if err != nil {
			t.Error(err.Error(), "-", relPath)
		}
		return nil
	})

	fmt.Println("found", count, "html files")
	fmt.Println("found", countLinks, "links")
	fmt.Println("found", countImages, "images")
}

// testURLs tests if we're not using absolute paths for URLs
// when pointing to local pages.
func testURLs(htmlBytes []byte) error {

	reader := bytes.NewReader(htmlBytes)

	z := html.NewTokenizer(reader)

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			return nil
		case html.StartTagToken:
			t := z.Token()

			url := ""

			// check tag types
			switch t.Data {
			case "a":
				countLinks++
				ok, href := getHref(t)
				// skip, it may just be an anchor
				if !ok {
					break
				}
				url = href

			case "img":
				countImages++
				ok, src := getSrc(t)
				if !ok {
					return errors.New("img with no src: " + t.String())
				}
				url = src
			}

			// there's an url to test!
			if url != "" {
				if strings.HasPrefix(url, "http://docs.docker.com") || strings.HasPrefix(url, "https://docs.docker.com") {
					return errors.New("found absolute link: " + t.String())
				}
			}
		}
	}

	return nil
}

// helpers

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func getSrc(t html.Token) (ok bool, src string) {
	for _, a := range t.Attr {
		if a.Key == "src" {
			src = a.Val
			ok = true
		}
	}
	return
}
