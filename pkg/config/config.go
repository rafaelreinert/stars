package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Config is the struct which carries all configurable values to startup the application
type Config struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	DBURI    string `env:"DB_URI" envDefault:"mongodb://localhost:27017"`
	SWAPIURL string `env:"SWAPI_URL" envDefault:"https://swapi.dev/api"`
}

// New return a New Config struct filled with the environment variables values or default values
func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "Error during the environment variables parse")
	}
	return cfg, nil
}
