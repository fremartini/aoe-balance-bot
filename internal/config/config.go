package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type config struct {
	Token          string
	LogLevel       uint
	SteamIdChannel string
}

const CONFIG_FILE = ".config"

func Read() (*config, error) {
	if _, err := os.Stat(CONFIG_FILE); errors.Is(err, os.ErrNotExist) {

		token := os.Getenv("token")
		logLevel := os.Getenv("logLevel")
		steamIdChannel := os.Getenv("steamIdChannel")

		if token == "" || logLevel == "" || steamIdChannel == "" {
			return nil, errors.New("token, logLevel and steamIdChannel must be supplied to run without config file")
		}

		level, err := strconv.ParseUint(logLevel, 10, 64)

		if err != nil {
			return nil, err
		}

		return &config{
			Token:          token,
			LogLevel:       uint(level),
			SteamIdChannel: steamIdChannel,
		}, nil
	}

	content, err := os.ReadFile(CONFIG_FILE)

	if err != nil {
		return nil, err
	}

	config := &config{}
	return config, json.Unmarshal(content, config)
}
