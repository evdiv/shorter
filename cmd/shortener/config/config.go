package config

import (
	"flag"
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

	flag.Func("a", "The hostname to bind the server to", func(flagValue string) error {
		setHost("Local", flagValue)
		return nil
	})

	flag.Func("b", "The result host name", func(flagValue string) error {
		setHost("Result", flagValue)
		return nil
	})

	flag.Parse()
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
