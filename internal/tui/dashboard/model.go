package dashboard

import (
	"github.com/AidanFogarty/pitwall/internal/tui/shared"
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
	Quit        key.Binding
	Timing      key.Binding
	RaceControl key.Binding
	Debug       key.Binding

	ScrollUp   key.Binding
	ScrollDown key.Binding
	Oldest     key.Binding
	Newest     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Timing, k.RaceControl, k.Debug}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Timing, k.RaceControl, k.Debug},
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
	RaceControl: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "race control"),
	),
	Debug: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "debug"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
	Oldest: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "oldest"),
	),
	Newest: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "newest"),
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

		m.debug = m.debug.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Timing):
			m.currentView = "timing"
			return m, nil
		case key.Matches(msg, m.keys.RaceControl):
			m.currentView = "racecontrol"
			return m, nil
		case key.Matches(msg, m.keys.Debug):
			m.currentView = "debug"
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case m.currentView == "racecontrol":
			switch {
			case key.Matches(msg, m.keys.ScrollUp):
				m.racecontrol = m.racecontrol.ScrollUp(1, m.contentHeight())
			case key.Matches(msg, m.keys.ScrollDown):
				m.racecontrol = m.racecontrol.ScrollDown(1, m.contentHeight())
			case key.Matches(msg, m.keys.Oldest):
				m.racecontrol = m.racecontrol.ScrollToStart(m.contentHeight())
			case key.Matches(msg, m.keys.Newest):
				m.racecontrol = m.racecontrol.ScrollToEnd(m.contentHeight())
			}
			return m, nil
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
	hint := ""
	if m.currentView == "racecontrol" {
		hint = m.racecontrol.ScrollHint()
	}
	helpView := shared.ShortcutBar(m.help.View(m.keys), hint, m.width)

	switch m.currentView {
	case "timing":
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.timingView(),
			helpView,
		)
	case "racecontrol":
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.racecontrol.FullView(m.width, m.contentHeight()),
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

	rightSide := m.information.SetHeight(lipgloss.Height(leftSide)).View()

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

func (m DashboardModel) contentHeight() int {
	return lipgloss.Height(m.timingView())
}
