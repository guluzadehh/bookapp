package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	GRPCServer GRPCServer `yaml:"grpc_server"`
	Postgres   Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	Port        int           `yaml:"port" env-default:"8080"`
}

type GRPCServer struct {
	Port int `yaml:"port" env-default:"50000"`
}

type Postgres struct {
	Url     string   `yaml:"-" env:"POSTGRES_URL" env-required:"true"`
	Options []string `yaml:"options"`
}

func (db *Postgres) DSN(options []string) string {
	if options == nil {
		options = db.Options
	}

	opts := strings.Join(options, "&")
	if len(opts) == 0 {
		return db.Url
	}

	return fmt.Sprintf("%s?%s", db.Url, opts)
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file `%s` does not exist", cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("can't read config file `%s` and env variables: %s", cfgPath, err)
	}

	return &cfg
}
