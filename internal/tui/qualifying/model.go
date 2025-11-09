package qualifying

import (
	"github.com/AidanFogarty/pitwall/internal/tui/shared/activity"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/information"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/qualifyingtable"
	"github.com/AidanFogarty/pitwall/internal/tui/shared/racecontrol"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QualifyingModel struct {
	width, height int

	qualifyingTable qualifyingtable.Model
	racecontrol     racecontrol.Model
	information     information.Model
	activity        activity.Model

	event int
}

func NewModel() tea.Model {
	model := &QualifyingModel{
		qualifyingTable: qualifyingtable.New(),
		racecontrol:     racecontrol.New(),
		information:     information.New(),
		activity:        activity.New(),
	}
	return model
}

func (m QualifyingModel) Init() tea.Cmd {
	return tea.Batch(
		m.qualifyingTable.Init(),
		m.racecontrol.Init(),
		m.information.Init(),
		m.activity.Init(),
	)
}

func (m QualifyingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		leftWidth := int(float64(msg.Width) * 0.8)
		rightWidth := msg.Width - leftWidth

		raceControlHeight := 10 // Approx height for race control section
		timingTableHeight := msg.Height - raceControlHeight

		m.qualifyingTable = m.qualifyingTable.SetSize(leftWidth, timingTableHeight)
		m.racecontrol = m.racecontrol.SetSize(leftWidth, raceControlHeight)
		m.information = m.information.SetSize(rightWidth, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.qualifyingTable, cmd = m.qualifyingTable.Update(msg)
	cmds = append(cmds, cmd)

	m.racecontrol, cmd = m.racecontrol.Update(msg)
	cmds = append(cmds, cmd)

	m.information, cmd = m.information.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m QualifyingModel) View() string {
	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		m.qualifyingTable.View(),
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
