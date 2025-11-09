package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/AidanFogarty/pitwall/internal/importer"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import historical f1 session data",
	Long:  "Import historical f1 session data to replay races, qualifying and more",
	RunE:  runImport,
}

func runImport(cmd *cobra.Command, args []string) error {
	year, err := cmd.Flags().GetInt("year")
	if err != nil {
		return err
	}

	meeting, err := cmd.Flags().GetInt("meeting")
	if err != nil {
		return err
	}

	session, err := cmd.Flags().GetInt("session")
	if err != nil {
		return err
	}

	config := importer.DefaultConfig()
	imp := importer.NewImporter(config, dataDir)

	if meeting == 0 || session == 0 {
		data, err := imp.GetAvailableMeetings(cmd.Context(), year)
		if err != nil {
			return err
		}

		printMeetingsTable(data.Meetings)
		return nil
	}

	return imp.ImportSession(cmd.Context(), year, meeting, session)
}

func printMeetingsTable(meetings []importer.Meeting) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Key\tName\tLocation\tP1\tP2\tP3\tSQ\tSprint\tQual\tRace")

	for _, meeting := range meetings {
		sessions := buildSessionMap(meeting.Sessions)
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			meeting.Key,
			meeting.Name,
			meeting.Location,
			formatSessionKey(sessions["P1"]),
			formatSessionKey(sessions["P2"]),
			formatSessionKey(sessions["P3"]),
			formatSessionKey(sessions["SQ"]),
			formatSessionKey(sessions["Sprint"]),
			formatSessionKey(sessions["Qual"]),
			formatSessionKey(sessions["Race"]))
	}
}

func buildSessionMap(sessions []importer.Session) map[string]int {
	sessionMap := make(map[string]int)
	for _, s := range sessions {
		switch s.Name {
		case "Practice 1", "Day 1":
			sessionMap["P1"] = s.Key
		case "Practice 2", "Day 2":
			sessionMap["P2"] = s.Key
		case "Practice 3", "Day 3":
			sessionMap["P3"] = s.Key
		case "Sprint Qualifying":
			sessionMap["SQ"] = s.Key
		case "Qualifying":
			sessionMap["Qual"] = s.Key
		case "Sprint":
			sessionMap["Sprint"] = s.Key
		case "Race":
			sessionMap["Race"] = s.Key
		}
	}
	return sessionMap
}

func formatSessionKey(key int) string {
	if key == 0 {
		return "-"
	}
	return strconv.Itoa(key)
}

func init() {
	year, _, _ := time.Now().Date()

	importCmd.Flags().Int("year", year, "the session to import")
	importCmd.Flags().Int("meeting", 0, "the meeting key of the session")
	importCmd.Flags().Int("session", 0, "the session key for the given meeting")

	rootCmd.AddCommand(importCmd)
}
