package funtranslations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Translator converts text
type Translator interface {
	// TranslateToYoda converts English to Yoda speak
	TranslateToYoda(ctx context.Context, text string) (string, error)

	// TranslateToShakespeare converts English to Shakespeare
	TranslateToShakespeare(ctx context.Context, text string) (string, error)
}

// TranslateToYoda converts English to Yoda speak
func (c *Client) TranslateToYoda(ctx context.Context, text string) (string, error) {
	return c.translate(ctx, text, "yoda")
}

// TranslateToShakespeare converts English to Shakespeare
func (c *Client) TranslateToShakespeare(ctx context.Context, text string) (string, error) {
	return c.translate(ctx, text, "shakespeare")
}

func (c *Client) translate(ctx context.Context, text string, to string) (string, error) {
	data, err := json.Marshal(TranslateRequest{Text: text})
	if err != nil {
		return "", errors.Wrapf(err, "error marshalling translate to %s request", to)
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf(path, c.url, to), bytes.NewBuffer(data))
	if err != nil {
		return "", errors.Wrap(err, "error creating http request")
	}

	raw := &TranslateResponse{}

	// requires nolint rule as body on this stage is already closed and linter doesn't pick that up
	res, err := c.execute(ctx, httpReq, raw) // nolint:bodyclose
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute translate to %s request", to)
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.Errorf("invalid response code: %s", res.Status)
	}

	return raw.Contents.Translated, nil
}
