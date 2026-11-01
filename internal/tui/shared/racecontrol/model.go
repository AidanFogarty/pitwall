package racecontrol

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/AidanFogarty/pitwall/internal/f1"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width, height int

	messages []f1.RaceControlMessage

	teamColours map[string]string

	cursor  int
	top     int
	tailing bool
}

func New() Model {
	return Model{
		messages:    []f1.RaceControlMessage{},
		teamColours: map[string]string{},
		tailing:     true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case f1.F1EventMsg:
		event := f1.F1Event(msg)
		m.handleEvent(event)
	}
	return m, nil
}

func (m Model) View() string {
	displayMessages := m.messages
	if len(displayMessages) > 7 {
		displayMessages = displayMessages[len(displayMessages)-7:]
	}

	var lines []string
	for _, msg := range displayMessages {
		lines = append(lines, msg.Message)
	}

	for len(lines) < 7 {
		lines = append(lines, "")
	}

	var styledLines []string
	lineStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(m.width - 2)

	for _, line := range lines {
		if line != "" {
			line = "> " + line
		}
		styledLines = append(styledLines, lineStyle.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, styledLines...)

	headerStyle := lipgloss.NewStyle().
		PaddingLeft(1)

	header := headerStyle.Render("Race Control Messages")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder())

	box := boxStyle.Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, header, box)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height

	return m
}

func (m *Model) handleEvent(event f1.F1Event) {
	switch event.Type {
	case "RaceControlMessages":
		m.handleRaceControlMessages(event)
	case "DriverList":
		m.handleDriverList(event)
	}
}

func (m *Model) handleDriverList(event f1.F1Event) error {
	var driverList f1.DriverList
	if err := json.Unmarshal(event.Data, &driverList); err != nil {
		return err
	}

	for racingNumber, driver := range driverList {
		if driver.TeamColour != "" {
			m.teamColours[racingNumber] = driver.TeamColour
		}
	}
	return nil
}

func (m *Model) handleRaceControlMessages(event f1.F1Event) error {
	var raceControlData f1.RaceControlMessages
	if err := json.Unmarshal(event.Data, &raceControlData); err != nil {
		return err
	}

	trimmed := bytes.TrimSpace(raceControlData.Messages)
	if len(trimmed) == 0 {
		return nil
	}

	switch trimmed[0] {
	case '[':
		var messagesArray []f1.RaceControlMessage
		if err := json.Unmarshal(raceControlData.Messages, &messagesArray); err != nil {
			return err
		}

		for _, msg := range messagesArray {
			if msg.Message != "" {
				m.messages = append(m.messages, msg)
			}
		}
	case '{':
		var messagesMap map[string]f1.RaceControlMessage
		if err := json.Unmarshal(raceControlData.Messages, &messagesMap); err != nil {
			return err
		}

		keys := make([]string, 0, len(messagesMap))
		for k := range messagesMap {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			a, errA := strconv.Atoi(keys[i])
			b, errB := strconv.Atoi(keys[j])
			if errA != nil || errB != nil {
				return keys[i] < keys[j]
			}
			return a < b
		})

		for _, k := range keys {
			if msg := messagesMap[k]; msg.Message != "" {
				m.messages = append(m.messages, msg)
			}
		}
	}

	return nil
}
