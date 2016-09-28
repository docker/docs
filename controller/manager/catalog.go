package manager

import (
	"strings"

	"github.com/docker/orca"
)

const (
	dockerIndexAddr = "https://index.docker.io"
)

// TODO: pull to central service for "curated" catalog?
var (
	defaultCatalog = []orca.CatalogItem{
		{
			Name:        "redis",
			Description: "Redis is an open source key-value store that functions as a data structure server.",
			Trusted:     false,
			Official:    true,
		},
		{
			Name:        "mongo",
			Description: "MongoDB document databases provide high availability and easy scalability.",
			Trusted:     false,
			Official:    true,
		},
		{
			Name:        "postgres",
			Description: "The PostgreSQL object-relational database system provides reliability and data integrity.",
			Trusted:     false,
			Official:    true,
		},
		{
			Name:        "tomcat",
			Description: "Apache Tomcat is an open source implementation of the Java Servlet and JavaServer Pages technologies.",
			Trusted:     false,
			Official:    true,
		},
		{
			Name:        "jenkins",
			Description: "Official Jenkins Docker image",
			Trusted:     false,
			Official:    true,
		},
		{
			Name:        "nginx",
			Description: "Official build of Nginx.",
			Trusted:     false,
			Official:    true,
		},
	}
)

func (m DefaultManager) GetDefaultCatalog() ([]orca.CatalogItem, error) {
	return defaultCatalog, nil
}

func (m DefaultManager) SearchCatalog(query string) ([]orca.CatalogItem, error) {
	items := []orca.CatalogItem{}
	for _, i := range defaultCatalog {
		if strings.Contains(i.Name, query) {
			items = append(items, i)
		}
	}

	return items, nil
}
