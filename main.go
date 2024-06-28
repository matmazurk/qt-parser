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

	const (
		widthTrackName  = "tracks width"
		heigthTrackName = "tracks heigth"
		audioFreqName   = "audio tracks frequency"
	)
	p := parser.NewBuilder().
		Find("moov/trak/tkhd", 84, 4, widthTrackName).
		Find("moov/trak/tkhd", 88, 4, heigthTrackName).
		Find("moov/trak/mdia/mdhd", 20, 4, audioFreqName).
		Build()

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("couldn't open '%s':%s\n", os.Args[1], err.Error())
		os.Exit(1)
	}
	defer f.Close()

	res, err := p.Parse(f)
	if err != nil {
		fmt.Printf("couldn't parse the file- might be corrupted:%s\n", err.Error())
		os.Exit(1)
	}

	widths, wOk := res[widthTrackName]
	heigths, hOk := res[heigthTrackName]
	frequencies, fOk := res[audioFreqName]
	if !wOk && !hOk && !fOk {
		fmt.Println("the file doesn't contain any video or audio tracks")
		return
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
