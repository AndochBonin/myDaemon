package process

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/AndochBonin/myDaemon/program"
)

var (
	s           *Scheduler
	once        sync.Once
	ErrSchedule = errors.New("invalid process scheduling")
)

type Process struct {
	Program   program.Program
	StartTime time.Time
	EndTime   time.Time
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
		previousProcessEndtime := scheduler.Schedule[insertIdx-1].EndTime
		if previousProcessEndtime.Equal(process.StartTime) || previousProcessEndtime.After(process.StartTime) {
			return ErrSchedule
		}
	}
	if insertIdx < len(scheduler.Schedule) {
		nextProcessStartTime := scheduler.Schedule[insertIdx].StartTime
		if process.EndTime.Equal(nextProcessStartTime) || process.EndTime.After(nextProcessStartTime) {
			return ErrSchedule
		}
	}
	scheduler.Schedule = slices.Insert(scheduler.Schedule, insertIdx, process)
	return nil
}

func (scheduler *Scheduler) RemoveProcess(processID int, endRecurrence bool) error {
	if processID < 0 || processID >= len(scheduler.Schedule) {
		return ErrSchedule
	}
	process := (scheduler.Schedule)[processID]
	scheduler.Schedule = slices.Delete(scheduler.Schedule, processID, processID+1)

	if process.IsRecurring && !endRecurrence {
		timeOffset := time.Hour * 24
		process.StartTime = process.StartTime.Add(timeOffset)
		process.EndTime = process.EndTime.Add(timeOffset)
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
	} else if time.Now().After(scheduler.Schedule[0].EndTime) {
		return scheduler.RemoveProcess(0, false)
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

func isValidProcess(process Process) bool {
	return process.StartTime.Before(process.EndTime)
}
