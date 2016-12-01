package frontparser

import (
	"bytes"
	"errors"

	"gopkg.in/yaml.v2"
)

// error messages
var ErrorFrontmatterNotFound error = errors.New("frontmatter header not found")

var FmYAMLDelimiter []byte = []byte("---")

// returns whether the provided content
func HasFrontmatterHeader(input []byte) bool {
	// remove heading and trailing spaces (and CR, LF, ...)
	input = bytes.TrimSpace(input)
	// test for frontmatter delimiter
	if !bytes.HasPrefix(input, FmYAMLDelimiter) {
		return false
	}
	// trim heading frontmatter delimiter
	input = bytes.TrimPrefix(input, FmYAMLDelimiter)
	// split on frontmatter delimiter to separate frontmatter from the rest
	elements := bytes.SplitN(input, FmYAMLDelimiter, 2)
	if len(elements) != 2 {
		// malformed input
		return false
	}
	// parse frontmatter to validate it is valid YAML
	var out map[string]interface{} = make(map[string]interface{})
	err := yaml.Unmarshal(elements[0], out)
	return err == nil
}

//
func ParseFrontmatter(input []byte, out interface{}) error {
	// remove heading and trailing spaces (and CR, LF, ...)
	input = bytes.TrimSpace(input)
	// test for frontmatter delimiter
	if !bytes.HasPrefix(input, FmYAMLDelimiter) {
		return errors.New("heading frontmatter delimiter not found")
	}
	// trim heading frontmatter delimiter
	input = bytes.TrimPrefix(input, FmYAMLDelimiter)
	// split on frontmatter delimiter to separate frontmatter from the rest
	elements := bytes.SplitN(input, FmYAMLDelimiter, 2)
	if len(elements) != 2 {
		// malformed input
		return errors.New("more than two frontmatter delimiters were found")
	}
	// parse frontmatter to validate it is valid YAML
	return yaml.Unmarshal(elements[0], out)
}

// input is frontmatter + markdown
// return values
// - frontmatter
// - markdown
// - error
func ParseFrontmatterAndContent(input []byte) (map[string]interface{}, []byte, error) {
	var resultFm map[string]interface{} = make(map[string]interface{})
	var resultRest []byte = make([]byte, 0)

	// remove heading and trailing spaces (and CR, LF, ...)
	input = bytes.TrimSpace(input)
	// test for frontmatter delimiter
	if !bytes.HasPrefix(input, FmYAMLDelimiter) {
		return resultFm, resultRest, errors.New("heading frontmatter delimiter not found")
	}

	// trim heading frontmatter delimiter
	input = bytes.TrimPrefix(input, FmYAMLDelimiter)

	// split on frontmatter delimiter to separate frontmatter from the rest
	elements := bytes.SplitN(input, FmYAMLDelimiter, 2)
	if len(elements) != 2 {
		// malformed input
		return resultFm, resultRest, errors.New("more than two frontmatter delimiters were found")
	}

	err := yaml.Unmarshal(elements[0], resultFm)
	if err != nil {
		return resultFm, resultRest, err
	}

	resultRest = elements[1]

	return resultFm, resultRest, nil
}
