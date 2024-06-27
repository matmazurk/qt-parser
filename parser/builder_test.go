package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilderBuild(t *testing.T) {
	p := NewBuilder().
		Find("moov/trak/tkhd", 10, 4, "some-finding").
		Find("moov/trak", 12, 4, "other-finding").
		Build()

	expected := &atom{
		typ: rootAtom,
		childs: map[atomType]*atom{
			"moov": {
				typ: "moov",
				childs: map[atomType]*atom{
					"trak": {
						typ: "trak",
						params: []searchParams{
							{
								offset:      12,
								bytesAmount: 4,
								findingName: "other-finding",
							},
						},
						childs: map[atomType]*atom{
							"tkhd": {
								typ:    "tkhd",
								childs: map[atomType]*atom{},
								params: []searchParams{
									{
										offset:      10,
										bytesAmount: 4,
										findingName: "some-finding",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	require.EqualValues(t, expected, p.ToFind())
}
