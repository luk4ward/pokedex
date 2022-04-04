//go:generate mockgen -destination=./mocks/pokeapi_mock.go -package=mocks pokedex/pkg/adapter/pokeapi PokemonFetcher
//go:generate mockgen -destination=./mocks/pokemonservice_mock.go -package=mocks pokedex/internal/service/pokemon DescriptionTranslator

package handler

import (
	"context"
	"net/http"
	"pokedex/internal/service/pokemon"
	"pokedex/pkg/adapter/pokeapi"
	"pokedex/pkg/http/response"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	nameParam = "name"
)

// GetPokemonByName fetches a Pokemon for a given name
func GetPokemonByName(pc pokeapi.PokemonFetcher) http.HandlerFunc {
	return response.Render(func(w http.ResponseWriter, r *http.Request) response.Renderer {
		ctx := r.Context()

		pok, err := getPokemon(ctx, pc, chi.URLParam(r, nameParam))
		if err != nil {
			return ErrorRenderer(err)
		}

		return response.JSON(http.StatusOK, pok)
	})
}

// GetPokemonByNameTranslated fetches a Pokemon with a translated description for a given name
func GetPokemonByNameTranslated(pc pokeapi.PokemonFetcher, dt pokemon.DescriptionTranslator) http.HandlerFunc {
	return response.Render(func(w http.ResponseWriter, r *http.Request) response.Renderer {
		ctx := r.Context()

		pok, err := getPokemon(ctx, pc, chi.URLParam(r, nameParam))
		if err != nil {
			return ErrorRenderer(err)
		}

		pok.Description = dt.TranslateDescription(ctx, pok)

		return response.JSON(http.StatusOK, pok)
	})
}

func getPokemon(ctx context.Context, pc pokeapi.PokemonFetcher, name string) (*pokeapi.Pokemon, error) {
	pok, err := pc.FetchByName(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch pokemon: %s", name)
	}

	if pok == nil {
		log.Info().Str("name", name).Msg(ErrPokemonNotFound.Error())
		return nil, ErrPokemonNotFound
	}

	return pok, nil
}
