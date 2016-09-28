package util

import (
	"fmt"
)

// MinutesToCron converts the given number of minutes into a cronspec string.
func MinutesToCron(minutes int) string {
	hours := minutes / 60
	if hours <= 1 {
		// This cron runs hourly at the top of the hour.
		return "@hourly"
	}
	if hours < 24 {
		// This cron runs at the top of the hour every `h` hours offset
		// from midnight.
		// If `h` is 5 it will run at midnight, 5am, 10am, 3pm, and
		// 8pm.
		return fmt.Sprintf("0 0 */%d * *", hours)
	}

	days := hours / 24
	if days == 1 {
		// This cron runs daily at midnight.
		return "@daily"
	}
	if days < 28 {
		// This cron runs at midnight every `d` days offset from the
		// first of the month.
		// If `d` is 7 it will run on the 1st, 8th, 15th, 22nd, and
		// 29th.
		return fmt.Sprintf("0 0 0 */%d *", days)
	}

	// If the interval is greater than or equal to 28 days, default to
	// running monthly at midnight the first day of every month.
	return "@monthly"
}
