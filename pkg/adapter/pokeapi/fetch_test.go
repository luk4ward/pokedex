package pokeapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestFetchByName_Error(t *testing.T) {
	t.Parallel()

	type testcase struct {
		server *httptest.Server
	}

	tests := map[string]testcase{
		"ReturnsErrorWhenFailedToExecute": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusOK)

				_, err := w.Write([]byte("not-a-json"))
				require.NoError(t, err)
			})),
		},
		"ReturnsErrorWhenInvalidResponseCode": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusTeapot)

				_, err := w.Write([]byte("{}"))
				require.NoError(t, err)
			})),
		},
	}

	for description, testCase := range tests {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			defer tc.server.Close()

			ctx := context.Background()
			client := Client{
				httpClient: tc.server.Client(),
				url:        tc.server.URL,
			}

			res, err := client.FetchByName(ctx, "mewtwo")

			require.Error(t, err)
			require.Empty(t, res)
		})
	}
}

func TestFetchByName_Success(t *testing.T) {
	t.Parallel()

	type testcase struct {
		server        *httptest.Server
		checkResponse func(t *testing.T, res *Pokemon)
	}

	tests := map[string]testcase{
		"ReturnsNilIfPokemonNotFound": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusNotFound)

				_, err := w.Write([]byte("{}"))
				require.NoError(t, err)
			})),
			checkResponse: func(t *testing.T, res *Pokemon) {
				require.Nil(t, res)
			},
		},
		"ReturnsPokemonIfPokemonFound": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusOK)

				_, err := w.Write([]byte(`
				{
					"name": "mewtwo",
					"is_legendary": true, 
					"habitat": { "name": "rare"},
					"flavor_text_entries": [ 
						{
							"flavor_text": "Il suo DNA è quasi uguale a quello di Mew. Ciò nonostante, sono agli antipodi per dimensioni e carattere",
							"language": { "name": "it" }
						},
						{
							"flavor_text": "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.",
							"language": { "name": "en" }
						}
					]
				}`))
				require.NoError(t, err)
			})),
			checkResponse: func(t *testing.T, res *Pokemon) {
				require.Equal(t, "mewtwo", res.Name)
				require.Equal(t, true, res.IsLegendary)
				require.Equal(t, "rare", res.Habitat)
				require.Equal(t, "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.", res.Description)
			},
		},
	}

	for description, testCase := range tests {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			defer tc.server.Close()

			ctx := context.Background()
			client := Client{
				httpClient: tc.server.Client(),
				url:        tc.server.URL,
			}

			res, err := client.FetchByName(ctx, "mewtwo")

			require.NoError(t, err)
			tc.checkResponse(t, res)
		})
	}
}
