package main_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testfiles
var videoTrakAtomHex embed.FS

func TestParseSmallAtom(t *testing.T) {
	tcs := []struct {
		name        string
		atomsPath   string
		expected    string
		offset      uint
		bytesAmount uint
		errContains string
	}{
		{
			name:        "should_successfully_parse_width_from_'trak/tkhd'",
			atomsPath:   "trak/tkhd",
			expected:    "00000500",
			offset:      128,
			bytesAmount: 4,
		},
		{
			name:        "should_return_err_when_trying_to_parse_from_nonexisting_atom",
			atomsPath:   "moov/trak/tkhd",
			errContains: "atom 'moov' not found",
		},
		{
			name:        "should_return_err_when_last_atom_in_path_is_not_data_atom",
			atomsPath:   "moov/trak",
			errContains: "atom 'moov/trak' is not a data atom",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			f, err := videoTrakAtomHex.Open("testfiles/video_trak_atom.bin")
			require.NoError(t, err)
			defer f.Close()

			findingName := "some finding"
			parser := NewParserBuilder{}.Find(tc.atomsPath, tc.offset, tc.bytesAmount, findingName).Build()
			res, err := parser.Parse(f)
			if len(tc.errContains) > 0 {
				require.ErrorContains(t, err, tc.errContains)
				return
			}
			require.NoError(t, err)

			findings, ok := res[findingName]
			require.True(t, ok)
			require.Contains(t, findings, tc.expected)
		})
	}
}
