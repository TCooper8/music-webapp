package main

import (
	"encoding/json"
	"errors"
	"os"
)

type configState struct {
	HttpPort     int
	HttpHostname string
	LogLevel     string
}

type Config struct {
	state configState
}

var config *Config

func init() {
	file, err := os.Open("config/default.json")
	if err != nil {
		panic(err)
	}

	config = new(Config)

	decoder := json.NewDecoder(file)
	config = new(Config)

	err = decoder.Decode(&config.state)
	if err != nil {
		panic(err)
	}

	config.GetLogLevel()
}

func GetConfig() *Config {
	return config
}

func (config *Config) GetHttpHostname() string {
	return config.state.HttpHostname
}

func (config *Config) GetHttpPort() int {
	return config.state.HttpPort
}

func (config *Config) GetLogLevel() int {
	switch config.state.LogLevel {
	case "FATAL":
		return LOG_FATAL
	case "ERROR":
		return LOG_ERROR
	case "WARN":
		return LOG_WARN
	case "INFO":
		return LOG_INFO
	case "DEBUG":
		return LOG_DEBUG
	case "TRACE":
		return LOG_TRACE
	default:
		panic(errors.New("Invalid logLevel"))
	}
}
