package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Config struct {
	App          AppConfig     `yaml:"app"`
	Bot          BotConfig     `yaml:"bot"`
	OpenAIConfig OpenAIConfig  `yaml:"openAIConfig"`
	Metrics      MetricsConfig `yaml:"metrics"`
	Tracing      TracingConfig `yaml:"tracing"`
}

type AppConfig struct {
	Id        string `yaml:"id" env:"APP_ID"`
	Name      string `yaml:"name" env:"APP_NAME"`
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL"`
	IsLogJSON bool   `yaml:"is_log_json" env:"IS_LOG_JSON"`
}

type BotConfig struct {
	Token   string        `yaml:"token" env:"BOT_TOKEN"`
	Timeout time.Duration `yaml:"timeout" env:"BOT_TIMEOUT"`
}

type OpenAIConfig struct {
	Enabled bool   `yaml:"enabled" env:"OPENAI_ENABLED"`
	ApiKey  string `yaml:"api_key" env:"OPENAI_API_KEY"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" env:"METRICS_ENABLED"`
	Host    string `yaml:"host" env:"METRICS_HOST"`
	Port    int    `yaml:"port" env:"METRICS_PORT"`
}

type TracingConfig struct {
	Enabled bool   `yaml:"enabled" env:"TRACING_ENABLED"`
	Host    string `yaml:"host" env:"TRACING_HOST"`
	Port    int    `yaml:"port" env:"TRACING_PORT"`
}

const (
	flagConfigPathName = "config"
	envConfigPathName  = "CONFIG_PATH"
)

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		var configPath string
		flag.StringVar(&configPath, flagConfigPathName, "", "config file path")
		flag.Parse()

		if path, ok := os.LookupEnv(envConfigPathName); ok {
			configPath = path
		}

		instance = &Config{}

		if readErr := cleanenv.ReadConfig(configPath, instance); readErr != nil {
			description, descErr := cleanenv.GetDescription(instance, nil)
			if descErr != nil {
				panic(descErr)
			}

			slog.Info(description)
			slog.Error("failed to parse config",
				slog.String("error", readErr.Error()),
				slog.String("path", configPath),
			)

			os.Exit(1)
		}
	})

	return instance
}
