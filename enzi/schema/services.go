package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	// ErrNoSuchService conveys that a service with the given id does not
	// exist.
	ErrNoSuchService = errors.New("no such service")
	// ErrServiceExists conveys that a service with the given ownerID and
	// name already exists.
	ErrServiceExists = errors.New("service already exists")
)

// Service represents an OpenID Connect Client, any service that uses eNZi for
// identity and authentication management.
type Service struct {
	PK                 string   `gorethink:"pk"`                 // Hash of ownerID and name. Primary key.
	ID                 string   `gorethink:"id"`                 // Randomly generated UUID.
	OwnerID            string   `gorethink:"ownerID"`            // Foreign key reference to owning account.
	Name               string   `gorethink:"name"`               // Name of the service, unique with ownerID.
	Description        string   `gorethink:"description"`        // Short description of the service.
	URL                string   `gorethink:"url"`                // Web address for the service.
	Privileged         bool     `gorethink:"privileged"`         // Indicates whether accounts are automatically opted into the service. Only system admins can set this.
	RedirectURIs       []string `gorethink:"redirectURIs"`       // List of Redirection URI values which MUST exactly match the redirect_uri parameter value in each Authorization Request.
	GrantTypes         []string `gorethink:"grantTypes"`         // List of Oauth 2.0 Grant types that the service will restrict itself to using.
	ResponseTypes      []string `gorethink:"responseTypes"`      // List of Oauth 2.0 response_type values that the service will restrict itself to using.
	JWKsURIs           []string `gorethink:"jwksURIs"`           // List of URLs for the service's JSON Web Key set, the set of public keys it uses to sign authentication tokens.
	ProviderIdentities []string `gorethink:"providerIdentities"` // List of identifiers that this service will use as the 'aud' (Audience) for its authentication tokens.
	CABundle           string   `gorethink:"caBundle"`           // PEM certificate bundle used to authenticate the JWKsURI endpoint.
}

var servicesTable = table{
	db:         dbName,
	name:       "services",
	primaryKey: "pk", // Guarantees uniqueness of (ownerID, Name). Quick lookups.
	secondaryIndexes: map[string][]string{
		"id":           nil,                 // For quick lookups by ID.
		"ownerID_name": {"ownerID", "name"}, // For quickly listing services for an owner ordered by name.
	},
}

func computeServicePK(ownerID, name string) string {
	hash := sha256.New()

	hash.Write([]byte(ownerID))
	hash.Write([]byte(name))

	return hex.EncodeToString(hash.Sum(nil))
}

// CreateService inserts a new service using the values supplied by the given
// service. The ID field of the service is set to a random UUID. The PK field
// of the service is set to a hash of the ownerID and service name. If a
// service already exists with the same ownerID and name the error will be
// ErrServiceExists.
func (m *manager) CreateService(service *Service) error {
	service.ID = uuid.NewV4().String()
	service.PK = computeServicePK(service.OwnerID, service.Name)

	if resp, err := servicesTable.Term().Insert(service).RunWrite(m.session); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrServiceExists
		}

		return fmt.Errorf("unable to create service: %s", err)
	}

	return nil
}

// getServiceByIndexVal queries the database for a service with the given value
// using the given index. If no such service exists the returned error will be
// ErrNoSuchService.
func (m *manager) getServiceByIndexVal(indexName string, val interface{}) (*Service, error) {
	cursor, err := servicesTable.Term().GetAllByIndex(indexName, val).Limit(1).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var service Service
	if err := cursor.One(&service); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchService
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &service, nil
}

// GetServiceByID retrieves the service with the given ID. If no such team
// exists the returned error will be ErrNoSuchService.
func (m *manager) GetServiceByID(id string) (*Service, error) {
	return m.getServiceByIndexVal("id", id)
}

// GetServiceByName retrieves the service with the given ownerID and name. If
// no such team exists the returned error will be ErrNoSuchService.
func (m *manager) GetServiceByName(ownerID, name string) (*Service, error) {
	return m.getServiceByIndexVal("pk", computeServicePK(ownerID, name))
}

// ListServicesForAccount returns a list services which are owned by the
// account with the given accountID. Services are ordered by name, starting
// from the given startName value. If limit is 0, all results are returned,
// otherwise returns at most limit results. Will return the startName
// of the next page or "" if no services remain.
func (m *manager) ListServicesForAccount(accountID, startName string, limit uint) (services []Service, nextName string, err error) {
	query := servicesTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "ownerID_name"},
	).Between(
		[]interface{}{accountID, startName},
		[]interface{}{accountID, rethink.MaxVal},
	)

	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(m.session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	services = []Service{}
	if err := cursor.All(&services); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(services)) > limit {
		nextName = services[limit].Name
		services = services[:limit]
	}

	return services, nextName, nil
}

// ServiceUpdateFields holds the fields of a service which are updatable. Any
// fields left nil will not be updated when used with the update method.
type ServiceUpdateFields struct {
	FullName           *string   `gorethink:"fullName,omitempty"`
	Description        *string   `gorethink:"description,omitempty"`
	URL                *string   `gorethink:"url,omitempty"`
	Privileged         *bool     `gorethink:"privileged,omitempty"`
	GrantTypes         *[]string `gorethink:"grantTypes,omitempty"`
	ResponseTypes      *[]string `gorethink:"responseTypes,omitempty"`
	RedirectURIs       *[]string `gorethink:"redirectURIs,omitempty"`
	JWKsURIs           *[]string `gorethink:"jwksURIs,omitempty"`
	ProviderIdentities *[]string `gorethink:"providerIdentities,omitempty"`
	CABundle           *string   `gorethink:"caBundle,omitempty"`
}

// UpdateService updates the service with the given ID using any set fields in
// the given updateFields struct. ServiceUpdateFields that are left null are
// unchanged.
func (m *manager) UpdateService(id string, updateFields ServiceUpdateFields) error {
	if _, err := servicesTable.Term().GetAllByIndex("id", id).Update(
		updateFields,
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to run update query: %s", err)
	}

	return nil
}

// DeleteService removes the service with the given ownerID and name.
func (m *manager) DeleteService(ownerID, name string) error {
	if _, err := servicesTable.Term().Get(computeServicePK(ownerID, name)).Delete().Run(m.session); err != nil {
		return fmt.Errorf("unable to delete service from database: %s", err)
	}

	return nil
}
