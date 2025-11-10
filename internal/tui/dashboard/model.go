package dashboard

import (
	"github.com/AidanFogarty/pitwall/internal/tui/shared/information"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/racecontrol"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/timingtable"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type DashboardModel struct {
	width, height int

	timingTable timingtable.Model
	racecontrol racecontrol.Model
	information information.Model

	keys keyMap
	help help.Model
}

func NewModel() tea.Model {
	model := &DashboardModel{
		timingTable: timingtable.New(),
		racecontrol: racecontrol.New(),
		information: information.New(),

		keys: keys,
		help: help.New(),
	}
	return model
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.timingTable.Init(),
		m.racecontrol.Init(),
		m.information.Init(),
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

	return m, tea.Batch(cmds...)
}

func (m DashboardModel) View() string {
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

	helpView := m.help.View(m.keys)
	helpView = lipgloss.NewStyle().PaddingLeft(1).Render(helpView)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		helpView,
	)
}
