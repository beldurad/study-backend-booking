package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Profile         string `yaml:"profile"`
	WebServerConfig `yaml:"web_server"`
	DatabaseConfig  `yaml:"database"`
	SecurityConfig  `yaml:"security"`
	WorkerConfig    `yaml:"worker"`
}

type WorkerConfig struct {
	Interval  time.Duration `yaml:"interval"`
	DaysAhead int           `yaml:"days_ahead"`
}

type WebServerConfig struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            uint16 `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DatabaseName    string `yaml:"name"`
	InitSqlFilepath string `yaml:"init_sql_filepath"`
}

type SecurityConfig struct {
	Secret string `yaml:"secret"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return cfg
}
