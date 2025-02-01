package config

import (
	"flag"
	"strings"
)

var Local struct {
	HostName string
	Port     string
}

var Result struct {
	HostName string
	Port     string
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
	h := strings.Split(flagValue, ":")

	domain := "http://localhost"
	if h[0] != "" {
		domain = "http://" + h[0]
	}

	port := ":8080"
	if h[1] != "" {
		port = ":" + h[1]
	}

	if typeOf == "Result" {
		Result.HostName = domain
		Result.Port = port
		return
	}
	Local.HostName = domain
	Local.Port = port
}
