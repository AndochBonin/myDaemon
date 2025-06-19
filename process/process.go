package process

import (
	"sync"
	"time"
	"slices"
	"github.com/AndochBonin/myDaemon/program"
)

var (
	scheduler *Scheduler
	once sync.Once
	SchedulerError error
)

type Process struct {
	Program program.Program
	StartTime time.Time
	DurationNanoseconds time.Duration
	IsRecurring bool
}

type Scheduler struct {
	CurrentProcess *Process
	Schedule []Process
}

func GetScheduler() *Scheduler {
	createScheduler := func() {
		scheduler = &Scheduler{}
	}
	once.Do(createScheduler)
	return scheduler
}

func (scheduler *Scheduler) AddProcess(process Process) error {
	insertIdx := 0

	for insertIdx < len(scheduler.Schedule) {
		scheduleProcess := scheduler.Schedule[insertIdx]

		if process.StartTime.Equal(scheduleProcess.StartTime) {
			return SchedulerError
		} else if  process.StartTime.After(scheduleProcess.StartTime) {
			insertIdx++
		} else {
			break
		}
	}

	if insertIdx < len(scheduler.Schedule) {
		nextProcessStartTime := scheduler.Schedule[insertIdx].StartTime
		processEndTime := process.StartTime.Add(process.DurationNanoseconds)

		if processEndTime.After(nextProcessStartTime) {
			return SchedulerError
		}
	}

	scheduler.Schedule = slices.Insert(scheduler.Schedule, insertIdx, process)
	return nil
}

func (scheduler *Scheduler) RemoveProcess(processID int) error {
	if processID < 0 || processID >= len(scheduler.Schedule) {
		return SchedulerError
	}
	scheduler.Schedule = slices.Delete(scheduler.Schedule, processID, processID + 1)
	return nil
}