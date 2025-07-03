package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) SchedulePage() string {
	pageTitle := "Schedule\n\n"
	processes := ""
	if m.scheduler == nil {
		return pageTitle
	}
	for i, process := range m.scheduler.Schedule {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		// reference time: Jan 2 15:04:05 2006 MST
		isRecurring := ""
		if process.IsRecurring {
			isRecurring = "(recurring)"
		}
		processes += cursor + process.Program.Name + ": " + process.StartTime.Format("02/01/2006 15:04") + " - " +
			process.Duration.Truncate(time.Minute).String() + " " + isRecurring + "\n"
	}
	return pageTitle + processes
}

func (m *Model) schedulePageKeyHandler(key string) tea.Cmd {
	switch key {
	case "ctrl+c", "q":
		return tea.Quit
	case "ctrl+r":
		m.err = nil
	case "p":
		m.page = programs
	case "h":
		m.page = help
	case "up":
		m.cursor = max(m.cursor-1, 0)
	case "down":
		m.cursor = min(m.cursor+1, len(m.scheduler.Schedule)-1)
	case "d":
		err := m.scheduler.RemoveProcess(m.cursor, true)
		if err != nil {
			m.err = err
		}
	}
	return nil
}
