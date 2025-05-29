package configuration

import (
	env "github.com/Netflix/go-env"
)

type (
	Configuration struct {
		AppHTTPPort int    `env:"APP_HTTP_PORT"`
		AppVersion  string `env:"APP_VERSION"`
		AppName     string `env:"APP_NAME"`
		Environment string `env:"ENVIRONMENT"`
		Database    Database
		Feature     Feature
	}

	Database struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
		Schema   string `env:"DB_SCHEMA"`
		Driver   string `env:"DB_DRIVER"`
		SSLMode  string `env:"DB_SSLMODE"`
	}

	Feature struct {
		FeatureHealthCheck bool `env:"FEATURE_HEALTH_CHECK"`
	}
)

func LoadConfig() (Configuration, error) {
	config := Configuration{}
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
