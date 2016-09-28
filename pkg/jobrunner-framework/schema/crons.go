package schema

import (
	"errors"
	"fmt"
	"io"

	"github.com/docker/dhe-deploy/rethinkutil"
	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	// ErrNoSuchCron conveys that a cron with the given name or id does not
	// exist.
	ErrNoSuchCron = errors.New("no such cron")
)

// A Cron is used to run a job on a regular schedule.
type Cron struct {
	ID          string            `gorethink:"id"`          // Randomly generated UUID for foreign references.
	Action      string            `gorethink:"action"`      // The action to be performed by jobs spawned from this cron.
	Schedule    string            `gorethink:"schedule"`    // The schedule for this cron as a cronspec string: (seconds) (minutes) (hours) (day of month) (month) (day of week) or @hourly, @weekly, etc.
	Retries     int               `gorethink:"retries"`     // Number of retries to create the job with
	Parameters  map[string]string `gorethink:"parameters"`  // Parameters to start the jobs with
	Deadline    string            `gorethink:"deadline"`    // After this amount of time has passed, a SIGTERM will be sent
	StopTimeout string            `gorethink:"stopTimeout"` // This long after SIGTERM is sent, SIGKILL will be sent if the proccess is still alive
}

var CronsTable = rethinkutil.Table{
	Name:       "crons",
	PrimaryKey: "action", // Guarantees uniqueness.
	SecondaryIndexes: map[string][]string{
		"id": nil, // For quick lookups by ID.
	},
}

// UpateCron creates or updates a cron with the given action and sets its
// schedule to the given schedule. The ID field of the cron is set to a random
// UUID.
func (m *jobrunnerManager) UpdateCron(cron *Cron) (*Cron, error) {
	cron.ID = uuid.NewV4().String()

	if _, err := CronsTable.Term(m.DB()).Get(cron.Action).Replace(cron).RunWrite(m.Session()); err != nil {
		return nil, err
	}

	return cron, nil
}

// ListCrons lists all scheduled crons
func (m *jobrunnerManager) ListCrons() ([]Cron, error) {
	cursor, err := CronsTable.Term(m.DB()).Run(m.Session())
	if err != nil {
		return nil, err
	}

	crons := []Cron{}
	if err := cursor.All(&crons); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}

	return crons, nil
}

// GetCron retrieves the cron for the given action. If no such cron exists the
// returned error will be ErrNoSuchCron.
func (m *jobrunnerManager) GetCron(action string) (*Cron, error) {
	cursor, err := CronsTable.Term(m.DB()).Get(action).Run(m.Session())
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var cron Cron
	if err := cursor.One(&cron); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchCron
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &cron, nil
}

// CronChange is used to deliver old and new values of a cron as part of a
// changes stream.
type CronChange struct {
	OldValue *Cron `gorethink:"old_val"`
	NewValue *Cron `gorethink:"new_val"`
}

// GetCronChanges begins listening for any changes to cron configs. Returns a
// channel on which the caller may receive a stream of CronChange objects and
// an io.Closer which performs necessary cleanup to end the stream's underlying
// goroutine. After closing, the changeStream should be checked for a possible
// remaining value.
func (m *jobrunnerManager) GetCronChanges() (changeStream <-chan CronChange, streamCloser io.Closer, err error) {
	cursor, err := CronsTable.Term(m.DB()).Changes(
		rethink.ChangesOpts{IncludeInitial: true},
	).Run(m.Session())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query db: %s", err)
	}

	changes := make(chan CronChange)
	cursor.Listen(changes)

	return changes, cursor, nil
}

// DeleteCron deletes the cron with the given action.
func (m *jobrunnerManager) DeleteCron(action string) error {
	if _, err := CronsTable.Term(m.DB()).Get(action).Delete().RunWrite(m.Session()); err != nil {
		return fmt.Errorf("unable to delete cron from database: %s", err)
	}

	return nil
}
