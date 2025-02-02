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

func LoadConfig() {
	//Use the system variables first
	set := useEnvForConfig()
	if !set {
		useFlagsForConfig()
	}
}

func GetHost(typeOf string) string {
	if typeOf == "Result" {
		return Result.HostName + Result.Port
	}
	return Local.HostName + Local.Port
}

func useFlagsForConfig() {

	// -a is a flag to set the Address of the server users send requests to
	flag.Func("a", "The hostname to bind the server to", func(flagValue string) error {
		setHost("Local", flagValue)
		return nil
	})

	// -b is a flag to set the URL that will be used as a base for generated short URLs
	flag.Func("b", "The result host name", func(flagValue string) error {
		setHost("Result", flagValue)
		return nil
	})
	flag.Parse()
}

func useEnvForConfig() bool {
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
