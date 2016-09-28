package main

// Maintainers defines the struct for a MAINTAINERS file
type Maintainers struct {
	Rules  map[string]Rule
	Roles  map[string]Role
	Org    map[string]*Org
	People map[string]Person
}

// Rule is a project rule
type Rule struct {
	Title string `toml:"title,omitempty"`
	Text  string `toml:"text,omitempty"`
}

// Role is a project role
type Role struct {
	Person string `toml:"person,omitempty"`
	Text   string `toml:"text,omitempty"`
}

// Org defines the organization within a project
type Org struct {
	People []string
}

// Person member of the project
type Person struct {
	Name   string
	Email  string
	GitHub string
}

// MaintainersDepreciated is an old struct for compatibility
// with the docker/docker maintainers file.
// TODO: delete this once the file in docker/docker repo is updated
type MaintainersDepreciated struct {
	Rules        map[string]Rule
	Organization Organization `toml:"Org"`
	People       map[string]Person
}

// Organization defines the project's organization
// TODO: delete this once MaintainersDepreciated is removed
type Organization struct {
	BDFL             string `toml:"bdfl"`
	ChiefArchitect   string `toml:"Chief Architect"`
	ChiefMaintainer  string `toml:"Chief Maintainer"`
	CommunityManager string `toml:"Community Manager"`
	CoreMaintainers  *Org   `toml:"Core maintainers"`
	DocsMaintainers  *Org   `toml:"Docs maintainers"`
	Curators         *Org   `toml:"Curators"`
}
