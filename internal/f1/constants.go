package f1

var BaseTopics = []string{
	"Heartbeat",
	"AudioStreams",
	"DriverList",
	"ExtrapolatedClock",
	"RaceControlMessages",
	"SessionInfo",
	"SessionStatus",
	"TeamRadio",
	"TimingAppData",
	"TimingStats",
	"TrackStatus",
	"WeatherData",
	"Position.z",
	"CarData.z",
	"ContentStreams",
	"SessionData",
	"TimingData",
	"TopThree",
}

var raceTopics = []string{"LapCount"}
var RaceTopics = append(BaseTopics, raceTopics...)

// For sessions that are not races, we get some other info
var AllTopics = append(BaseTopics, raceTopics...)

var NumberOfDrivers int = 20

var (
	F1SignalRURL = "https://livetiming.formula1.com/signalr"
	HubName      = "Streaming"
)
