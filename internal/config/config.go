package config

import (
	"encoding/json"
	"os"
)

type config struct {
	Token    string
	LogLevel uint
}

const CONFIG_FILE = ".config"

func Read() (*config, error) {
	content, err := os.ReadFile(CONFIG_FILE)

	if err != nil {
		return nil, err
	}

	config := &config{}
	return config, json.Unmarshal(content, config)
}
