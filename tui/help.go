package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) HelpPage() string {
	title := "Help\n\n"

	description := "myDaemon: A process manager. Not a todo app.\n\n"

	programs := "Program: A set of whitelisted applications\n\n"

	processes := "Process: The execution of a Program for a specified duration.\n" +
		"         During this period myDaemon kills all applications not listed in the Program whitelist.\n\n"

	schedule := "Schedule: A time sorted sequence of non-overlapping Processes\n"

	return title + description + programs + processes + schedule
}

func (m *Model) helpPageKeyHandler(key string) tea.Cmd {
	switch key {
	case "ctrl+c", "q":
		return tea.Quit
	case "ctrl+r":
		m.err = nil
	case "s":
		m.page = schedule
		m.cursor = 0
	case "p":
		m.page = programs
		m.cursor = 0
	}
	return nil
}
