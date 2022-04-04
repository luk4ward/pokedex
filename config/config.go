package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigPath = "../../config.yaml"
)

type (
	// Config variables for the application
	Config struct {
		Service    Service    `yaml:"service"`
		ThirdParty ThirdParty `yaml:"third_party"`
	}

	// Service represents service configuration
	Service struct {
		Port string `yaml:"port"`
	}

	// ThirdParty represents configuration for third party services
	ThirdParty struct {
		Funtranslations API `yaml:"funtranslations"`
		PokeAPI         API `yaml:"pokeapi"`
	}

	// API variables for an external API
	API struct {
		Url string `yaml:"url"`
	}
)

// Load loads the configuration for the application
func Load() (*Config, error) {
	config := &Config{}

	file, err := os.Open(defaultConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
