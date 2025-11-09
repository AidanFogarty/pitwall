package shared

import "github.com/charmbracelet/lipgloss"

var pit = `
▗▄▄▖▗▄▄▄▖▗▄▄▄▖
▐▌ ▐▌ █    █
▐▛▀▘  █    █
▐▌  ▗▄█▄▖  █
`

var wall = `
▗▖ ▗▖ ▗▄▖ ▗▖   ▗▖
▐▌ ▐▌▐▌ ▐▌▐▌   ▐▌
▐▌ ▐▌▐▛▀▜▌▐▌   ▐▌
▐▙█▟▌▐▌ ▐▌▐▙▄▄▖▐▙▄▄▖
`

func Logo() string {
	f1RedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#DB2A20"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	logo := lipgloss.JoinHorizontal(
		lipgloss.Center,
		f1RedStyle.Render(pit),
		normalStyle.Render(wall),
	)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Width(lipgloss.Width(logo)).
		Align(lipgloss.Left)

	desc := descStyle.Render("A F1 live timing client for your terminal")

	logoAndDesc := lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		desc,
	)

	lines := []string{logoAndDesc, ""}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return content
}
