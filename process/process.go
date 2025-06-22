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
	// start function timer -> tell main to run this process
	// end function timer -> tell main to end this process
	IsRecurring bool
}

type Scheduler struct {
	Schedule    []Process
	ProcessChan chan *Process
	StopChan    chan bool
}

func GetScheduler() *Scheduler {
	createSchedule := func() {
		s = &Scheduler{}
	}
	once.Do(createSchedule)
	return s
}

func (scheduler *Scheduler) AddProcess(process Process) error {
	insertIdx := 0

	for insertIdx < len(scheduler.Schedule) {
		scheduleProcess := scheduler.Schedule[insertIdx]

		if process.StartTime.Equal(scheduleProcess.StartTime) {
			return ErrSchedule
		} else if process.StartTime.After(scheduleProcess.StartTime) {
			insertIdx++
		} else {
			break
		}
	}
	if insertIdx > len(scheduler.Schedule) {
		return ErrSchedule
	}

	if insertIdx > 0 {
		previousProcessEndtime := scheduler.Schedule[insertIdx-1].EndTime
		if previousProcessEndtime.Equal(process.EndTime) || previousProcessEndtime.After(process.StartTime) {
			return ErrSchedule
		}
	}
	if insertIdx < len(scheduler.Schedule) {
		nextProcessStartTime := scheduler.Schedule[insertIdx].StartTime
		if process.EndTime.After(nextProcessStartTime) {
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
		scheduler.AddProcess(process)
	}
	return nil
}

func (scheduler *Scheduler) RunSchedule() {
	go func() {
		for {
			if <-scheduler.StopChan {
				return
			}
			if time.Now().After(scheduler.Schedule[0].EndTime) {
				s.RemoveProcess(0, false)
				scheduler.ProcessChan <- nil
			} else if time.Now().Equal(scheduler.Schedule[0].StartTime) || time.Now().After(scheduler.Schedule[0].StartTime) {
				scheduler.ProcessChan <- &scheduler.Schedule[0]
			}
		}
	}()
}
