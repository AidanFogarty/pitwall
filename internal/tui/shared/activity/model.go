package activity

import (
	"fmt"

	"github.com/AidanFogarty/pitwall/internal/f1"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	event int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case f1.F1EventMsg:
		m.event++
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("Events processed: %d", m.event)
}
