package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var countLinks = 0
var countImages = 0
var htmlContentRootPath = "/usr/src/app/allvbuild"

// TestURLs tests if we're not using absolute paths for URLs
// when pointing to local pages.
func TestURLs(t *testing.T) {
	count := 0

	filepath.Walk(htmlContentRootPath, func(path string, info os.FileInfo, err error) error {

		relPath := strings.TrimPrefix(path, htmlContentRootPath)

		isArchive, err := regexp.MatchString(`^/v[0-9]+\.[0-9]+/.*`, relPath)
		if err != nil {
			t.Error(err.Error(), "-", relPath)
		}
		// skip archives for now, only test URLs in current version
		// TODO: test archives
		if isArchive {
			return nil
		}

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

		err = testURLs(htmlBytes, path)
		if err != nil {
			t.Error(relPath + err.Error())
		}
		return nil
	})

	fmt.Println("found", count, "html files (excluding archives)")
	fmt.Println("found", countLinks, "links (excluding archives)")
	fmt.Println("found", countImages, "images (excluding archives)")
}

// testURLs tests if we're not using absolute paths for URLs
// when pointing to local pages.
func testURLs(htmlBytes []byte, htmlPath string) error {

	reader := bytes.NewReader(htmlBytes)

	z := html.NewTokenizer(reader)

	urlErrors := ""
	// fmt.Println("urlErrors:", urlErrors)
	done := false

	for !done {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			done = true
		case html.StartTagToken:
			t := z.Token()

			urlStr := ""

			// check tag types
			switch t.Data {
			case "a":
				countLinks++
				ok, href := getHref(t)
				// skip, it may just be an anchor
				if !ok {
					break
				}
				urlStr = href

			case "img":
				countImages++
				ok, src := getSrc(t)
				if !ok {
					urlErrors += "\nimg with no src: " + t.String()
					break
				}
				urlStr = src
			}

			// there's an url to test!
			if urlStr != "" {
				u, err := url.Parse(urlStr)
				if err != nil {
					urlErrors += "\ncan't parse url: " + t.String()
					break
					// return errors.New("can't parse url: " + t.String())
				}
				// test with github.com
				if u.Scheme != "" && u.Host == "docs.docker.com" {
					urlErrors += "\nabsolute: " + t.String()
					break
				}

				// relative link
				if u.Scheme == "" {

					resourcePath := ""
					resourcePathIsAbs := false

					if filepath.IsAbs(u.Path) {
						resourcePath = filepath.Join(htmlContentRootPath, mdToHtmlPath(u.Path))
						resourcePathIsAbs = true
					} else {
						resourcePath = filepath.Join(filepath.Dir(htmlPath), mdToHtmlPath(u.Path))
					}

					if _, err := os.Stat(resourcePath); os.IsNotExist(err) {

						fail := true

						// index.html could mean there's a corresponding index.md meaning built the correct path
						// but Jekyll actually creates index.html files for all md files.
						// foo.md -> foo/index.html
						// it does this to prettify urls, content of foo.md would then be rendered here:
						// http://domain.com/foo/ (instead of http://domain.com/foo.html)
						// so if there's an error, let's see if index.md exists, otherwise retry from parent folder
						// (only if the resource path is not absolute)
						if !resourcePathIsAbs && filepath.Base(htmlPath) == "index.html" {
							// retry from parent folder
							resourcePath = filepath.Join(filepath.Dir(htmlPath), "..", mdToHtmlPath(u.Path))
							if _, err := os.Stat(resourcePath); err == nil {
								fail = false
							}
						}

						if fail {
							urlErrors += "\nbroken: " + t.String()
							break
						}
					}
				}
			}
		}
	}

	// fmt.Println("urlErrors:", urlErrors)
	if urlErrors != "" {
		return errors.New(urlErrors)
	}
	return nil
}

func mdToHtmlPath(mdPath string) string {
	if strings.HasSuffix(mdPath, ".md") == false {
		// file is not a markdown, don't change anything
		return mdPath
	}
	if strings.HasSuffix(mdPath, "index.md") {
		return strings.TrimSuffix(mdPath, "md") + "html"
	}
	return strings.TrimSuffix(mdPath, ".md") + "/index.html"
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
