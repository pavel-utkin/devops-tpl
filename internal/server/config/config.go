package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type StoreConfig struct {
	Interval    time.Duration `env:"STORE_INTERVAL"`
	DatabaseDSN string        `env:"DATABASE_DSN"`
	File        string        `env:"STORE_FILE"`
	Restore     bool          `env:"RESTORE"`
}

type Config struct {
	ServerAddr       string `env:"ADDRESS"`
	ProfilingAddr    string
	TemplatesAbsPath string `env:"TEMPLATES_ABS_PATH"`
	SignKey          string `env:"KEY"`
	Store            StoreConfig
}

func newConfig() *Config {
	config := Config{}
	config.initDefaultValues()

	return &config
}

func (config *Config) initDefaultValues() {
	config.ServerAddr = "127.0.0.1:8080"
	config.ProfilingAddr = "127.0.0.1:8090"
	config.TemplatesAbsPath = "./templates"
	config.Store = StoreConfig{
		Interval: time.Duration(300) * time.Second,
		File:     "/tmp/devops-metrics-db.json",
		Restore:  true,
	}
}

func (config *Config) parseEnv() error {
	return env.Parse(config)
}

func (config *Config) parseFlags() {
	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "server address (host:port)")
	flag.StringVar(&config.SignKey, "k", config.SignKey, "sign key")
	flag.BoolVar(&config.Store.Restore, "r", config.Store.Restore, "restoring metrics from file")
	flag.StringVar(&config.Store.DatabaseDSN, "d", config.Store.DatabaseDSN, "Database DSN")
	flag.DurationVar(&config.Store.Interval, "i", config.Store.Interval, "store interval (example: 10s)")
	flag.StringVar(&config.Store.File, "f", config.Store.File, "path to file for storage metrics")
	flag.Parse()
}

func LoadConfig() Config {
	config := newConfig()

	config.parseFlags()
	err := config.parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	return *config
}
