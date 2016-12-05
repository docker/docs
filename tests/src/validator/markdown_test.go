package main

import (
	"errors"
	"github.com/gdevillele/frontparser"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFrontmatterTitle tests if there's a title present in all
// published markdown frontmatters.
func TestFrontmatterTitle(t *testing.T) {
	filepath.Walk("/docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		published, mdBytes, err := isPublishedMarkdown(path)
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		if published == false {
			return nil
		}
		err = testFrontmatterTitle(mdBytes)
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		return nil
	})
}

// testFrontmatterTitle tests if there's a title present in
// given markdown file bytes
func testFrontmatterTitle(mdBytes []byte) error {
	fm, _, err := frontparser.ParseFrontmatterAndContent(mdBytes)
	if err != nil {
		return err
	}
	if _, exists := fm["title"]; exists == false {
		return errors.New("can't find title in frontmatter")
	}
	return nil
}

// TestFrontMatterKeywords tests if keywords are present and correctly
// formatted in all published markdown frontmatters.
func TestFrontMatterKeywords(t *testing.T) {
	filepath.Walk("/docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		published, mdBytes, err := isPublishedMarkdown(path)
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		if published == false {
			return nil
		}
		err = testFrontMatterKeywords(mdBytes)
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		return nil
	})
}

// testFrontMatterKeywords tests if if keywords are present and correctly
// formatted in given markdown file bytes
func testFrontMatterKeywords(mdBytes []byte) error {
	fm, _, err := frontparser.ParseFrontmatterAndContent(mdBytes)
	if err != nil {
		return err
	}

	keywords, exists := fm["keywords"]

	// it's ok to have a page without keywords
	if exists == false {
		return nil
	}

	if _, ok := keywords.(string); !ok {
		return errors.New("keywords should be a comma separated string")
	}

	return nil
}

//-----------------
// utils
//-----------------

// isPublishedMarkdown returns wether a file is a published markdown or not
// as a convenience it also returns the markdown bytes to avoid reading files twice
func isPublishedMarkdown(path string) (bool, []byte, error) {
	if strings.HasSuffix(path, ".md") {
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return false, nil, err
		}
		if frontparser.HasFrontmatterHeader(fileBytes) {
			fm, _, err := frontparser.ParseFrontmatterAndContent(fileBytes)
			if err != nil {
				return false, nil, err
			}
			// skip markdowns that are not published
			if published, exists := fm["published"]; exists {
				if publishedBool, ok := published.(bool); ok {
					if publishedBool {
						// file is markdown, has frontmatter and is published
						return true, fileBytes, nil
					}
				}
			} else {
				// if "published" field is missing, it means published == true
				return true, fileBytes, nil
			}
		}
	}
	return false, nil, nil
}
