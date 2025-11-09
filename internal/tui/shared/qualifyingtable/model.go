package qualifyingtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/AidanFogarty/pitwall/internal/f1"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var qualifyingTableHeaders = []string{
	"Pos",
	"Driver",
	"Lap Time",
	"Best Lap",
	"Sector 1",
	"Sector 2",
	"Sector 3",
	"Pit",
	"Tyre",
	"S1",
	"S2",
	"S3",
}

const (
	SegmentStatusNormal         = 2048
	SegmentStatusPersonalBest   = 2049
	SegmentStatusOverallFastest = 2051
)

const (
	ColorPurple = "#800080"
	ColorGreen  = "#00FF00"
	ColorYellow = "#FFFF00"
	ColorGray   = "#808080"
)

type DriverState struct {
	RacingNumber string
	Tla          string
	TeamColour   string

	Position int

	LastLap string
	BestLap string

	S1 string
	S2 string
	S3 string

	S1Segments []int
	S2Segments []int
	S3Segments []int

	PitStatus string

	Tyre    string
	TyreAge int
}

type Model struct {
	table *table.Table

	drivers map[string]*DriverState

	width  int
	height int
}

func New() Model {
	return Model{
		drivers: make(map[string]*DriverState),
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
	if m.table == nil {
		return ""
	}
	return m.table.String()
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height

	if m.table == nil {
		m.rebuildTable()
	}

	return m
}

func (m *Model) handleEvent(event f1.F1Event) {
	switch event.Type {
	case "DriverList":
		m.handleDriverList(event)
	case "TimingData":
		m.handleTimingData(event)
	case "TimingAppData":
		m.handleTimingAppData(event)
	}
}

func (m *Model) handleDriverList(event f1.F1Event) error {
	var driverListData f1.DriverList
	if err := json.Unmarshal(event.Data, &driverListData); err != nil {
		return err
	}

	for racingNumber, data := range driverListData {
		driver := m.getOrCreateDriver(racingNumber)

		if data.Tla != "" {
			driver.Tla = data.Tla
		}
		if data.TeamColour != "" {
			driver.TeamColour = data.TeamColour
		}

		if data.Line != 0 {
			driver.Position = data.Line
		}
	}

	m.rebuildTable()
	return nil
}

func (m *Model) handleTimingData(event f1.F1Event) error {
	// First parse to get the raw data for detecting explicit false values
	var rawTimingData struct {
		Lines map[string]json.RawMessage `json:"Lines"`
	}
	if err := json.Unmarshal(event.Data, &rawTimingData); err != nil {
		return err
	}

	var timingData f1.TimingData
	if err := json.Unmarshal(event.Data, &timingData); err != nil {
		return err
	}

	for racingNumber, data := range timingData.Lines {
		driver := m.getOrCreateDriver(racingNumber)

		if data.Line != 0 {
			driver.Position = data.Line
		}

		if data.LastLapTime.Value != "" {
			driver.LastLap = data.LastLapTime.Value
		}

		if data.BestLapTime.Value != "" {
			driver.BestLap = data.BestLapTime.Value
		}

		// This is a hack for now until I change the pitout to use boolean pointers.
		rawLine := rawTimingData.Lines[racingNumber]
		pitOutExplicitlyFalse := bytes.Contains(rawLine, []byte(`"PitOut":false`))

		if data.InPit {
			driver.PitStatus = "Pit"
		} else if data.PitOut {
			driver.PitStatus = "Out"
		} else if pitOutExplicitlyFalse && driver.PitStatus == "Out Lap" {
			driver.PitStatus = ""
		}

		if len(data.Sectors) > 0 {
			m.updateSectors(driver, data.Sectors)
		}
	}

	m.rebuildTable()
	return nil
}

func (m *Model) updateSectors(driver *DriverState, sectorsRaw json.RawMessage) {
	trimmed := bytes.TrimSpace(sectorsRaw)
	if len(trimmed) == 0 {
		return
	}

	var sectorsMap map[string]json.RawMessage
	if err := json.Unmarshal(sectorsRaw, &sectorsMap); err != nil {
		return
	}

	for sectorKey, sectorRaw := range sectorsMap {
		var sectorData f1.SectorData
		if err := json.Unmarshal(sectorRaw, &sectorData); err != nil {
			continue
		}

		switch sectorKey {
		case "0":
			if sectorData.Value != "" {
				driver.S1 = sectorData.Value
				driver.S2 = ""
				driver.S3 = ""
			}
		case "1":
			if sectorData.Value != "" {
				driver.S2 = sectorData.Value
			}
		case "2":
			if sectorData.Value != "" {
				driver.S3 = sectorData.Value
			}
		}

		if len(sectorData.Segments) == 0 {
			continue
		}

		var existingSegments []int
		switch sectorKey {
		case "0":
			existingSegments = driver.S1Segments
		case "1":
			existingSegments = driver.S2Segments
		case "2":
			existingSegments = driver.S3Segments
		}

		segments := m.parseSegments(sectorData.Segments, existingSegments)

		switch sectorKey {
		case "0":
			driver.S1Segments = segments
		case "1":
			driver.S2Segments = segments
		case "2":
			driver.S3Segments = segments
		}
	}
}

func (m *Model) parseSegments(segmentsRaw json.RawMessage, existingSegments []int) []int {
	trimmed := bytes.TrimSpace(segmentsRaw)
	if len(trimmed) == 0 {
		return existingSegments
	}

	if trimmed[0] == '[' {
		var segmentsArray []f1.SegmentData
		if err := json.Unmarshal(segmentsRaw, &segmentsArray); err != nil {
			return existingSegments
		}

		statuses := make([]int, len(segmentsArray))
		for i, seg := range segmentsArray {
			statuses[i] = seg.Status
		}
		return statuses
	}

	if trimmed[0] == '{' {
		var segmentsMap map[string]f1.SegmentData
		if err := json.Unmarshal(segmentsRaw, &segmentsMap); err != nil {
			return existingSegments
		}

		maxIndex := len(existingSegments) - 1
		for key := range segmentsMap {
			if idx, err := strconv.Atoi(key); err == nil && idx > maxIndex {
				maxIndex = idx
			}
		}

		requiredSize := maxIndex + 1
		statuses := make([]int, requiredSize)
		copy(statuses, existingSegments)

		for key, seg := range segmentsMap {
			if idx, err := strconv.Atoi(key); err == nil && idx < len(statuses) {
				statuses[idx] = seg.Status
			}
		}
		return statuses
	}

	return existingSegments
}

func (m *Model) handleTimingAppData(event f1.F1Event) error {
	var timingAppData f1.TimingAppData
	if err := json.Unmarshal(event.Data, &timingAppData); err != nil {
		return err
	}

	for racingNumber, data := range timingAppData.Lines {
		driver := m.getOrCreateDriver(racingNumber)

		if len(data.Stints) > 0 {
			m.updateStints(driver, data.Stints)
		}
	}

	return nil
}

func (m *Model) updateStints(driver *DriverState, stints json.RawMessage) {
	trimmed := bytes.TrimSpace(stints)
	if len(trimmed) == 0 {
		return
	}

	switch stints[0] {
	case '[':
		var stintsArray []f1.StintData
		if err := json.Unmarshal(stints, &stintsArray); err != nil {
			return
		}

		if len(stintsArray) == 0 {
			return
		}

		latestStint := stintsArray[len(stintsArray)-1]

		if latestStint.Compound != "" {
			driver.Tyre = cases.Title(language.English, cases.Compact).String(latestStint.Compound)
		}

		if latestStint.TotalLaps != nil {
			driver.TyreAge = *latestStint.TotalLaps
		}

	case '{':
		var stintsMap map[string]f1.StintData
		if err := json.Unmarshal(stints, &stintsMap); err != nil {
			return
		}

		if len(stintsMap) == 0 {
			return
		}

		var keys []string
		for key := range stintsMap {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i, j int) bool {
			first, _ := strconv.Atoi(keys[i])
			second, _ := strconv.Atoi(keys[j])
			return first < second
		})

		latestStintNumber := keys[len(keys)-1]
		latestStint := stintsMap[latestStintNumber]

		if latestStint.Compound != "" {
			driver.Tyre = cases.Title(language.English, cases.Compact).String(latestStint.Compound)
		}

		if latestStint.TotalLaps != nil {
			driver.TyreAge = *latestStint.TotalLaps
		}
	}
}

func (m *Model) getOrCreateDriver(racingNumber string) *DriverState {
	driver, exists := m.drivers[racingNumber]
	if !exists {
		driver = &DriverState{
			RacingNumber: racingNumber,
		}
		m.drivers[racingNumber] = driver
	}
	return driver
}

func (m *Model) buildSortedDriverList() []*DriverState {
	drivers := make([]*DriverState, 0, len(m.drivers))
	for _, driver := range m.drivers {
		if driver.Position > 0 {
			drivers = append(drivers, driver)
		}
	}

	sort.Slice(drivers, func(i, j int) bool {
		return drivers[i].Position < drivers[j].Position
	})

	return drivers
}

func renderMiniSectors(segments []int) string {
	if len(segments) == 0 {
		return ""
	}

	var result string
	for i, status := range segments {
		var color lipgloss.Color
		switch status {
		case SegmentStatusOverallFastest:
			color = lipgloss.Color(ColorPurple)
		case SegmentStatusPersonalBest:
			color = lipgloss.Color(ColorGreen)
		case SegmentStatusNormal:
			color = lipgloss.Color(ColorYellow)
		default:
			color = lipgloss.Color(ColorGray)
		}

		style := lipgloss.NewStyle().
			Foreground(color).
			Bold(true)

		result += style.Render("■")

		if i < len(segments)-1 {
			result += " "
		}
	}

	return result
}

func (m *Model) buildTableRows(drivers []*DriverState) [][]string {
	baseStyle := lipgloss.NewStyle().Padding(0, 1)

	if len(drivers) == 0 {
		rows := make([][]string, f1.NumberOfDrivers)
		emptyRow := []string{"", "", "", "", "", "", "", "", "", "", "", ""}
		for i := range rows {
			rows[i] = emptyRow
		}
		return rows
	}

	rows := make([][]string, len(drivers))

	for i, driver := range drivers {
		driverStyle := baseStyle
		if driver.TeamColour != "" {
			driverStyle = baseStyle.Foreground(lipgloss.Color("#" + driver.TeamColour))
		}

		rows[i] = []string{
			baseStyle.Render(fmt.Sprintf("%d", driver.Position)),
			driverStyle.Render(fmt.Sprintf("%s %s", driver.RacingNumber, driver.Tla)),
			baseStyle.Render(driver.LastLap),
			baseStyle.Render(driver.BestLap),
			baseStyle.Render(renderMiniSectors(driver.S1Segments)),
			baseStyle.Render(renderMiniSectors(driver.S2Segments)),
			baseStyle.Render(renderMiniSectors(driver.S3Segments)),
			baseStyle.Render(driver.PitStatus),
			baseStyle.Render(fmt.Sprintf("%s %d", driver.Tyre, driver.TyreAge)),
			baseStyle.Render(driver.S1),
			baseStyle.Render(driver.S2),
			baseStyle.Render(driver.S3),
		}
	}

	return rows
}

func (m *Model) rebuildTable() {
	drivers := m.buildSortedDriverList()
	rows := m.buildTableRows(drivers)

	headerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("252")).
		Bold(true)

	m.table = table.New().
		Border(lipgloss.NormalBorder()).
		Headers(qualifyingTableHeaders...).
		Rows(rows...).
		Width(m.width).
		Height(m.height).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == -1 {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})
}
