package main

import (
	"github.com/gdevillele/frontparser"
	"io/ioutil"
	"os"
	"strings"
)

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

// isHTML returns wether a file is an html file or not
// as a convenience it also returns the markdown bytes to avoid reading files twice
func isHTML(path string) (bool, []byte, error) {
	if strings.HasSuffix(path, ".html") {
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return false, nil, err
		}
		return true, fileBytes, nil
	}
	return false, nil, nil
}

// fileExists returns true if the given file exists
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
