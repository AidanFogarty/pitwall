package debug

import (
	"encoding/json"
	"fmt"

	"github.com/AidanFogarty/pitwall/internal/f1"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxEvents    = 20
	maxDataChars = 100
)

type Model struct {
	width, height int

	name        string
	sessionType string
	location    string

	meetingKey int
	sessionKey int

	event  int
	events []f1.F1Event
}

func New() Model {
	return Model{}
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
	baseStyle := lipgloss.NewStyle().
		PaddingLeft(1)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1)

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Width(14)

	locationBaseStyle := lipgloss.NewStyle().
		PaddingLeft(6)

	title := titleStyle.Render("Debug Mode\n")

	sessionInfo := lipgloss.JoinVertical(lipgloss.Left,
		baseStyle.Render(fmt.Sprintf("%s %d", keyStyle.Render("Events:"), m.event)),
		baseStyle.Render(fmt.Sprintf("%s %d", keyStyle.Render("Meeting Key:"), m.meetingKey)),
		baseStyle.Render(fmt.Sprintf("%s %d", keyStyle.Render("Session Key:"), m.sessionKey)),
	)

	locationInfo := lipgloss.JoinVertical(lipgloss.Left,
		locationBaseStyle.Render(fmt.Sprintf("%s %s", keyStyle.Render("Meeting Name:"), m.name)),
		locationBaseStyle.Render(fmt.Sprintf("%s %s", keyStyle.Render("Session Type:"), m.sessionType)),
		locationBaseStyle.Render(fmt.Sprintf("%s %s", keyStyle.Render("Location:"), m.location)),
	)

	stats := lipgloss.JoinHorizontal(lipgloss.Left, sessionInfo, locationInfo)

	var eventLines []string
	for _, event := range m.events {
		dataStr := string(event.Data)
		if len(dataStr) > maxDataChars {
			dataStr = dataStr[:maxDataChars] + "..."
		}

		eventLine := fmt.Sprintf("%s | %s | %s",
			keyStyle.Render(event.Type),
			event.Timestamp.Format("15:04:05.000"),
			dataStr,
		)
		eventLines = append(eventLines, baseStyle.Render(eventLine))
	}

	eventContent := lipgloss.JoinVertical(lipgloss.Left, eventLines...)

	eventLog := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width - 4)

	debugInfo := []string{
		title,
		stats,
		eventLog.Render(eventContent),
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		debugInfo...,
	)
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func (m *Model) handleEvent(event f1.F1Event) {
	m.event++

	m.events = append(m.events, event)
	if len(m.events) > maxEvents {
		m.events = m.events[len(m.events)-maxEvents:]
	}

	switch event.Type {
	case "SessionInfo":
		m.handleSessionInfo(event)
	}
}

func (m *Model) handleSessionInfo(event f1.F1Event) error {
	var sessionInfoData f1.SessionInfo
	if err := json.Unmarshal(event.Data, &sessionInfoData); err != nil {
		return err
	}

	meeting := sessionInfoData.Meeting
	if meeting.OfficialName != "" {
		m.name = meeting.OfficialName
	}

	if sessionInfoData.Name != "" {
		m.sessionType = sessionInfoData.Name
	}

	if meeting.Location != "" && meeting.Country.Name != "" {
		m.location = fmt.Sprintf("%s, %s", meeting.Location, meeting.Country.Name)
	}

	if meeting.Key != 0 {
		m.meetingKey = meeting.Key
	}

	if sessionInfoData.Key != 0 {
		m.sessionKey = sessionInfoData.Key
	}

	return nil
}
