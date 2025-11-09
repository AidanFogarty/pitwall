package sessions

import (
	"fmt"
	"io"
	"strings"

	"github.com/AidanFogarty/pitwall/internal/tui/shared"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().Padding(0).Margin(0)
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DB2A20"))
)

type item struct {
	dirName     string
	displayName string
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.displayName)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i item) FilterValue() string { return "" }

type ReplayListSessionModel struct {
	width, height int

	list list.Model

	selected string
	quitting bool
}

func NewModel(dirNames []string) tea.Model {
	var items []list.Item
	for _, dirName := range dirNames {
		details := strings.Split(dirName, "_")

		year := details[0]
		session := details[1]
		location := details[2]
		sessionType := details[3]

		// example: 2025_Singapore-Grand-Prix_Marina-Bay_Race -> 2025 Singapore Grand Prix, Marina Bay, Race
		displayName := fmt.Sprintf("%s %s, %s, %s", year, strings.ReplaceAll(session, "-", " "), strings.ReplaceAll(location, "-", " "), sessionType)

		items = append(items, item{
			dirName:     dirName,
			displayName: displayName,
		})
	}

	title := "Select a session to replay"
	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = title
	l.Styles.Title = titleStyle
	l.Styles.TitleBar = lipgloss.NewStyle()
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	model := &ReplayListSessionModel{
		list: l,
	}
	return model
}

func (m ReplayListSessionModel) Init() tea.Cmd {
	return nil
}

func (m ReplayListSessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.dirName
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		listHeight := 1 + len(m.list.Items()) + 1
		m.list.SetSize(msg.Width, listHeight)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ReplayListSessionModel) View() string {
	content := shared.Logo() + "\n" + m.list.View()
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m ReplayListSessionModel) IsQuiting() bool {
	return m.quitting
}

func (m ReplayListSessionModel) Selected() string {
	return m.selected
}
