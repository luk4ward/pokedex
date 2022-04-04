package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"pokedex/internal/handler/mocks"
	"pokedex/pkg/adapter/pokeapi"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	pokemonName = "mewtwo"
)

func TestGetPokemonByName(t *testing.T) {
	t.Parallel()

	type testcase struct {
		req            func(t *testing.T) *http.Request
		pokemonFetcher func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}

	tests := map[string]testcase{
		"ReturnsNotFoundWhenPokemonNotFound": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, pokemonName)

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), pokemonName).
					Return(nil, nil)

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		"ReturnsNotFoundWhenPokemonNameEmpty": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, "")

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), "").
					Return(nil, nil)

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		"ReturnsInternalServerErrorWhenFailedToFetchPokemon": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, pokemonName)

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), pokemonName).
					Return(nil, errors.New("foo"))

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		"ReturnsOKAndPokemonWhenPokemonFound": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, pokemonName)

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), pokemonName).
					Return(&pokeapi.Pokemon{
						Name: pokemonName,
					}, nil)

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)

				body, err := ioutil.ReadAll(rec.Body)
				require.NoError(t, err)

				pokemon := &pokeapi.Pokemon{}
				err = json.Unmarshal(body, pokemon)
				require.NoError(t, err)

				require.Equal(t, pokemonName, pokemon.Name)
			},
		},
	}

	for description, testCase := range tests {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rec := httptest.NewRecorder()
			req := tc.req(t)

			handler := GetPokemonByName(tc.pokemonFetcher(t, ctrl))
			handler.ServeHTTP(rec, req)

			tc.checkResponse(t, rec)
		})
	}
}

func TestGetPokemonByNameTranslated(t *testing.T) {
	t.Parallel()

	type testcase struct {
		req            func(t *testing.T) *http.Request
		pokemonFetcher func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher
		translator     func(t *testing.T, c *gomock.Controller) *mocks.MockDescriptionTranslator
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}

	tests := map[string]testcase{
		"ReturnsNotFoundWhenPokemonNotFound": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, pokemonName)

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), pokemonName).
					Return(nil, nil)

				return m
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockDescriptionTranslator {
				m := mocks.NewMockDescriptionTranslator(c)

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		"ReturnsNotFoundWhenPokemonNameEmpty": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, "")

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), "").
					Return(nil, nil)

				return m
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockDescriptionTranslator {
				m := mocks.NewMockDescriptionTranslator(c)

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		"ReturnsInternalServerErrorWhenFailedToFetchPokemon": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, pokemonName)

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), pokemonName).
					Return(nil, errors.New("foo"))

				return m
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockDescriptionTranslator {
				m := mocks.NewMockDescriptionTranslator(c)

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		"ReturnsOKAndPokemonWithTranslatedDescriptionWhenPokemonFound": {
			req: func(t *testing.T) *http.Request {
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add(nameParam, pokemonName)

				req := &http.Request{}

				return req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, ctx))
			},
			pokemonFetcher: func(t *testing.T, c *gomock.Controller) *mocks.MockPokemonFetcher {
				m := mocks.NewMockPokemonFetcher(c)

				m.EXPECT().
					FetchByName(gomock.Any(), pokemonName).
					Return(&pokeapi.Pokemon{
						Name:        pokemonName,
						Description: "old pokemon description",
					}, nil)

				return m
			},
			translator: func(t *testing.T, c *gomock.Controller) *mocks.MockDescriptionTranslator {
				m := mocks.NewMockDescriptionTranslator(c)

				m.EXPECT().
					TranslateDescription(gomock.Any(), gomock.AssignableToTypeOf(&pokeapi.Pokemon{})).
					Return("translated description")

				return m
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)

				body, err := ioutil.ReadAll(rec.Body)
				require.NoError(t, err)

				pokemon := &pokeapi.Pokemon{}
				err = json.Unmarshal(body, pokemon)
				require.NoError(t, err)

				require.Equal(t, pokemonName, pokemon.Name)
				require.Equal(t, "translated description", pokemon.Description)
			},
		},
	}

	for description, testCase := range tests {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rec := httptest.NewRecorder()
			req := tc.req(t)

			handler := GetPokemonByNameTranslated(tc.pokemonFetcher(t, ctrl), tc.translator(t, ctrl))
			handler.ServeHTTP(rec, req)

			tc.checkResponse(t, rec)
		})
	}
}
