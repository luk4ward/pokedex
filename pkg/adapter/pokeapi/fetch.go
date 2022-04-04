package pokeapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// DocumentFetcher fetches Pokemons from PokeAPI
type PokemonFetcher interface {
	//FetchByName returns Pokemon details by a given name
	FetchByName(ctx context.Context, name string) (*Pokemon, error)
}

// FetchByName returns Pokemon details by a given name.
func (c *Client) FetchByName(ctx context.Context, name string) (*Pokemon, error) {
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf(path, c.url, "pokemon-species", name), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating http request")
	}

	var p *Pokemon
	raw := &PokemonResponse{}

	// requires nolint rule as body on this stage is already closed and linter doesn't pick that up
	res, err := c.execute(ctx, httpReq, raw) // nolint:bodyclose
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute fetch by name request")
	}

	if res.StatusCode == http.StatusNotFound {
		return p, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid response code: %s", res.Status)
	}

	p = &Pokemon{
		Name:        raw.Name,
		Habitat:     raw.Habitat.Name,
		IsLegendary: raw.IsLegendary,
		Description: raw.GetEnglishDescription(),
	}

	return p, nil
}
