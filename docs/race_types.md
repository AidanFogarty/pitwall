# F1 Live Timing API Event Types

This doc outlines the available the json stream events and their meaning

## SessionInfo

Just gives a brief overview of the session information, such as the circuit name, country etc. Usually gets streamed 3 times, each with a different session status, inactive -> started -> finalized

```json
{
  "offset": 0,
  "timestamp": "2025-10-05T11:08:36.6258196Z",
  "type": "SessionInfo",
  "data": {
    "ArchiveStatus": {
      "Status": "Generating"
    },
    "EndDate": "2025-10-05T22:00:00",
    "GmtOffset": "08:00:00",
    "Key": 9896,
    "Meeting": {
      "Circuit": {
        "Key": 61,
        "ShortName": "Singapore"
      },
      "Country": {
        "Code": "SGP",
        "Key": 157,
        "Name": "Singapore"
      },
      "Key": 1270,
      "Location": "Marina Bay",
      "Name": "Singapore Grand Prix",
      "Number": 18,
      "OfficialName": "FORMULA 1 SINGAPORE AIRLINES SINGAPORE GRAND PRIX 2025"
    },
    "Name": "Race",
    "Path": "2025/2025-10-05_Singapore_Grand_Prix/2025-10-05_Race/",
    "SessionStatus": "Inactive",
    "StartDate": "2025-10-05T20:00:00",
    "Type": "Race"
  }
}
```

## TrackStatus

Simply just outlines whether the track is all clear or is in a yellow flag state. Status 1 = AllClear, Status 2 = Yellow. No indiciation on which sector the status is applied to.

```json
{
  "offset": 2064,
  "timestamp": "2025-10-05T11:08:38.6898196Z",
  "type": "TrackStatus",
  "data": {
    "Message": "Yellow",
    "Status": "2"
  }
}
```

```json
{
  "offset": 293525,
  "timestamp": "2025-10-05T11:13:30.1508196Z",
  "type": "TrackStatus",
  "data": {
    "Message": "AllClear",
    "Status": "1"
  }
}
```

## SessionData

Session data seems to hold info such as track status information, similar to the TrackStatus, but under nested fields called StatusSeries, which I am not sure on the meaning. Contains info the current lap as well. When Series is empty, its just an empty array, but the series number seems to increment each time this event is streamed.

There is event with Series for each lap for this event type.

```json
{
  "offset": 2064,
  "timestamp": "2025-10-05T11:08:38.6898196Z",
  "type": "SessionData",
  "data": {
    "Series": [],
    "StatusSeries": [
      {
        "TrackStatus": "Yellow",
        "Utc": "2025-10-05T11:08:36.526Z"
      }
    ]
  }
}
```

```json
{
  "offset": 4455,
  "timestamp": "2025-10-05T11:08:41.0808196Z",
  "type": "SessionData",
  "data": {
    "Series": {
      "0": {
        "Lap": 1,
        "Utc": "2025-10-05T11:08:38.917Z"
      }
    }
  }
}
```

```json
{
  "offset": 293525,
  "timestamp": "2025-10-05T11:13:30.1508196Z",
  "type": "SessionData",
  "data": {
    "StatusSeries": {
      "1": {
        "TrackStatus": "AllClear",
        "Utc": "2025-10-05T11:13:27.987Z"
      }
    }
  }
}
```

```json
{
  "offset": 3297735,
  "timestamp": "2025-10-05T12:03:34.3608196Z",
  "type": "SessionData",
  "data": {
    "StatusSeries": {
      "10": {
        "SessionStatus": "Started",
        "Utc": "2025-10-05T12:03:32.197Z"
      }
    }
  }
}
```

```json
{
  "offset": 3401836,
  "timestamp": "2025-10-05T12:05:18.4618196Z",
  "type": "SessionData",
  "data": {
    "Series": {
      "1": {
        "Lap": 2,
        "Utc": "2025-10-05T12:05:16.298Z"
      }
    }
  }
}
```

## ExtrapolatedClock

Not too sure on what this is.

```json
{
  "offset": 2549,
  "timestamp": "2025-10-05T11:08:39.1748196Z",
  "type": "ExtrapolatedClock",
  "data": {
    "Extrapolating": false,
    "Remaining": "00:00:00",
    "Utc": "2025-10-05T11:08:37.011Z"
  }
}
```

```json
{
  "offset": 7549,
  "timestamp": "2025-10-05T11:08:44.1748196Z",
  "type": "ExtrapolatedClock",
  "data": {
    "Remaining": "02:00:00",
    "Utc": "2025-10-05T11:08:42.011Z"
  }
}
```

```json
{
  "offset": 3298536,
  "timestamp": "2025-10-05T12:03:35.1618196Z",
  "type": "ExtrapolatedClock",
  "data": {
    "Extrapolating": true,
    "Remaining": "01:59:59",
    "Utc": "2025-10-05T12:03:32.998Z"
  }
}
```

## Position.z

This just contains position data about the cars. Needs to processed as the data is encoded.

```json
{
  "offset": 3029,
  "timestamp": "2025-10-05T11:08:39.6548196Z",
  "type": "Position.z",
  "data": "7ZWxDoIwEIbf5WYkvba0pbuzJjIoxoEYBmIAA3UifXfRF7A3yXDLJU2+4e/dfbkFjuPchW4cwF8XqLq+nUPTP8GDFLLYodiJokL0wnllcqtQOWVqyGA/hKlrZ/AL4KecQhNe6xMOQzU198eKnMGLDC7fWq81ZqDSUZ2OoiCwhLRIyWAIrEtnJeFvUhJYwigkoQ/SEraBMAtNWR1C3qJIZw0hgyX0wSX3Icbst6VliUprZEvZUrZ0o5baHFdLUVm2lC1lSzdrqS5RGse3lC1lS/9k6S2+AQ=="
}
```

## LapCount

Thankfully an easy one, just tracks the current lap versus the total laps

```json
{
  "offset": 4455,
  "timestamp": "2025-10-05T11:08:41.0808196Z",
  "type": "LapCount",
  "data": {
    "CurrentLap": 1,
    "TotalLaps": 62
  }
}
```

## TimingAppData

From initial look, it tracks the driver position. It has a nested object `Lines`, which has map of nested objects where the key is a string number, which actually maps to the driver number. In the first example, you can see the number "1", which maps to max verstappen at the time of writing, and in that object it outlines his starting grid position, which is same as the `Line`, and then his racing number which will always match the key.

This event type is streamed quite a lot, 1000ish times across a race session.

```json
{
  "offset": 6424,
  "timestamp": "2025-10-05T11:08:43.0498196Z",
  "type": "TimingAppData",
  "data": {
    "Lines": {
      "1": {
        "GridPos": "2",
        "Line": 2,
        "RacingNumber": "1"
      },
      "4": {
        "GridPos": "5",
        "Line": 5,
        "RacingNumber": "4"
      },
      "5": {
        "GridPos": "14",
        "Line": 14,
        "RacingNumber": "5"
      },
      "6": {
        "GridPos": "8",
        "Line": 8,
        "RacingNumber": "6"
      },
      "10": {
        "GridPos": "19",
        "Line": 19,
        "RacingNumber": "10"
      },
      "12": {
        "GridPos": "4",
        "Line": 4,
        "RacingNumber": "12"
      },
      "14": {
        "GridPos": "10",
        "Line": 10,
        "RacingNumber": "14"
      },
      "16": {
        "GridPos": "7",
        "Line": 7,
        "RacingNumber": "16"
      },
      "18": {
        "GridPos": "15",
        "Line": 15,
        "RacingNumber": "18"
      },
      "22": {
        "GridPos": "13",
        "Line": 13,
        "RacingNumber": "22"
      },
      "23": {
        "GridPos": "20",
        "Line": 20,
        "RacingNumber": "23"
      },
      "27": {
        "GridPos": "11",
        "Line": 11,
        "RacingNumber": "27"
      },
      "30": {
        "GridPos": "12",
        "Line": 12,
        "RacingNumber": "30"
      },
      "31": {
        "GridPos": "17",
        "Line": 17,
        "RacingNumber": "31"
      },
      "43": {
        "GridPos": "16",
        "Line": 16,
        "RacingNumber": "43"
      },
      "44": {
        "GridPos": "6",
        "Line": 6,
        "RacingNumber": "44"
      },
      "55": {
        "GridPos": "18",
        "Line": 18,
        "RacingNumber": "55"
      },
      "63": {
        "GridPos": "1",
        "Line": 1,
        "RacingNumber": "63"
      },
      "81": {
        "GridPos": "3",
        "Line": 3,
        "RacingNumber": "81"
      },
      "87": {
        "GridPos": "9",
        "Line": 9,
        "RacingNumber": "87"
      }
    }
  }
}
```

This event type also seems have another nested field "Stints", which seems to track the driver tyre stint, its type, number of laps, whether its a new tyre etc. I think `LapFlags` is the number of laps under a yellow flag that the tyre completed. Some stints are empty.

```json
{
  "offset": 2790461,
  "timestamp": "2025-10-05T11:55:07.0868196Z",
  "type": "TimingAppData",
  "data": {
    "Lines": {
      "1": {
        "Stints": [
          {
            "Compound": "SOFT",
            "LapFlags": 0,
            "New": "false",
            "StartLaps": 3,
            "TotalLaps": 3,
            "TyresNotChanged": "0"
          }
        ]
      },
      "4": {
        "Stints": []
      },
      "5": {
        "Stints": []
      },
      "6": {
        "Stints": []
      },
      "10": {
        "Stints": [
          {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        ]
      },
      "12": {
        "Stints": []
      },
      "14": {
        "Stints": []
      },
      "16": {
        "Stints": []
      },
      "18": {
        "Stints": []
      },
      "22": {
        "Stints": [
          {
            "Compound": "SOFT",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        ]
      },
      "23": {
        "Stints": [
          {
            "Compound": "HARD",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        ]
      },
      "27": {
        "Stints": []
      },
      "30": {
        "Stints": []
      },
      "31": {
        "Stints": []
      },
      "43": {
        "Stints": [
          {
            "Compound": "SOFT",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        ]
      },
      "44": {
        "Stints": []
      },
      "55": {
        "Stints": [
          {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        ]
      },
      "63": {
        "Stints": []
      },
      "81": {
        "Stints": []
      },
      "87": {
        "Stints": []
      }
    }
  }
}
```

The stints also can have a nested key, which in the below are zero, which I am understanding as the first stint.

```json
{
  "offset": 2795554,
  "timestamp": "2025-10-05T11:55:12.1798196Z",
  "type": "TimingAppData",
  "data": {
    "Lines": {
      "4": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "5": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "6": {
        "Stints": {
          "0": {
            "Compound": "SOFT",
            "LapFlags": 0,
            "New": "false",
            "StartLaps": 2,
            "TotalLaps": 2,
            "TyresNotChanged": "0"
          }
        }
      },
      "12": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "14": {
        "Stints": {
          "0": {
            "Compound": "SOFT",
            "LapFlags": 0,
            "New": "false",
            "StartLaps": 3,
            "TotalLaps": 3,
            "TyresNotChanged": "0"
          }
        }
      },
      "16": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "18": {
        "Stints": {
          "0": {
            "Compound": "SOFT",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "23": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM"
          }
        }
      },
      "27": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "30": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "31": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "44": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "63": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "81": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      },
      "87": {
        "Stints": {
          "0": {
            "Compound": "MEDIUM",
            "LapFlags": 0,
            "New": "true",
            "StartLaps": 0,
            "TotalLaps": 0,
            "TyresNotChanged": "0"
          }
        }
      }
    }
  }
}
```

Then there is cases where there is way less data, I am thinking this might be an overtake. Below, Driver number 4, moved to position 3, starting from position 5 which outlined above. Which I think is correct because that is what happened in the race.

```json
{
  "offset": 3331586,
  "timestamp": "2025-10-05T12:04:08.2118196Z",
  "type": "TimingAppData",
  "data": {
    "Lines": {
      "4": {
        "Line": 3
      },
      "12": {
        "Line": 5
      },
      "81": {
        "Line": 4
      }
    }
  }
}
```

```json
{
  "offset": 3405609,
  "timestamp": "2025-10-05T12:05:22.2348196Z",
  "type": "TimingAppData",
  "data": {
    "Lines": {
      "1": {
        "Stints": {
          "0": {
            "TotalLaps": 4
          }
        }
      },
      "4": {
        "Stints": {
          "0": {
            "TotalLaps": 1
          }
        }
      },
      "16": {
        "Stints": {
          "0": {
            "TotalLaps": 1
          }
        }
      },
      "63": {
        "Stints": {
          "0": {
            "TotalLaps": 1
          }
        }
      },
      "81": {
        "Stints": {
          "0": {
            "TotalLaps": 1
          }
        }
      }
    }
  }
}
```

## TimingStats

This is one of the most important events. Useful for tracking fastest laps, best sectors, speed trap values

```json
{
  "offset": 6424,
  "timestamp": "2025-10-05T11:08:43.0498196Z",
  "type": "TimingStats",
  "data": {
    "Lines": {
      "1": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 1,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "1"
      },
      "4": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 2,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "4"
      },
      "5": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 3,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "5"
      },
      "6": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 4,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "6"
      },
      "10": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 5,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "10"
      },
      "12": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 6,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "12"
      },
      "14": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 7,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "14"
      },
      "16": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 8,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "16"
      },
      "18": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 9,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "18"
      },
      "22": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 10,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "22"
      },
      "23": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 11,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "23"
      },
      "27": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 12,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "27"
      },
      "30": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 13,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "30"
      },
      "31": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 14,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "31"
      },
      "43": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 15,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "43"
      },
      "44": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 16,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "44"
      },
      "55": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 17,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "55"
      },
      "63": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 18,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "63"
      },
      "81": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 19,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "81"
      },
      "87": {
        "BestSectors": [
          {
            "Value": ""
          },
          {
            "Value": ""
          },
          {
            "Value": ""
          }
        ],
        "BestSpeeds": {
          "FL": {
            "Value": ""
          },
          "I1": {
            "Value": ""
          },
          "I2": {
            "Value": ""
          },
          "ST": {
            "Value": ""
          }
        },
        "Line": 20,
        "PersonalBestLapTime": {
          "Value": ""
        },
        "RacingNumber": "87"
      }
    },
    "SessionType": "Race",
    "Withheld": false
  }
}
```

```json
{
  "offset": 3303237,
  "timestamp": "2025-10-05T12:03:39.8628196Z",
  "type": "TimingStats",
  "data": {
    "Lines": {
      "63": {
        "BestSpeeds": {
          "ST": {
            "Position": 1,
            "Value": "185"
          }
        }
      }
    }
  }
}
```

```json
{
  "offset": 3303585,
  "timestamp": "2025-10-05T12:03:40.2108196Z",
  "type": "TimingStats",
  "data": {
    "Lines": {
      "1": {
        "BestSpeeds": {
          "ST": {
            "Position": 1,
            "Value": "190"
          }
        }
      },
      "63": {
        "BestSpeeds": {
          "ST": {
            "Position": 2
          }
        }
      }
    }
  }
}
```

```json
{
  "offset": 3303595,
  "timestamp": "2025-10-05T12:03:40.2208196Z",
  "type": "TimingStats",
  "data": {
    "Lines": {
      "1": {
        "BestSpeeds": {
          "ST": {
            "Position": 2
          }
        }
      },
      "63": {
        "BestSpeeds": {
          "ST": {
            "Position": 3
          }
        }
      },
      "81": {
        "BestSpeeds": {
          "ST": {
            "Position": 1,
            "Value": "194"
          }
        }
      }
    }
  }
}
```

## TimingData

TimingData tracks things such as the interval to the leader, interval to the next driver. The first TimingData is usually zero values. Subsequent events usually have gaps for a subset of drivers. Probably the most common event, 60k+ occurances.

```json
{
  "offset": 6424,
  "timestamp": "2025-10-05T11:08:43.0498196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "1": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 2,
        "PitOut": false,
        "Position": "2",
        "RacingNumber": "1",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "4": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 5,
        "PitOut": false,
        "Position": "5",
        "RacingNumber": "4",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "5": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 14,
        "PitOut": false,
        "Position": "14",
        "RacingNumber": "5",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "6": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 8,
        "PitOut": false,
        "Position": "8",
        "RacingNumber": "6",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "10": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 19,
        "PitOut": false,
        "Position": "19",
        "RacingNumber": "10",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "12": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 4,
        "PitOut": false,
        "Position": "4",
        "RacingNumber": "12",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "14": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 10,
        "PitOut": false,
        "Position": "10",
        "RacingNumber": "14",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "16": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 7,
        "PitOut": false,
        "Position": "7",
        "RacingNumber": "16",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "18": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 15,
        "PitOut": false,
        "Position": "15",
        "RacingNumber": "18",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "22": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 13,
        "PitOut": false,
        "Position": "13",
        "RacingNumber": "22",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "23": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 20,
        "PitOut": false,
        "Position": "20",
        "RacingNumber": "23",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "27": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 11,
        "PitOut": false,
        "Position": "11",
        "RacingNumber": "27",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "30": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 12,
        "PitOut": false,
        "Position": "12",
        "RacingNumber": "30",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "31": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 17,
        "PitOut": false,
        "Position": "17",
        "RacingNumber": "31",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "43": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 16,
        "PitOut": false,
        "Position": "16",
        "RacingNumber": "43",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "44": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 6,
        "PitOut": false,
        "Position": "6",
        "RacingNumber": "44",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "55": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 18,
        "PitOut": false,
        "Position": "18",
        "RacingNumber": "55",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "63": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "LAP 1",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": "LAP 1"
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 1,
        "PitOut": false,
        "Position": "1",
        "RacingNumber": "63",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "81": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 3,
        "PitOut": false,
        "Position": "3",
        "RacingNumber": "81",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      },
      "87": {
        "BestLapTime": {
          "Value": ""
        },
        "GapToLeader": "",
        "InPit": true,
        "IntervalToPositionAhead": {
          "Catching": false,
          "Value": ""
        },
        "LastLapTime": {
          "OverallFastest": false,
          "PersonalFastest": false,
          "Status": 0,
          "Value": ""
        },
        "Line": 9,
        "PitOut": false,
        "Position": "9",
        "RacingNumber": "87",
        "Retired": false,
        "Sectors": [
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          },
          {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Stopped": false,
            "Value": ""
          }
        ],
        "ShowPosition": true,
        "Speeds": {
          "FL": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I1": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "I2": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          },
          "ST": {
            "OverallFastest": false,
            "PersonalFastest": false,
            "Status": 0,
            "Value": ""
          }
        },
        "Status": 80,
        "Stopped": false
      }
    },
    "Withheld": false
  }
}
```

```json
{
  "offset": 689183,
  "timestamp": "2025-10-05T11:20:05.8088196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "12": {
        "InPit": false,
        "Status": 64
      }
    }
  }
}
```

```json
{"offset":702463,"timestamp":"2025-10-05T11:20:19.0888196Z","type":"TimingData","data":{"Lines":{"12":{"Sectors":{"0":{"Segments":{"0":{"Status":2064}}}}}}}}
{"offset":702495,"timestamp":"2025-10-05T11:20:19.1208196Z","type":"TimingData","data":{"Lines":{"12":{"Sectors":{"0":{"Segments":{"1":{"Status":2064}}}}}}}}
{"offset":702495,"timestamp":"2025-10-05T11:20:19.1208196Z","type":"TimingData","data":{"Lines":{"12":{"Sectors":{"0":{"Segments":{"2":{"Status":2064}}}}}}}}
```

Example of the diff between drivers

```json
{
  "offset": 3302942,
  "timestamp": "2025-10-05T12:03:39.5678196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "1": {
        "GapToLeader": "+0.240",
        "IntervalToPositionAhead": {
          "Value": "+0.240"
        }
      }
    }
  }
}
```

```json
{
  "offset": 3303264,
  "timestamp": "2025-10-05T12:03:39.8898196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "12": {
        "GapToLeader": "+0.635",
        "IntervalToPositionAhead": {
          "Value": "+0.267"
        }
      }
    }
  }
}
```

```json
{
  "offset": 3303393,
  "timestamp": "2025-10-05T12:03:40.0188196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "63": {
        "Speeds": {
          "ST": {
            "OverallFastest": true,
            "PersonalFastest": true,
            "Value": "185"
          }
        }
      }
    }
  }
}
```

```json
{
  "offset": 3304278,
  "timestamp": "2025-10-05T12:03:40.9038196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "4": {
        "Speeds": {
          "ST": {
            "OverallFastest": false
          }
        }
      },
      "16": {
        "Speeds": {
          "ST": {
            "OverallFastest": true,
            "PersonalFastest": true,
            "Value": "215"
          }
        }
      }
    }
  }
}
```

```json
{
  "offset": 9112755,
  "timestamp": "2025-10-05T13:40:29.3808196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "14": {
        "LastLapTime": {
          "Value": "1:36.004"
        },
        "NumberOfLaps": 59,
        "Sectors": {
          "2": {
            "Value": "26.737"
          }
        },
        "Speeds": {
          "FL": {
            "Value": "255"
          }
        }
      },
      "81": {
        "Sectors": {
          "1": {
            "PreviousValue": "40.579"
          }
        }
      }
    }
  }
}
```

```json
{
  "offset": 3373007,
  "timestamp": "2025-10-05T12:04:49.6328196Z",
  "type": "TimingData",
  "data": {
    "Lines": {
      "63": {
        "Sectors": {
          "1": {
            "OverallFastest": true,
            "PersonalFastest": true,
            "Value": "42.718"
          }
        },
        "Speeds": {
          "I2": {
            "OverallFastest": true,
            "PersonalFastest": true,
            "Value": "258"
          }
        }
      }
    }
  }
}
```

## TopThree

Similar to TimingData but with extra information about the top three. Can be used to populate the first row with diff to ahead of the lap number.

```json
{
  "offset": 6424,
  "timestamp": "2025-10-05T11:08:43.0498196Z",
  "type": "TopThree",
  "data": {
    "Lines": [
      {
        "BroadcastName": "G RUSSELL",
        "DiffToAhead": "LAP 1",
        "DiffToLeader": "LAP 1",
        "FirstName": "George",
        "FullName": "George RUSSELL",
        "LapState": 80,
        "LapTime": "",
        "LastName": "Russell",
        "OverallFastest": false,
        "PersonalFastest": false,
        "Position": "1",
        "RacingNumber": "63",
        "Reference": "GEORUS01",
        "ShowPosition": true,
        "Team": "Mercedes",
        "TeamColour": "00D7B6",
        "Tla": "RUS"
      },
      {
        "BroadcastName": "M VERSTAPPEN",
        "DiffToAhead": "",
        "DiffToLeader": "",
        "FirstName": "Max",
        "FullName": "Max VERSTAPPEN",
        "LapState": 80,
        "LapTime": "",
        "LastName": "Verstappen",
        "OverallFastest": false,
        "PersonalFastest": false,
        "Position": "2",
        "RacingNumber": "1",
        "Reference": "MAXVER01",
        "ShowPosition": true,
        "Team": "Red Bull Racing",
        "TeamColour": "4781D7",
        "Tla": "VER"
      },
      {
        "BroadcastName": "O PIASTRI",
        "DiffToAhead": "",
        "DiffToLeader": "",
        "FirstName": "Oscar",
        "FullName": "Oscar PIASTRI",
        "LapState": 80,
        "LapTime": "",
        "LastName": "Piastri",
        "OverallFastest": false,
        "PersonalFastest": false,
        "Position": "3",
        "RacingNumber": "81",
        "Reference": "OSCPIA01",
        "ShowPosition": true,
        "Team": "McLaren",
        "TeamColour": "F47600",
        "Tla": "PIA"
      }
    ],
    "Withheld": false
  }
}
```

```json
{
  "offset": 695026,
  "timestamp": "2025-10-05T11:20:11.6518196Z",
  "type": "TopThree",
  "data": {
    "Lines": {
      "0": {
        "LapState": 64
      }
    }
  }
}
```

```json
{
  "offset": 3302942,
  "timestamp": "2025-10-05T12:03:39.5678196Z",
  "type": "TopThree",
  "data": {
    "Lines": {
      "1": {
        "DiffToAhead": "+0.240",
        "DiffToLeader": "+0.240"
      }
    }
  }
}
```

## DriverList

Just gives a rundown of all the drivers in the session. Can be used to build the initial list of driver positions on the timing tower. Think it also tracks whenever a change of position happens. The subsequent events seems to overlap with TimingAppData

```json
{
  "offset": 6424,
  "timestamp": "2025-10-05T11:08:43.0498196Z",
  "type": "DriverList",
  "data": {
    "1": {
      "BroadcastName": "M VERSTAPPEN",
      "FirstName": "Max",
      "FullName": "Max VERSTAPPEN",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/M/MAXVER01_Max_Verstappen/maxver01.png.transform/1col/image.png",
      "LastName": "Verstappen",
      "Line": 2,
      "PublicIdRight": "common/f1/2025/redbullracing/maxver01/2025redbullracingmaxver01right",
      "RacingNumber": "1",
      "Reference": "MAXVER01",
      "TeamColour": "4781D7",
      "TeamName": "Red Bull Racing",
      "Tla": "VER"
    },
    "4": {
      "BroadcastName": "L NORRIS",
      "FirstName": "Lando",
      "FullName": "Lando NORRIS",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/L/LANNOR01_Lando_Norris/lannor01.png.transform/1col/image.png",
      "LastName": "Norris",
      "Line": 5,
      "PublicIdRight": "common/f1/2025/mclaren/lannor01/2025mclarenlannor01right",
      "RacingNumber": "4",
      "Reference": "LANNOR01",
      "TeamColour": "F47600",
      "TeamName": "McLaren",
      "Tla": "NOR"
    },
    "5": {
      "BroadcastName": "G BORTOLETO",
      "FirstName": "Gabriel",
      "FullName": "Gabriel BORTOLETO",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/G/GABBOR01_Gabriel_Bortoleto/gabbor01.png.transform/1col/image.png",
      "LastName": "Bortoleto",
      "Line": 14,
      "PublicIdRight": "common/f1/2025/kicksauber/gabbor01/2025kicksaubergabbor01right",
      "RacingNumber": "5",
      "Reference": "GABBOR01",
      "TeamColour": "01C00E",
      "TeamName": "Kick Sauber",
      "Tla": "BOR"
    },
    "6": {
      "BroadcastName": "I HADJAR",
      "FirstName": "Isack",
      "FullName": "Isack HADJAR",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/I/ISAHAD01_Isack_Hadjar/isahad01.png.transform/1col/image.png",
      "LastName": "Hadjar",
      "Line": 8,
      "PublicIdRight": "common/f1/2025/racingbulls/isahad01/2025racingbullsisahad01right",
      "RacingNumber": "6",
      "Reference": "ISAHAD01",
      "TeamColour": "6C98FF",
      "TeamName": "Racing Bulls",
      "Tla": "HAD"
    },
    "10": {
      "BroadcastName": "P GASLY",
      "FirstName": "Pierre",
      "FullName": "Pierre GASLY",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/P/PIEGAS01_Pierre_Gasly/piegas01.png.transform/1col/image.png",
      "LastName": "Gasly",
      "Line": 19,
      "PublicIdRight": "common/f1/2025/alpine/piegas01/2025alpinepiegas01right",
      "RacingNumber": "10",
      "Reference": "PIEGAS01",
      "TeamColour": "00A1E8",
      "TeamName": "Alpine",
      "Tla": "GAS"
    },
    "12": {
      "BroadcastName": "K ANTONELLI",
      "FirstName": "Kimi",
      "FullName": "Kimi ANTONELLI",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/A/ANDANT01_Andrea%20Kimi_Antonelli/andant01.png.transform/1col/image.png",
      "LastName": "Antonelli",
      "Line": 4,
      "PublicIdRight": "common/f1/2025/mercedes/andant01/2025mercedesandant01right",
      "RacingNumber": "12",
      "Reference": "ANDANT01",
      "TeamColour": "00D7B6",
      "TeamName": "Mercedes",
      "Tla": "ANT"
    },
    "14": {
      "BroadcastName": "F ALONSO",
      "FirstName": "Fernando",
      "FullName": "Fernando ALONSO",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/F/FERALO01_Fernando_Alonso/feralo01.png.transform/1col/image.png",
      "LastName": "Alonso",
      "Line": 10,
      "PublicIdRight": "common/f1/2025/astonmartin/feralo01/2025astonmartinferalo01right",
      "RacingNumber": "14",
      "Reference": "FERALO01",
      "TeamColour": "229971",
      "TeamName": "Aston Martin",
      "Tla": "ALO"
    },
    "16": {
      "BroadcastName": "C LECLERC",
      "FirstName": "Charles",
      "FullName": "Charles LECLERC",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/C/CHALEC01_Charles_Leclerc/chalec01.png.transform/1col/image.png",
      "LastName": "Leclerc",
      "Line": 7,
      "PublicIdRight": "common/f1/2025/ferrari/chalec01/2025ferrarichalec01right",
      "RacingNumber": "16",
      "Reference": "CHALEC01",
      "TeamColour": "ED1131",
      "TeamName": "Ferrari",
      "Tla": "LEC"
    },
    "18": {
      "BroadcastName": "L STROLL",
      "FirstName": "Lance",
      "FullName": "Lance STROLL",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/L/LANSTR01_Lance_Stroll/lanstr01.png.transform/1col/image.png",
      "LastName": "Stroll",
      "Line": 15,
      "PublicIdRight": "common/f1/2025/astonmartin/lanstr01/2025astonmartinlanstr01right",
      "RacingNumber": "18",
      "Reference": "LANSTR01",
      "TeamColour": "229971",
      "TeamName": "Aston Martin",
      "Tla": "STR"
    },
    "22": {
      "BroadcastName": "Y TSUNODA",
      "FirstName": "Yuki",
      "FullName": "Yuki TSUNODA",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/Y/YUKTSU01_Yuki_Tsunoda/yuktsu01.png.transform/1col/image.png",
      "LastName": "Tsunoda",
      "Line": 13,
      "PublicIdRight": "common/f1/2025/redbullracing/yuktsu01/2025redbullracingyuktsu01right",
      "RacingNumber": "22",
      "Reference": "YUKTSU01",
      "TeamColour": "4781D7",
      "TeamName": "Red Bull Racing",
      "Tla": "TSU"
    },
    "23": {
      "BroadcastName": "A ALBON",
      "FirstName": "Alexander",
      "FullName": "Alexander ALBON",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/A/ALEALB01_Alexander_Albon/alealb01.png.transform/1col/image.png",
      "LastName": "Albon",
      "Line": 20,
      "PublicIdRight": "common/f1/2025/williams/alealb01/2025williamsalealb01right",
      "RacingNumber": "23",
      "Reference": "ALEALB01",
      "TeamColour": "1868DB",
      "TeamName": "Williams",
      "Tla": "ALB"
    },
    "27": {
      "BroadcastName": "N HULKENBERG",
      "FirstName": "Nico",
      "FullName": "Nico HULKENBERG",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/N/NICHUL01_Nico_Hulkenberg/nichul01.png.transform/1col/image.png",
      "LastName": "Hulkenberg",
      "Line": 11,
      "PublicIdRight": "common/f1/2025/kicksauber/nichul01/2025kicksaubernichul01right",
      "RacingNumber": "27",
      "Reference": "NICHUL01",
      "TeamColour": "01C00E",
      "TeamName": "Kick Sauber",
      "Tla": "HUL"
    },
    "30": {
      "BroadcastName": "L LAWSON",
      "FirstName": "Liam",
      "FullName": "Liam LAWSON",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/L/LIALAW01_Liam_Lawson/lialaw01.png.transform/1col/image.png",
      "LastName": "Lawson",
      "Line": 12,
      "PublicIdRight": "common/f1/2025/racingbulls/lialaw01/2025racingbullslialaw01right",
      "RacingNumber": "30",
      "Reference": "LIALAW01",
      "TeamColour": "6C98FF",
      "TeamName": "Racing Bulls",
      "Tla": "LAW"
    },
    "31": {
      "BroadcastName": "E OCON",
      "FirstName": "Esteban",
      "FullName": "Esteban OCON",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/E/ESTOCO01_Esteban_Ocon/estoco01.png.transform/1col/image.png",
      "LastName": "Ocon",
      "Line": 17,
      "PublicIdRight": "common/f1/2025/haas/estoco01/2025haasestoco01right",
      "RacingNumber": "31",
      "Reference": "ESTOCO01",
      "TeamColour": "9C9FA2",
      "TeamName": "Haas F1 Team",
      "Tla": "OCO"
    },
    "43": {
      "BroadcastName": "F COLAPINTO",
      "FirstName": "Franco",
      "FullName": "Franco COLAPINTO",
      "LastName": "Colapinto",
      "Line": 16,
      "PublicIdRight": "common/f1/2025/alpine/fracol01/2025alpinefracol01right",
      "RacingNumber": "43",
      "Reference": "FRACOL01",
      "TeamColour": "00A1E8",
      "TeamName": "Alpine",
      "Tla": "COL"
    },
    "44": {
      "BroadcastName": "L HAMILTON",
      "FirstName": "Lewis",
      "FullName": "Lewis HAMILTON",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/L/LEWHAM01_Lewis_Hamilton/lewham01.png.transform/1col/image.png",
      "LastName": "Hamilton",
      "Line": 6,
      "PublicIdRight": "common/f1/2025/ferrari/lewham01/2025ferrarilewham01right",
      "RacingNumber": "44",
      "Reference": "LEWHAM01",
      "TeamColour": "ED1131",
      "TeamName": "Ferrari",
      "Tla": "HAM"
    },
    "55": {
      "BroadcastName": "C SAINZ",
      "FirstName": "Carlos",
      "FullName": "Carlos SAINZ",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/C/CARSAI01_Carlos_Sainz/carsai01.png.transform/1col/image.png",
      "LastName": "Sainz",
      "Line": 18,
      "PublicIdRight": "common/f1/2025/williams/carsai01/2025williamscarsai01right",
      "RacingNumber": "55",
      "Reference": "CARSAI01",
      "TeamColour": "1868DB",
      "TeamName": "Williams",
      "Tla": "SAI"
    },
    "63": {
      "BroadcastName": "G RUSSELL",
      "FirstName": "George",
      "FullName": "George RUSSELL",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/G/GEORUS01_George_Russell/georus01.png.transform/1col/image.png",
      "LastName": "Russell",
      "Line": 1,
      "PublicIdRight": "common/f1/2025/mercedes/georus01/2025mercedesgeorus01right",
      "RacingNumber": "63",
      "Reference": "GEORUS01",
      "TeamColour": "00D7B6",
      "TeamName": "Mercedes",
      "Tla": "RUS"
    },
    "81": {
      "BroadcastName": "O PIASTRI",
      "FirstName": "Oscar",
      "FullName": "Oscar PIASTRI",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/O/OSCPIA01_Oscar_Piastri/oscpia01.png.transform/1col/image.png",
      "LastName": "Piastri",
      "Line": 3,
      "PublicIdRight": "common/f1/2025/mclaren/oscpia01/2025mclarenoscpia01right",
      "RacingNumber": "81",
      "Reference": "OSCPIA01",
      "TeamColour": "F47600",
      "TeamName": "McLaren",
      "Tla": "PIA"
    },
    "87": {
      "BroadcastName": "O BEARMAN",
      "FirstName": "Oliver",
      "FullName": "Oliver BEARMAN",
      "HeadshotUrl": "https://media.formula1.com/d_driver_fallback_image.png/content/dam/fom-website/drivers/O/OLIBEA01_Oliver_Bearman/olibea01.png.transform/1col/image.png",
      "LastName": "Bearman",
      "Line": 9,
      "PublicIdRight": "common/f1/2025/haas/olibea01/2025haasolibea01right",
      "RacingNumber": "87",
      "Reference": "OLIBEA01",
      "TeamColour": "9C9FA2",
      "TeamName": "Haas F1 Team",
      "Tla": "BEA"
    }
  }
}
```

```json
{
  "offset": 3331586,
  "timestamp": "2025-10-05T12:04:08.2118196Z",
  "type": "DriverList",
  "data": {
    "4": {
      "Line": 3
    },
    "12": {
      "Line": 5
    },
    "81": {
      "Line": 4
    }
  }
}
```

## SessionStatus

```json
{
  "offset": 8064,
  "timestamp": "2025-10-05T11:08:44.6898196Z",
  "type": "SessionStatus",
  "data": {
    "Started": "Inactive",
    "Status": "Inactive"
  }
}
```

```json
{
  "offset": 3297735,
  "timestamp": "2025-10-05T12:03:34.3608196Z",
  "type": "SessionStatus",
  "data": { "Started": "Started", "Status": "Started" }
}
```

```json
{
  "offset": 9320249,
  "timestamp": "2025-10-05T13:43:56.8748196Z",
  "type": "SessionStatus",
  "data": { "Started": "Finished", "Status": "Finished" }
}
```

```json
{
  "offset": 9640947,
  "timestamp": "2025-10-05T13:49:17.5728196Z",
  "type": "SessionStatus",
  "data": { "Started": "Finished", "Status": "Finalised" }
}
```

```json
{
  "offset": 9642852,
  "timestamp": "2025-10-05T13:49:19.4778196Z",
  "type": "SessionStatus",
  "data": { "Started": "Finished", "Status": "Ends" }
}
```

## WeatherData

Just contains info such as rainfall, windspeed etc.

```json
{
  "offset": 38510,
  "timestamp": "2025-10-05T11:09:15.1358196Z",
  "type": "WeatherData",
  "data": {
    "AirTemp": "28.5",
    "Humidity": "75.0",
    "Pressure": "1009.9",
    "Rainfall": "0",
    "TrackTemp": "34.0",
    "WindDirection": "151",
    "WindSpeed": "1.6"
  }
}
```

```json
{
  "offset": 98519,
  "timestamp": "2025-10-05T11:10:15.1448196Z",
  "type": "WeatherData",
  "data": {
    "AirTemp": "28.6",
    "Humidity": "75.0",
    "Pressure": "1009.8",
    "Rainfall": "0",
    "TrackTemp": "34.4",
    "WindDirection": "156",
    "WindSpeed": "1.2"
  }
}
```
