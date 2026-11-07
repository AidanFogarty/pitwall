package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AidanFogarty/pitwall/internal/f1"
	"github.com/AidanFogarty/pitwall/internal/logging"
	"github.com/AidanFogarty/pitwall/internal/tui/dashboard"
	"github.com/AidanFogarty/pitwall/internal/tui/qualifying"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var liveLogFile string

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "connect to a live f1 session",
	Long:  "view detailed timing stats, dashboards, driver info for a live Formula 1 session",
	RunE:  runLive,
}

func runLive(cmd *cobra.Command, args []string) error {
	var model tea.Model
	var program *tea.Program

	client := f1.NewF1LiveClient(f1.AllTopics)
	if liveLogFile != "" {
		if err := client.SetLogFile(liveLogFile); err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	initialState, err := client.Start(cmd.Context())
	if err != nil {
		return err
	}

	if initialState.SessionInfo == nil {
		return fmt.Errorf("no session info in initial state")
	}

	sessionType := initialState.SessionInfo.Type
	switch sessionType {
	case "Race":
		model = dashboard.NewModel()
		program = tea.NewProgram(model, tea.WithAltScreen())
	case "Qualifying", "Practice":
		model = qualifying.NewModel()
		program = tea.NewProgram(model, tea.WithAltScreen())
	default:
		panic("Unknown session type.")
	}

	client.SetHandler(func(topic string, data any, timestamp string) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		program.Send(f1.F1EventMsg{
			Type: topic,
			Data: jsonData,
		})

		return nil
	})

	go client.Watch(cmd.Context(), func(state *f1.InitialLiveState) {
		populateInitialState(program, state)
	}, func(state string) {
		program.Send(f1.ConnectionMsg(state))
	})

	// TODO: I am marshalling and unmarshalling data everywhere, I need to fix this later
	go populateInitialState(program, initialState)
	go program.Send(f1.LoginMsg(loginStatus()))
	_, err = program.Run()

	return err
}

func init() {
	liveCmd.Flags().StringVar(&liveLogFile, "log-file", logging.DefaultPath("live.log"), "connection diagnostics log file (empty to disable)")
	rootCmd.AddCommand(liveCmd)
}

func populateInitialState(program *tea.Program, state *f1.InitialLiveState) {
	if state.SessionInfo != nil {
		data, _ := json.Marshal(state.SessionInfo)
		program.Send(f1.F1EventMsg{
			Type: "SessionInfo",
			Data: data,
		})
	}

	if state.WeatherData != nil {
		data, _ := json.Marshal(state.WeatherData)
		program.Send(f1.F1EventMsg{
			Type: "WeatherData",
			Data: data,
		})
	}

	if state.TrackStatus != nil {
		data, _ := json.Marshal(state.TrackStatus)
		program.Send(f1.F1EventMsg{
			Type: "TrackStatus",
			Data: data,
		})
	}

	if state.LapCount != nil {
		data, _ := json.Marshal(state.LapCount)
		program.Send(f1.F1EventMsg{
			Type: "LapCount",
			Data: data,
		})
	}

	if state.DriverList != nil {
		data, _ := json.Marshal(state.DriverList)
		program.Send(f1.F1EventMsg{
			Type: "DriverList",
			Data: data,
		})
	}

	if state.TimingData != nil {
		data, _ := json.Marshal(state.TimingData)
		program.Send(f1.F1EventMsg{
			Type: "TimingData",
			Data: data,
		})
	}

	if state.TopThree != nil {
		data, _ := json.Marshal(state.TopThree)
		program.Send(f1.F1EventMsg{
			Type: "TopThree",
			Data: data,
		})
	}

	if state.TimingAppData != nil {
		data, _ := json.Marshal(state.TimingAppData)
		program.Send(f1.F1EventMsg{
			Type: "TimingAppData",
			Data: data,
		})
	}

	if state.RaceControlMessages != nil {
		data, _ := json.Marshal(state.RaceControlMessages)
		program.Send(f1.F1EventMsg{
			Type: "RaceControlMessages",
			Data: data,
		})
	}
}
