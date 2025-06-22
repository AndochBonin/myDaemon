package process

import (
	"slices"
	"sync"
	"time"

	"github.com/AndochBonin/myDaemon/program"
)

var (
	s      *Schedule
	once          sync.Once
	ScheduleError error
)

type Process struct {
	Program     program.Program
	StartTime   time.Time
	EndTime     time.Time
	// start function timer -> tell main to run this process
	// end function timer -> tell main to end this process
	IsRecurring bool
}

type Schedule []Process

func GetSchedule() *Schedule {
	createSchedule := func() {
		s = &Schedule{}
	}
	once.Do(createSchedule)
	return s
}

func (schedule *Schedule) AddProcess(process Process) error {
	insertIdx := 0

	for insertIdx < len(*schedule) {
		scheduleProcess := (*schedule)[insertIdx]

		if process.StartTime.Equal(scheduleProcess.StartTime) {
			return ScheduleError
		} else if process.StartTime.After(scheduleProcess.StartTime) {
			insertIdx++
		} else {
			break
		}
	}

	if insertIdx < len(*schedule) {
		nextProcessStartTime := (*schedule)[insertIdx].StartTime

		if process.EndTime.After(nextProcessStartTime) {
			return ScheduleError
		}
	} else if insertIdx > len(*schedule) {
		return ScheduleError
	}
	*schedule = slices.Insert(*schedule, insertIdx, process)
	return nil
}

func (schedule *Schedule) RemoveProcess(processID int, endRecurrence bool) error {
	if processID < 0 || processID >= len(*schedule) {
		return ScheduleError
	}
	process := (*schedule)[processID]
	*schedule = slices.Delete(*schedule, processID, processID+1)

	if process.IsRecurring && !endRecurrence {
		timeOffset := time.Hour * 24
		process.StartTime = process.StartTime.Add(timeOffset)
		process.EndTime = process.EndTime.Add(timeOffset)
		schedule.AddProcess(process)
	}
	return nil
}

func (schedule *Schedule) RunSchedule(stop chan bool, process chan *Process) {
	go func() {
		for {
			if <-stop {
				return
			}
			if time.Now().After((*schedule)[0].EndTime) {
				s.RemoveProcess(0, false)
				process<-nil
			} else if time.Now().Equal((*schedule)[0].StartTime) || time.Now().After((*schedule)[0].StartTime) {
				process<-&(*schedule)[0]
			}
		}	
	}()
}