package main_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testfiles
var videoTrakAtomHex embed.FS

func TestParser(t *testing.T) {
	t.Run("should_successfully_parse_width_from_'trak/tkhd'", func(t *testing.T) {
		f, err := videoTrakAtomHex.Open("testfiles/video_trak_atom.bin")
		require.NoError(t, err)
		defer f.Close()

		var (
			offset      = 128
			bytesAmount = 4
			findingName = "video track width"
		)
		parser := NewParserBuilder{}.Find("trak/tkhd", offset, bytesAmount, findingName).Build()
		res, err := parser.Parse(f)
		require.NoError(t, err)

		findings, ok := res[findingName]
		require.True(t, ok)
		expected := "00000500"
		require.Contains(t, findings, expected)
	})
}
