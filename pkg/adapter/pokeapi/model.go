package pokeapi

import "strings"

type (
	// PokemonResponse raw response from PokeAPI
	PokemonResponse struct {
		Name              string            `json:"name"`
		FlavorTextEntries FlavorTextEntries `json:"flavor_text_entries"`
		Habitat           Habitat           `json:"habitat"`
		IsLegendary       bool              `json:"is_legendary"`
	}

	// FlavorTextEntries an arry of FlavorTextEntry
	FlavorTextEntries []FlavorTextEntry

	// FlavorTextEntry contains single text entry with a language information
	FlavorTextEntry struct {
		FlavorText string   `json:"flavor_text"`
		Language   Language `json:"language"`
	}

	// Language defines the language of a text entry
	Language struct {
		Name string `json:"name"`
	}

	// Habitat contains details about Pokemon's habitat
	Habitat struct {
		Name string `json:"name"`
	}

	// Pokemon contains basic information about Pokemon
	Pokemon struct {
		Name        string
		Description string
		Habitat     string
		IsLegendary bool
	}
)

// GetEnglishDescription returns the first english description available
func (r *PokemonResponse) GetEnglishDescription() string {
	var desc string

	for _, e := range r.FlavorTextEntries {
		if e.Language.Name == "en" {
			desc = e.FlavorText

			replacer := strings.NewReplacer("\n", " ", "\f", "")
			desc = replacer.Replace(desc)
		}
	}

	return desc
}
