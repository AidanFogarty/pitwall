package information

import (
	"encoding/json"
	"fmt"

	"github.com/AidanFogarty/pitwall/internal/f1"
	"github.com/AidanFogarty/pitwall/internal/tui/shared"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width, height int

	officialName string
	sessionType  string
	location     string

	airTemperture   string
	trackTemperture string
	humidity        string
	rainfall        string
	windSpeed       string
	windDirection   string

	trackStatus string
	currentLap  int
	totalLaps   int
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
	logo := shared.Logo()
	logoStyle := lipgloss.NewStyle().
		Width(m.width - 4).
		Align(lipgloss.Left)
	styledLogo := logoStyle.Render(logo)

	style := lipgloss.NewStyle()

	titleStyle := lipgloss.NewStyle().Bold(true)

	officialName := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Official Name"),
		style.Render(m.officialName),
	)

	sessionType := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Session Type"),
		style.Render(m.sessionType),
	)

	location := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Location"),
		style.Render(m.location),
	)

	airTemp := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Air Temperature"),
		style.Render(m.airTemperture),
	)

	trackTemp := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Track Temperature"),
		style.Render(m.trackTemperture),
	)

	humidity := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Humidity"),
		style.Render(m.humidity),
	)

	rainfall := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Rainfall"),
		style.Render(m.rainfall),
	)

	windSpeed := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Wind Speed"),
		style.Render(m.windSpeed),
	)

	windDirection := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Wind Direction"),
		style.Render(m.windDirection),
	)

	trackStatus := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Track Status"),
		style.Render(m.trackStatus),
	)

	currentLap := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Current Lap"),
		style.Render(fmt.Sprintf("%d / %d", m.currentLap, m.totalLaps)),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styledLogo,
		officialName,
		sessionType,
		location,
		airTemp,
		trackTemp,
		humidity,
		rainfall,
		windSpeed,
		windDirection,
		trackStatus,
		currentLap,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(m.width - 2).
		Height(30).
		PaddingLeft(1).
		Render(content)

	return box
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func (m *Model) handleEvent(event f1.F1Event) {
	switch event.Type {
	case "SessionInfo":
		m.handleSessionInfo(event)
	case "WeatherData":
		m.handleWeatherData(event)
	case "TrackStatus":
		m.handleTrackStatus(event)
	case "LapCount":
		m.handleLapCount(event)
	}
}

func (m *Model) handleSessionInfo(event f1.F1Event) error {
	var sessionInfoData f1.SessionInfo
	if err := json.Unmarshal(event.Data, &sessionInfoData); err != nil {
		return err
	}

	meeting := sessionInfoData.Meeting
	if meeting.OfficialName != "" {
		m.officialName = meeting.OfficialName
	}

	if sessionInfoData.Name != "" {
		m.sessionType = sessionInfoData.Name
	}

	if meeting.Location != "" && meeting.Country.Name != "" {
		m.location = fmt.Sprintf("%s, %s", meeting.Location, meeting.Country.Name)
	}

	return nil
}

func (m *Model) handleWeatherData(event f1.F1Event) error {
	var weatherData f1.WeatherData
	if err := json.Unmarshal(event.Data, &weatherData); err != nil {
		return err
	}

	if weatherData.AirTemp != "" {
		m.airTemperture = fmt.Sprintf("%s°C", weatherData.AirTemp)
	}

	if weatherData.TrackTemp != "" {
		m.trackTemperture = fmt.Sprintf("%s°C", weatherData.TrackTemp)
	}

	if weatherData.Humidity != "" {
		m.humidity = fmt.Sprintf("%s%%", weatherData.Humidity)
	}

	if weatherData.Rainfall != "" {
		m.rainfall = fmt.Sprintf("%smm", weatherData.Rainfall)
	}

	if weatherData.WindSpeed != "" {
		m.windSpeed = fmt.Sprintf("%smph", weatherData.WindSpeed)
	}

	if weatherData.WindDirection != "" {
		m.windDirection = fmt.Sprintf("%s", weatherData.WindDirection)
	}

	return nil
}

func (m *Model) handleTrackStatus(event f1.F1Event) error {
	var trackStatusData f1.TrackStatus
	if err := json.Unmarshal(event.Data, &trackStatusData); err != nil {
		return err
	}

	if trackStatusData.Message != "" {
		m.trackStatus = trackStatusData.Message
	}

	return nil
}

func (m *Model) handleLapCount(event f1.F1Event) error {
	var lapCountData f1.LapCount
	if err := json.Unmarshal(event.Data, &lapCountData); err != nil {
		return err
	}

	if lapCountData.CurrentLap != 0 {
		m.currentLap = lapCountData.CurrentLap
	}

	if lapCountData.TotalLaps != 0 {
		m.totalLaps = lapCountData.TotalLaps
	}

	return nil
}
