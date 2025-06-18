package process

import (
	"time"

	"github.com/AndochBonin/myDaemon/program"
)

type Process struct {
	Program program.Program
	StartTime time.Time
	DurationNanoseconds time.Duration
	IsRecurring bool
}

type Schedule []Process

// Define 