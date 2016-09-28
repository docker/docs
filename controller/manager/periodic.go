package manager

import (
	"time"
)

// Implement a basic polling algorithm to run periodic tasks
// Each periodic task is responsible for enabling/disabling itself
// based on its own configuration settings
func (m DefaultManager) startPeriodicTasks() {
	doChecks := func(m DefaultManager) {
		// If this list gets long, we might want to make it an actual list
		m.reportUsage()
		m.periodicLicenseCheck()
		m.periodicAuthSync()
		m.periodicCheckForUpdates()
	}

	ticker := time.NewTicker(PeriodicInterval)
	go func(m DefaultManager) {

		doChecks(m) // Run once at startup to initialize everything

		for range ticker.C {
			doChecks(m)
		}
	}(m)
}
