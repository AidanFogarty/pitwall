package logmerge

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type fileReader struct {
	scanner   *bufio.Scanner
	timestamp time.Time
	line      string
	topic     string
}

type fileHeap struct {
	readers []*fileReader
	writer  *bufio.Writer
}

func (h *fileHeap) Len() int {
	return len(h.readers)
}

func (h *fileHeap) Less(i, j int) bool {
	return h.readers[i].timestamp.Before(h.readers[j].timestamp)
}

func (h *fileHeap) Pop() any {
	old := h.readers
	n := len(old)
	item := old[n-1]
	h.readers = old[:n-1]

	return item
}

func (h *fileHeap) Push(val any) {
	reader := val.(*fileReader)
	h.readers = append(h.readers, reader)
}

func (h *fileHeap) Swap(i, j int) {
	h.readers[i], h.readers[j] = h.readers[j], h.readers[i]
}

func parseTimestamp(line string, startTime time.Time) (time.Time, error) {
	if len(line) < 12 {
		return time.Time{}, fmt.Errorf("line too short: %s", line)
	}

	timeStr := line[:12]

	// Parse using regex: HH:MM:SS.mmm
	re := regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2})\.(\d{3})`)
	matches := re.FindStringSubmatch(timeStr)
	if len(matches) != 5 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	milliseconds, _ := strconv.Atoi(matches[4])

	offset := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond

	return startTime.Add(offset), nil
}

type OutputLine struct {
	Offset    int64  `json:"offset"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Data      any    `json:"data"`
}

func merge(scanners []*bufio.Scanner, topics []string, writer *bufio.Writer, startTime time.Time) error {
	h := &fileHeap{}
	heap.Init(h)

	for i, scanner := range scanners {
		if scanner.Scan() {
			line := scanner.Text()

			timestamp, err := parseTimestamp(line, startTime)
			if err != nil {
				continue
			}

			reader := &fileReader{
				scanner:   scanner,
				timestamp: timestamp,
				line:      line,
				topic:     topics[i],
			}
			heap.Push(h, reader)
		}
	}

	for h.Len() > 0 {
		reader := heap.Pop(h).(*fileReader)

		outputLine, err := createOutputLine(reader.line, reader.topic, reader.timestamp, startTime)
		if err != nil {
			if reader.scanner.Scan() {
				line := reader.scanner.Text()
				timestamp, err := parseTimestamp(line, startTime)
				if err == nil {
					reader.line = line
					reader.timestamp = timestamp
					heap.Push(h, reader)
				}
			}
			continue
		}

		jsonBytes, err := json.Marshal(outputLine)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		if _, err := writer.WriteString(string(jsonBytes) + "\n"); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}

		if reader.scanner.Scan() {
			line := reader.scanner.Text()

			timestamp, err := parseTimestamp(line, startTime)
			if err != nil {
				continue
			}

			reader.line = line
			reader.timestamp = timestamp

			heap.Push(h, reader)
		}
	}
	return writer.Flush()
}

func createOutputLine(line, topic string, timestamp time.Time, startTime time.Time) (*OutputLine, error) {
	if len(line) < 12 {
		return nil, fmt.Errorf("line too short")
	}
	dataStr := line[12:]

	var data any

	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		// If it's not valid JSON, treat it as raw string (likely compressed data)
		data = strings.Trim(dataStr, "\"")
	}

	offsetMs := timestamp.Sub(startTime).Milliseconds()

	return &OutputLine{
		Offset:    offsetMs,
		Type:      topic,
		Data:      data,
		Timestamp: timestamp.Format(time.RFC3339Nano),
	}, nil
}

func Merge(srcPaths []string, outputPath string, startTime time.Time) error {
	var scanners []*bufio.Scanner
	var topics []string

	for _, path := range srcPaths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		filename := filepath.Base(path)
		topic := strings.TrimSuffix(filename, ".jsonStream")
		topics = append(topics, topic)

		scanners = append(scanners, bufio.NewScanner(file))
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	return merge(scanners, topics, writer, startTime)
}
