package main

import (
	"errors"
	"fmt"
	"github.com/gdevillele/frontparser"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFrontmatterTitle tests if there's a title present in all
// markdown frontmatters.
func TestFrontmatterTitle(t *testing.T) {
	filepath.Walk("/docs", func(path string, info os.FileInfo, err error) error {
		err = testFrontmatterTitle(path)
		if err != nil {
			fmt.Println(err.Error(), "-", path)
			t.Fail()
		}
		return nil
	})
}

// testFrontmatterTitle tests if there's a title present in
// markdown file at given path
func testFrontmatterTitle(path string) error {
	if strings.HasSuffix(path, ".md") {
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// if file has frontmatter
		if frontparser.HasFrontmatterHeader(fileBytes) {
			fm, _, err := frontparser.ParseFrontmatterAndContent(fileBytes)
			if err != nil {
				return err
			}

			// skip markdowns that are not published
			if published, exists := fm["published"]; exists {
				if publishedBool, ok := published.(bool); ok {
					if publishedBool == false {
						return nil
					}
				}
			}

			if _, exists := fm["title"]; exists == false {
				return errors.New("can't find title in frontmatter")
			}
		} else {
			// no frontmatter is not an error
			// markdown files without frontmatter won't be considered
			return nil
		}
	}
	return nil
}
