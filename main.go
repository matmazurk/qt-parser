package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/matmazurk/qt-parser/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: qt-parser filepath")
		os.Exit(1)
	}

	widthTrackName := "tracks width"
	heigthTrackName := "tracks heigth"
	audioFreqName := "audio tracks frequency"
	p := parser.NewBuilder().
		Find("moov/trak/tkhd", 84, 4, widthTrackName).
		Find("moov/trak/tkhd", 88, 4, heigthTrackName).
		Find("moov/trak/mdia/mdhd", 20, 4, audioFreqName).
		Build()

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	res, err := p.Parse(f)
	if err != nil {
		panic(err)
	}

	widths, ok := res[widthTrackName]
	if !ok {
		panic("couldn't find video widths within the file")
	}

	heigths, ok := res[heigthTrackName]
	if !ok {
		panic("couldn't find video heigths within the file")
	}

	frequencies, ok := res[audioFreqName]
	if !ok {
		panic("couldn't find audio frequencies within the file")
	}

	for i := range len(widths) {
		width := parseFixedPoint32(widths[i])
		heigth := parseFixedPoint32(heigths[i])
		if width != 0 && heigth != 0 {
			fmt.Printf("track %d: video; dimensions %.5fx%.5f\n", i+1, width, heigth)
			continue
		}

		freq := binary.BigEndian.Uint32(frequencies[i])
		fmt.Printf("track %d: audio; sampling frequency %dHz\n", i+1, freq)
	}
}

func parseFixedPoint32(data []byte) float64 {
	if len(data) != 4 {
		panic("data must be 4 bytes long")
	}

	value := binary.BigEndian.Uint32(data)
	integerPart := int32(value >> 16)
	fractionalPart := value & 0xFFFF
	fraction := float64(fractionalPart) / 65536.0

	return float64(integerPart) + fraction
}
