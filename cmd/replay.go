package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AidanFogarty/pitwall/internal/f1"
	"github.com/AidanFogarty/pitwall/internal/tui/dashboard"
	"github.com/AidanFogarty/pitwall/internal/tui/qualifying"
	"github.com/AidanFogarty/pitwall/internal/tui/replay/sessions"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var skipEventCount int

var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "replay a previous f1 session",
	Long:  "replay previous races, qualifying and view how the session played out.",
	RunE:  runReplay,
}

func runReplay(cmd *cobra.Command, args []string) error {
	items, err := getImportedSessions()

	model := sessions.NewModel(items)
	program := tea.NewProgram(model, tea.WithAltScreen())

	final, err := program.Run()
	if err != nil {
		return err
	}

	if m, ok := final.(sessions.ReplayListSessionModel); ok {
		if m.IsQuiting() {
			return nil
		}

		sessionType := getSessionType(m.Selected())

		var program *tea.Program
		switch sessionType {
		case "Race":
			model = dashboard.NewModel()
			program = tea.NewProgram(model, tea.WithAltScreen())
		case "Qualifying":
			model = qualifying.NewModel()
			program = tea.NewProgram(model, tea.WithAltScreen())
		default:
			panic("unknown session type")
		}

		config := f1.NewF1JsonConfig(dataDir, m.Selected())
		config.SkipDelayEventCount = skipEventCount

		client := f1.NewF1JsonClient(config, func(ctx context.Context, data f1.F1Event) error {
			program.Send(f1.F1EventMsg(data))
			return nil
		})

		go client.Start(cmd.Context())

		_, err = program.Run()
		if err != nil {
			return err
		}
	}

	return err
}

func init() {
	rootCmd.AddCommand(replayCmd)

	replayCmd.Flags().IntVar(&skipEventCount, "skip", -1, "number of events to fast-forward through immediately for debugging")
}

func getImportedSessions() ([]string, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return []string{}, err
	}

	items := []string{}
	for _, entry := range entries {
		dirName := entry.Name()
		items = append(items, dirName)
	}

	return items, nil
}

func getSessionType(selected string) string {
	sessionInfoFilePath := filepath.Join(dataDir, selected, "session_info.json")

	sessionInfo, err := os.ReadFile(sessionInfoFilePath)
	if err != nil {
		panic(err)
	}

	var sessionInfoData f1.SessionInfo
	err = json.Unmarshal(sessionInfo, &sessionInfoData)
	if err != nil {
		panic(err)
	}

	return sessionInfoData.Type
}
