package cmd

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
)

// TickJob is the worker subcommand for the tick job. It's useful as a test.
var TickJob = cli.Command{
	Name:   "tick",
	Usage:  "Tick Tock; Don't Stop",
	Action: runTick,
}

const timeFormat = "2006-01-02 15:04:05 MST"

func runTick(*cli.Context) error {
	for tick := true; true; tick = !tick {
		time.Sleep(time.Second)

		msg := "tick"
		if !tick {
			msg = "tock"
		}

		fmt.Println(time.Now().UTC().Format(timeFormat), msg)
	}

	return nil
}
