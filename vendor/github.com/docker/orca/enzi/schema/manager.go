package schema

import (
	"io"
	"time"

	rethink "gopkg.in/dancannon/gorethink.v2"
)

// Manager exports CRUDy methods for Accounts and Teams.
type Manager interface {
	// GetProperty looks up the property with the given key and unmarshals
	// its value into the given val interface which must be a pointer to a
	// type which can be unmarshalled from JSON.
	GetProperty(key string, val interface{}) error
	// SetProperty sets the property with the given key to the given value.
	// The given value must be a type which can be marshalled to JSON.
	SetProperty(key string, val interface{}) error
	// GetPropertyChanges begins listening for any changes to properties.
	// Returns a channel on which the caller may receive a stream of
	// PropertyChange objects and an io.Closer which performs necessary
	// cleanup to end the stream's underlying goroutine. Only changes for
	// the given keys are streamed. If no keys are specified, all property
	// changes are streamed. After closing, the changeStream should be
	// checked for a possible remaining value.
	GetPropertyChanges(keys ...string) (changeStream <-chan PropertyChange, streamCloser io.Closer, err error)
	// DeleteProperty deletes the property with the given key if it exists.
	DeleteProperty(key string) error

	// CreateAccount creates a new account using the values supplied by the
	// given Account. The ID field of the account is set to a random UUID.
	CreateAccount(acct *Account) error
	// GetAccountByName gets an account with the given name. If no such
	// account exists the returned error will be ErrNoSuchAccount.
	GetAccountByName(name string) (*Account, error)
	// GetAccountByID gets an account with the given id. If no such account
	// exists the returned error will be ErrNoSuchAccount.
	GetAccountByID(id string) (*Account, error)
	// GetUserByName gets a user account with the given name. If no such
	// user account exists the returned error will be ErrNoSuchAccount.
	GetUserByName(name string) (*Account, error)
	// GetUserByID gets a user account with the given id. If no such user
	// account exists the returned error will be ErrNoSuchAccount.
	GetUserByID(id string) (*Account, error)
	// GetUserByLdapDN gets a user account with the given ldapDN. If no
	// such user account exists the returned error will be
	// ErrNoSuchAccount.
	GetUserByLdapDN(ldapDN string) (*Account, error)
	// GetOrgByName gets an organization account with the given name. If no
	// such organization account exists the returned error will be
	// ErrNoSuchAccount.
	GetOrgByName(name string) (*Account, error)
	// GetOrgByID gets an organization account with the given id. If no
	// such organization account exists the returned error will be
	// ErrNoSuchAccount.
	GetOrgByID(id string) (*Account, error)
	// ListAccounts gets a slice of all accounts ordered by name, starting
	// from the given startName value. If limit is 0, all results
	// are returned, otherwise returns at most limit results. Will return
	// the startName of the next page or "" if no accounts remain.
	ListAccounts(startName string, limit uint) (accts []Account, nextName string, err error)
	// ListUsers gets a slice of all user accounts ordered by name,
	// starting from the given startName value. If limit is 0, all results
	// are returned, otherwise returns at most limit results. Will return
	// the startName of the next page or "" if no user accounts remain.
	ListUsers(startName string, limit uint) (accts []Account, nextName string, err error)
	// ListOrgs queries for a slice of all organization accounts ordered by
	// name, starting from the given startName value. If limit is 0, all
	// results are returned, otherwise returns at most limit results. Will
	// return the startName of the next page or "" if no organization
	// accounts remain.
	ListOrgs(startName string, limit uint) (accts []Account, nextName string, err error)
	// ListAdmins queries for a slice of all accounts marked as admins
	// ordered by name, starting from the given startName value. If limit
	// is 0, all results are returned, otherwise returns at most limit
	// results. Will return the startName of the next page or "" if no
	// admin accounts remain.
	ListAdmins(startName string, limit uint) (accts []Account, nextName string, err error)
	// ListNonAdmins queries for a slice of all accounts not marked as
	// admins ordered by name, starting from the given startName value. If
	// limit is 0, all results are returned, otherwise returns at most
	// limit results. Will return the startName of the next page or "" if
	// no non-admin accounts remain.
	ListNonAdmins(startName string, limit uint) (accts []Account, nextName string, err error)
	// UpdateAccount updates fields for the account with the given id.
	UpdateAccount(id string, updateFields AccountUpdateFields) error
	// DeleteAccount removes the account with the given id.
	DeleteAccount(id string) error

	// AddOrgMembership updates or creates an org membership record for the
	// given userID and orgID. If either isAdmin or isPublic are nil,
	// those attributes will not be updated or set.
	AddOrgMembership(orgID, userID string, isAdmin, isPublic *bool) error
	// GetOrgMembership gets the org membership info for the given orgID
	// and userID. Returns (nil, nil) if the user is not a member of the
	// org.
	GetOrgMembership(orgID, userID string) (*OrgMembership, error)
	// ListOrgsForUser gets a slice of all orgs that the user with the
	// given userID is a member of, ordered by id, starting from the given
	// startID value. If limit is 0, all results are returned, otherwise
	// returns at most limit results. Will return the startID of the next
	// page or "" if no orgs remain.
	ListOrgsForUser(userID, startID string, limit uint) (memberOrgs []MemberOrg, nextID string, err error)
	// ListOrgMembers gets a slice of all user accounts which are members
	// of the org with the given orgID. Members are ordered by id, starting
	// from the given startID value. If limit is 0, all results are
	// returned, otherwise returns at most limit results. Will return the
	// startID of the next page or "" if no user accounts remain.
	ListOrgMembers(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error)
	// ListPublicOrgMembers gets a slice of all user accounts which are
	// public members of the org with the given orgID. Members are ordered
	// by id, starting from the given startID value. If limit is 0, all
	// results are returned, otherwise returns at most limit results. Will
	// return the startID of the next page or "" if no public members
	// remain.
	ListPublicOrgMembers(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error)
	// ListOrgAdmins queries for a slice of all accounts marked as admins
	// of the org with the given orgID. Admins are ordered by ID, starting
	// from the given startID value. If limit is 0, all results are
	// returned, otherwise returns at most limit results. Will return the
	// startID of the next page or "" if no org admins remain.
	ListOrgAdmins(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error)
	// ListOrgNonAdmins queries for a slice of all accounts not marked as
	// admins of the org with the given orgID. Non-admins are ordered by
	// id, starting from the given startID value. If limit is 0, all
	// results are returned, otherwise returns at most limit results. Will
	// return the startID of the next page or "" if no non-admins remain.
	ListOrgNonAdmins(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error)
	// DeleteOrgMembership removes the org membership for the given orgID
	// and userID.
	DeleteOrgMembership(orgID, userID string) error

	// CreateTeam creates a new team using the values supplied by the given
	// team. The ID field of the team is set to a random UUID. The PK field
	// of the team is set to a hash of the organization ID and team name.
	CreateTeam(team *Team) error
	// GetTeamByName gets a team with the given name in an org with the
	// given orgID. If no such team exists the returned error will be
	// ErrNoSuchTeam.
	GetTeamByName(orgID, name string) (*Team, error)
	// GetTeamByID gets a team with the given id. If no such team exists
	// the returned error will be ErrNoSuchTeam.
	GetTeamByID(id string) (*Team, error)
	// ListTeamsInOrg gets a slice of all teams with the given orgID,
	// ordered by name, starting from the given startName value. If limit
	// is 0, all results are returned, otherwise returns at most limit
	// results. Will return the startName of the next page or "" if no
	// teams remain.
	ListTeamsInOrg(orgID, startName string, limit uint) (teams []Team, nextName string, err error)
	// LisTeamMembersInOrg gets a slice of distinct userIDs of members of
	// all teams with the given orgID, in no particular order. This method
	// is used as a convenience function for syncing org membership with
	// the set of all teams' memberships.
	ListTeamMembersInOrg(orgID string) (userIDs []string, err error)
	// ListLDAPSyncTeamsInOrg returns a slice of teams with the given orgID
	// which are configured to be synced with an LDAP group. Teams are
	// ordered by name, starting from the given startName value. If limit
	// is 0, all results are returned, otherwise returns at most limit
	// results. Will return the startName of the next page or "" if no
	// teams remain.
	ListLDAPSyncTeamsInOrg(orgID, startName string, limit uint) (teams []Team, nextName string, err error)
	// UpdateTeam updates the fields for the team with the given id.
	UpdateTeam(id string, updateFields TeamUpdateFields) error
	// RenameTeam changes the name of the given team to the given new name.
	// This operation is more advanced than simply updating the team record
	// as changing the PK of the record is not possible. A duplicate team
	// with the same ID must be made and the original must then be deleted.
	// This operation is not atomic, but if the copy succeeds and the
	// delete fails, the delete of the old team by name can be retried. On
	// success, the given team object will have its Name and PK field
	// changed.
	RenameTeam(team *Team, newName string) error
	// DeleteTeam removes the team with the given id.
	DeleteTeam(ordID, name string) error

	// AddTeamMembership updates or creates a team membership record for
	// the given userID and teamID. If either isAdmin or isPublic are nil,
	// those attributes will not be updated or set. If the user is not
	// already a member of the org, they will be added.
	AddTeamMembership(orgID, teamID, userID string, isAdmin, isPublic *bool) error
	// GetTeamMembership gets the team membership info for the given teamID
	// and userID. Returns (nil, nil) if the user is not a member of the
	// team.
	GetTeamMembership(teamID, userID string) (*TeamMembership, error)
	// ListTeamMembers gets a slice of all user accounts which are members
	// of the team with the given teamID. Members are ordered by id,
	// starting from the given startID value. If limit is 0, all results
	// are returned, otherwise returns at most limit results. Will return
	// the startID of the next page or "" if no user accounts remain.
	ListTeamMembers(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error)
	// ListPublicTeamMembers gets a slice of all user accounts which are
	// public members of the team with the given teamID. Members are
	// ordered by id, starting from the given startID value. If limit is 0,
	// all results are returned, otherwise returns at most limit results.
	// Will return the startID of the next page or "" if no public members
	// remain.
	ListPublicTeamMembers(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error)
	// ListTeamAdmins queries for a slice of all accounts marked as admins
	// of the team with the given teamID. Admins are ordered by ID,
	// starting from the given startID value. If limit is 0, all results
	// are returned, otherwise returns at most limit results. Will return
	// the startID of the next page or "" if no team admins remain.
	ListTeamAdmins(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error)
	// ListTeamNonAdmins queries for a slice of all accounts not marked as
	// admins of the team with the given teamID. Non-admins are ordered by
	// id, starting from the given startID value. If limit is 0, all
	// results are returned, otherwise returns at most limit results. Will
	// return the startID of the next page or "" if no non-admins remain.
	ListTeamNonAdmins(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error)
	// ListTeamsInOrgForUser gets a slice of all teams with the given orgID
	// that the user with the given userID is a member of. Teams are
	// ordered by id, starting from the given startID value. If limit is 0,
	// all results are returned, otherwise returns at most limit results.
	// Will return the startID of the next page or "" if no teams remain.
	ListTeamsInOrgForUser(orgID, userID, startID string, limit uint) (memberTeams []MemberTeam, nextID string, err error)
	// DeleteTeamMembership removes the team membership for the given
	// teamID and userID.
	DeleteTeamMembership(teamID, userID string) error

	// CreateSession creates a new session using the given user ID and
	// duration. A secret and CSRF token should be generated as random
	// values. The ID field of the session is a hash of the secret which
	// should not be stored in the backend. Returns the new session object.
	// The Secret field of the session is only available once.
	CreateSession(userID string, duration time.Duration) (*Session, error)
	// GetSession retrieves the session corresponding to the given ID. If
	// no such session exists, the error will be ErrNoSuchSession.
	GetSession(id string) (*Session, error)
	// ExtendSession extends the given session's expiration to the given
	// additional duration of time from now. Sets the new expiration on the
	// given session.
	ExtendSession(session *Session, duration time.Duration) error
	// DeleteSession removes the session with the given ID.
	DeleteSession(id string) error
	// DeleteSessionsForUser deletes all sessions for the given userID
	// except for the session which matches the given excludeID.
	DeleteSessionsForUser(userID string, excludeID string) error
	// DeleteExpiredSessions removes any sessions with an expiration time
	// which is in the past.
	DeleteExpiredSessions() error

	// CreateService inserts a new service using the values supplied by the
	// given service. The ID field of the service is set to a random UUID.
	// The PK field of the service is set to a hash of the ownerID and
	// service name. If a service already exists with the same ownerID and
	// name the error will be ErrServiceExists.
	CreateService(service *Service) error
	// GetServiceByID retrieves the service with the given ID. If no such
	// service exists the returned error will be ErrNoSuchService.
	GetServiceByID(id string) (*Service, error)
	// GetServiceByName retrieves the service with the given ownerID and
	// name. If no such team exists the returned error will be
	// ErrNoSuchService.
	GetServiceByName(ownerID, name string) (*Service, error)
	// ListServicesForAccount returns a list services which are owned by
	// the account with the given accountID. Services are ordered by name,
	// starting from the given startName value. If limit is 0, all results
	// are returned, otherwise returns at most limit results. Will return
	// the startName of the next page or "" if no services remain.
	ListServicesForAccount(accountID, startName string, limit uint) (services []Service, nextName string, err error)
	// UpdateService updates the service with the given ID using any set
	// fields in the given updateFields struct. ServiceUpdateFields that
	// are left null are unchanged.
	UpdateService(id string, updateFields ServiceUpdateFields) error
	// DeleteService removes the service with the given ownerID and name.
	DeleteService(ownerID, name string) error

	// CreateServiceAuthCode creates a new service authcode using the
	// values from the given serviceAuthCode. To prevent secrets from being
	// stored in the backend, the service auth code ID is set to a hash of
	// a random secret. This secret is not stored in the backend and can
	// only be produced once by this function which sets it as the Code
	// field on the given serviceAuthCode value.
	CreateServiceAuthCode(serviceAuthCode *ServiceAuthCode) error
	// GetServiceAuthCode returns the service auth code with the given
	// code. If no such service auth code exists the error will be
	// ErrNoSuchServiceAuthCode.
	GetServiceAuthCode(code string) (*ServiceAuthCode, error)
	// DeleteServiceAuthCode deletes the service auth code with the given
	// ID. This should be done to ensure that auth codes are not used more
	// than once.
	DeleteServiceAuthCode(id string) error
	// DeleteExpiredServiceAuthCodes removes any service auth codes with an
	// expiration time which is in the past.
	DeleteExpiredServiceAuthCodes() error

	// CreateServiceSession creates a new service session linked to the
	// given root sessionID for the service with the given serviceID. To
	// prevent secrets from being stored in the backend, the service
	// session ID is set to a hash of a random secret UUID. This secret is
	// not stored in the backend and can only be produced once by this
	// function. The CSRF token is also set to a random UUID but is not
	// kept secret.
	CreateServiceSession(sessionID, serviceID string) (*ServiceSession, error)
	// GetServiceSession retrieves the service session corresponding to the
	// given ID. If no such session exists, the error will be
	// ErrNoSuchSession.
	GetServiceSession(id string) (*ServiceSession, error)
	// DeleteServiceSessions removes the service sessions associated with
	// the given root session ID.
	DeleteServiceSessions(sessionID string) error

	// SaveSigningKey stores the given signing key.
	SaveSigningKey(jwk JWK) error
	// ListSigningKeys returns a list of all currently used signing keys.
	ListSigningKeys() (keys []JWK, err error)
	// GetSigningKey returns the signing key with the given ID. If no such
	// key exists, the error will be ErrNoSuchSigningKey.
	GetSigningKey(id string) (*JWK, error)
	// ExtendSigningKeyExpiration updates the expiration time of the
	// signing key with the given id to be now plus the given additional
	// duration of time.
	ExtendSigningKeyExpiration(id string, duration time.Duration) error
	// DeleteExpiredSigningKeys removes any signing keys with an expiration
	// time which is in the past.
	DeleteExpiredSigningKeys() error

	// SaveServiceKeys stores the given keys for the given service. Keys
	// are updated if they already exist, extending key expiration in the
	// process.
	SaveServiceKeys(serviceID string, keys ...JWK) error
	// GetServiceKey retrieves the service key with the given serviceID and
	// keyID. If no such service key exists the error will be
	// ErrNoSuchServiceKey.
	GetServiceKey(serviceID, keyID string) (*ServiceKey, error)
	// DeleteExpiredServiceKeys removes any service keys with an expiration time
	// which is in the past.
	DeleteExpiredServiceKeys() error

	// CreateOrUpdateWorker creates a worker with the given ID and address.
	// If a worker with the given ID already exists, its address is set to
	// the given address value.
	CreateOrUpdateWorker(id, address string) error
	// GetWorker retrieves the worker object with the given ID.
	GetWorker(id string) (*Worker, error)
	// ListWorkers returns a slice of all registered workers.
	ListWorkers() ([]Worker, error)
	// DeleteWorker deletes the worker with the given ID.
	DeleteWorker(id string) error

	// CreateJob creates a new job using the values supplied by the given
	// job. The ID field of the job is set to a random UUID. The PK field
	// of the job is set to a hash of the cronID and scheduledAt time from
	// the job. If a job with the same PK already exists, the returned
	// error will be ErrJobExists.
	CreateJob(job *Job) error
	// ClaimJob attempts to claim the job with the given jobID for the
	// worker with the given workerID by performing a conditional update on
	// the job. If the job is still unclaimed, the job's workerID will be
	// set to the given workerID and its status will be set to "running"
	// and it lastUpdated time will be set to the current UTC time. If
	// claiming the job was successful, the job will be returned. If it is
	// nil, then the job has already been claimed by another worker.
	ClaimJob(jobID, workerID string) (*Job, error)
	// UpdateJobStatus updates the job with the given jobID by setting its
	// status to the given status value.
	UpdateJobStatus(jobID, status string) error
	// GetJob retrieves the job with the given jobID. If no such job exists
	// the returned error will be ErrNoSuchJob.
	GetJob(jobID string) (*Job, error)
	// GetMostRecentlyScheduledJobs retreives a slice of the most recently
	// scheduled jobs. The length of the slice will be at most limit
	GetMostRecentlyScheduledJobs(offset, limit uint) (jobs []Job, err error)
	// GetMostRecentlyScheduledJobsForWorker returns a slice of the most
	// recently scheduled jobs which are claimed by the worker with the
	// given workerID. The length of the slice will be at most limit. If
	// limit is 0, all jobs are returned.
	GetMostRecentlyScheduledJobsForWorker(workerID string, offset, limit uint) (jobs []Job, err error)
	// GetMostRecentlyScheduledJobsWithAction retreives a slice of the most
	// recently scheduled jobs with the given action. The length of the
	// slice will be at most limit. If limit is 0, all jobs are returned.
	GetMostRecentlyScheduledJobsWithAction(action string, offset, limit uint) (jobs []Job, err error)
	// GetMostRecentlyScheduledJobsForWorker returns a slice of the most
	// recently scheduled jobs which are claimed by the worker with the
	// given workerID. The length of the slice will be at most limit. If
	// limit is 0, all jobs are returned.
	GetMostRecentlyScheduledJobsForWorkerWithAction(workerID, action string, offset, limit uint) (jobs []Job, err error)
	// CountJobsWithActionStatus returns the number jobs with the given
	// action and status. This is a convenience method used to determine if
	// there are multiple jobs running which are performing the same action.
	CountJobsWithActionStatus(action, status string) (count uint, err error)
	// GetUnclaimedJobChanges begins listening for any changes to jobs with
	// an empty workerID field. Returns a channel on which the caller may
	// receive a stream of JobChange objects and an io.Closer which
	// performs necessary cleanup to end the stream's underlying goroutine.
	// After closing, the changeStream should be checked for a possible
	// remaining value.
	GetUnclaimedJobChanges() (changeStream <-chan JobChange, streamCloser io.Closer, err error)
	// DeleteJob deletes the job with the given jobID.
	DeleteJob(jobID string) error
}

type manager struct {
	session *rethink.Session
}

var _ Manager = &manager{}

// NewRethinkDBManager returns a new schema manager which connects to a
// RethinkDB cluster for storing data.
func NewRethinkDBManager(session *rethink.Session) Manager {
	return &manager{
		session: session,
	}
}
