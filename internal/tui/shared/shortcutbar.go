package shared

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func ShortcutBar(help, hint string, width int) string {
	left := lipgloss.NewStyle().PaddingLeft(1).Render(help)
	if hint == "" {
		return left
	}

	right := lipgloss.NewStyle().
		Foreground(ColorDim).
		PaddingRight(1).
		Render(hint)

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		return left
	}
	return left + strings.Repeat(" ", gap) + right
}
