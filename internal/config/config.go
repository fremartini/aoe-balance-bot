package config

import (
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

		token := os.Getenv("token")
		logLevel := os.Getenv("level")

		if token == "" || logLevel == "" {
			return nil, errors.New("both token and logLevel must be supplied to run without config file")
		}

		level, err := strconv.ParseUint(logLevel, 10, 64)

		if err != nil {
			return nil, err
		}

		return &config{
			Token:    token,
			LogLevel: uint(level),
		}, nil
	}

	content, err := os.ReadFile(CONFIG_FILE)

	if err != nil {
		return nil, err
	}

	config := &config{}
	return config, json.Unmarshal(content, config)
}
