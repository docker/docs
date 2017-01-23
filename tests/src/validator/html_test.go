package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var countLinks = 0
var countImages = 0

// testTocDotYaml tests if there are no broken links in
// _data/toc.yaml
func testTocDotYaml() error {
	tocBytes, err := ioutil.ReadFile(filepath.Join(docsSource, "_data/toc.yaml"))
	if err != nil {
		return err
	}

	toc := make(map[string]interface{})
	err = yaml.Unmarshal(tocBytes, &toc)
	if err != nil {
		return err
	}

	// top level map
	topLevelArray, ok := toc["toc"].([]interface{})
	if !ok {
		return errors.New("top level object in toc.yaml can't be parsed into []interface{}")
	}

	err = tocLinkChecker(topLevelArray)
	if err != nil {
		return err
	}

	return nil
}

// TestURLs checks different things regarding urls (in that order):
// - no links in _data/toc.yaml are broken
// - we don't have absolute links to docs.docker.com pointing to local ressources
// - no relative links to local ressources are broken (warning for working redirects)
func TestURLs(t *testing.T) {
	count := 0

	err := testTocDotYaml()
	if err != nil {
		t.Error("broken links in _data/toc.yaml:", err.Error())
		return
	}

	csv := ""

	filepath.Walk(docsHtmlWithoutRedirects, func(path string, info os.FileInfo, err error) error {

		relPath := strings.TrimPrefix(path, docsHtmlWithoutRedirects)

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

		err, redirects := testURLs(htmlBytes, relPath)

		// build csv
		if err != nil || redirects != "" {
			csv += "\n" + relPath
		}

		if err != nil {
			t.Error(relPath + " (broken links)\n" + err.Error())
			csv += ",-- broken --"
			links := strings.Split(err.Error(), ",")
			for _, link := range links {
				csv += "\n," + link
			}
		}
		if redirects != "" {
			if err != nil {
				csv += "\n,-- redirects --"
			} else {
				csv += ",-- redirects --"
			}
			links := strings.Split(redirects, ",")
			for _, link := range links {
				csv += "\n," + link
			}
			// t.Log(relPath + " (redirects)\n" + redirects)
		}

		return nil
	})

	fmt.Println("found", count, "html files (excluding archives)")
	fmt.Println("found", countLinks, "links (excluding archives)")
	fmt.Println("found", countImages, "images (excluding archives)")
	if csv != "" {
		fmt.Println("CSV:\n", csv)
	}
}

// testURLs tests if we're not using absolute paths for URLs
// when pointing to local pages.
func testURLs(htmlBytes []byte, htmlPath string) (err error, redirects string) {

	reader := bytes.NewReader(htmlBytes)

	z := html.NewTokenizer(reader)

	urlErrors := make([]string, 0)
	urlRedirects := make([]string, 0)

	// if there's a <base> tag, we should use it to check
	// if resources pointed by relative links do exist
	baseTagValue := ""

	done := false

	for !done {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			done = true

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			urlStr := ""

			// check tag types
			switch t.Data {
			case "base":
				ok, href := getHref(t)
				// ignore if href can't be found
				if !ok {
					break
				}
				baseTagValue = href
				break
			case "a":
				countLinks++
				ok, href := getHref(t)
				// skip, it may just be an anchor
				if !ok {
					break
				}
				urlStr = href
				break
			case "img":
				countImages++
				ok, src := getSrc(t)
				if !ok {
					urlErrors = append(urlErrors, "image-with-no-src")
					break
				}
				urlStr = src
			}

			// there's an url to test!
			if urlStr != "" {
				u, err := url.Parse(urlStr)
				if err != nil {
					urlErrors = append(urlErrors, urlStr)
					break
					// return errors.New("can't parse url: " + t.String())
				}
				// test with github.com
				if u.Scheme != "" && u.Host == "docs.docker.com" {
					urlErrors = append(urlErrors, urlStr)
					break
				}

				// relative link
				if u.Scheme == "" {
					origin := htmlPath
					if baseTagValue != "" {
						origin = baseTagValue
					}

					err := checkLinkToLocalResourceValid(cleanPath(u.Path), docsHtmlWithoutRedirects, origin)
					if err != nil {
						// maybe it's a redirect
						err := checkLinkToLocalResourceValid(cleanPath(u.Path), docsHtmlWithRedirects, origin)
						if err != nil {
							urlErrors = append(urlErrors, urlStr)
						} else {
							urlRedirects = append(urlRedirects, urlStr)
						}
					}
					break // exit switch
				}
			}
		}
	}

	redirects = ""
	if len(urlRedirects) > 0 {
		redirects = strings.Join(urlRedirects, ",")
	}

	if len(urlErrors) > 0 {
		return errors.New(strings.Join(urlErrors, ",")), redirects
	}
	return nil, redirects
}

// cleanPath takes a path in parameter and returns path
// to corresponding index.html.
// It deals with different kind of situations:
// - markdown paths:
// 		- /foo/bar.md -> /foo/bar/index.html
//		- /foo/index.md -> /foo/index.html
// - folders:
// 		- /foo/bar/baz/ -> /foo/bar/baz/index.html
// cleanPath also removes queries and fragments
// - /foo/bar?query#fragment -> /foo/bar/index.html

func cleanPath(path string) string {

	if strings.HasPrefix(path, "#") {
		return ""
	}

	// remove queries & fragments
	parts := strings.Split(path, "#")
	if len(parts) > 1 {
		path = strings.Join(parts[:len(parts)-1], "#")
	}
	parts = strings.Split(path, "?")
	if len(parts) > 1 {
		path = strings.Join(parts[:len(parts)-1], "?")
	}

	if strings.HasSuffix(path, ".md") {
		if strings.HasSuffix(path, "index.md") {
			return strings.TrimSuffix(path, "md") + "html"
		}
		return strings.TrimSuffix(path, ".md") + "/index.html"
	}
	if strings.HasSuffix(path, "/") {
		return path + "index.html"
	}
	return path
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

// isLinkToLocalRessourceBroken returns wether given path does
// point to existing local resource or not
// path: path to be checked
// rootFolder: the root for absolute paths (like /foo/bar)
// origin: where the link is made
//
// example:
// isLinkToLocalResourceValid("../baz", "/www", "foo/bar")
// will look for "/www/foo/baz"
//
// if path is absolute, origin is not used:
// isLinkToLocalResourceValid("/baz", "/www", "foo/bar")
// will look for "/www/baz"
func checkLinkToLocalResourceValid(path string, rootFolder string, origin string) error {

	var absPath string
	if filepath.IsAbs(path) {
		absPath = filepath.Join(rootFolder, path)
	} else {
		// if origin is not a folder, we should select parent directory
		if filepath.Ext(origin) != "" {
			origin = filepath.Dir(origin)
		}

		absPath = filepath.Join(rootFolder, origin, path)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		// fmt.Println(err.Error())
		return err
	}

	return nil

	// TODO: see if that is necessary
	// // index.html could mean there's a corresponding index.md meaning built the correct path
	// // but Jekyll actually creates index.html files for all md files.
	// // foo.md -> foo/index.html
	// // it does this to prettify urls, content of foo.md would then be rendered here:
	// // http://domain.com/foo/ (instead of http://domain.com/foo.html)
	// // so if there's an error, let's see if index.md exists, otherwise retry from parent folder
	// // (only if the resource path is not absolute)
	// if !resourcePathIsAbs && filepath.Base(htmlPath) == "index.html" {
	// 	// retry from parent folder
	// 	resourcePath = filepath.Join(filepath.Dir(htmlPath), "..", cleanPath(u.Path))
	// 	if _, err := os.Stat(resourcePath); err == nil {
	// 		fail = false
	// 	}
	// }

}

// TOC

type TocLink struct {
	Path  string
	Title string
}

func TopLinkFromMSI(msi map[string]interface{}) (TocLink, error) {
	t := TocLink{}
	ok := false

	if msi["path"] != nil && msi["title"] != nil {
		if t.Path, ok = msi["path"].(string); !ok {
			return TocLink{}, errors.New("path is not a string")
		}
		if t.Title, ok = msi["title"].(string); !ok {
			return TocLink{}, errors.New("title is not a string")
		}
		return t, nil
	}
	return t, errors.New("given msi does not represent a TocLink")
}

type TocSection struct {
	SectionTitle string
	Section      []interface{}
}

func TopSectionFromMSI(msi map[string]interface{}) (TocSection, error) {
	t := TocSection{}
	ok := false

	if msi["section"] != nil && msi["title"] != nil {
		if t.SectionTitle, ok = msi["sectiontitle"].(string); !ok {
			return TocSection{}, errors.New("sectiontitle is not a string")
		}
		if t.Section, ok = msi["section"].([]interface{}); !ok {
			return TocSection{}, errors.New("section is not a []interface{}")
		}
		return t, nil
	}
	return t, errors.New("given msi does not represent a TocSection")
}

func tocLinkChecker(items []interface{}) error {

	urlErrors := make([]string, 0)

	// TODO: gather errors, don't exit on first one
	for _, item := range items {

		mii, ok := item.(map[interface{}]interface{})
		if !ok {
			fmt.Printf("%#v\n", item)
			return errors.New("item is not a map[interface{}]interface{}")
		}

		for key, value := range mii {
			keyStr, ok := key.(string)
			if !ok {
				return errors.New("map key is not a string")
			}

			if keyStr == "section" {
				arr, ok := value.([]interface{})
				if !ok {
					return errors.New("section value is not a []interface{}")
				}

				err := tocLinkChecker(arr)
				if err != nil {
					return err
				}
				continue
			}

			if keyStr != "path" {
				continue
			}

			path, ok := value.(string)
			if !ok {
				return errors.New("path is not a string")
			}

			// don't accept links to external resources
			u, err := url.Parse(path)
			if err != nil {
				return err
			}
			if u.Scheme != "" {
				urlErrors = append(urlErrors, path)
			}

			err = checkLinkToLocalResourceValid(cleanPath(path), docsHtmlWithoutRedirects, "/all/links/should/be/absolute")

			if err != nil {
				urlErrors = append(urlErrors, path)
			}
		}
	}

	if len(urlErrors) > 0 {
		return errors.New(strings.Join(urlErrors, ","))
	}
	return nil
}
