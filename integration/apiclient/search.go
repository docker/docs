package apiclient

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

const (
	autocompletePath = "/api/v0/index/autocomplete"
	reindexPath      = "/api/v0/index/reindex"
)

type SearchOptions struct {
	forms.SearchOptions
	Start int
}

func (c *apiClient) Autocomplete(opts SearchOptions) (*responses.Autocomplete, error) {
	v := opts.SearchOptions.GetURLValues()
	v.Set("start", strconv.Itoa(opts.Start))
	resp, err := c.makeRequest("GET", url.URL{Path: autocompletePath, RawQuery: v.Encode()}, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}
	var autocompleteResp responses.Autocomplete
	if err := json.NewDecoder(resp.Body).Decode(&autocompleteResp); err != nil {
		return nil, err
	} else {
		return &autocompleteResp, nil
	}
}

func (c *apiClient) Reindex() error {
	if resp, err := c.makeRequest("POST", url.URL{Path: reindexPath}, nil); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		return validateStatusCode(resp, http.StatusAccepted)
	}
}
