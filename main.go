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
	Identifier   string
	Minutes      int
	Seconds      int
	Milliseconds int
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

	m_s_ms := int32(b4)<<16 | int32(b5)<<8 | int32(b6)

	rkgHeader.Minutes = int(m_s_ms >> 17)
	rkgHeader.Seconds = int((m_s_ms >> 10) & 0x7f) // 0x7f = 01111111
	rkgHeader.Milliseconds = int(m_s_ms & 0x3ff)   // 0x3ff = 001111111111

	// fmt.Printf("%b\n", m_s_ms)
	// fmt.Printf("%b\n", b5)
	// fmt.Printf("%d\n", int(b5<<2))
	fmt.Println(rkgHeader)
}
