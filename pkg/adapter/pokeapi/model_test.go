package pokeapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEnglishDescription(t *testing.T) {
	t.Parallel()

	type testcase struct {
		res                *PokemonResponse
		englishDescription string
	}

	tests := map[string]testcase{
		"ReturnsEmptyWhenNoEnglishDescription": {
			res: &PokemonResponse{
				FlavorTextEntries: FlavorTextEntries{
					FlavorTextEntry{
						FlavorText: "polish description",
						Language:   Language{"pl"},
					},
				},
			},
			englishDescription: "",
		},
		"ReturnsEnglishDescription": {
			res: &PokemonResponse{
				FlavorTextEntries: FlavorTextEntries{
					FlavorTextEntry{
						FlavorText: "polish description",
						Language:   Language{"pl"},
					},
					FlavorTextEntry{
						FlavorText: "english description",
						Language:   Language{"en"},
					},
					FlavorTextEntry{
						FlavorText: "italian description",
						Language:   Language{"it"},
					},
				},
			},
			englishDescription: "english description",
		},
		"ReturnsEnglishDescriptionWithoutEscapeCharacters": {
			res: &PokemonResponse{
				FlavorTextEntries: FlavorTextEntries{
					FlavorTextEntry{
						FlavorText: "english\ndescription with some escape char\fact\fers",
						Language:   Language{"en"},
					},
				},
			},
			englishDescription: "english description with some escape characters",
		},
	}

	for description, testCase := range tests {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			desc := tc.res.GetEnglishDescription()

			require.Equal(t, tc.englishDescription, desc)
		})
	}
}
