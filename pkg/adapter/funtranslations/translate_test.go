package funtranslations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTranslate_Error(t *testing.T) {
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

			res, err := client.TranslateToYoda(ctx, "text to translate")

			require.Error(t, err)
			require.Empty(t, res)
		})
	}
}

func TestTranslate_Success(t *testing.T) {
	t.Parallel()

	type testcase struct {
		server         *httptest.Server
		translatedText string
	}

	tests := map[string]testcase{
		"ReturnsTranslationWhenNoErrorOccured": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusOK)

				_, err := w.Write([]byte(`{"contents": { "translated": "text translated to Yoda"}}`))
				require.NoError(t, err)
			})),
			translatedText: "text translated to Yoda",
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

			res, err := client.TranslateToYoda(ctx, "text to translate")

			require.NoError(t, err)
			require.Equal(t, tc.translatedText, res)
		})
	}
}

func TestTranslateToYoda_Success(t *testing.T) {
	t.Parallel()

	type testcase struct {
		server *httptest.Server
	}

	tests := map[string]testcase{
		"ReturnsTranslationWhenNoErrorOccured": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				require.True(t, strings.Contains(r.RequestURI, "yoda"))
				w.WriteHeader(http.StatusOK)

				_, err := w.Write([]byte(`{"contents": { "translated": "text translated to Yoda"}}`))
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

			_, err := client.TranslateToYoda(ctx, "text to translate")

			require.NoError(t, err)
		})
	}
}

func TestTranslateToShakespeare_Success(t *testing.T) {
	t.Parallel()

	type testcase struct {
		server *httptest.Server
	}

	tests := map[string]testcase{
		"ReturnsTranslationWhenNoErrorOccured": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				require.True(t, strings.Contains(r.RequestURI, "shakespeare"))
				w.WriteHeader(http.StatusOK)

				_, err := w.Write([]byte(`{"contents": { "translated": "text translated to Shakespeare"}}`))
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

			_, err := client.TranslateToShakespeare(ctx, "text to translate")

			require.NoError(t, err)
		})
	}
}
