package tui

import (
	"os"
	"path/filepath"

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
	myDaemonStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true)
	navStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Faint(true)
	pageTitleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true)
	focusedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	textContentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))
	errStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("124"))
)

type Model struct {
	programDetails struct {
		programName      textinput.Model
		programWhitelist textinput.Model
		URLWhitelist     textinput.Model
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

// var programListFile string = "./program/programList.json"
var exePath, _ = os.Executable()
var programListFile string = filepath.Join(filepath.Dir(exePath), "program", "programList.json")

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
		switch m.page {
		case schedule:
			keyCmd = m.schedulePageKeyHandler(msg.String())
		case programs:
			keyCmd = m.programsPageKeyHandler(msg.String())
			if keyCmd != nil { // doing this to consume "n" / "e" presses to avoid them showing up in the text input
				return m, keyCmd
			}
		case help:
			keyCmd = m.helpPageKeyHandler(msg.String())
		case addProgram, editProgram:
			keyCmd = m.programDetailsPageKeyHandler(msg.String())
		case addProcess:
			keyCmd = m.processDetailsPageKeyHandler(msg.String())
		}
	}
	var nameUpdate, appWhitelistUpdate, urlWhitelistUpdate, startTimeUpdate, durationUpdate, isRecurringUpdate tea.Cmd
	m.programDetails.programName, nameUpdate = m.programDetails.programName.Update(msg)
	m.programDetails.programWhitelist, appWhitelistUpdate = m.programDetails.programWhitelist.Update(msg)
	m.programDetails.URLWhitelist, urlWhitelistUpdate = m.programDetails.URLWhitelist.Update(msg)

	m.processDetails.startTime, startTimeUpdate = m.processDetails.startTime.Update(msg)
	m.processDetails.duration, durationUpdate = m.processDetails.duration.Update(msg)
	m.processDetails.isRecurring, isRecurringUpdate = m.processDetails.isRecurring.Update(msg)
	cmd := tea.Batch(nameUpdate, appWhitelistUpdate, urlWhitelistUpdate, startTimeUpdate, durationUpdate, isRecurringUpdate, keyCmd)
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
	return Title() + NavHeader(m.page) + view + ErrMessage(m.err)
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

func Title() string {
	return myDaemonStyle.Render("myDaemon") + "\n\n"
}

func NavHeader(page int) string {
	mainNav := navStyle.Render("schedule [s]   programs [p]   help [h]   quit [q / ctrl+c]")
	backNav := navStyle.Render("back [esc]   quit [ctrl+c]")
	if page == schedule || page == programs || page == help {
		return mainNav + "\n\n"
	}
	return backNav + "\n\n"
}

func ErrMessage(err error) string {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}
	return "\n" + errStyle.Render(errMessage)
}
