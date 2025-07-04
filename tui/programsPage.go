package tui

import (
	"slices"
	"strings"

	"github.com/AndochBonin/myDaemon/program"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) ProgramsPage() string {
	pageTitle := "Programs"
	pageKeys := "schedule [enter]   new [n]   edit [e]   delete [d]"
	programs := ""

	for i, program := range m.programList {
		var style lipgloss.Style = textContentStyle
		cursor := " "
		if i == m.cursor {
			cursor = "> "
			style = focusedStyle.Bold(true)
		}
		var whitelist string
		for j, uri := range program.URIWhitelist {
			whitelist += uri
			if j < len(program.URIWhitelist)-1 {
				whitelist += ", "
			}
		}
		programs += style.Render(cursor + program.Name + ": " + whitelist) + "\n"
	}
	return pageTitleStyle.Render(pageTitle)  + "\n\n" + navStyle.Render(pageKeys) + "\n\n" + programs
}

func (m *Model) ProgramDetailsPage() string {
	pageTitle := "Program Details"

	return pageTitleStyle.Render(pageTitle) + "\n\n" + 
		   focusedStyle.Render("Program Name: ") + "\n" + m.programDetails.programName.View() + "\n\n" + 
		   focusedStyle.Render("Program Whitelist: ") + "\n" + m.programDetails.programWhitelist.View()
}

func (m *Model) initProgramDetailsInput(program program.Program) tea.Cmd {
	programName := textinput.New()
	programName.Placeholder = program.Name
	programName.Cursor.Style = textContentStyle
	programName.PromptStyle = textContentStyle
	programName.TextStyle = textContentStyle
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
	programWhitelist.Cursor.Style = textContentStyle
	programWhitelist.PromptStyle = textContentStyle
	programWhitelist.TextStyle = textContentStyle
	programWhitelist.CharLimit = 156
	programWhitelist.Width = 20
	m.programDetails.programWhitelist = programWhitelist

	return m.programDetails.programName.Focus()
}

func (m *Model) programsPageKeyHandler(key string) tea.Cmd {
	switch key {
	case "q", "ctrl+c":
		return tea.Quit
	case "ctrl+r":
		m.err = nil
	case "s":
		m.page = schedule
		m.cursor = 0
	case "h":
		m.page = help
	case "up":
		m.cursor = max(m.cursor-1, 0)
	case "down":
		m.cursor = min(m.cursor+1, len(m.programList)-1)
	case "d":
		if len(m.programList) == 0 {
			return nil
		}
		err := program.DeleteProgram(programListFile, m.cursor)
		if err != nil {
			m.err = err
			break
		}
		program.ReadPrograms(programListFile, &m.programList)
		m.cursor = max(m.cursor-1, 0)
	case "n":
		m.page = addProgram
		cmd := m.initProgramDetailsInput(program.Program{})
		return cmd
	case "e":
		if len(m.programList) == 0 {
			return nil
		}
		m.page = editProgram
		cmd := m.initProgramDetailsInput(m.programList[m.cursor])
		return cmd
	case "enter":
		if len(m.programList) == 0 {
			break
		}
		m.page = addProcess
		cmd := m.initProcessDetailsInput()
		return cmd
	}
	return nil
}

func (m *Model) programDetailsPageKeyHandler(key string) tea.Cmd {
	switch key {
	case "ctrl+c":
		return tea.Quit
	case "esc":
		m.page = programs
	case "enter":
		switch m.programDetails.focused {
		case 0:
			m.programDetails.focused = 1
			cmd := m.programDetails.programWhitelist.Focus()
			m.programDetails.programWhitelist.PromptStyle = focusedStyle
			m.programDetails.programWhitelist.TextStyle = focusedStyle

			m.programDetails.programName.Blur()
			return cmd
		case 1:
			name := m.programDetails.programName.Value()
			whitelist := strings.Split(m.programDetails.programWhitelist.Value(), ",")
			for i := 0; i < len(whitelist); i++ {
				whitelist[i] = strings.Trim(whitelist[i], " ")
				if whitelist[i] == "" {
					whitelist = slices.Delete(whitelist, i, i+1)
					i--
				}
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
	}
	return nil
}
