package funtranslations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNew_Error(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		url string
		err string
	}{
		"ReturnsErrorWhenURLEmpty": {
			url: "",
			err: "url: invalid parameter",
		},
	}
	for description, testCase := range cases {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			c, err := New(tc.url)
			require.EqualError(t, err, testCase.err)
			require.Nil(t, c)
		})
	}
}

func TestNew_Success(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		url    string
		client *http.Client
	}{
		"ReturnsClientWhenURLSet": {
			url: "http://example.com",
		},
	}
	for description, testCase := range cases {
		tc := testCase

		t.Run(description, func(t *testing.T) {
			t.Parallel()

			c, err := New(tc.url)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}
}

func TestExecute_Error(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		checkErr func(t *testing.T, err error)
		server   *httptest.Server
	}{
		"ReturnsErrorWhenInvalidJSONResponse": {
			checkErr: func(t *testing.T, err error) {
				require.Error(t, err)
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("not-a-json"))
				require.NoError(t, err)

				w.WriteHeader(http.StatusTeapot)
			})),
		},
	}

	for description, testCase := range cases {
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

			req, err := http.NewRequest(http.MethodPost, tc.server.URL, nil)
			require.NoError(t, err)

			var val TranslateResponse
			res, err := client.execute(ctx, req, &val)

			res.Body.Close()

			tc.checkErr(t, err)
		})
	}
}

func TestExecute_Success(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		checkErr func(t *testing.T, err error)
		server   *httptest.Server
	}{
		"ReturnsNilWhenValidJSONResponse": {
			checkErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(`{"name":"mewtwo"}`))
				require.NoError(t, err)

				w.WriteHeader(http.StatusTeapot)
			})),
		},
	}

	for description, testCase := range cases {
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

			req, err := http.NewRequest(http.MethodPost, tc.server.URL, nil)
			require.NoError(t, err)

			var val TranslateResponse
			res, err := client.execute(ctx, req, &val)

			res.Body.Close()

			tc.checkErr(t, err)
		})
	}
}
