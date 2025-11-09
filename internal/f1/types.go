package f1

import (
	"encoding/json"
	"time"
)

type F1Event struct {
	Offset    int64           `json:"offset"`
	Timestamp time.Time       `json:"timestamp"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
}

type F1EventMsg F1Event

type ArchiveStatus struct {
	Status string `json:"Status"`
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
	Key          int     `json:"Key"`
	Name         string  `json:"Name"`
	OfficialName string  `json:"OfficialName"`
	Location     string  `json:"Location"`
	Number       int     `json:"Number"`
	Country      Country `json:"Country"`
	Circuit      Circuit `json:"Circuit"`
}

type SessionInfo struct {
	Meeting       Meeting       `json:"Meeting"`
	SessionStatus string        `json:"SessionStatus"`
	ArchiveStatus ArchiveStatus `json:"ArchiveStatus"`
	Key           int           `json:"Key"`
	Type          string        `json:"Type"`
	Name          string        `json:"Name"`
	StartDate     string        `json:"StartDate"`
	EndDate       string        `json:"EndDate"`
	GmtOffset     string        `json:"GmtOffset"`
	Path          string        `json:"Path"`
}

type Driver struct {
	BroadcastName string `json:"BroadcastName,omitempty"`
	FirstName     string `json:"FirstName,omitempty"`
	FullName      string `json:"FullName,omitempty"`
	LastName      string `json:"LastName,omitempty"`
	Line          int    `json:"Line,omitempty"`
	RacingNumber  string `json:"RacingNumber,omitempty"`
	TeamColour    string `json:"TeamColour,omitempty"`
	TeamName      string `json:"TeamName,omitempty"`
	Tla           string `json:"Tla,omitempty"`
}

type DriverList map[string]Driver

type TopThreeLine struct {
	BroadcastName   string `json:"BroadcastName,omitempty"`
	DiffToAhead     string `json:"DiffToAhead,omitempty"`
	DiffToLeader    string `json:"DiffToLeader,omitempty"`
	FirstName       string `json:"FirstName,omitempty"`
	FullName        string `json:"FullName,omitempty"`
	LapTime         string `json:"LapTime,omitempty"`
	LastName        string `json:"LastName,omitempty"`
	OverallFastest  bool   `json:"OverallFastest,omitempty"`
	PersonalFastest bool   `json:"PersonalFastest,omitempty"`
	Position        string `json:"Position,omitempty"`
	RacingNumber    string `json:"RacingNumber,omitempty"`
	Team            string `json:"Team,omitempty"`
	TeamColour      string `json:"TeamColour,omitempty"`
	Tla             string `json:"Tla,omitempty"`
}

type TopThree struct {
	Lines    json.RawMessage `json:"Lines"`
	Withheld bool            `json:"Withheld,omitempty"`
}

type TimingData struct {
	Lines    map[string]TimingDataLine `json:"Lines"`
	Withheld bool                      `json:"Withheld,omitempty"`
}

type TimingDataLine struct {
	GapToLeader             string       `json:"GapToLeader,omitempty"`
	IntervalToPositionAhead IntervalData `json:"IntervalToPositionAhead"`
	Position                string       `json:"Position,omitempty"`
	RacingNumber            string       `json:"RacingNumber,omitempty"`
	Line                    int          `json:"Line,omitempty"`

	InPit   bool `json:"InPit,omitempty"`
	PitOut  bool `json:"PitOut,omitempty"`
	Retired bool `json:"Retired,omitempty"`
	Status  int  `json:"Status,omitempty"`

	BestLapTime LapTimeData          `json:"BestLapTime"`
	LastLapTime LapTimeData          `json:"LastLapTime"`
	Sectors     json.RawMessage      `json:"Sectors,omitempty"`
	Speeds      map[string]SpeedData `json:"Speeds,omitempty"`
}

type IntervalData struct {
	Catching bool   `json:"Catching,omitempty"`
	Value    string `json:"Value,omitempty"`
}

type LapTimeData struct {
	Value           string `json:"Value,omitempty"`
	Status          int    `json:"Status,omitempty"`
	OverallFastest  bool   `json:"OverallFastest,omitempty"`
	PersonalFastest bool   `json:"PersonalFastest,omitempty"`
}

type SegmentData struct {
	Status int `json:"Status,omitempty"`
}

type SectorData struct {
	Value           string          `json:"Value,omitempty"`
	Status          int             `json:"Status,omitempty"`
	OverallFastest  bool            `json:"OverallFastest,omitempty"`
	PersonalFastest bool            `json:"PersonalFastest,omitempty"`
	Stopped         bool            `json:"Stopped,omitempty"`
	Segments        json.RawMessage `json:"Segments,omitempty"`
}

type SpeedData struct {
	Value           string `json:"Value,omitempty"`
	Status          int    `json:"Status,omitempty"`
	OverallFastest  bool   `json:"OverallFastest,omitempty"`
	PersonalFastest bool   `json:"PersonalFastest,omitempty"`
}

type TimingAppData struct {
	Lines    map[string]TimingAppDataLine `json:"Lines"`
	Withheld bool                         `json:"Withheld,omitempty"`
}

type TimingAppDataLine struct {
	GridPos      string          `json:"GridPos,omitempty"`
	Line         int             `json:"Line,omitempty"`
	RacingNumber string          `json:"RacingNumber,omitempty"`
	Stints       json.RawMessage `json:"Stints,omitempty"`
}

type StintData struct {
	Compound        string `json:"Compound,omitempty"`
	LapFlags        *int   `json:"LapFlags,omitempty"`
	New             string `json:"New,omitempty"`
	StartLaps       *int   `json:"StartLaps,omitempty"`
	TotalLaps       *int   `json:"TotalLaps,omitempty"`
	TyresNotChanged string `json:"TyresNotChanged,omitempty"`
}

type RaceControlMessages struct {
	Messages json.RawMessage `json:"Messages"`
}

type RaceControlMessage struct {
	Category string `json:"Category,omitempty"`
	Flag     string `json:"Flag,omitempty"`
	Lap      int    `json:"Lap,omitempty"`
	Message  string `json:"Message,omitempty"`
	Scope    string `json:"Scope,omitempty"`
	Sector   int    `json:"Sector,omitempty"`
	Utc      string `json:"Utc,omitempty"`
}

type WeatherData struct {
	AirTemp       string `json:"AirTemp,omitempty"`
	Humidity      string `json:"Humidity,omitempty"`
	Pressure      string `json:"Pressure,omitempty"`
	Rainfall      string `json:"Rainfall,omitempty"`
	TrackTemp     string `json:"TrackTemp,omitempty"`
	WindDirection string `json:"WindDirection,omitempty"`
	WindSpeed     string `json:"WindSpeed,omitempty"`
}

type TrackStatus struct {
	Message string `json:"Message,omitempty"`
}

type LapCount struct {
	CurrentLap int `json:"CurrentLap,omitempty"`
	TotalLaps  int `json:"TotalLaps,omitempty"`
}

type SessionStatus struct {
	Status string `json:"Status,omitempty"`
}

type InitialLiveState struct {
	SessionInfo *SessionInfo `json:"SessionInfo,omitempty"`
	WeatherData *WeatherData `json:"WeatherData,omitempty"`
	TrackStatus *TrackStatus `json:"TrackStatus,omitempty"`
	LapCount    *LapCount    `json:"LapCount,omitempty"`

	DriverList    *DriverList    `json:"DriverList,omitempty"`
	TimingData    *TimingData    `json:"TimingData,omitempty"`
	TopThree      *TopThree      `json:"TopThree,omitempty"`
	TimingAppData *TimingAppData `json:"TimingAppData,omitempty"`

	RaceControlMessages *RaceControlMessages `json:"RaceControlMessages,omitempty"`
}
