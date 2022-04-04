package main

import (
	"fmt"
	"net/http"
	"pokedex/config"
	"pokedex/internal/handler"
	"pokedex/internal/router"
	"pokedex/internal/service/pokemon"
	"pokedex/pkg/adapter/funtranslations"
	"pokedex/pkg/adapter/pokeapi"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := setup(); err != nil {
		log.Fatal().Err(err).Msg("faild to start the service")
		panic(err)
	}
}

func setup() error {
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	pokeClient, err := pokeapi.New(cfg.ThirdParty.PokeAPI.Url)
	if err != nil {
		return errors.Wrap(err, "failed to create new pokeapi client")
	}
	funtranslationsClient, err := funtranslations.New(cfg.ThirdParty.Funtranslations.Url)
	if err != nil {
		return errors.Wrap(err, "failed to create new funtranslations client")
	}

	translateService, err := pokemon.NewTranslateService(funtranslationsClient)
	if err != nil {
		return errors.Wrap(err, "failed to create new pokemon service")
	}

	router := router.New(router.Handlers{
		HealthCheck: handler.HealthCheck,

		GetPokemonByName:           handler.GetPokemonByName(pokeClient),
		GetPokemonByNameTranslated: handler.GetPokemonByNameTranslated(pokeClient, translateService),
	})

	log.Info().Str("port", cfg.Service.Port).Msg("starting service")
	defer log.Info().Msg("stopped service")

	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.Service.Port), router)
}
