package process

import (
	"errors"
	"github.com/AndochBonin/myDaemon/program"
	"slices"
	"sync"
	"time"
)

var (
	s           *Scheduler
	once        sync.Once
	ErrSchedule = errors.New("invalid process scheduling")
)

type Process struct {
	Program     program.Program
	StartTime   time.Time
	Duration    time.Duration
	IsRecurring bool
}

type Scheduler struct {
	Schedule []Process
}

func GetScheduler() *Scheduler {
	createSchedule := func() {
		s = &Scheduler{}
	}
	once.Do(createSchedule)
	return s
}

func (scheduler *Scheduler) AddProcess(process Process) error {
	insertIdx, match := slices.BinarySearchFunc(scheduler.Schedule, process, func(E Process, T Process) int {
		if E.StartTime.Equal(T.StartTime) {
			return 0
		} else if E.StartTime.Before(T.StartTime) {
			return -1
		} else {
			return 1
		}
	})

	if match || insertIdx > len(scheduler.Schedule) || insertIdx < 0 {
		return ErrSchedule
	}

	if insertIdx > 0 {
		previousProcess := scheduler.Schedule[insertIdx-1]
		previousProcessEndtime := previousProcess.StartTime.Add(previousProcess.Duration)
		if previousProcessEndtime.Equal(process.StartTime) || previousProcessEndtime.After(process.StartTime) {
			return ErrSchedule
		}
	}
	if insertIdx < len(scheduler.Schedule) {
		nextProcessStartTime := scheduler.Schedule[insertIdx].StartTime
		if process.StartTime.Add(process.Duration).Equal(nextProcessStartTime) ||
			process.StartTime.Add(process.Duration).After(nextProcessStartTime) {
			return ErrSchedule
		}
	}
	scheduler.Schedule = slices.Insert(scheduler.Schedule, insertIdx, process)
	return nil
}

func (scheduler *Scheduler) RemoveProcess(processID int, endRecurrence bool) error {
	if processID < 0 || processID >= len(scheduler.Schedule) {
		return nil
	}
	process := (scheduler.Schedule)[processID]
	scheduler.Schedule = slices.Delete(scheduler.Schedule, processID, processID+1)

	if process.IsRecurring && !endRecurrence {
		timeOffset := time.Hour * 24
		process.StartTime = process.StartTime.Add(timeOffset)
		err := scheduler.AddProcess(process)
		if err != nil {
			return err
		}
	}
	return nil
}

func (scheduler *Scheduler) UpdateSchedule() error {
	if len(scheduler.Schedule) == 0 {
		return nil
	} else if process := scheduler.Schedule[0]; time.Now().After(process.StartTime.Add(process.Duration)) {
		return scheduler.RemoveProcess(0, !process.IsRecurring)
	}
	return nil
}

func (scheduler *Scheduler) GetCurrentProcess() *Process {
	if len(scheduler.Schedule) == 0 {
		return nil
	} else if time.Now().After(scheduler.Schedule[0].StartTime) {
		return &scheduler.Schedule[0]
	}
	return nil
}
