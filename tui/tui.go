package tui

import (
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
		startTime   textinput.Model
		duration    textinput.Model
		isRecurring textinput.Model
		focused     int
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
	var keyCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// new implementation:
		// switch page:
		// pass msg.String() to [page]KeyHandler
		switch m.page {
		case schedule:
			keyCmd = m.schedulePageKeyHandler(msg.String())
		case programs:
			keyCmd = m.programsPageKeyHandler(msg.String())
		case help:
			keyCmd = m.helpPageKeyHandler(msg.String())
		case addProgram, editProgram:
			keyCmd = m.programDetailsPageKeyHandler(msg.String())
		case addProcess:
			keyCmd = m.processDetailsPageKeyHandler(msg.String())
		}
	}
	var cmd1, cmd2, cmd3, cmd4, cmd5 tea.Cmd
	m.programDetails.programName, cmd1 = m.programDetails.programName.Update(msg)
	m.programDetails.programWhitelist, cmd2 = m.programDetails.programWhitelist.Update(msg)

	m.processDetails.startTime, cmd3 = m.processDetails.startTime.Update(msg)
	m.processDetails.duration, cmd4 = m.processDetails.duration.Update(msg)
	m.processDetails.isRecurring, cmd5 = m.processDetails.isRecurring.Update(msg)
	cmd := tea.Batch(cmd1, cmd2, cmd3, cmd4, cmd5, keyCmd)
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
