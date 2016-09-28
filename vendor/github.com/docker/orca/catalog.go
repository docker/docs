package orca

type (
	CatalogItem struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Trusted     bool   `json:"is_trusted,omitempty"`
		Official    bool   `json:"is_official,omitempty"`
		ImageURL    string `json:"image_url,omitempty"`
	}
)
