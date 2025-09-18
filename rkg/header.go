package rkg

type Header struct {
	Identifier string `json:"identifier,omitempty"`

	FinishTime *RaceTime `json:"finishTime,omitempty"`

	TrackID     int `json:"trackId,omitempty"`
	VehicleID   int `json:"vehicleId,omitempty"`
	CharacterID int `json:"characterId,omitempty"`

	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`

	ControllerID int `json:"controllerId,omitempty"`
	Compressed   int `json:"compressed,omitempty"`
	GhostType    int `json:"ghostType,omitempty"`
	DriftType    int `json:"driftType,omitempty"`

	DataLength int `json:"dataLength,omitempty"`
	LapCount   int `json:"lapCount,omitempty"`

	Laps []*RaceTime `json:"laps,omitempty"`

	CountryCode  int `json:"countryCode,omitempty"`
	StateCode    int `json:"stateCode,omitempty"`
	LocationCode int `json:"locationCode,omitempty"`
}

type ReadableHeader struct {
	Identifier string `json:"identifier,omitempty"`

	FinishTime *RaceTime `json:"finishTime,omitempty"`

	Track     string `json:"track,omitempty"`
	Vehicle   string `json:"vehicle,omitempty"`
	Character string `json:"character,omitempty"`

	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`

	Controller string `json:"controller,omitempty"`
	Compressed bool   `json:"compressed,omitempty"`
	GhostType  string `json:"ghostType,omitempty"`
	DriftType  string `json:"driftType,omitempty"`

	DataLength int `json:"dataLength,omitempty"`
	LapCount   int `json:"lapCount,omitempty"`

	Laps []*RaceTime `json:"laps,omitempty"`

	Country string `json:"country,omitempty"`
}

func ConvertHeader(header Header) *ReadableHeader {
	readableHeader := &ReadableHeader{
		Identifier: header.Identifier,
		FinishTime: header.FinishTime,
		Track:      GetStringValue(header.TrackID, TrackIDs),
		Vehicle:    GetStringValue(header.VehicleID, VehicleIDs),
		Character:  GetStringValue(header.CharacterID, CharacterIDs),
		Year:       header.Year,
		Month:      header.Month,
		Day:        header.Day,
		Controller: GetStringValue(header.ControllerID, ControllerIDs),
		Compressed: GetBoolValue(header.Compressed, Compressed),
		GhostType:  GetStringValue(header.GhostType, GhostTypes),
		DriftType:  GetStringValue(header.DriftType, DriftTypes),
		DataLength: header.DataLength,
		LapCount:   header.LapCount,
		Laps:       header.Laps,
		Country:    GetStringValue(header.CountryCode, CountryCodes),
	}

	return readableHeader
}

func ParseLaps(rkgHeader *Header, header []byte) {
	for i := 0x11; i+2 < 0x11+(rkgHeader.LapCount*3); i = i + 3 {
		lap := &RaceTime{}

		b1 := header[i]
		b2 := header[i+1]
		b3 := header[i+2]

		lap.Minutes, lap.Seconds,
			lap.Milliseconds = ParseTime(b1, b2, b3)

		rkgHeader.Laps = append(rkgHeader.Laps, lap)
	}
}

func ParseTime(minutes byte, seconds byte, ms byte) (int, int, int) {
	minutes_seconds_ms := int32(minutes)<<16 | int32(seconds)<<8 | int32(ms)

	min := int(minutes_seconds_ms >> 17)
	sec := int((minutes_seconds_ms >> 10) & 0x7f) // 0x7f = 01111111
	millis := int(minutes_seconds_ms & 0x3ff)     // 0x3ff = 001111111111

	return min, sec, millis
}
