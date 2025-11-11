package dashboard

import (
	"github.com/AidanFogarty/pitwall/internal/tui/shared/debug"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/information"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/racecontrol"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/timingtable"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Quit   key.Binding
	Timing key.Binding
	Debug  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Timing, k.Debug}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Timing, k.Debug},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Timing: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "timing"),
	),
	Debug: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "debug"),
	),
}

type DashboardModel struct {
	width, height int

	timingTable timingtable.Model
	racecontrol racecontrol.Model
	information information.Model
	debug       debug.Model

	keys keyMap
	help help.Model

	currentView string
}

func NewModel() tea.Model {
	model := &DashboardModel{
		timingTable: timingtable.New(),
		racecontrol: racecontrol.New(),
		information: information.New(),

		keys: keys,
		help: help.New(),

		currentView: "timing",
	}
	return model
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.timingTable.Init(),
		m.racecontrol.Init(),
		m.information.Init(),
		m.debug.Init(),
	)
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		leftWidth := int(float64(msg.Width) * 0.8)
		rightWidth := msg.Width - leftWidth

		raceControlHeight := 10 // Approx height for race control section
		timingTableHeight := msg.Height - raceControlHeight

		m.timingTable = m.timingTable.SetSize(leftWidth, timingTableHeight)
		m.racecontrol = m.racecontrol.SetSize(leftWidth, raceControlHeight)
		m.information = m.information.SetSize(rightWidth, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Timing):
			m.currentView = "timing"
			return m, nil
		case key.Matches(msg, m.keys.Debug):
			m.currentView = "debug"
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.timingTable, cmd = m.timingTable.Update(msg)
	cmds = append(cmds, cmd)

	m.racecontrol, cmd = m.racecontrol.Update(msg)
	cmds = append(cmds, cmd)

	m.information, cmd = m.information.Update(msg)
	cmds = append(cmds, cmd)

	m.debug, cmd = m.debug.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DashboardModel) View() string {
	helpView := m.help.View(m.keys)
	helpView = lipgloss.NewStyle().PaddingLeft(1).Render(helpView)

	switch m.currentView {
	case "timing":
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.timingView(),
			helpView,
		)
	case "debug":
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.debug.View(),
			helpView,
		)
	default:
		return ""
	}
}

func (m DashboardModel) timingView() string {
	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		m.timingTable.View(),
		m.racecontrol.View(),
	)

	rightSide := m.information.View()

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftSide,
		rightSide,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
	)
}
