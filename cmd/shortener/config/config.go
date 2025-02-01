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

type envConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
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

	// Set the Address of the server
	// It is the address the users send requests to
	setServerAddress()

	//Set the URL that will return to the user
	//The base domain is used as a base for generated short URLs
	setBaseAddress()
}

func GetHost(typeOf string) string {
	if typeOf == "Result" {
		return Result.HostName + Result.Port
	}
	return Local.HostName + Local.Port
}

func setServerAddress() {
	//Use the system variable first
	var ec envConfig
	env.Parse(&ec)

	if ec.ServerAddress != "" {
		setHost("Local", ec.ServerAddress)
		return
	}

	//Check the flag in the command line
	flag.Func("a", "The hostname to bind the server to", func(flagValue string) error {
		setHost("Local", flagValue)
		return nil
	})
	flag.Parse()
}

func setBaseAddress() {
	//Use the system variable first
	var ec envConfig
	env.Parse(&ec)

	if ec.ServerAddress != "" {
		setHost("Local", ec.BaseURL)
		return
	}

	//check the flag in the command line
	flag.Func("b", "The result host name", func(flagValue string) error {
		setHost("Result", flagValue)
		return nil
	})
	flag.Parse()
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
