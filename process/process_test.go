package process

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/AndochBonin/myDaemon/program"
)

func TestGetScheduler(t *testing.T) {
	// test that GetScheduler always returns the same instance (singleton)
	scheduler := GetScheduler()
	anotherScheduler := GetScheduler()
	if anotherScheduler != scheduler {
		t.Fatal("GetScheduler created a new scheduler")
	}
	// test concurrent calls to GetScheduler return the same instance (thread safety)
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if GetScheduler() != scheduler {
				t.Error("GetScheduler created a new scheduler in a separater goroutine")
			}
		}()
	}
	wg.Wait()
}

func checkSchedule(t *testing.T, actual, expected []Process) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected schedule:\n%v\nGot:\n%v", expected, actual)
	}
}

func newMockProcess(t *testing.T, name string, startOffset time.Duration, duration time.Duration, isRecurring bool) Process {
	t.Helper()
	start := time.Now().Add(startOffset)
	return Process{
		Program:     program.Program{Name: name},
		StartTime:   start,
		EndTime:     start.Add(duration),
		IsRecurring: isRecurring,
	}
}


func TestAddProcess(t *testing.T) {
	scheduler := GetScheduler()
	scheduler.Schedule = nil

	processStartNow := newMockProcess(t, "NOW", 0, time.Hour, false)
	processStartEarlier := newMockProcess(t, "EARLIER", -time.Hour, time.Minute*59, false)
	processStartLater := newMockProcess(t, "LATER", 2*time.Hour, time.Hour, false)
	processOverlapLater := newMockProcess(t, "OVERLAP LATER", 2*time.Hour+time.Minute, 2*time.Hour, false)
	processSameTimeLater := Process{
		Program:     program.Program{Name: "SAME TIME LATER"},
		StartTime:   processStartLater.StartTime,
		EndTime:     processStartLater.EndTime.Add(time.Hour),
		IsRecurring: false,
	}

	t.Run("adds initial process", func(t *testing.T) {
		err := scheduler.AddProcess(processStartNow)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		checkSchedule(t, scheduler.Schedule, []Process{processStartNow})
	})

	t.Run("adds later process in order", func(t *testing.T) {
		err := scheduler.AddProcess(processStartLater)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		checkSchedule(t, scheduler.Schedule, []Process{processStartNow, processStartLater})
	})

	t.Run("adds earlier process in order", func(t *testing.T) {
		err := scheduler.AddProcess(processStartEarlier)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		checkSchedule(t, scheduler.Schedule, []Process{processStartEarlier, processStartNow, processStartLater})
	})

	t.Run("rejects overlapping process", func(t *testing.T) {
		err := scheduler.AddProcess(processOverlapLater)
		if err != ErrSchedule {
			t.Errorf("Expected ErrSchedule, got: %v", err)
		}
		checkSchedule(t, scheduler.Schedule, []Process{processStartEarlier, processStartNow, processStartLater})
	})

	t.Run("rejects process with duplicate start time", func(t *testing.T) {
		err := scheduler.AddProcess(processSameTimeLater)
		if err != ErrSchedule {
			t.Errorf("Expected ErrSchedule, got: %v", err)
		}
		checkSchedule(t, scheduler.Schedule, []Process{processStartEarlier, processStartNow, processStartLater})
	})
}

func TestDeleteProcess(t *testing.T) {
	// test deleting a process by valid index
	// test deleting a recurring process and keeping recurrence (re-added with +24h offset)
	// test deleting a recurring process and ending recurrence (not re-added)
	// test deleting with invalid index (negative or out of bounds)
	// test recurrence is re-added in correct sorted position
	// test recurrence fails if it overlaps with an existing process
}

func TestRunScheduler(t *testing.T) {
	// test that scheduler stops when StopChan receives true
	// test process switching: remove old, send nil, send new when time passes
	// test that process is sent through ProcessChan at correct start time
	// test behavior when Schedule is empty (should not panic)
	// test concurrency behavior during schedule execution (e.g., run and remove at same time)
	// test using short durations or time mocking to avoid long waits
}
