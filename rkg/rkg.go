package rkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type RKG struct {
	Header *Header `json:"header,omitempty"`
	Mii    []byte  `json:"mii,omitempty"`
	CRC    []byte  `json:"crc,omitempty"`
	Data   []byte  `json:"data,omitempty"`
}

type ReadbleRKG struct {
	Header *ReadableHeader `json:"header,omitempty"`
	Mii    []byte          `json:"mii,omitempty"`
	CRC    []byte          `json:"crc,omitempty"`
	Data   []byte          `json:"data,omitempty"`
}

func ReadFile(filepath string) []byte {
	file, err := os.ReadFile(filepath)
	if err != nil {
		println(err.Error())
	}

	return file
}

func ExportToJsonRaw(filepath string) {
	rkg := ParseRKG(ReadFile(filepath))

	filename := strings.Split(filepath, ".")[0]

	jsonBytes, _ := json.MarshalIndent(rkg, "", "\t")
	WriteFile(fmt.Sprintf("%s-raw-values.json", filename), jsonBytes)
}

func ExportToJsonReadable(filepath string) {
	rkg := ParseRKG(ReadFile(filepath))
	readable := ConvertRkg(rkg)

	filename := strings.Split(filepath, ".")[0]

	jsonBytes, _ := json.MarshalIndent(readable, "", "\t")
	WriteFile(fmt.Sprintf("%s-readable.json", filename), jsonBytes)
}

func ExportMii(filepath string) {
	bytes := ReadFile(filepath)

	mii := bytes[0x3c:0x86]
	filename := strings.Split(filepath, ".")[0]

	WriteFile(fmt.Sprintf("%s.miigx", filename), mii)
}

func WriteFile(fileName string, bytes []byte) {
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	file.Write(bytes)
}

func ConvertRkg(rkg *RKG) *ReadbleRKG {
	readbleRkg := &ReadbleRKG{}
	readbleRkg.Data = rkg.Data
	readbleRkg.Mii = rkg.Mii
	readbleRkg.CRC = rkg.CRC

	readbleRkg.Header = ConvertHeader(*rkg.Header)

	return readbleRkg
}

func ParseRKG(bytes []byte) *RKG {
	rkg := &RKG{}
	rkgHeader := &Header{}
	header := bytes[:0x88]
	data := bytes[0x88:]
	mii := bytes[0x3c:0x86]
	crc := bytes[0x86:0x88]

	rkg.Header = rkgHeader
	rkg.Data = data
	rkg.Mii = mii
	rkg.CRC = crc

	rkgHeader.Identifier = string(header[:0x04])

	b4 := header[0x04]
	b5 := header[0x05]
	b6 := header[0x06]

	finshTime := &RaceTime{}
	finshTime.Minutes, finshTime.Seconds,
		finshTime.Milliseconds = ParseTime(b4, b5, b6)
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

	ParseLaps(rkgHeader, header)

	b19 := header[0x34]

	rkgHeader.CountryCode = int(int8(b19))

	b20 := header[0x35]

	rkgHeader.StateCode = int(int8(b20))

	b21 := header[0x36]
	b22 := header[0x37]

	rkgHeader.LocationCode = int(int16(b21)<<8 | int16(b22))

	return rkg
}
