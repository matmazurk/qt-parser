package parser_test

import (
	"embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matmazurk/qt-parser/parser"
)

//go:embed testfiles
var testfiles embed.FS

func TestParse(t *testing.T) {
	tcs := []struct {
		filename        string
		expectedWidths  [][]byte
		expectedHeigths [][]byte
		expectedFreqs   [][]byte
	}{
		{
			filename: "testfiles/sample.mp4",
			expectedWidths: [][]byte{
				{0x01, 0x40, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedHeigths: [][]byte{
				{0x00, 0xb4, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedFreqs: [][]byte{
				{0x00, 0x00, 0x28, 0x00},
				{0x00, 0x00, 0xac, 0x44},
			},
		},
		{
			filename: "testfiles/sample_fragmented.mp4",
			expectedWidths: [][]byte{
				{0x05, 0x00, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedHeigths: [][]byte{
				{0x02, 0xd0, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedFreqs: [][]byte{
				{0x00, 0x1, 0x5f, 0x90},
				{0x00, 0x00, 0xac, 0x44},
			},
		},
		{
			filename: "testfiles/Clouds.mov",
			expectedWidths: [][]byte{
				{0x2, 0xd0, 0x00, 0x00},
			},
			expectedHeigths: [][]byte{
				{0x01, 0xe6, 0x00, 0x00},
			},
			expectedFreqs: [][]byte{
				{0x00, 0x00, 0x00, 0x1e},
			},
		},
		{
			filename: "testfiles/example.mp4",
			expectedWidths: [][]byte{
				{0x2, 0x80, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedHeigths: [][]byte{
				{0x01, 0x68, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedFreqs: [][]byte{
				{0x00, 0x00, 0x0, 0x1e},
				{0x00, 0x00, 0xbb, 0x80},
			},
		},
		{
			filename: "testfiles/example2.mp4",
			expectedWidths: [][]byte{
				{0x1, 0xe0, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedHeigths: [][]byte{
				{0x01, 0xe, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedFreqs: [][]byte{
				{0x00, 0x00, 0x0, 0x1e},
				{0x00, 0x00, 0xbb, 0x80},
			},
		},
		{
			filename: "testfiles/example3.mp4",
			expectedWidths: [][]byte{
				{0x1, 0x40, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedHeigths: [][]byte{
				{0x00, 0xf0, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
			expectedFreqs: [][]byte{
				{0x00, 0x00, 0x3c, 0x00},
				{0x00, 0x00, 0xbb, 0x80},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("should_correctly_parse_%s", tc.filename), func(t *testing.T) {
			f, err := testfiles.Open(tc.filename)
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
			require.Equal(t, tc.expectedWidths, widths)

			heigths, ok := res[heigthTrackName]
			require.True(t, ok)
			require.Equal(t, tc.expectedHeigths, heigths)

			frequencies, ok := res[audioFreqName]
			require.True(t, ok)
			require.Equal(t, tc.expectedFreqs, frequencies)
		})
	}
}

func TestNonExistingAtomPath(t *testing.T) {
	f, err := testfiles.Open("testfiles/sample.mp4")
	require.NoError(t, err)
	defer f.Close()

	someSearch := "some search"
	p := parser.NewBuilder().
		Find("maav/trak/tkhd", 84, 4, someSearch).
		Build()
	res, err := p.Parse(f)
	require.NoError(t, err)
	require.Empty(t, res)
}
