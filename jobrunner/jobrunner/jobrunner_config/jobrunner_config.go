package jobrunner_config

import (
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
)

// The currently available actions to perform in a job.
const (
	ActionGC           = "gc"
	ActionGCRecovery   = "gc_recovery"
	ActionSleep        = "sleep"
	ActionTagMigration = "tagmigration"
)

// RegisteredActions is the set of valid actions available to run in a job and
// their corresponding configuration.
var RegisteredActions = map[string]worker.ActionInfo{
	ActionGC: {
		Command:   "gc",
		Exclusive: true,
		CleanupJob: &schema.Job{
			Action: ActionGCRecovery,
		},
	},
	ActionGCRecovery: {
		Command:   "gc",
		Args:      []string{"set-rw-mode"},
		Exclusive: true,
	},
	ActionSleep: {
		Command:   "sleep",
		Args:      []string{"60"},
		Exclusive: true,
	},
	ActionTagMigration: {
		Command:   "tagmigration",
		Exclusive: true,
	},
}
