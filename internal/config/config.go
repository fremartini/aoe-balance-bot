package config

import (
	"aoe-bot/internal/logger"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type config struct {
	Token    *string
	LogLevel *uint
	Cache    *cache
}

type cache struct {
	ExpiryHours *uint
	MaxSize     *uint
}

const (
	CONFIG_FILE                 = ".config"
	DEF_LOG_LEVEL          uint = logger.INFO
	DEF_CACHE_EXPIRY_HOURS uint = 24
	DEF_CACHE_SIZE         uint = 20
)

func Read() (*config, error) {
	if _, err := os.Stat(CONFIG_FILE); errors.Is(err, os.ErrNotExist) {
		return readConfigFromEnv()
	}

	return readConfigFromFile()
}

func readConfigFromEnv() (*config, error) {
	tokenEnv := os.Getenv("token")
	logLevelEnv := os.Getenv("logLevel")
	cacheExpiryHoursEnv := os.Getenv("cacheExpiryHours")
	cacheMaxSizeEnv := os.Getenv("cacheMaxSize")

	if tokenEnv == "" {
		return nil, errors.New("token not supplied")
	}

	logLevel, err := uintOrDefault(logLevelEnv, DEF_LOG_LEVEL)

	if err != nil {
		return nil, err
	}

	cacheExpiryHours, err := uintOrDefault(cacheExpiryHoursEnv, DEF_CACHE_EXPIRY_HOURS)

	if err != nil {
		return nil, err
	}

	cacheMaxSize, err := uintOrDefault(cacheMaxSizeEnv, DEF_CACHE_SIZE)

	if err != nil {
		return nil, err
	}

	return &config{
		Token:    &tokenEnv,
		LogLevel: &logLevel,
		Cache: &cache{
			ExpiryHours: &cacheExpiryHours,
			MaxSize:     &cacheMaxSize,
		},
	}, nil
}

func readConfigFromFile() (*config, error) {
	content, err := os.ReadFile(CONFIG_FILE)

	if err != nil {
		return nil, err
	}

	config := &config{}

	err = json.Unmarshal(content, config)

	if config.Token == nil {
		return nil, errors.New("token not supplied")
	}

	if config.Cache.ExpiryHours == nil {
		hrs := DEF_CACHE_EXPIRY_HOURS
		config.Cache.ExpiryHours = &hrs
	}

	if config.Cache.MaxSize == nil {
		size := DEF_CACHE_SIZE
		config.Cache.MaxSize = &size
	}

	if config.LogLevel == nil {
		lvl := DEF_LOG_LEVEL
		config.LogLevel = &lvl
	}

	return config, err
}

func uintOrDefault(v string, def uint) (uint, error) {
	if v == "" {
		// no env value provided
		return def, nil
	}

	// env value is provided
	l, err := parseUint(v)

	if err != nil {
		return def, err
	}

	return l, err
}

func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return 0, err
	}

	return uint(v), err
}
