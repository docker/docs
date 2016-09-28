package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

const (
	RepositoryVisibilityPublic  = "public"
	RepositoryVisibilityPrivate = "private"
)

var (
	ErrRepositoryExists = errors.New("repository already exists")
	ErrNoSuchRepository = errors.New("no such repository")
)

// Repository represents a repository in the system.
type Repository struct {
	// ID is a randomly generated uuid for foreign references.
	ID string `gorethink:"id"`
	// PK is the hash of the namespace's account ID and the repository name.
	PK string `gorethink:"pk"`
	// NamespaceAccountID is the ID in eNZi of the account that the repository
	// belongs to.
	NamespaceAccountID string `gorethink:"namespaceAccountID"`
	// Name is the repository name.
	Name string `gorethink:"name"`
	// ShortDescription is the short description of the repository.
	ShortDescription string `gorethink:"shortDescription,omitempty"`
	// LongDescription is the long description of the repository.
	LongDescription string `gorethink:"longDescription,omitempty"`
	// Visibility is the visibility of the repository, e.g. "public" or
	// "private".
	Visibility string `gorethink:"visibility"`
}

func (r Repository) PublicID() string {
	return r.NamespaceAccountID + "/" + r.Name
}

// deconstructPublicID deconstructs the publicID into (NamespaceAccountID, Name)
func deconstructPublicID(publicID string) (string, string) {
	parts := strings.Split(publicID, "/")
	if len(parts) != 2 {
		return "", ""
	}
	namespaceAccountID, name := parts[0], parts[1]
	return namespaceAccountID, name
}

var repositoriesTable = table{
	db:         deploy.DTRDBName,
	name:       "repositories",
	primaryKey: "pk", // Guarantees uniqueness of (namespaceAccountID, name), and allows quick lookups by pk
	secondaryIndexes: map[string][]string{
		"id":   nil, // For quick lookups by ID.
		"name": nil, // For autocomplete by name.
		"namespaceAccountID_name": {"namespaceAccountID", "name"}, // For quickly listing repositories under a namespace, ordered by name
	},
}

// RepositoryManager exports CRUDy methods for Repositories in the database.
type RepositoryManager struct {
	session *rethink.Session
}

func NewRepositoryManager(session *rethink.Session) *RepositoryManager {
	return &RepositoryManager{session}
}

func (m *RepositoryManager) CreateRepository(repo *Repository) error {
	repo.ID = uuid.NewV4().String()
	repo.PK = computeRepositoryPK(repo.NamespaceAccountID, repo.Name)
	if resp, err := repositoriesTable.Term().Insert(repo).RunWrite(m.session); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrRepositoryExists
		}
		return fmt.Errorf("unable to create repository in database: %s", err)
	}
	return nil
}

func (m *RepositoryManager) GetRepositoryByName(namespaceID, name string) (*Repository, error) {
	pk := computeRepositoryPK(namespaceID, name)
	var repo Repository
	if err := repositoriesTable.getRowByIndexVal(m.session, "pk", pk, &repo, ErrNoSuchRepository); err != nil {
		return nil, err
	}
	return &repo, nil
}

type RepositoryUpdateFields struct {
	ShortDescription *string `gorethink:"shortDescription,omitempty"`
	LongDescription  *string `gorethink:"longDescription,omitempty"`
	Visibility       *string `gorethink:"visibility,omitempty"`
}

func (m *RepositoryManager) UpdateRepository(id string, updateFields RepositoryUpdateFields) error {
	if _, err := repositoriesTable.Term().GetAllByIndex("id", id).Update(
		updateFields,
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to run update query: %s", err)
	}
	return nil
}

func (m *RepositoryManager) DeleteRepository(id string) error {
	if _, err := repositoriesTable.Term().GetAllByIndex("id", id).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete repository from database: %s", err)
	}
	return nil
}

// DeleteRepositoryByPK uses the PK, and not the secondary UUID, to delete a repo...
func (m *RepositoryManager) DeleteRepositoryByPK(pk string) error {
	if _, err := repositoriesTable.Term().GetAllByIndex("pk", pk).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete repository from database: %s", err)
	}
	return nil
}

func (m *RepositoryManager) ListAllRepositories(startPublicID string, limit uint) (repositories []*Repository, nextPublicID string, err error) {
	namespaceAccountID, name := deconstructPublicID(startPublicID)
	query := repositoriesTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "namespaceAccountID_name"},
	).Between(
		[]interface{}{namespaceAccountID, name},
		rethink.MaxVal,
	)
	return m.paginateRepositoryQuery(query, limit)
}

func (m *RepositoryManager) CountPublicRepositories() (count int, err error) {
	publicFilter := rethink.Row.Field("visibility").Eq(RepositoryVisibilityPublic)
	cursor, err := repositoriesTable.Term().Filter(publicFilter).Count().Run(m.session)
	if err != nil {
		return 0, fmt.Errorf("unable to query the db: %s", err)
	}

	var result []int
	if err := cursor.All(&result); err != nil {
		return 0, fmt.Errorf("unable to scan query results: %s", err)
	}

	return result[0], nil
}

func (m *RepositoryManager) CountPrivateRepositories() (count int, err error) {
	privateFilter := rethink.Row.Field("visibility").Eq(RepositoryVisibilityPrivate)
	cursor, err := repositoriesTable.Term().Filter(privateFilter).Count().Run(m.session)
	if err != nil {
		return 0, fmt.Errorf("unable to query the db: %s", err)
	}

	var result []int
	if err := cursor.All(&result); err != nil {
		return 0, fmt.Errorf("unable to scan query results: %s", err)
	}

	return result[0], nil
}

func (m *RepositoryManager) AutocompleteAllRepositories(repoPrefix string, limit uint, namespaceAccountID string) (repositories []*Repository, err error) {
	var nameAutocompleteEnd interface{}
	if repoPrefix != "" {
		nameAutocompleteEndBytes := []byte(repoPrefix)
		nameAutocompleteEndBytes[len(repoPrefix)-1]++
		nameAutocompleteEnd = string(nameAutocompleteEndBytes)
	} else {
		nameAutocompleteEnd = rethink.MaxVal
	}
	query := repositoriesTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "name"},
	).Between(repoPrefix, nameAutocompleteEnd)
	if namespaceAccountID != "" {
		query = query.Filter(rethink.Row.Field("namespaceAccountID").Eq(namespaceAccountID))
	}
	repos, _, err := m.paginateRepositoryQuery(query, limit)
	return repos, err
}

func (m *RepositoryManager) ListRepositoriesInNamespace(namespaceAccountID, startPublicID string, limit uint) (repositories []*Repository, nextPublicID string, err error) {
	_, name := deconstructPublicID(startPublicID)
	query := repositoriesTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "namespaceAccountID_name"},
	).Between(
		[]interface{}{namespaceAccountID, name},
		[]interface{}{namespaceAccountID, rethink.MaxVal},
	)
	return m.paginateRepositoryQuery(query, limit)
}

func (m *RepositoryManager) ListRepositoriesInNamespaces(namespaceAccountIDs []string, startPublicID string, limit uint) (repositories []*Repository, nextPublicID string, err error) {
	if len(namespaceAccountIDs) == 0 {
		return nil, "", nil
	}

	namespaceAccountID, name := deconstructPublicID(startPublicID)
	// if we don't find the start account ID then just start from the
	// beginning
	namespaceIndex := 0
	for i, accID := range namespaceAccountIDs {
		if accID == namespaceAccountID {
			namespaceIndex = i
		}
	}
	var firstQuery rethink.Term
	queries := make([]interface{}, len(namespaceAccountIDs)-namespaceIndex-1)
	for i, namespaceAccountID := range namespaceAccountIDs[namespaceIndex:] {
		startRepoName := ""
		if namespaceIndex == i {
			startRepoName = name
		}
		q := repositoriesTable.Term().OrderBy(
			rethink.OrderByOpts{Index: "namespaceAccountID_name"},
		).
			Between(
				[]interface{}{namespaceAccountID, startRepoName},
				[]interface{}{namespaceAccountID, rethink.MaxVal},
			)
		if i == 0 {
			firstQuery = q
		} else {
			queries[i-1] = q
		}
	}
	if len(queries) > 1 {
		firstQuery = firstQuery.Union(queries[1:]...)
	}
	return m.paginateRepositoryQuery(firstQuery, limit)
}

func (m *RepositoryManager) ListRepositoriesInNamespacesOrPublic(namespaceAccountIDs []string, startPublicID string, limit uint) (repositories []*Repository, nextPublicID string, err error) {
	namespaceAccountID, name := deconstructPublicID(startPublicID)
	nsOrPublicFilter := rethink.Row.Field("visibility").Eq(RepositoryVisibilityPublic)
	for _, namespaceAccountID := range namespaceAccountIDs {
		nsOrPublicFilter = nsOrPublicFilter.Or(rethink.Row.Field("namespaceAccountID").Eq(namespaceAccountID))
	}
	query := repositoriesTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "namespaceAccountID_name"},
	).
		Between(
			[]interface{}{namespaceAccountID, name},
			rethink.MaxVal,
		).
		Filter(nsOrPublicFilter)
	return m.paginateRepositoryQuery(query, limit)
}

func (m *RepositoryManager) paginateRepositoryQuery(query rethink.Term, limit uint) (repositories []*Repository, nextKey string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(m.session)
	if err != nil {
		return nil, "", fmt.Errorf("RepositoryManager: unable to query db: %s", err)
	}

	if err := cursor.All(&repositories); err != nil {
		return nil, "", fmt.Errorf("RepositoryManager: unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(repositories)) > limit {
		nextKey = repositories[limit].PublicID()
		repositories = repositories[:limit]
	}

	return repositories, nextKey, nil
}

func computeRepositoryPK(namespaceID, name string) string {
	hash := sha256.New()
	hash.Write([]byte(namespaceID))
	hash.Write([]byte(name))
	return hex.EncodeToString(hash.Sum(nil))
}
