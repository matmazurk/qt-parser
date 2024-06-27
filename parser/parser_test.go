package parser_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matmazurk/qt-parser/parser"
)

//go:embed testfiles
var testfiles embed.FS

func TestParseSample(t *testing.T) {
	f, err := testfiles.Open("testfiles/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	widthTrackName := "track width"
	heigthTrackName := "track heigth"
	audioFreqName := "audio track frequency"
	p := parser.NewBuilder().
		Find("moov/trak/tkhd", 84, 4, widthTrackName).
		Find("moov/trak/tkhd", 88, 4, heigthTrackName).
		Find("moov/trak/mdia/mdhd", 20, 4, audioFreqName).
		Build()
	res, err := p.Parse(f)
	require.NoError(t, err)

	widths, ok := res[widthTrackName]
	require.True(t, ok)
	expectedWidths := [][]byte{
		{0x01, 0x40, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
	}
	require.Equal(t, expectedWidths, widths)

	heigths, ok := res[heigthTrackName]
	require.True(t, ok)
	expectedHeigths := [][]byte{
		{0x00, 0xb4, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
	}
	require.Equal(t, expectedHeigths, heigths)

	frequencies, ok := res[audioFreqName]
	require.True(t, ok)
	expectedFrequencies := [][]byte{
		{0x00, 0x00, 0x28, 0x00},
		{0x00, 0x00, 0xac, 0x44},
	}
	require.Equal(t, expectedFrequencies, frequencies)
}

func TestParseSampleFragmented(t *testing.T) {
	f, err := testfiles.Open("testfiles/sample_fragmented.mp4")
	require.NoError(t, err)
	defer f.Close()

	widthTrackName := "track width"
	heigthTrackName := "track heigth"
	audioFreqName := "audio track frequency"
	p := parser.NewBuilder().
		Find("moov/trak/tkhd", 84, 4, widthTrackName).
		Find("moov/trak/tkhd", 88, 4, heigthTrackName).
		Find("moov/trak/mdia/mdhd", 20, 4, audioFreqName).
		Build()
	res, err := p.Parse(f)
	require.NoError(t, err)

	widths, ok := res[widthTrackName]
	require.True(t, ok)
	expectedWidths := [][]byte{
		{0x05, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
	}
	require.Equal(t, expectedWidths, widths)

	heigths, ok := res[heigthTrackName]
	require.True(t, ok)
	expectedHeigths := [][]byte{
		{0x02, 0xd0, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
	}
	require.Equal(t, expectedHeigths, heigths)

	frequencies, ok := res[audioFreqName]
	require.True(t, ok)
	expectedFrequencies := [][]byte{
		{0x00, 0x1, 0x5f, 0x90},
		{0x00, 0x00, 0xac, 0x44},
	}
	require.Equal(t, expectedFrequencies, frequencies)
}

func TestParseCoulds(t *testing.T) {
	f, err := testfiles.Open("testfiles/Clouds.mov")
	require.NoError(t, err)
	defer f.Close()

	widthTrackName := "track width"
	heigthTrackName := "track heigth"
	audioFreqName := "audio track frequency"
	p := parser.NewBuilder().
		Find("moov/trak/tkhd", 84, 4, widthTrackName).
		Find("moov/trak/tkhd", 88, 4, heigthTrackName).
		Find("moov/trak/mdia/mdhd", 20, 4, audioFreqName).
		Build()
	res, err := p.Parse(f)
	require.NoError(t, err)

	widths, ok := res[widthTrackName]
	require.True(t, ok)
	expectedWidths := [][]byte{
		{0x2, 0xd0, 0x00, 0x00},
	}
	require.Equal(t, expectedWidths, widths)

	heigths, ok := res[heigthTrackName]
	require.True(t, ok)
	expectedHeigths := [][]byte{
		{0x01, 0xe6, 0x00, 0x00},
	}
	require.Equal(t, expectedHeigths, heigths)

	frequencies, ok := res[audioFreqName]
	require.True(t, ok)
	expectedFrequencies := [][]byte{
		{0x00, 0x00, 0x00, 0x1e},
	}
	require.Equal(t, expectedFrequencies, frequencies)
}
