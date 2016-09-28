package schema

import (
	"errors"
	"fmt"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/constants"
	"github.com/docker/dhe-deploy/rethinkutil"
	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// ActionConfigs are always looked up immediately before running a job
// The parameters act as defaults for the parameters of the job. If the job is
// created by a cron, the parameters from the cron will take precedence

var (
	// ErrNoSuchActionConfig conveys that an action config with the given action doesn't exist.
	ErrNoSuchActionConfig = errors.New("no such action config")
)

type ActionConfig struct {
	ID               string            `gorethink:"id"`
	Action           string            `gorethink:"action"`
	MaxJobsPerWorker int               `gorethink:"maxJobsPerWorker"`
	HeartbeatTimeout string            `gorethink:"heartbeatTimeout"`
	Parameters       map[string]string `gorethink:"parameters"`
}

var ActionConfigsTable = rethinkutil.Table{
	Name:       "action_configs",
	PrimaryKey: "action", // Guarantees uniqueness.
	SecondaryIndexes: map[string][]string{
		"id": nil, // For quick lookups by ID.
	},
}

// UpateActionConfig creates or updates a actionConfig with the given action. The ID
// field of the actionConfig is set to a random UUID.
func (m *jobrunnerManager) UpdateActionConfig(actionConfig *ActionConfig) (*ActionConfig, error) {
	actionConfig.ID = uuid.NewV4().String()

	if _, err := ActionConfigsTable.Term(m.DB()).Get(actionConfig.Action).Replace(actionConfig).RunWrite(m.Session()); err != nil {
		return nil, err
	}

	return actionConfig, nil
}

// ListActionConfigs lists all scheduled actionConfigs
func (m *jobrunnerManager) ListActionConfigs() ([]ActionConfig, error) {
	cursor, err := ActionConfigsTable.Term(m.DB()).Run(m.Session())
	if err != nil {
		return nil, err
	}

	actionConfigs := []ActionConfig{}
	if err := cursor.All(&actionConfigs); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}
	fmt.Println("action configs got")
	fmt.Println(actionConfigs)

	return actionConfigs, nil
}

// GetActionConfig retrieves the action config for the given action. If no such action config exists the
// returned error will be ErrNoSuchActionConfig.
func (m *jobrunnerManager) GetActionConfig(action string) (*ActionConfig, error) {
	cursor, err := ActionConfigsTable.Term(m.DB()).Get(action).Run(m.Session())
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var actionConfig ActionConfig
	if err := cursor.One(&actionConfig); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchActionConfig
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &actionConfig, nil
}

func (m *jobrunnerManager) SafeGetActionConfig(action string) (*ActionConfig, error) {
	ac, err := m.GetActionConfig(action)
	if err != nil {
		ac = &ActionConfig{
			Action:           action,
			HeartbeatTimeout: constants.DefaultHeartbeatTimeout.String(),
		}
	}
	if err == ErrNoSuchActionConfig {
		err = nil
	}
	return ac, err
}

// DeleteActionConfig deletes the action config with the given action.
func (m *jobrunnerManager) DeleteActionConfig(action string) error {
	if _, err := ActionConfigsTable.Term(m.DB()).Get(action).Delete().RunWrite(m.Session()); err != nil {
		return fmt.Errorf("unable to delete actionConfig from database: %s", err)
	}

	return nil
}
