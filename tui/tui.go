package tui

import (
	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/program"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

type Model struct {
	page string
	scheduler *process.Scheduler
	programList program.ProgramList
}

var programListFile string = "./program/programList.json"

func initialModel() (Model, error) {
	var m Model
	m.scheduler = process.GetScheduler()
	m.page = SchedulePage(m.scheduler)
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
			m.page = SchedulePage(m.scheduler)
		case "p":
			m.page = ProgramsPage(m.programList)
		case "h":
			m.page = HelpPage()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return Header() + m.page + Footer()
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

func SchedulePage(scheduler *process.Scheduler) string {
	s := "Schedule\n\n"
	//shows current process (and progress), rest of schedule (remove process)
	processes := ""
	if scheduler == nil {
		return s
	}
	for _, process := range(scheduler.Schedule) {
		processes += process.Program.Name + "\n"
	}
	return s + processes
}

func ProgramsPage(programList program.ProgramList) string {
	s := "Programs\n\n"
	// programs: shows list of programs programs (actions: add process, delete program, edit programs, create program)
	programs := ""
	for _, program := range(programList) {
		programs += program.Name + "\n"
	}
	return s + programs
}

func HelpPage() string {
	title := "Help\n\n"
	// help: explains myDaemon, programs, processes, and the scheduler
	description := "myDaemon is a process manager. Not a todo app.\n\n" 

	programs := "Program: a set of whitelisted applications\n\n"

	processes := "Process: the execution of a programs for a specified duration\n" + 
				 "during this period myDaemon kills all applications not listed in the program whitelist.\n\n"

	schedule := "Schedule: processes can be scheduled and can be recurring. The scheduler manages this." 

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
