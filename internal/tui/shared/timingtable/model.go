package timingtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/AidanFogarty/pitwall/internal/f1"
	"github.com/AidanFogarty/pitwall/internal/tui/shared"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var timingTableHeaders = []string{
	"Pos",
	"Driver",
	"Gap",
	"Leader",
	"Best Lap",
	"Last Lap",
	"S1",
	"S2",
	"S3",
	"Pit",
	"Tyre",
}

type DriverState struct {
	RacingNumber string
	Tla          string
	TeamColour   string

	Position int

	Gap     string
	Leader  string
	BestLap string
	LastLap string

	S1 string
	S2 string
	S3 string

	S1PersonalFastest bool
	S2PersonalFastest bool
	S3PersonalFastest bool

	PitStatus string

	Tyre    string
	TyreAge int
}

type Model struct {
	table *table.Table

	drivers map[string]*DriverState

	s1OverallFastestHolder string
	s2OverallFastestHolder string
	s3OverallFastestHolder string

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
	case "TopThree":
		m.handleTopThree(event)
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

func (m *Model) handleTopThree(event f1.F1Event) error {
	var topThreeData f1.TopThree
	if err := json.Unmarshal(event.Data, &topThreeData); err != nil {
		return err
	}

	lines := bytes.TrimSpace(topThreeData.Lines)
	if len(lines) == 0 {
		return nil
	}

	if lines[0] == '[' {
		var linesArray []f1.TopThreeLine
		if err := json.Unmarshal(topThreeData.Lines, &linesArray); err != nil {
			return err
		}

		for i, data := range linesArray {
			if data.RacingNumber == "" {
				continue
			}

			driver := m.getOrCreateDriver(data.RacingNumber)
			driver.Position = i + 1

			if data.Tla != "" {
				driver.Tla = data.Tla
			}
			if data.TeamColour != "" {
				driver.TeamColour = data.TeamColour
			}

			if data.DiffToAhead != "" {
				driver.Gap = data.DiffToAhead
			}
			if data.DiffToLeader != "" {
				driver.Leader = data.DiffToLeader
			}
		}

	} else if lines[0] == '{' {
		var linesMap map[string]f1.TopThreeLine
		if err := json.Unmarshal(topThreeData.Lines, &linesMap); err != nil {
			return err
		}

		for positionStr, data := range linesMap {
			positionIndex, err := strconv.Atoi(positionStr)
			if err != nil {
				continue
			}
			position := positionIndex + 1

			var targetDriver *DriverState
			for _, driver := range m.drivers {
				if driver.Position == position {
					targetDriver = driver
					break
				}
			}

			if targetDriver == nil {
				continue
			}

			if data.DiffToLeader != "" {
				targetDriver.Leader = data.DiffToLeader
			}
			if data.DiffToAhead != "" {
				targetDriver.Gap = data.DiffToAhead
			}
		}
	}

	m.rebuildTable()
	return nil
}

func (m *Model) handleTimingData(event f1.F1Event) error {
	var timingData f1.TimingData
	if err := json.Unmarshal(event.Data, &timingData); err != nil {
		return err
	}

	for racingNumber, data := range timingData.Lines {
		driver := m.getOrCreateDriver(racingNumber)

		if data.Line != 0 {
			driver.Position = data.Line
		}

		if data.IntervalToPositionAhead.Value != "" {
			driver.Gap = data.IntervalToPositionAhead.Value
		}
		if data.GapToLeader != "" {
			driver.Leader = data.GapToLeader
		}

		if data.BestLapTime.Value != "" {
			driver.BestLap = data.BestLapTime.Value
		}
		if data.LastLapTime.Value != "" {
			driver.LastLap = data.LastLapTime.Value
		}

		if data.InPit {
			driver.PitStatus = "Pit"
		} else if data.PitOut {
			driver.PitStatus = "Out"
		} else {
			driver.PitStatus = ""
		}

		if len(data.Sectors) > 0 {
			m.updateSectors(driver, data.Sectors)
		}
	}

	m.rebuildTable()
	return nil
}

func (m *Model) updateSectors(driver *DriverState, sectors json.RawMessage) {
	trimmed := bytes.TrimSpace(sectors)
	if len(trimmed) == 0 {
		return
	}

	switch trimmed[0] {
	case '[':
		var sectorsArray []f1.SectorData
		if err := json.Unmarshal(sectors, &sectorsArray); err != nil {
			return
		}

		if len(sectorsArray) > 0 && sectorsArray[0].Value != "" {
			driver.S1 = sectorsArray[0].Value

			driver.S1PersonalFastest = sectorsArray[0].PersonalFastest
			if sectorsArray[0].OverallFastest {
				m.s1OverallFastestHolder = driver.RacingNumber
			}

			driver.S2 = ""
			driver.S3 = ""
			driver.S2PersonalFastest = false
			driver.S3PersonalFastest = false
		}
		if len(sectorsArray) > 1 && sectorsArray[1].Value != "" {
			driver.S2 = sectorsArray[1].Value

			driver.S2PersonalFastest = sectorsArray[1].PersonalFastest
			if sectorsArray[1].OverallFastest {
				m.s2OverallFastestHolder = driver.RacingNumber
			}
		}
		if len(sectorsArray) > 2 && sectorsArray[2].Value != "" {
			driver.S3 = sectorsArray[2].Value

			driver.S3PersonalFastest = sectorsArray[2].PersonalFastest
			if sectorsArray[2].OverallFastest {
				m.s3OverallFastestHolder = driver.RacingNumber
			}
		}

	case '{':
		var sectorsMap map[string]f1.SectorData
		if err := json.Unmarshal(sectors, &sectorsMap); err != nil {
			return
		}

		if s1, exists := sectorsMap["0"]; exists && s1.Value != "" {
			driver.S1 = s1.Value

			driver.S1PersonalFastest = s1.PersonalFastest
			if s1.OverallFastest {
				m.s1OverallFastestHolder = driver.RacingNumber
			}

			driver.S2 = ""
			driver.S3 = ""
			driver.S2PersonalFastest = false
			driver.S3PersonalFastest = false
		}
		if s2, exists := sectorsMap["1"]; exists && s2.Value != "" {
			driver.S2 = s2.Value

			driver.S2PersonalFastest = s2.PersonalFastest
			if s2.OverallFastest {
				m.s2OverallFastestHolder = driver.RacingNumber
			}
		}
		if s3, exists := sectorsMap["2"]; exists && s3.Value != "" {
			driver.S3 = s3.Value

			driver.S3PersonalFastest = s3.PersonalFastest
			if s3.OverallFastest {
				m.s3OverallFastestHolder = driver.RacingNumber
			}
		}
	}
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

func (m *Model) buildTableRows(drivers []*DriverState) [][]string {
	baseStyle := lipgloss.NewStyle().Padding(0, 1)

	personalBestStyle := baseStyle.Foreground(shared.ColorSectorPersonalFastest)
	overallBestStyle := baseStyle.Foreground(shared.ColorSectorOverallFastest)

	if len(drivers) == 0 {
		rows := make([][]string, f1.NumberOfDrivers)
		emptyRow := []string{"", "", "", "", "", "", "", "", "", "", ""}
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

		var tyreDisplay string
		if driver.Tyre != "" {
			tyreChar := string(driver.Tyre[0])
			tyreStyle := baseStyle

			switch tyreChar {
			case "S":
				tyreStyle = baseStyle.Foreground(shared.ColorTyreSoft)
			case "M":
				tyreStyle = baseStyle.Foreground(shared.ColorTyreMedium)
			case "H":
				tyreStyle = baseStyle.Foreground(shared.ColorTyreHard)
			}

			tyreDisplay = tyreStyle.Render(fmt.Sprintf("%s %d", tyreChar, driver.TyreAge))
		} else {
			tyreDisplay = baseStyle.Render("")
		}

		s1Style := baseStyle
		if driver.RacingNumber == m.s1OverallFastestHolder {
			s1Style = overallBestStyle
		} else if driver.S1PersonalFastest {
			s1Style = personalBestStyle
		}

		s2Style := baseStyle
		if driver.RacingNumber == m.s2OverallFastestHolder {
			s2Style = overallBestStyle
		} else if driver.S2PersonalFastest {
			s2Style = personalBestStyle
		}

		s3Style := baseStyle
		if driver.RacingNumber == m.s3OverallFastestHolder {
			s3Style = overallBestStyle
		} else if driver.S3PersonalFastest {
			s3Style = personalBestStyle
		}

		rows[i] = []string{
			baseStyle.Render(fmt.Sprintf("%d", driver.Position)),
			driverStyle.Render(fmt.Sprintf("%s %s", driver.RacingNumber, driver.Tla)),
			baseStyle.Render(driver.Gap),
			baseStyle.Render(driver.Leader),
			baseStyle.Render(driver.BestLap),
			baseStyle.Render(driver.LastLap),
			s1Style.Render(driver.S1),
			s2Style.Render(driver.S2),
			s3Style.Render(driver.S3),
			baseStyle.Render(driver.PitStatus),
			baseStyle.Render(tyreDisplay),
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
		Headers(timingTableHeaders...).
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
