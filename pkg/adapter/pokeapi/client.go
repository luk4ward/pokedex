package pokeapi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	path = "%s/api/v2/%s/%s"
)

var ErrInvalidParam = errors.New("invalid parameter")

type (
	// Client represents a PokeAPI client for making request.
	Client struct {
		httpClient *http.Client
		url        string
	}
)

// New creates a new PokeAPI client
func New(url string) (*Client, error) {
	if url == "" {
		return nil, errors.Wrap(ErrInvalidParam, "url")
	}

	return &Client{
		httpClient: &http.Client{},
		url:        url,
	}, nil
}

func (c *Client) execute(ctx context.Context, req *http.Request, responseVal interface{}) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending request")
	}

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "error reading body")
		}

		if err := json.Unmarshal(body, responseVal); err != nil {
			return res, errors.Wrap(err, "error unmarshaling response")
		}
	}

	return res, nil
}
