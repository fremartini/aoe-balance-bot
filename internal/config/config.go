package config

import (
	"aoe-bot/internal/logger"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

const (
	ConfigFile               = ".config"
	DefLogLevel         uint = logger.INFO
	DefCacheExpiryHours uint = 24
	DefCacheSize        uint = 20
)

var (
	ErrTokenNotSupplied = errors.New("token not supplied")
)

type config struct {
	Token    string
	LogLevel uint
	Cache    *struct {
		ExpiryHours uint
		MaxSize     uint
	}
	Port *uint
}

func Read() (*config, error) {
	if _, err := os.Stat(ConfigFile); errors.Is(err, os.ErrNotExist) {
		return readFromEnv()
	}

	return readFromFile()
}

func readFromEnv() (*config, error) {
	token := envValueOrDefault("token", nil, func(s string) *string { return &s })
	logLevel := envValueOrDefault("logLevel", DefLogLevel, parseUint)
	cacheExpiryHours := envValueOrDefault("cacheExpiryHours", DefCacheExpiryHours, parseUint)
	cacheMaxSize := envValueOrDefault("cacheMaxSize", DefCacheSize, parseUint)
	port := envValueOrDefault("PORT", nil, func(s string) *uint {
		r := parseUint(s)
		return &r
	})

	if token == nil {
		return nil, ErrTokenNotSupplied
	}

	return &config{
		Token:    *token,
		LogLevel: logLevel,
		Cache: &struct {
			ExpiryHours uint
			MaxSize     uint
		}{
			ExpiryHours: cacheExpiryHours,
			MaxSize:     cacheMaxSize,
		},
		Port: port,
	}, nil
}

type fileConfig struct {
	Token    *string
	LogLevel *uint
	Cache    *struct {
		ExpiryHours *uint
		MaxSize     *uint
	}
	Port *uint
}

func readFromFile() (*config, error) {
	content, err := os.ReadFile(ConfigFile)

	if err != nil {
		return nil, err
	}

	fileConfig := &fileConfig{}

	err = json.Unmarshal(content, fileConfig)

	if err != nil {
		return nil, err
	}

	if fileConfig.Token == nil {
		return nil, ErrTokenNotSupplied
	}

	return &config{
		Token:    *fileConfig.Token,
		LogLevel: valueOrDefault(fileConfig.LogLevel, DefLogLevel),
		Cache: &struct {
			ExpiryHours uint
			MaxSize     uint
		}{
			ExpiryHours: valueOrDefault(fileConfig.Cache.ExpiryHours, DefCacheExpiryHours),
			MaxSize:     valueOrDefault(fileConfig.Cache.MaxSize, DefCacheSize),
		},
		Port: fileConfig.Port,
	}, nil
}

func envValueOrDefault[K any](env string, def K, f func(string) K) K {
	v := os.Getenv(env)

	if v == "" {
		return def
	}

	return f(v)
}

func valueOrDefault[K any](a *K, def K) K {
	if a == nil {
		return def
	}

	return *a
}

func parseUint(s string) uint {
	v, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		panic(err)
	}

	return uint(v)
}
