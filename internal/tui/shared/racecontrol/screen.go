package racecontrol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AidanFogarty/pitwall/internal/f1"
	"github.com/AidanFogarty/pitwall/internal/tui/shared"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type kind int

const (
	kindInfo kind = iota
	kindFlag
	kindSafetyCar
	kindDRS
	kindIncident
)

const (
	colorYellow = lipgloss.Color("#ffd230")
	colorRed    = lipgloss.Color("#ff4d4d")
	colorGreen  = lipgloss.Color("#00ff00")
	colorOrange = lipgloss.Color("#ff9f1a")
	colorSelBg  = lipgloss.Color("#1c2433")
	colorAccent = lipgloss.Color("#4aa3ff")
	colorSelFg  = lipgloss.Color("#e5e7eb")
)

var driverRef = regexp.MustCompile(`\b(\d{1,2}) \([A-Z]{3}\)`)

var incidentKeywords = []string{
	"INCIDENT", "DELETED", "PENALTY", "INVESTIGATION", "STEWARDS", "NOTED", "REPRIMAND",
}

func classify(msg f1.RaceControlMessage) kind {
	text := strings.ToUpper(msg.Message)

	switch {
	case msg.Category == "SafetyCar" || strings.Contains(text, "SAFETY CAR"):
		return kindSafetyCar
	case msg.Category == "Flag":
		return kindFlag
	case msg.Category == "Drs" || strings.HasPrefix(text, "DRS "):
		return kindDRS
	}

	for _, kw := range incidentKeywords {
		if strings.Contains(text, kw) {
			return kindIncident
		}
	}
	return kindInfo
}

func (k kind) badge() string {
	switch k {
	case kindFlag:
		return "FLAG"
	case kindSafetyCar:
		return "SC"
	case kindDRS:
		return "DRS"
	case kindIncident:
		return "STEWARD"
	default:
		return "INFO"
	}
}

func messageStyle(msg f1.RaceControlMessage, k kind, selected bool) lipgloss.Style {
	plain := lipgloss.NewStyle()
	dim := lipgloss.NewStyle().Foreground(shared.ColorDim)
	if selected {
		plain = plain.Foreground(colorSelFg)
		dim = dim.Foreground(colorSelFg)
	}

	text := strings.ToUpper(msg.Message)

	switch k {
	case kindFlag:
		switch strings.ToUpper(msg.Flag) {
		case "YELLOW":
			return plain.Foreground(colorYellow)
		case "DOUBLE YELLOW":
			return plain.Foreground(colorYellow).Bold(true)
		case "RED":
			return plain.Foreground(lipgloss.Color("#ffffff")).Background(colorRed).Bold(true)
		case "GREEN":
			return plain.Foreground(colorGreen)
		case "CHEQUERED":
			if selected {
				return plain
			}
			return plain.Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#ffffff")).Bold(true)
		case "BLUE", "CLEAR":
			return dim
		}
	case kindSafetyCar:
		if strings.Contains(text, "ENDING") || strings.Contains(text, "ENDED") {
			return dim
		}
		return plain.Foreground(colorOrange).Bold(true)
	case kindDRS:
		return dim
	case kindIncident:
		switch {
		case strings.Contains(text, "PENALTY") && !strings.Contains(text, "SERVED"):
			return plain.Foreground(colorRed).Bold(true)
		case strings.Contains(text, "NO FURTHER") || strings.Contains(text, "REVIEWED"):
			return dim
		}
	}
	return plain
}

const (
	timeW  = 10
	lapW   = 4
	badgeW = 7
	scopeW = 9

	chromeRows = 4
)

func viewRows(height int) int {
	return max(height-chromeRows, 1)
}

func (m Model) selected() int {
	if m.tailing {
		return max(len(m.messages)-1, 0)
	}
	return min(m.cursor, max(len(m.messages)-1, 0))
}

func (m Model) ScrollUp(n, height int) Model {
	return m.moveTo(m.selected()-n, height)
}

func (m Model) ScrollDown(n, height int) Model {
	return m.moveTo(m.selected()+n, height)
}

func (m Model) ScrollToStart(height int) Model {
	return m.moveTo(0, height)
}

func (m Model) ScrollToEnd(height int) Model {
	return m.moveTo(len(m.messages)-1, height)
}

func (m Model) moveTo(index, height int) Model {
	last := max(len(m.messages)-1, 0)
	rows := viewRows(height)

	if m.tailing {
		m.top = max(len(m.messages)-rows, 0)
	}

	index = min(max(index, 0), last)
	m.cursor = index
	m.tailing = index == last

	if index < m.top {
		m.top = index
	} else if index >= m.top+rows {
		m.top = index - rows + 1
	}
	return m
}

func formatTime(utc string) string {
	if i := strings.Index(utc, "T"); i >= 0 && len(utc) >= i+9 {
		return utc[i+1 : i+9]
	}
	return ""
}

func formatScope(msg f1.RaceControlMessage) string {
	switch {
	case strings.EqualFold(msg.Scope, "Sector") && msg.Sector > 0:
		return fmt.Sprintf("SECTOR %d", msg.Sector)
	case msg.Scope != "":
		return strings.ToUpper(msg.Scope)
	}
	return ""
}

func (m Model) renderMessage(text string, base lipgloss.Style) string {
	matches := driverRef.FindAllStringSubmatchIndex(text, -1)

	var out strings.Builder
	last := 0
	for _, loc := range matches {
		colour, ok := m.teamColours[text[loc[2]:loc[3]]]
		if !ok {
			continue
		}

		out.WriteString(base.Render(text[last:loc[0]]))
		out.WriteString(base.Foreground(lipgloss.Color("#" + colour)).Bold(true).Render(text[loc[0]:loc[1]]))
		last = loc[1]
	}
	out.WriteString(base.Render(text[last:]))
	return out.String()
}

func (m Model) tableRow(msg f1.RaceControlMessage, selected bool) []string {
	k := classify(msg)

	lap := ""
	if msg.Lap > 0 {
		lap = fmt.Sprintf("L%d", msg.Lap)
	}

	timeStyle := lipgloss.NewStyle().Foreground(shared.ColorDim)
	msgStyle := messageStyle(msg, k, selected)
	marker := " "

	if selected {
		timeStyle = lipgloss.NewStyle().Foreground(colorSelFg).Background(colorSelBg)
		if _, none := msgStyle.GetBackground().(lipgloss.NoColor); none {
			msgStyle = msgStyle.Background(colorSelBg)
		}
		marker = lipgloss.NewStyle().Foreground(colorAccent).Background(colorSelBg).Render("▌")
	}

	return []string{
		marker + timeStyle.Render(" "+formatTime(msg.Utc)),
		lap,
		k.badge(),
		formatScope(msg),
		m.renderMessage(msg.Message, msgStyle),
	}
}

func (m Model) FullView(width, height int) string {
	rows := viewRows(height)
	total := len(m.messages)
	selected := m.selected()

	top := m.top
	if m.tailing {
		top = max(total-rows, 0)
	}
	top = min(max(top, 0), max(total-rows, 0))

	var data [][]string
	for i := top; i < min(top+rows, total); i++ {
		data = append(data, m.tableRow(m.messages[i], i == selected))
	}
	for len(data) < rows {
		data = append(data, []string{"", "", "", "", ""})
	}

	columnWidths := []int{timeW, lapW, badgeW, scopeW}

	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(shared.ColorDim)).
		Headers("  TIME", "LAP", "TYPE", "SCOPE", "MESSAGE").
		Rows(data...).
		Width(width).
		Wrap(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)
			if col < len(columnWidths) {
				style = style.Width(columnWidths[col] + 2)
			}

			if row == table.HeaderRow {
				return style.Foreground(colorSelFg).Bold(true)
			}

			if col < len(columnWidths) {
				style = style.Foreground(shared.ColorDim)
			}
			if row == selected-top {
				style = style.Background(colorSelBg)
				if col < len(columnWidths) {
					style = style.Foreground(colorSelFg)
				}
			}
			return style
		}).
		String()
}

func (m Model) ScrollHint() string {
	hint := "↑/↓ scroll · g/G oldest/newest"
	if !m.tailing {
		hint = fmt.Sprintf("▼ %d newer · %s", len(m.messages)-1-m.selected(), hint)
	}
	return hint
}
