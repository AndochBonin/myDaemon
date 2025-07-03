package tui

import (
	"errors"
	"time"

	"github.com/AndochBonin/myDaemon/process"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) ProcessDetailsPage() string {
	pageTitle := "Process Details\n\n"

	return pageTitle + m.programList[m.cursor].Name +
		"\n\nStart Time:\n" + m.processDetails.startTime.View() +
		"\n\nDuration:\n" + m.processDetails.duration.View() +
		"\n\nIs Recurring?:\n" + m.processDetails.isRecurring.View()
}

func (m *Model) initProcessDetailsInput() tea.Cmd {
	startTime := textinput.New()
	startTime.Placeholder = "00:00"
	startTime.Cursor.Style = cursorStyle
	startTime.PromptStyle = focusedStyle
	startTime.TextStyle = focusedStyle
	startTime.CharLimit = 10
	startTime.Width = 10
	m.processDetails.startTime = startTime

	duration := textinput.New()
	duration.Placeholder = "0h0m"
	duration.Cursor.Style = cursorStyle
	duration.PromptStyle = focusedStyle
	duration.TextStyle = focusedStyle
	duration.CharLimit = 10
	duration.Width = 10
	m.processDetails.duration = duration

	isRecurring := textinput.New()
	isRecurring.Placeholder = "N"
	isRecurring.Cursor.Style = cursorStyle
	isRecurring.PromptStyle = focusedStyle
	isRecurring.TextStyle = focusedStyle
	isRecurring.CharLimit = 1
	isRecurring.Width = 1
	m.processDetails.isRecurring = isRecurring

	return m.processDetails.startTime.Focus()
}

func (m *Model) processDetailsPageKeyHandler(key string) tea.Cmd {
	switch key {
	case "ctrl+c":
		return tea.Quit
	case "esc":
		m.page = programs
	case "enter":
		switch m.processDetails.focused {
		case 0:
			m.processDetails.focused = 1
			cmd := m.processDetails.duration.Focus()
			m.processDetails.duration.PromptStyle = focusedStyle
			m.processDetails.duration.TextStyle = focusedStyle

			m.processDetails.startTime.Blur()
			m.processDetails.startTime.PromptStyle = noStyle
			m.processDetails.startTime.TextStyle = noStyle
			return cmd
		case 1:
			m.processDetails.focused = 2
			cmd := m.processDetails.isRecurring.Focus()
			m.processDetails.isRecurring.PromptStyle = focusedStyle
			m.processDetails.isRecurring.TextStyle = focusedStyle

			m.processDetails.duration.Blur()
			m.processDetails.duration.PromptStyle = noStyle
			m.processDetails.duration.TextStyle = noStyle
			return cmd
		case 2:
			startTime, startErr := time.Parse("15:04", m.processDetails.startTime.Value())
			duration, durationErr := time.ParseDuration(m.processDetails.duration.Value())
			startTime = startTime.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)

			if startErr != nil || durationErr != nil {
				m.err = errors.Join(startErr, durationErr)
			} else {
				newProcess := process.Process{Program: m.programList[m.cursor], StartTime: startTime, Duration: duration}
				if m.processDetails.isRecurring.Value() == "Y" {
					newProcess.IsRecurring = true
				}
				scheduleErr := m.scheduler.AddProcess(newProcess)
				if scheduleErr != nil {
					m.err = scheduleErr
				}
			}
			m.page = programs
			m.processDetails.focused = 0
		}
	}
	return nil
}
