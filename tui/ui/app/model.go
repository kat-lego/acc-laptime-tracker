package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type Model struct {
	term      string
	profile   string
	width     int
	height    int
	bg        string
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
}

func New(renderer *lipgloss.Renderer, pty *ssh.Pty) Model {

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	return Model{
		term:      pty.Term,
		profile:   renderer.ColorProfile().Name(),
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		bg:        bg,
		txtStyle:  txtStyle,
		quitStyle: quitStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
