package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type HTTPClientConfig struct {
	RetryCount       int           `env:"RETRY_CONN_COUNT"`
	RetryWaitTime    time.Duration `env:"RETRY_CONN_WAIT_TIME"`
	RetryMaxWaitTime time.Duration `env:"RETRY_CONN_MAX_WAIT_TIME"`
	ServerAddr       string        `env:"ADDRESS"`
}

type Config struct {
	PollInterval         time.Duration `env:"POLL_INTERVAL"`
	ReportInterval       time.Duration `env:"REPORT_INTERVAL"`
	SignKey              string        `env:"KEY"`
	HTTPClientConnection HTTPClientConfig
}

func (config *Config) initDefaultValues() {
	config.PollInterval = time.Duration(2) * time.Second
	config.ReportInterval = time.Duration(10) * time.Second

	config.HTTPClientConnection = HTTPClientConfig{
		RetryCount:       2,
		RetryWaitTime:    time.Duration(10) * time.Second,
		RetryMaxWaitTime: time.Duration(90) * time.Second,
		ServerAddr:       "127.0.0.1:8080",
	}
}

func newConfig() *Config {
	config := Config{}
	config.initDefaultValues()

	return &config
}

func (config *Config) parseEnv() error {
	return env.Parse(config)
}

func (config *Config) parseFlags() {
	flag.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "report interval (example: 10s)")
	flag.DurationVar(&config.PollInterval, "p", config.PollInterval, "poll interval (example: 10s)")
	flag.StringVar(&config.HTTPClientConnection.ServerAddr, "a", config.HTTPClientConnection.ServerAddr, "server address (host:port)")
	flag.StringVar(&config.SignKey, "k", config.SignKey, "sign key")
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
