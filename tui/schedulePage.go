package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//var scheduleFile string = "./storage/schedule.json" //for testing

func (m *Model) SchedulePage() string {
	pageTitle := "Schedule"
	pageKeys := "delete [d]"
	processes := ""
	if m.scheduler == nil {
		return pageTitleStyle.Render(pageTitle) + "\n\n" + navStyle.Render(pageKeys) + "\n\n" + navStyle.Render("nothing here yet.")
	}
	for i, process := range m.scheduler.Schedule {
		var style lipgloss.Style = textContentStyle
		cursor := " "
		if m.cursor == i {
			style = focusedStyle.Bold(true)
			cursor = "> "
		}
		// reference time: Jan 2 15:04:05 2006 MST
		isRecurring := ""
		if process.IsRecurring {
			isRecurring = "(R)"
		}
		processes += style.Render(cursor + process.Program.Name + ": " + process.StartTime.Format("02/01/2006 15:04") + " - " +
			process.Duration.Truncate(time.Second).String() + " " + isRecurring) + "\n"
	}
	if processes == "" {
		processes = navStyle.Render("nothing yet.")
	}
	return pageTitleStyle.Render(pageTitle) + "\n\n" + navStyle.Render(pageKeys) + "\n\n" + processes
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
		err := m.scheduler.RemoveProcess(m.cursor, true, scheduleFile)
		if err != nil {
			m.err = err
		}
	}
	return nil
}
