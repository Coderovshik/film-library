package config

import (
	"log"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        string `env:"SERVER_PORT" env-default:"8080"`
	SigningKey  string `env:"SIGNING_KEY" env-required:"true"`
	DatabaseURI string `env:"DATABASE_URI" env-required:"true"`
	Docs        string `env:"DOCS" env-required:"true"`
}

func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func New() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
