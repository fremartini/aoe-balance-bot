package config

import (
	"aoe-bot/internal/logger"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type config struct {
	Token    string
	LogLevel uint
}

const CONFIG_FILE = ".config"

func Read() (*config, error) {
	if _, err := os.Stat(CONFIG_FILE); errors.Is(err, os.ErrNotExist) {
		return readConfigFromEnv()
	}

	return readConfigFromFile()
}

func readConfigFromEnv() (*config, error) {
	token := os.Getenv("token")
	logLevel := os.Getenv("logLevel")

	if token == "" {
		return nil, errors.New("token supplied")
	}

	var level uint = logger.INFO
	if logLevel != "" {
		l, err := parseUint(logLevel)

		if err != nil {
			return nil, err
		}

		level = l
	}

	return &config{
		Token:    token,
		LogLevel: level,
	}, nil
}

func readConfigFromFile() (*config, error) {
	content, err := os.ReadFile(CONFIG_FILE)

	if err != nil {
		return nil, err
	}

	config := &config{}
	return config, json.Unmarshal(content, config)
}

func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return 0, err
	}

	return uint(v), err
}
