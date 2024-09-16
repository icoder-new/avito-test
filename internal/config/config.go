package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	HTTPServer
	Database
}

type HTTPServer struct {
	Host        string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port        string        `env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT" env-default:"45s"`
}

type Database struct {
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"0.0.0.0"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USERNAME" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" env-default:"postgres"`
}

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("cannot load env: %s", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
