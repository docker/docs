package schema

import (
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type MetadataManager struct {
	*ManifestRegistryManager
	*TagRegistryManager
	*eventManager
}

func NewMetadataManager(s *rethink.Session) *MetadataManager {
	return &MetadataManager{
		&ManifestRegistryManager{s},
		&TagRegistryManager{s},
		&eventManager{s},
	}
}
