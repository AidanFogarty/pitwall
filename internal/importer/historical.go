package importer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AidanFogarty/pitwall/internal/f1"
	logmerge "github.com/AidanFogarty/pitwall/internal/util"
)

type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	DataDir    string
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL:    "https://livetiming.formula1.com/static",
		HTTPClient: http.DefaultClient,
	}
}

type Session struct {
	Key       int    `json:"Key"`
	Type      string `json:"Type"`
	Number    int    `json:"Number"`
	Name      string `json:"Name"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	GmtOffset string `json:"GmtOffset"`
	Path      string `json:"Path"`
}

type Country struct {
	Key  int    `json:"Key"`
	Code string `json:"Code"`
	Name string `json:"Name"`
}

type Circuit struct {
	Key       int    `json:"Key"`
	ShortName string `json:"ShortName"`
}

type Meeting struct {
	Sessions     []Session `json:"Sessions"`
	Key          int       `json:"Key"`
	Code         string    `json:"Code"`
	Number       int       `json:"Number"`
	Location     string    `json:"Location"`
	OfficialName string    `json:"OfficialName"`
	Name         string    `json:"Name"`
	Country      Country   `json:"Country"`
	Circuit      Circuit   `json:"Circuit"`
}

type IndexResponse struct {
	Year     int       `json:"Year"`
	Meetings []Meeting `json:"Meetings"`
}

type Importer struct {
	config  *Config
	dataDir string
}

func NewImporter(config *Config, dataDir string) *Importer {
	return &Importer{config: config, dataDir: dataDir}
}

func (imp *Importer) GetAvailableMeetings(ctx context.Context, year int) (*IndexResponse, error) {
	url := fmt.Sprintf("%s/%d/Index.json", imp.config.BaseURL, year)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating index request: %w", err)
	}

	resp, err := imp.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching index from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading index response: %w", err)
	}

	// Remove BOM if present
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	var index IndexResponse
	if err := json.Unmarshal(body, &index); err != nil {
		return nil, fmt.Errorf("parsing index JSON: %w", err)
	}

	return &index, nil
}

func (imp *Importer) ImportSession(ctx context.Context, year int, meetingKey int, sessionKey int) error {
	index, err := imp.fetchIndex(ctx, year)
	if err != nil {
		return err
	}

	meeting, err := findMeeting(index, meetingKey)
	if err != nil {
		return err
	}

	session, err := findSession(meeting, sessionKey)
	if err != nil {
		return err
	}

	outputDir, err := imp.setupOutputDirectories(year, meeting, session)
	if err != nil {
		return err
	}

	streamPaths, err := imp.downloadStreams(ctx, session, outputDir)
	if err != nil {
		return err
	}

	sessionInfoPath := filepath.Join(outputDir, "raw", "SessionInfo.jsonStream")
	outputSessionInfoPath := filepath.Join(outputDir, "session_info.json")
	err = createSessionInfoJson(sessionInfoPath, outputSessionInfoPath)
	if err != nil {
		return err
	}

	heartbeatPath := filepath.Join(outputDir, "raw", "Heartbeat.jsonStream")
	startTime, err := parseStartTimeFromHeartbeat(heartbeatPath)
	if err != nil {
		return err
	}

	mergedFile := filepath.Join(outputDir, "live.txt")
	if err := logmerge.Merge(streamPaths, mergedFile, startTime); err != nil {
		return err
	}

	return nil
}

func (imp *Importer) fetchIndex(ctx context.Context, year int) (*IndexResponse, error) {
	url := fmt.Sprintf("%s/%d/Index.json", imp.config.BaseURL, year)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating index request: %w", err)
	}

	resp, err := imp.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching index from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading index response: %w", err)
	}

	// Remove BOM if present
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	var index IndexResponse
	if err := json.Unmarshal(body, &index); err != nil {
		return nil, fmt.Errorf("parsing index JSON: %w", err)
	}

	return &index, nil
}

func (imp *Importer) setupOutputDirectories(year int, meeting *Meeting, session *Session) (string, error) {
	fullName := fmt.Sprintf("%d_%s_%s_%s", year, meeting.Name, meeting.Location, session.Name)
	normalized := strings.ReplaceAll(fullName, " ", "-")

	outputPath := filepath.Join(imp.dataDir, normalized)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return "", fmt.Errorf("creating directory %s: %w", outputPath, err)
	}

	rawPath := filepath.Join(outputPath, "raw")
	if err := os.MkdirAll(rawPath, 0755); err != nil {
		return "", fmt.Errorf("creating directory %s: %w", rawPath, err)
	}

	return outputPath, nil
}

func (imp *Importer) downloadStreams(ctx context.Context, session *Session, outputDir string) ([]string, error) {
	streamBasePath := fmt.Sprintf("%s/%s", imp.config.BaseURL, session.Path)
	var streamPaths []string

	var topics []string
	if session.Type == "Race" {
		topics = f1.RaceTopics
	} else {
		topics = f1.BaseTopics
	}

	for _, topic := range topics {
		streamURL := fmt.Sprintf("%s%s.jsonStream", streamBasePath, topic)
		fileName := fmt.Sprintf("%s.jsonStream", topic)
		filePath := filepath.Join(outputDir, "raw", fileName)

		err := imp.downloadStream(ctx, streamURL, filePath)
		if err != nil {
			return nil, err
		}

		streamPaths = append(streamPaths, filePath)
	}

	return streamPaths, nil
}

func (imp *Importer) downloadStream(ctx context.Context, url, filePath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request for %s: %w", url, err)
	}

	resp, err := imp.config.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response from %s: %w", url, err)
	}

	// Remove BOM if present
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	if err := os.WriteFile(filePath, body, 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", filePath, err)
	}

	return nil
}

func findMeeting(index *IndexResponse, key int) (*Meeting, error) {
	for _, meeting := range index.Meetings {
		if meeting.Key == key {
			return &meeting, nil
		}
	}
	return nil, fmt.Errorf("meeting with key %d not found", key)
}

func findSession(meeting *Meeting, key int) (*Session, error) {
	for _, session := range meeting.Sessions {
		if session.Key == key {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("session %d not found in meeting %d", key, meeting.Key)
}

func createSessionInfoJson(sessionInfoPath string, outputPath string) error {
	file, err := os.Open(sessionInfoPath)
	if err != nil {
		return fmt.Errorf("opening session info file %s: %w", sessionInfoPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 12 {
			continue
		}

		data := line[12:]

		var sessionInfo f1.SessionInfo
		if err := json.Unmarshal([]byte(data), &sessionInfo); err != nil {
			continue
		}

		jsonData, err := json.MarshalIndent(sessionInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling session info: %w", err)
		}

		if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
			return fmt.Errorf("writing session info file: %w", err)
		}

		return nil
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading session info file: %w", err)
	}

	return fmt.Errorf("no valid session info found in %s", sessionInfoPath)
}

func parseStartTimeFromHeartbeat(heartbeatPath string) (time.Time, error) {
	file, err := os.Open(heartbeatPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("opening heartbeat file %s: %w", heartbeatPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 12 {
			continue
		}

		// Parse the JSON part
		data := line[12:]
		var heartbeat struct {
			Utc string `json:"Utc"`
		}

		if err := json.Unmarshal([]byte(data), &heartbeat); err != nil {
			continue
		}

		timeStr := line[:12]
		offset, err := parseTimeOffset(timeStr)
		if err != nil {
			continue
		}

		// Parse the heartbeat UTC time
		heartbeatTime, err := time.Parse("2006-01-02T15:04:05.9999999Z", heartbeat.Utc)
		if err != nil {
			continue
		}

		startTime := heartbeatTime.Add(-offset)
		return startTime, nil
	}

	if err := scanner.Err(); err != nil {
		return time.Time{}, fmt.Errorf("reading heartbeat file: %w", err)
	}

	return time.Time{}, fmt.Errorf("no valid heartbeat found in %s", heartbeatPath)
}

func parseTimeOffset(timeStr string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2})\.(\d{3})`)
	matches := re.FindStringSubmatch(timeStr)
	if len(matches) != 5 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	milliseconds, _ := strconv.Atoi(matches[4])

	offset := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond

	return offset, nil
}
