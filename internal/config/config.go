package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string           `yaml:"env" env-default:"local`
	Storage    StorageConfig    `yaml: storage env-required:"true"`
	HTTPServer HTTPServerConfig `yaml: "http_server"`
	SSO        SSOsrt           `yaml:"sso"`
}

type StorageConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml: "user" env-required:"true"`
	Password    string        `yaml: "password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type SSOsrt struct {
	Address      string        `yaml:"address" env-default:"localhost:11011"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	RetriesCount int           `yaml:"retriesCount" env-default:"3"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config: %s", err)
	}

	return &cfg
}
