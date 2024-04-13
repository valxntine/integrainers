package config

import "github.com/caarlos0/env/v6"

type AppConfig struct {
	Port        string `env:"PORT" envDefault:"8088"`
	LibraryHost string `env:"LIBRARY_URL" envDefault:"http://library-mock.app.internal:1080"`
	DB          DatabaseConfig
}

type DatabaseConfig struct {
	Host           string `env:"DB_HOST" envDefault:"book-db.app.internal"`
	Name           string `env:"DB_NAME" envDefault:"book_db"`
	Password       string `env:"DB_PASSWORD" envDefault:"book_db"`
	Port           int    `env:"DB_PORT" envDefault:"3306"`
	Type           string `env:"DB_TYPE" envDefault:"mysql"`
	User           string `env:"DB_USER" envDefault:"book_db"`
	SSLMode        string `env:"DB_SSL_MODE" envDefault:"disable"`
	TimeoutSeconds int    `env:"DB_TIMEOUT_SECONDS" envDefault:"10"`
}

func NewAppConfig() (AppConfig, error) {
	var cfg AppConfig
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
