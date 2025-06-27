package tui

import (
	"fmt"

	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/program"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

const (
	schedule = iota
	programs
	help
)

type Model struct {
	page int
	cursor int
	scheduler *process.Scheduler
	programList program.ProgramList
}

var programListFile string = "./program/programList.json"

func initialModel() (Model, error) {
	var m Model
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
			return m, tea.Quit
		case "s":
			if m.page == schedule {
				break
			}
			m.page = schedule
			m.cursor = 0
		case "p":
			if m.page == programs {
				break
			}
			m.page = programs
			m.cursor = 0
		case "h":
			if m.page == help {
				break
			}
			m.page = help
		case "up":
			if m.page == programs {
				m.cursor = max(m.cursor - 1, 0)
			}
			if m.page == schedule {
				m.cursor = max(m.cursor - 1, 0)
			}
		case "down":
			if m.page == programs {
				m.cursor = min(m.cursor + 1, len(m.programList) - 1)
			}
			if m.page == schedule {
				m.cursor = max(m.cursor + 1, len(m.scheduler.Schedule))
			}
		case "d":
			if m.page == programs {
				err := program.DeleteProgram(programListFile, m.cursor)
				if err != nil {
					fmt.Println("\nCould not delete program")
					break
				}
				program.ReadPrograms(programListFile, &m.programList)
				m.cursor = max(m.cursor - 1, 0)
			}
			if m.page == schedule {
				err := m.scheduler.RemoveProcess(m.cursor, true)
				if err != nil {
					fmt.Println("\nCould not delete process")
				}
			}
		}
	}
	return m, nil
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
	}
	return Header() + view + Footer()
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

func (m Model) SchedulePage() string {
	pageTitle := "Schedule\n\n"
	//shows current process (and progress), rest of schedule (remove process)
	processes := ""
	if m.scheduler == nil {
		return pageTitle
	}
	for i, process := range(m.scheduler.Schedule) {
		cursor := " "
		if (m.cursor == i) {
			cursor = ">"
		} 
		processes += cursor + process.Program.Name + " " + process.StartTime.String() + " - " + process.EndTime.String() + "\n"
	}
	return pageTitle + processes
}

func (m Model) ProgramsPage() string {
	pageTitle := "Programs\n"
	pageDescription := "schedule process [enter] / new program [n] / edit program [e] / delete program [d]\n\n"
	// programs: shows list of programs programs (actions: add process, delete program, edit programs, create program)
	programs := ""
	for i, program := range(m.programList) {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		programs += cursor + program.Name + "\n"
	}
	return pageTitle + pageDescription + programs
}

func (m Model) HelpPage() string {
	title := "Help\n\n"
	// help: explains myDaemon, programs, processes, and the scheduler
	description := "myDaemon: A process manager. Not a todo app.\n\n" 

	programs := "Program: A set of whitelisted applications\n\n"

	processes := "Process: The execution of a Program for a specified duration.\n" + 
				 "         During this period myDaemon kills all applications not listed in the Program whitelist.\n\n"

	schedule := "Schedule: A time sorted sequence of non-overlapping Processes\n" 

	return title + description + programs + processes + schedule
}

func Header() string {
    s := "myDaemon\n\n" + "Schedule [s] / Programs [p] / Help [h]\n\n"
	return s
}

func Footer() string {
	s := "\n\npress q or ctrl+c to exit"
	return s
}
