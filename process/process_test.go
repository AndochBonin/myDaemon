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
				t.Error("GetScheduler created a new scheduler in a separate goroutine")
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
		Duration: duration,
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

func TestRemoveProcess(t *testing.T) {
	scheduler := GetScheduler()
	scheduler.Schedule = nil

	t.Run("remove invalid index", func(t *testing.T) {
		scheduler.Schedule = nil
		process := Process{}
		err := scheduler.AddProcess(process)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		err = scheduler.RemoveProcess(-1, true)
		if err != nil {
			t.Errorf("Expected %v got %v instead", ErrSchedule, err)
		}
		err = scheduler.RemoveProcess(1, true)
		if err != nil {
			t.Errorf("Expected %v got %v instead", ErrSchedule, err)
		}
		checkSchedule(t, scheduler.Schedule, []Process{process})
	})

	t.Run("remove valid index, end recurrence", func(t *testing.T) {
		scheduler.Schedule = nil
		addErr := scheduler.AddProcess(Process{})
		removeErr := scheduler.RemoveProcess(0, true)
		if addErr != nil || removeErr != nil {
			t.Errorf("Unexpected errors: %v, %v", addErr, removeErr)
		}
		checkSchedule(t, scheduler.Schedule, []Process{})
	})

	t.Run("remove valid index, keep recurrence", func(t *testing.T) {
		scheduler.Schedule = nil
		process := newMockProcess(t, "recurring", time.Second, time.Hour, true)
		addErr := scheduler.AddProcess(process)
		removeErr := scheduler.RemoveProcess(0, false)
		if addErr != nil || removeErr != nil {
			t.Errorf("Unexpected errors: %v, %v", addErr, removeErr)
		}
		process.StartTime = process.StartTime.Add(time.Hour * 24)
		checkSchedule(t, scheduler.Schedule, []Process{process})
	})

	t.Run("process with recurrence reinserted in correct order", func(t *testing.T) {
		scheduler.Schedule = nil
		testProcess := newMockProcess(t, "recurring", time.Second, time.Hour, true)
		beforeTestProcess := newMockProcess(t, "earlier", -time.Hour, time.Minute*30, false)
		afterTestProcess := newMockProcess(t, "later", time.Hour*2, time.Minute*30, false)
		afterRecurredProcess := newMockProcess(t, "later", time.Hour*27, time.Minute*30, false)

		scheduler.AddProcess(testProcess)
		scheduler.AddProcess(beforeTestProcess)
		scheduler.AddProcess(afterTestProcess)
		scheduler.AddProcess(afterRecurredProcess)
		// now schedule should look like: {earlier, recurring, later, muchlater}
		removeErr := scheduler.RemoveProcess(1, false)
		if removeErr != nil {
			t.Errorf("Unexpected error: %v", removeErr)
		}
		testProcess.StartTime = testProcess.StartTime.Add(time.Hour * 24)

		// recurring should slot between later and much later
		expected := []Process{beforeTestProcess, afterTestProcess, testProcess, afterRecurredProcess}
		checkSchedule(t, scheduler.Schedule, expected)
	})

	t.Run("recurring process reinsertion fails when overlapping with existing process", func(t *testing.T) {
		scheduler.Schedule = nil
		testProcess := newMockProcess(t, "recurring", time.Second, time.Hour, true)
		beforeTestProcess := newMockProcess(t, "earlier", -time.Hour, time.Minute*30, false)
		afterTestProcess := newMockProcess(t, "later", time.Hour*2, time.Minute*30, false)
		duringRecurredProcess := newMockProcess(t, "later", time.Hour*24, time.Hour, false)

		scheduler.AddProcess(testProcess)
		scheduler.AddProcess(beforeTestProcess)
		scheduler.AddProcess(afterTestProcess)
		scheduler.AddProcess(duringRecurredProcess)
		// now schedule should look like: {earlier, recurring, later, muchlater}
		removeErr := scheduler.RemoveProcess(1, false)
		if removeErr != ErrSchedule {
			t.Errorf("Expected: %v got %v instead", ErrSchedule, removeErr)
		}
		expected := []Process{beforeTestProcess, afterTestProcess, duringRecurredProcess}
		checkSchedule(t, scheduler.Schedule, expected)
	})
}

func TestUpdateSchedule(t *testing.T) {
	scheduler := GetScheduler()
	scheduler.Schedule = nil
	t.Run("returns nil when schedule is empty", func(t *testing.T) {
		err := scheduler.UpdateSchedule()
		scheduleLength := len(scheduler.Schedule)
		if err != nil {
			t.Errorf("Unexpected error. Expected %v got %v", nil, err)
		}
		if scheduleLength != 0 {
			t.Errorf("Expected schedule length 0 but got length %v", scheduleLength)
		}
	})

	t.Run("leaves schedule as is if time is within start and end time of schedule[0]", func(t *testing.T) {
		scheduler.Schedule = nil
		ongoingProcess := newMockProcess(t, "ongoing process", -time.Minute, time.Hour, false)
		scheduler.AddProcess(ongoingProcess)
		err := scheduler.UpdateSchedule()
		if err != nil {
		t.Errorf("Unexpected error. Expected %v got %v", nil, err)
		}
		expected := []Process{ongoingProcess}
		if !reflect.DeepEqual(scheduler.Schedule, expected) {
			t.Errorf("Expected %v got %v", expected, scheduler.Schedule)
		}
	})

	t.Run("removes process at index 0 if time is past end time of process", func(t *testing.T) {
		scheduler.Schedule = nil
		completedProcess := newMockProcess(t, "completed process", -time.Hour, time.Minute, false)
		scheduler.AddProcess(completedProcess)
		err := scheduler.UpdateSchedule()
		if err != nil {
		t.Errorf("Unexpected error. Expected %v got %v", nil, err)
		}
		expected := []Process{}
		if !reflect.DeepEqual(scheduler.Schedule, expected) {
			t.Errorf("Expected %v got %v", expected, scheduler.Schedule)
		}
	})

	t.Run("leaves schedule as is if time is before start time of schedule[0]", func(t *testing.T) {
		scheduler.Schedule = nil
		pendingProcess := newMockProcess(t, "pending process", time.Hour, time.Minute, false)
		scheduler.AddProcess(pendingProcess)
		err := scheduler.UpdateSchedule()
		if err != nil {
		t.Errorf("Unexpected error. Expected %v got %v", nil, err)
		}
		expected := []Process{pendingProcess}
		if !reflect.DeepEqual(scheduler.Schedule, expected) {
			t.Errorf("Expected %v got %v", expected, scheduler.Schedule)
		}
	})
}

func TestGetCurrentProcess(t *testing.T) {
	// test return nil when schedule is empty
	scheduler := GetScheduler()
	t.Run("return nil when schedule is empty", func(t *testing.T) {
		scheduler.Schedule = nil
		process := scheduler.GetCurrentProcess()
		if process != nil {
			t.Errorf("Expected nil got %v", process)
		}
	})
	// test return nil when process hasn't started
	t.Run("return nil when process at index 0 hasn't started", func(t *testing.T) {
		scheduler.Schedule = nil
		pendingProcess := newMockProcess(t, "pending process", time.Hour, time.Minute, false)
		scheduler.AddProcess(pendingProcess)
		process := scheduler.GetCurrentProcess()
		if process != nil {
			t.Errorf("Expected nil got %v", process)
		}
	})
	// test return process when time is within start and end time
	t.Run("return process at index 0 when time is within start and end time", func(t *testing.T) {
		scheduler.Schedule = nil
		ongoingProcess := newMockProcess(t, "ongoing process", -time.Minute, time.Hour, false)
		scheduler.AddProcess(ongoingProcess)
		process := scheduler.GetCurrentProcess()
		if process == nil ||!reflect.DeepEqual(*process, ongoingProcess) {
			t.Errorf("Expected %v got %v", ongoingProcess, process)
		}
	})
}