package main

import (
	"fmt"
	"os"
)

type RKG struct {
	header Header
	data   []byte
}

type Header struct {
	Identifier string

	Minutes      int8
	Seconds      int8
	Milliseconds int8

	TrackID     int8
	VehicleID   int8
	CharacterID int8

	Year  int
	Month int
	Day   int
}

func main() {
	filename := "./2m02s964.rkg"

	data, err := os.ReadFile(filename)
	if err != nil {
		println(err.Error())
	}

	rkgHeader := &Header{}
	header := data[:0x88]

	rkgHeader.Identifier = string(header[:0x04])

	b4 := header[0x04]
	b5 := header[0x05]
	b6 := header[0x06]
	minutes_seconds_ms := int32(b4)<<16 | int32(b5)<<8 | int32(b6)

	rkgHeader.Minutes = int8(minutes_seconds_ms >> 17)
	rkgHeader.Seconds = int8((minutes_seconds_ms >> 10) & 0x7f) // 0x7f = 01111111
	rkgHeader.Milliseconds = int8(minutes_seconds_ms & 0x3ff)   // 0x3ff = 001111111111

	b7 := header[0x07]

	rkgHeader.TrackID = int8(b7 >> 2)

	b8 := header[0x08]

	rkgHeader.VehicleID = int8(b8 >> 2)

	b9 := header[0x08]
	b10 := header[0x09]
	vehicleID_characterID := int16(b9)<<8 | int16(b10)

	rkgHeader.CharacterID = int8((vehicleID_characterID >> 4) & 0x3F)

	b11 := header[0x09]
	b12 := header[0x0a]
	b13 := header[0x0b]

	year_month_day := int32(b11)<<16 | int32(b12)<<8 | int32(b13)

	rkgHeader.Year = 2000 + int((year_month_day>>13)&0x7f)
	rkgHeader.Month = int((year_month_day >> 9) & 0x0f)
	rkgHeader.Day = int((year_month_day >> 4) & 0x1f)
	//00101101 00010011 00001000
	fmt.Printf("%b\n", header[0x0a:0x0d])
	fmt.Printf("%b\n", (year_month_day>>13)&0x7f)
	fmt.Printf("%b\n", year_month_day)

	fmt.Println(rkgHeader)
}
