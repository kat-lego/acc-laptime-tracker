package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyUpdate(msg, msg.String())
	case tea.WindowSizeMsg:
		return m.handleWindowSizeUpdate(msg.Width, msg.Height)
	}
	return m, nil
}

func (m Model) handleKeyUpdate(msg tea.Msg, key string) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	cmds = append(cmds, cmd)

	switch key {
	case "ctrl+c", "q":
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleWindowSizeUpdate(width int, height int) (tea.Model, tea.Cmd) {
	m.width = width
	m.height = height

	return m, nil
}
