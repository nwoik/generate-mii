package main

import (
	"fmt"
	"os"
)

type RKG struct {
	header *Header
	data   []byte
}

type Header struct {
	Identifier string

	FinishTime *RaceTime

	TrackID     int
	VehicleID   int
	CharacterID int

	Year  int
	Month int
	Day   int

	ControllerID int
	Compressed   int
	GhostType    int
	DriftType    int

	DataLength int
	LapCount   int

	Laps []*RaceTime

	CountryCode  int
	StateCode    int
	LocationCode int
}

type RaceTime struct {
	Minutes      int
	Seconds      int
	Milliseconds int
}

func main() {
	filename := "./decomp.rkg"

	file, err := os.ReadFile(filename)
	if err != nil {
		println(err.Error())
	}

	rkg := &RKG{}
	rkgHeader := &Header{}
	header := file[:0x88]
	data := file[0x88:]

	rkg.header = rkgHeader
	rkg.data = data

	rkgHeader.Identifier = string(header[:0x04])

	b4 := header[0x04]
	b5 := header[0x05]
	b6 := header[0x06]

	finshTime := &RaceTime{}
	finshTime.Minutes, finshTime.Seconds,
		finshTime.Milliseconds = parseTime(b4, b5, b6)
	rkgHeader.FinishTime = finshTime

	b7 := header[0x07]

	rkgHeader.TrackID = int(b7 >> 2)

	b8 := header[0x08]

	rkgHeader.VehicleID = int(b8 >> 2)

	b9 := header[0x08]
	b10 := header[0x09]
	vehicleID_characterID := int16(b9)<<8 | int16(b10)

	rkgHeader.CharacterID = int((vehicleID_characterID >> 4) & 0x3F)

	b11 := header[0x09]
	b12 := header[0x0a]
	b13 := header[0x0b]

	year_month_day := int32(b11)<<16 | int32(b12)<<8 | int32(b13)

	rkgHeader.Year = 2000 + int((year_month_day>>13)&0x7f)
	rkgHeader.Month = int((year_month_day >> 9) & 0x0f)
	rkgHeader.Day = int((year_month_day >> 4) & 0x1f)

	rkgHeader.ControllerID = int(int8(b13) & 0x0f)

	b14 := header[0x0c]

	rkgHeader.Compressed = int((int8(b14) >> 3) & 0x01)

	b15 := header[0x0d]
	ghost_drift_type := int16(b14)<<8 | int16(b15)

	rkgHeader.GhostType = int((ghost_drift_type >> 2) & 0x7f)
	rkgHeader.DriftType = int((ghost_drift_type >> 1) & 0x01)

	b16 := header[0x0e]
	b17 := header[0x0f]
	input_data_length := int16(b16)<<8 | int16(b17)

	rkgHeader.DataLength = int(input_data_length)

	b18 := header[0x10]

	rkgHeader.LapCount = int(b18)

	parseLaps(rkgHeader, header)

	b19 := header[0x34]

	rkgHeader.CountryCode = int(int8(b19))

	b20 := header[0x35]

	rkgHeader.StateCode = int(int8(b20))

	b21 := header[0x36]
	b22 := header[0x37]

	rkgHeader.LocationCode = int(int16(b21)<<8 | int16(b22))

	fmt.Println(rkg)
}

func parseLaps(rkgHeader *Header, header []byte) {
	for i := 0x11; i < 0x11+(rkgHeader.LapCount*3); i = i + 3 {
		lap := &RaceTime{}

		b1 := header[i]
		b2 := header[i+1]
		b3 := header[i+2]

		lap.Minutes, lap.Seconds,
			lap.Milliseconds = parseTime(b1, b2, b3)

		rkgHeader.Laps = append(rkgHeader.Laps, lap)
	}
}

func parseTime(minutes byte, seconds byte, ms byte) (int, int, int) {
	minutes_seconds_ms := int32(minutes)<<16 | int32(seconds)<<8 | int32(ms)

	min := int(minutes_seconds_ms >> 17)
	sec := int((minutes_seconds_ms >> 10) & 0x7f) // 0x7f = 01111111
	millis := int(minutes_seconds_ms & 0x3ff)     // 0x3ff = 001111111111

	return min, sec, millis
}
