package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"net/url"
	"strings"
)

type Config struct {
	LocalHost    string `env:"LOCAL_ADDRESS"`
	ResultHost   string `env:"RESULT_ADDRESS"`
	StoragePath  string `env:"FILE_STORAGE_PATH"`
	DBConnection string `env:"DATABASE_DSN"`
	LoadedFrom   map[string]string
}

var AppConfig = Config{
	LocalHost:    "",
	ResultHost:   "",
	StoragePath:  "",
	DBConnection: "",
	LoadedFrom:   make(map[string]string),
}

// NewConfig - loads configs in the required order
func NewConfig(loaders ...func()) *Config {
	for _, loader := range loaders {
		// Load each loader, because not every loader may have all required parameters.
		loader()
	}
	return &AppConfig
}

// LoadFromEnv - loads from Environment variables
func LoadFromEnv() {

	if err := env.Parse(&AppConfig); err != nil {
		return
	}
	AppConfig.LocalHost = addPrefix(AppConfig.LocalHost)
	AppConfig.ResultHost = addPrefix(AppConfig.ResultHost)
	AppConfig.StoragePath = strings.TrimSpace(AppConfig.StoragePath)
	AppConfig.DBConnection = strings.TrimSpace(AppConfig.DBConnection)
}

// LoadFromFlags - loads from command-line flags
func LoadFromFlags() {

	if AppConfig.LocalHost == "" {
		flag.Func("a", "The hostname to bind the server to", func(value string) error {
			AppConfig.LocalHost = addPrefix(value)
			return nil
		})
	}

	if AppConfig.ResultHost == "" {
		flag.Func("b", "The result host name", func(value string) error {
			AppConfig.ResultHost = addPrefix(value)
			return nil
		})
	}

	if AppConfig.StoragePath == "" {
		flag.Func("f", "The path for storing a file", func(value string) error {
			AppConfig.StoragePath = strings.TrimSpace(value)
			return nil
		})
	}

	if AppConfig.DBConnection == "" {
		flag.Func("d", "Database connection string", func(value string) error {
			AppConfig.DBConnection = strings.TrimSpace(value)
			return nil
		})
	}
	flag.Parse()
}

func LoadDefault() {
	if AppConfig.LocalHost == "" {
		AppConfig.LocalHost = "http://localhost:8080"
	}
	if AppConfig.ResultHost == "" {
		AppConfig.ResultHost = "http://localhost:8080"
	}
	if AppConfig.StoragePath == "" {
		AppConfig.StoragePath = "./tmp/data.txt"
	}
	//if AppConfig.DBConnection == "" {
	//	//AppConfig.DBConnection = "postgres://postgres:55555@localhost:5432/postgres"
	//}
}

func GetPort(typeOf string) string {
	if typeOf == "Local" {
		return extractPort(AppConfig.LocalHost)
	}
	return extractPort(AppConfig.ResultHost)
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
