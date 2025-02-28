package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"net/url"
	"strings"
)

type Config struct {
	LoadedFrom  string
	LocalHost   string `env:"LOCAL_ADDRESS"`
	ResultHost  string `env:"RESULT_ADDRESS"`
	StoragePath string `env:"FILE_STORAGE_PATH"`
	DbDSN       string `env:"DATABASE_DSN"`
}

var AppConfig = Config{
	LoadedFrom:  "",
	LocalHost:   "",
	ResultHost:  "",
	StoragePath: "",
	DbDSN:       "",
}

// LoadFromEnv - loads from Environment variables
func LoadFromEnv() bool {

	if err := env.Parse(&AppConfig); err != nil {
		return false
	}

	AppConfig.LocalHost = addPrefix(AppConfig.LocalHost)
	AppConfig.ResultHost = addPrefix(AppConfig.ResultHost)

	if AppConfig.LocalHost != "" && AppConfig.ResultHost != "" && AppConfig.StoragePath != "" {
		AppConfig.LoadedFrom = "environment"
		return true
	}
	return false
}

// LoadFromFlags - loads from command-line flags
func LoadFromFlags() bool {
	flag.Func("a", "The hostname to bind the server to", func(value string) error {
		AppConfig.LocalHost = addPrefix(value)
		return nil
	})

	flag.Func("b", "The result host name", func(value string) error {
		AppConfig.ResultHost = addPrefix(value)
		return nil
	})

	flag.Func("f", "The path for storing a file", func(value string) error {
		setPath(value)
		return nil
	})

	flag.Func("d", "Database connection string", func(value string) error {
		AppConfig.DbDSN = value
		return nil
	})

	flag.Parse()
	AppConfig.LoadedFrom = "flags"
	return true
}

func LoadDefault() bool {
	AppConfig.LoadedFrom = "default"
	AppConfig.LocalHost = "http://localhost:8080"
	AppConfig.ResultHost = "http://localhost:8080"
	AppConfig.StoragePath = "./tmp/data.txt"
	AppConfig.DbDSN = "postgres://username:password@localhost:5432/mydb"
	return true
}

// NewConfig load configs in the required order
func NewConfig(loaders ...func() bool) {
	for _, loader := range loaders {
		success := loader()
		if success {
			return // Stop at the first successful loader
		}
	}
	fmt.Println("No valid configuration found, using defaults")
}

func GetPort(typeOf string) string {
	if typeOf == "Local" {
		return extractPort(AppConfig.LocalHost)
	}
	return extractPort(AppConfig.ResultHost)
}

func GetHost(typeOf string) string {
	if typeOf == "Local" {
		return extractHost(AppConfig.LocalHost)
	}
	return extractHost(AppConfig.ResultHost)
}

func addPrefix(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return host
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}
	return host
}

func extractPort(address string) string {
	parsed, err := url.Parse(address)
	if err != nil {
		return ""
	}
	if parsed.Port() != "" {
		return ":" + parsed.Port()
	}
	return ""
}

func extractHost(address string) string {
	parsed, err := url.Parse(address)
	if err != nil {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Hostname()
}

func setPath(path string) {
	path = strings.TrimSpace(path)
	AppConfig.StoragePath = path
}
