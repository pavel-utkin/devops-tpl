package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

// HTTPClientConfig используется для хранения конфигурации агента, связанной с настройкой http клиента.
type HTTPClientConfig struct {
	// RetryCount - количество попыток отправки запроса (default: 2)
	RetryCount int `env:"RETRY_CONN_COUNT"`
	// RetryWaitTime - время ожидания между попытками отправки запроса (default: 10s)
	RetryWaitTime time.Duration `env:"RETRY_CONN_WAIT_TIME"`
	// RetryMaxWaitTime - макс. время для попыток отправки запроса (default: 90s)
	RetryMaxWaitTime time.Duration `env:"RETRY_CONN_MAX_WAIT_TIME"`
	// ServerAddr - адрес сервера (default: 127.0.0.1:8080)
	ServerAddr string `env:"ADDRESS"`
}

// Config используется для хранения конфигурации агента.
type Config struct {
	// PollInterval - интервал между считыванием метрик (flag: p; default: 2s)
	PollInterval time.Duration `env:"POLL_INTERVAL"`
	// ReportInterval - интервал между отправки метрик (flag: r; default: 2s)
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	// SignKey - ключ для подписи сообщений (flag: k)
	SignKey string `env:"KEY"`
	// RateLimit
	RateLimit int `env:"RATE_LIMIT"`
	// LogFile - лог файл (flag: l)
	LogFile string `env:"LOG_FILE"`
	// DebugMode - debug мод (flag: d)
	DebugMode            bool `env:"DEBUG"`
	HTTPClientConnection HTTPClientConfig
}

// initDefaultValues - значения конфига по умолчанию.
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
	flag.IntVar(&config.RateLimit, "l", config.RateLimit, "number of concurrent requests to the server")
	flag.BoolVar(&config.DebugMode, "d", config.DebugMode, "debug mode \"\"")
	flag.Parse()
}

func LoadConfig() Config {
	config := newConfig()

	config.parseFlags()
	err := config.parseEnv()

	if config.RateLimit == 0 {
		config.RateLimit = 1
	}

	if err != nil {
		log.Fatal(err)
	}

	return *config
}
