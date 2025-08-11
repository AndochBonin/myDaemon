package tui

import (
	"os"
	"path/filepath"
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
		webHostBlocklist     textinput.Model
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

//var programListFile string = "./storage/programList.json" //for testing
var exePath, _ = os.Executable()
var programListFile string = filepath.Join(filepath.Dir(exePath), "storage", "programList.json")
var scheduleFile string = filepath.Join(filepath.Dir(exePath), "storage", "schedule.json")

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
	
	return tea.Batch(tea.SetWindowTitle("myDaemon"), tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var keyCmd tea.Cmd
	var timeCmd tea.Cmd
	switch msg := msg.(type) {
	case time.Time:
		timeCmd = tick() 
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
	for _, textModel := range [](*textinput.Model){
		&m.programDetails.programName,
		&m.programDetails.programWhitelist,
		&m.programDetails.webHostBlocklist,
		&m.processDetails.startTime,
		&m.processDetails.duration,
		&m.processDetails.isRecurring,
	} {
		var cmd tea.Cmd
		*textModel, cmd = textModel.Update(msg)
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, keyCmd, timeCmd)
	cmd := tea.Batch(cmds...)
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

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}
