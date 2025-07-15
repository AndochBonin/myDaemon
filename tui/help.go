package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) HelpPage() string {
	title := "Help"

	description := textContentStyle.Render("myDaemon (v. 1.1): A process manager. Not a todo app.")

	programs := textContentStyle.Render("Program: A set of whitelisted applications and blacklisted webhosts.")

	processes := textContentStyle.Render("Process: The execution of a Program for a specified duration.\n" +
		"         During this period myDaemon kills all applications not listed in the Program whitelist.")

	schedule := textContentStyle.Render("Schedule: A time sorted sequence of non-overlapping Processes.")

	return pageTitleStyle.Render(title) + "\n\n" + description + "\n\n" + programs + "\n\n" + processes + "\n\n" + schedule
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
