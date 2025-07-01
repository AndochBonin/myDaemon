package tui

import (
	"errors"
	"strings"
	"time"

	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/program"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	schedule = iota
	programs
	help
	addProcess
	addProgram
	editProgram
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
	cursorStyle  = focusedStyle
)

type Model struct {
	programDetails struct {
		programName      textinput.Model
		programWhitelist textinput.Model
		focused          int
	}
	processDetails struct {
		startTime textinput.Model
		duration  textinput.Model
		isRecurring textinput.Model
		focused   int
	}
	page        int
	cursor      int
	err         error
	scheduler   *process.Scheduler
	programList program.ProgramList
}

var programListFile string = "./program/programList.json"

func initialModel() (Model, error) {
	var m Model
	m.programDetails.focused = 0
	m.page = schedule
	m.cursor = 0
	m.scheduler = process.GetScheduler()
	m.scheduler = process.GetScheduler()
	err := program.ReadPrograms(programListFile, &m.programList)

	return m, err
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("myDaemon")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.page == editProgram || m.page == addProgram || m.page == addProcess {
				break
			}
			return m, tea.Quit
		case "ctrl+r":
			m.err = nil
			return m, nil
		case "s":
			if m.page == schedule || m.page == editProgram || m.page == addProgram || m.page == addProcess {
				break
			}
			m.page = schedule
			m.cursor = 0
		case "p":
			if m.page == programs || m.page == editProgram || m.page == addProgram || m.page == addProcess {
				break
			}
			m.page = programs
			m.cursor = 0
		case "h":
			if m.page == help || m.page == editProgram || m.page == addProgram || m.page == addProcess {
				break
			}
			m.page = help
		case "up":
			if m.page == programs {
				m.cursor = max(m.cursor-1, 0)
			}
			if m.page == schedule {
				m.cursor = max(m.cursor-1, 0)
			}
		case "down":
			if m.page == programs {
				m.cursor = min(m.cursor+1, len(m.programList)-1)
			}
			if m.page == schedule {
				m.cursor = min(m.cursor+1, len(m.scheduler.Schedule)-1)
			}
		case "d":
			if m.page == programs {
				if len(m.programList) == 0 {
					return m, nil
				}
				err := program.DeleteProgram(programListFile, m.cursor)
				if err != nil {
					m.err = err
					break
				}
				program.ReadPrograms(programListFile, &m.programList)
				m.cursor = max(m.cursor-1, 0)
			}
			if m.page == schedule {
				err := m.scheduler.RemoveProcess(m.cursor, true)
				if err != nil {
					m.err = err
				}
			}
		case "n":
			if m.page == programs {
				m.page = addProgram
				cmd := m.initProgramDetailsInput(program.Program{})
				return m, cmd
			}
		case "e":
			if m.page == programs {
				if len(m.programList) == 0 {
					return m, nil
				}
				m.page = editProgram
				cmd := m.initProgramDetailsInput(m.programList[m.cursor])
				return m, cmd
			}
		case "esc":
			if m.page == editProgram || m.page == addProgram || m.page == addProcess {
				m.page = programs
			}
		case "enter":
			switch m.page {
			case editProgram, addProgram:
				switch m.programDetails.focused {
				case 0:
					m.programDetails.focused = 1
					cmd := m.programDetails.programWhitelist.Focus()
					m.programDetails.programWhitelist.PromptStyle = focusedStyle
					m.programDetails.programWhitelist.TextStyle = focusedStyle

					m.programDetails.programName.Blur()
					m.programDetails.programName.PromptStyle = noStyle
					m.programDetails.programName.TextStyle = noStyle
					return m, cmd
				case 1:
					name := m.programDetails.programName.Value()
					whitelist := strings.Split(m.programDetails.programWhitelist.Value(), ",")
					for i, uri := range whitelist {
						whitelist[i] = strings.Trim(uri, " ")
					}
					newProgram := program.Program{Name: name, URIWhitelist: whitelist}
					var err error
					switch m.page {
					case editProgram:
						err = program.UpdateProgram(programListFile, m.cursor, newProgram)
					case addProgram:
						err = program.CreateProgram(programListFile, newProgram)
					}
					if err != nil {
						m.err = err
					}
					program.ReadPrograms(programListFile, &m.programList)
					m.page = programs
					m.programDetails.focused = 0
				}
			case programs:
				if len(m.programList) == 0 {
					break
				}
				m.page = addProcess
				cmd := m.initProcessDetailsInput()
				return m, cmd
			case addProcess:
				switch m.processDetails.focused {
				case 0:
					m.processDetails.focused = 1
					cmd := m.processDetails.duration.Focus()
					m.processDetails.duration.PromptStyle = focusedStyle
					m.processDetails.duration.TextStyle = focusedStyle

					m.processDetails.startTime.Blur()
					m.processDetails.startTime.PromptStyle = noStyle
					m.processDetails.startTime.TextStyle = noStyle
					return m, cmd
				case 1:
					m.processDetails.focused = 2
					cmd := m.processDetails.isRecurring.Focus()
					m.processDetails.isRecurring.PromptStyle = focusedStyle
					m.processDetails.isRecurring.TextStyle = focusedStyle

					m.processDetails.duration.Blur()
					m.processDetails.duration.PromptStyle = noStyle
					m.processDetails.duration.TextStyle = noStyle
					return m, cmd
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

		}
	}
	var cmd1, cmd2, cmd3, cmd4, cmd5 tea.Cmd
	m.programDetails.programName, cmd1 = m.programDetails.programName.Update(msg)
	m.programDetails.programWhitelist, cmd2 = m.programDetails.programWhitelist.Update(msg)

	m.processDetails.startTime, cmd3 = m.processDetails.startTime.Update(msg)
	m.processDetails.duration, cmd4 = m.processDetails.duration.Update(msg)
	m.processDetails.isRecurring, cmd5 = m.processDetails.isRecurring.Update(msg)
	cmd := tea.Batch(cmd1, cmd2, cmd3, cmd4, cmd5)
	return m, cmd
}

func (m Model) View() string {
	view := ""
	switch m.page {
	case schedule:
		view = m.SchedulePage()
	case programs:
		view = m.ProgramsPage()
	case help:
		view = m.HelpPage()
	case addProgram, editProgram:
		view = m.ProgramDetailsPage()
	case addProcess:
		view = m.ProcessDetailsPage()
	}
	return Header() + view + Footer(m.err)
}

func Run() error {
	m, initErr := initialModel()
	if initErr != nil {
		return initErr
	}
	p := tea.NewProgram(m)
	_, runErr := p.Run()
	return runErr
}

func Header() string {
	s := "myDaemon\n\n" + "Schedule [s] / Programs [p] / Help [h]\n\n"
	return s
}

func Footer(err error) string {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}
	s := "\n\npress q or ctrl+c to exit"
	return "\n" + errMessage + s
}

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

func (m *Model) ProgramsPage() string {
	pageTitle := "Programs\n"
	pageDescription := "schedule process [enter] / new program [n] / edit program [e] / delete program [d]\n\n"

	programs := ""
	for i, program := range m.programList {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		var whitelist string
		for j, uri := range program.URIWhitelist {
			whitelist += uri
			if j < len(program.URIWhitelist)-1 {
				whitelist += ", "
			}
		}
		programs += cursor + program.Name + ": " + whitelist + "\n"
	}
	return pageTitle + pageDescription + programs
}

func (m *Model) HelpPage() string {
	title := "Help\n\n"

	description := "myDaemon: A process manager. Not a todo app.\n\n"

	programs := "Program: A set of whitelisted applications\n\n"

	processes := "Process: The execution of a Program for a specified duration.\n" +
		"         During this period myDaemon kills all applications not listed in the Program whitelist.\n\n"

	schedule := "Schedule: A time sorted sequence of non-overlapping Processes\n"

	return title + description + programs + processes + schedule
}

func (m *Model) ProgramDetailsPage() string {
	pageTitle := "Program Details\n\n"

	return pageTitle + "\nProgram Name:\n" + m.programDetails.programName.View() +
		"\n\nProgram Whitelist:\n" + m.programDetails.programWhitelist.View()
}

func (m *Model) initProgramDetailsInput(program program.Program) tea.Cmd {
	programName := textinput.New()
	programName.Placeholder = program.Name
	programName.Cursor.Style = cursorStyle
	programName.PromptStyle = focusedStyle
	programName.TextStyle = focusedStyle
	programName.CharLimit = 156
	programName.Width = 20
	m.programDetails.programName = programName

	programWhitelist := textinput.New()
	programWhitelist.Placeholder = ""
	for i, uri := range program.URIWhitelist {
		programWhitelist.Placeholder += uri
		if i < len(program.URIWhitelist)-1 {
			programWhitelist.Placeholder += ", "
		}
	}
	programWhitelist.Cursor.Style = cursorStyle
	programWhitelist.CharLimit = 156
	programWhitelist.Width = 20
	m.programDetails.programWhitelist = programWhitelist

	return m.programDetails.programName.Focus()
}

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
