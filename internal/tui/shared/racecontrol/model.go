package racecontrol

import (
	"bytes"
	"encoding/json"

	"github.com/AidanFogarty/pitwall/internal/f1"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width, height int

	messages []string
}

func New() Model {
	return Model{
		messages: []string{},
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
		lines = append(lines, msg)
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
	}
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
				m.messages = append(m.messages, msg.Message)
			}
		}
	case '{':
		var messagesMap map[string]f1.RaceControlMessage
		if err := json.Unmarshal(raceControlData.Messages, &messagesMap); err != nil {
			return err
		}

		for _, msg := range messagesMap {
			if msg.Message != "" {
				m.messages = append(m.messages, msg.Message)
			}
		}
	}

	return nil
}
