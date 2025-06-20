package process

import (
	"github.com/AndochBonin/myDaemon/program"
	"slices"
	"sync"
	"time"
)

var (
	scheduler      *Scheduler
	once           sync.Once
	SchedulerError error
)

type Process struct {
	Program             program.Program
	StartTime           time.Time
	DurationNanoseconds time.Duration
	IsRecurring         bool
}

type Scheduler struct {
	Schedule         []Process
	CurrentProcess   *Process
	currentProcessID int
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
		} else if process.StartTime.After(scheduleProcess.StartTime) {
			insertIdx++
		} else {
			break
		}
	}

	if insertIdx < len(scheduler.Schedule) {
		processEndTime := process.StartTime.Add(process.DurationNanoseconds)
		nextProcessStartTime := scheduler.Schedule[insertIdx].StartTime

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
	scheduler.Schedule = slices.Delete(scheduler.Schedule, processID, processID+1)
	return nil
}

func (scheduler *Scheduler) UpdateCurrentProcess() bool {
	currentTime := time.Now().UTC()

	if scheduler.CurrentProcess != nil {
		currentProcessEndTime := scheduler.CurrentProcess.StartTime.Add(scheduler.CurrentProcess.DurationNanoseconds)

		if currentTime.Before(currentProcessEndTime) {
			return false
		}

		if !scheduler.CurrentProcess.IsRecurring {
			scheduler.RemoveProcess(scheduler.currentProcessID)
		} else {
			scheduler.currentProcessID++
		}
		if scheduler.currentProcessID == len(scheduler.Schedule) {
			scheduler.currentProcessID = 0
		}
		scheduler.CurrentProcess = nil
		scheduler.UpdateCurrentProcess()
		return true
	} else {
		nextProcessStartTime := scheduler.Schedule[scheduler.currentProcessID].StartTime

		if currentTime.Equal(nextProcessStartTime) || currentTime.After(nextProcessStartTime) {
			scheduler.CurrentProcess = &scheduler.Schedule[scheduler.currentProcessID]
			return true
		}
	}
	return false
}

func (scheduler *Scheduler) RunSchedule(done chan bool, process chan *Process) {
	go func() {
		for {
			if <-done {
				return
			}
			if scheduler.UpdateCurrentProcess() { // update ? send next process through channel
				process <- scheduler.CurrentProcess
			}
		}
	}()
}
