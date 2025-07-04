package tui

import (
	"errors"
	"time"

	"github.com/AndochBonin/myDaemon/process"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) ProcessDetailsPage() string {
	pageTitle := "Process Details"

	return pageTitleStyle.Render(pageTitle + ": " + m.programList[m.cursor].Name) + "\n\n" + 
		   focusedStyle.Render("Start Time:") + "\n" + m.processDetails.startTime.View() + "\n\n" +
		   focusedStyle.Render("Duration:") + "\n" + m.processDetails.duration.View() + "\n\n" +
		   focusedStyle.Render("Is Recurring?:") + "\n" + m.processDetails.isRecurring.View()
}

func (m *Model) initProcessDetailsInput() tea.Cmd {
	startTime := textinput.New()
	startTime.Placeholder = "00:00"
	startTime.Cursor.Style = textContentStyle
	startTime.PromptStyle = textContentStyle
	startTime.TextStyle = textContentStyle
	startTime.CharLimit = 10
	startTime.Width = 10
	m.processDetails.startTime = startTime

	duration := textinput.New()
	duration.Placeholder = "0h0m"
	duration.Cursor.Style = textContentStyle
	duration.PromptStyle = textContentStyle
	duration.TextStyle = textContentStyle
	duration.CharLimit = 10
	duration.Width = 10
	m.processDetails.duration = duration

	isRecurring := textinput.New()
	isRecurring.Placeholder = "N"
	isRecurring.Cursor.Style = textContentStyle
	isRecurring.PromptStyle = textContentStyle
	isRecurring.TextStyle = textContentStyle
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
			return cmd
		case 1:
			m.processDetails.focused = 2
			cmd := m.processDetails.isRecurring.Focus()
			m.processDetails.isRecurring.PromptStyle = focusedStyle
			m.processDetails.isRecurring.TextStyle = focusedStyle

			m.processDetails.duration.Blur()
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
