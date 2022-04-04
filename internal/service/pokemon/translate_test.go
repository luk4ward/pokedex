package pokemon

import (
	"context"
	"pokedex/internal/service/pokemon/mocks"
	"pokedex/pkg/adapter/funtranslations"
	"pokedex/pkg/adapter/pokeapi"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewTranslateService_Error(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		traslateSerivce funtranslations.Translator
		err             string
	}{
		"ReturnsErrortranslatorNil": {
			traslateSerivce: nil,
			err:             "translator: invalid parameter",
		},
	}
	for description, testCase := range cases {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			s, err := NewTranslateService(tc.traslateSerivce)
			require.EqualError(t, err, tc.err)
			require.Nil(t, s)
		})
	}
}

func TestNewTranslateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	traslateSerivceMock := mocks.NewMockTranslator(ctrl)

	c, err := NewTranslateService(traslateSerivceMock)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestTranslateDescription(t *testing.T) {
	t.Parallel()

	defaultDesc := "it is default description"
	yodaDesc := "yoda description, it is"
	shakespeareDesc := "shakespear babble description"

	type testcase struct {
		pokemon       *pokeapi.Pokemon
		translator    func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator
		checkResponse func(t *testing.T, desc string)
	}

	tests := map[string]testcase{
		"ReturnYodaTranslationWhenLegendary": {
			pokemon: &pokeapi.Pokemon{
				Name:        "mewtwo",
				Habitat:     "rare",
				Description: defaultDesc,
				IsLegendary: true,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToYoda(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return yodaDesc, nil
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, yodaDesc, desc)
			},
		},
		"ReturnYodaTranslationWhenCaveHabitat": {
			pokemon: &pokeapi.Pokemon{
				Name:        "zubat",
				Habitat:     "cave",
				Description: defaultDesc,
				IsLegendary: false,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToYoda(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return yodaDesc, nil
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, yodaDesc, desc)
			},
		},
		"ReturnYodaTranslationWhenLegendaryAndCave": {
			pokemon: &pokeapi.Pokemon{
				Name:        "mewtwo",
				Habitat:     "cave",
				Description: defaultDesc,
				IsLegendary: true,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToYoda(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return yodaDesc, nil
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, yodaDesc, desc)
			},
		},
		"ReturnDefaultDescriptionWhenLegendaryAndFaildToTranslate": {
			pokemon: &pokeapi.Pokemon{
				Name:        "mewtwo",
				Habitat:     "rare",
				Description: defaultDesc,
				IsLegendary: true,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToYoda(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return "", errors.New("foo")
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, defaultDesc, desc)
			},
		},
		"ReturnDefaultDescriptionWhenCaveHabitatAndFaildToTranslate": {
			pokemon: &pokeapi.Pokemon{
				Name:        "zubat",
				Habitat:     "cave",
				Description: defaultDesc,
				IsLegendary: false,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToYoda(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return "", errors.New("foo")
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, defaultDesc, desc)
			},
		},
		"ReturnShakespeareTranslationWhenNotLegendaryOrCaveHabitat": {
			pokemon: &pokeapi.Pokemon{
				Name:        "eevee",
				Habitat:     "normal",
				Description: defaultDesc,
				IsLegendary: false,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToShakespeare(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return shakespeareDesc, nil
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, shakespeareDesc, desc)
			},
		},
		"ReturnDefaultDescriptionWhenNotLegendaryOrCaveHabitatAndFailedToTranslate": {
			pokemon: &pokeapi.Pokemon{
				Name:        "eevee",
				Habitat:     "normal",
				Description: defaultDesc,
				IsLegendary: false,
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockTranslator {
				m := mocks.NewMockTranslator(c)
				m.EXPECT().
					TranslateToShakespeare(gomock.Any(), defaultDesc).
					DoAndReturn(func(ctx context.Context, desc string) (string, error) {
						return "", errors.New("foo")
					})

				return m
			},
			checkResponse: func(t *testing.T, desc string) {
				require.Equal(t, defaultDesc, desc)
			},
		},
	}

	for description, testCase := range tests {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			svc := TranslatService{
				translator: tc.translator(t, ctrl),
			}

			res := svc.TranslateDescription(ctx, tc.pokemon)

			tc.checkResponse(t, res)
		})
	}
}
