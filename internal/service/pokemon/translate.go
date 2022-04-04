//go:generate mockgen -destination=./mocks/funtranslations_mock.go -package=mocks pokedex/pkg/adapter/funtranslations Translator

package pokemon

import (
	"context"
	"pokedex/pkg/adapter/funtranslations"
	"pokedex/pkg/adapter/pokeapi"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	caveHabitat = "cave"
)

var ErrInvalidParam = errors.New("invalid parameter")

type (
	// DescriptionTranslator alows description translations
	DescriptionTranslator interface {
		// TranslateDescription translates description of a given Pokemon
		TranslateDescription(ctx context.Context, pok *pokeapi.Pokemon) string
	}

	// TranslatService service that allows transations on the original Pokemon data
	TranslatService struct {
		translator funtranslations.Translator
	}
)

// NewTranslateService creates a new Pokemon translate service
func NewTranslateService(translator funtranslations.Translator) (*TranslatService, error) {
	if translator == nil {
		return nil, errors.Wrap(ErrInvalidParam, "translator")
	}

	return &TranslatService{
		translator: translator,
	}, nil
}

// TranslateDescription translates description of a given Pokemon
func (t *TranslatService) TranslateDescription(ctx context.Context, pok *pokeapi.Pokemon) string {
	if pok.IsLegendary || pok.Habitat == caveHabitat {
		desc, err := t.translator.TranslateToYoda(ctx, pok.Description)
		if err != nil {
			log.Info().Err(err).Msg("failed to translate; defauling description")

			return pok.Description
		}

		return desc
	}

	desc, err := t.translator.TranslateToShakespeare(ctx, pok.Description)
	if err != nil {
		log.Info().Err(err).Msg("failed to translate; defauling description")

		return pok.Description
	}

	return desc
}
