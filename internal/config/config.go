package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"strings"
)

type Config struct {
	LocalHost   string
	LocalPort   string
	ResultHost  string
	ResultPort  string
	FileName    string
	StoragePath string
}

var AppConfig = Config{
	LocalHost:   "http://localhost",
	LocalPort:   ":8080",
	ResultHost:  "http://localhost",
	ResultPort:  ":8080",
	StoragePath: "./storage/",
	FileName:    "data.txt",
}

// Load configuration from environment variables
func loadFromEnv() bool {
	var envVars struct {
		LocalAddr   string `env:"LOCAL_ADDRESS"`
		ResultAddr  string `env:"RESULT_ADDRESS"`
		StoragePath string `env:"FILE_STORAGE_PATH"`
	}

	if err := env.Parse(&envVars); err != nil {
		return false
	}

	if envVars.LocalAddr != "" {
		setHost("Local", envVars.LocalAddr)
	}
	if envVars.ResultAddr != "" {
		setHost("Result", envVars.ResultAddr)
	}
	if envVars.StoragePath != "" {
		AppConfig.StoragePath = envVars.StoragePath
	}

	return envVars.LocalAddr != "" && envVars.ResultAddr != "" && envVars.StoragePath != ""
}

// Load configuration from command-line flags
func loadFromFlags() bool {
	flag.Func("a", "The local hostname and port", func(value string) error {
		setHost("Local", value)
		return nil
	})

	flag.Func("b", "The result hostname and port", func(value string) error {
		setHost("Result", value)
		return nil
	})

	flag.Func("f", "The path for storing a file", func(value string) error {
		setPath(value)
		return nil
	})

	flag.Parse()
	return true
}

// Initialize configuration with priority: environment -> flags
func InitConfig() {
	if !loadFromEnv() {
		loadFromFlags()
	}
}

func GetHost(typeOf string) string {
	if typeOf == "Result" {
		return AppConfig.ResultHost + AppConfig.ResultPort
	}
	return AppConfig.LocalHost + AppConfig.LocalPort
}

func setHost(typeOf string, flagValue string) {
	flagValue = strings.TrimPrefix(flagValue, "http://")
	flagValue = strings.TrimPrefix(flagValue, "https://")

	h := strings.Split(flagValue, ":")

	if len(h) == 0 || h[0] == "" || h[1] == "" {
		return
	}

	domain := "http://" + h[0]
	port := ":" + h[1]

	if typeOf == "Result" {
		AppConfig.ResultHost = domain
		AppConfig.ResultPort = port
		return
	}
	AppConfig.LocalHost = domain
	AppConfig.LocalPort = port
}

func setPath(flagValue string) {
	// Trim any whitespace
	path := strings.TrimSpace(flagValue)

	// If the path doesn't end with a '/', append it
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	AppConfig.StoragePath = path
}
