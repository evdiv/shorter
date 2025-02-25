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
	LocalHost   string
	ResultHost  string
	StoragePath string
}

var AppConfig = Config{
	LoadedFrom:  "default",
	LocalHost:   "http://localhost:8080",
	ResultHost:  "http://localhost:8080",
	StoragePath: "./tmp/data.txt",
}

// LoadFromEnv - loads from Environment variables
func LoadFromEnv() bool {

	var envVars struct {
		LocalAddr   string `env:"LOCAL_ADDRESS"`
		ResultAddr  string `env:"RESULT_ADDRESS"`
		StoragePath string `env:"FILE_STORAGE_PATH"`
	}

	if err := env.Parse(&envVars); err != nil {
		return false
	}

	if envVars.LocalAddr != "" {
		AppConfig.LocalHost = addPrefix(envVars.LocalAddr)
	}

	if envVars.ResultAddr != "" {
		AppConfig.ResultHost = addPrefix(envVars.ResultAddr)
	}

	if envVars.StoragePath != "" {
		AppConfig.StoragePath = envVars.StoragePath
	}

	if envVars.LocalAddr != "" && envVars.ResultAddr != "" && envVars.StoragePath != "" {
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

	flag.Parse()
	AppConfig.LoadedFrom = "flags"
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

//func NewConfig() {
//	if !LoadFromEnv() {
//		LoadFromFlags()
//	}
//}

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
