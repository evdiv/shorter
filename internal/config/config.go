package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"strings"
)

type Config struct {
	HostName string
	Port     string
}

var Local = Config{
	HostName: "http://localhost",
	Port:     ":8080",
}

var Result = Config{
	HostName: "http://localhost",
	Port:     ":8080",
}

// ConfigLoader is an interface for configuration loading strategies
type ConfigLoader interface {
	Load() bool
}

// EnvConfigLoader loads configuration from environment variables
type EnvConfigLoader struct{}

func (e EnvConfigLoader) Load() bool {
	type envConfig struct {
		ServerAddress string `env:"SERVER_ADDRESS"`
		BaseURL       string `env:"BASE_URL"`
	}

	var ec envConfig
	err := env.Parse(&ec)
	if err != nil {
		return false
	}

	if ec.ServerAddress != "" && ec.BaseURL != "" {
		setHost("Local", ec.ServerAddress)
		setHost("Result", ec.BaseURL)
		return true
	}
	return false
}

// FlagConfigLoader loads configuration from command-line flags
type FlagConfigLoader struct{}

func (f FlagConfigLoader) Load() bool {
	flag.Func("a", "The hostname to bind the server to", func(flagValue string) error {
		setHost("Local", flagValue)
		return nil
	})

	flag.Func("b", "The result host name", func(flagValue string) error {
		setHost("Result", flagValue)
		return nil
	})
	flag.Parse()
	return true
}

// NewConfig uses dependency injection to load configuration
func NewConfig(loaders ...ConfigLoader) {
	for _, loader := range loaders {
		if loader.Load() {
			return
		}
	}
}

func GetHost(typeOf string) string {
	if typeOf == "Result" {
		return Result.HostName + Result.Port
	}
	return Local.HostName + Local.Port
}

func setHost(typeOf string, flagValue string) {
	// Remove "http://" or "https:// from the flag if exists
	flagValue = strings.TrimPrefix(flagValue, "http://")
	flagValue = strings.TrimPrefix(flagValue, "https://")

	h := strings.Split(flagValue, ":")

	if len(h) == 0 || h[0] == "" || h[1] == "" {
		return
	}

	domain := "http://" + h[0]
	port := ":" + h[1]

	if typeOf == "Result" {
		Result.HostName = domain
		Result.Port = port
		return
	}
	Local.HostName = domain
	Local.Port = port
}
