package main

import (
	"errors"
	"github.com/gdevillele/frontparser"
	"os"
	"path/filepath"
	"testing"
)

// TestFrontmatterTitle tests if there's a title present in all
// published markdown frontmatters.
func TestFrontMatterTitle(t *testing.T) {
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
		err = testFrontMatterTitle(mdBytes)
		if err != nil {
			t.Error(err.Error(), "-", path)
		}
		return nil
	})
}

// testFrontmatterTitle tests if there's a title present in
// given markdown file bytes
func testFrontMatterTitle(mdBytes []byte) error {
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
