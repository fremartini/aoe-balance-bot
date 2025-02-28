package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

var (
	ErrTokenNotSupplied = errors.New("token not supplied")

	tokenEntry = configEntry[*string]{
		Key:          "token",
		DefaultValue: nil,
		ParseFunc:    func(s string) *string { return &s },
	}

	cacheExpiryHoursEntry = configEntry[uint]{
		Key:          "cacheExpiryHours",
		DefaultValue: 24,
		ParseFunc:    parseUint,
	}

	cacheMaxSizeEntry = configEntry[uint]{
		Key:          "cacheMaxSize",
		DefaultValue: 20,
		ParseFunc:    parseUint,
	}

	portEntry = configEntry[*uint]{
		Key:          "port",
		DefaultValue: nil,
		ParseFunc: func(s string) *uint {
			r := parseUint(s)
			return &r
		},
	}

	trustInsecureCertificatesEntry = configEntry[bool]{
		Key:          "trustInsecureCertificates",
		DefaultValue: false,
		ParseFunc: func(s string) bool {
			r, err := strconv.ParseBool(s)

			if err != nil {
				return false
			}

			return r
		},
	}

	prefixEntry = configEntry[string]{
		Key:          "prefix",
		DefaultValue: "!",
		ParseFunc:    func(s string) string { return s },
	}

	whitelistedChannelsEntry = configEntry[[]uint]{
		Key:          "whitelistedChannels",
		DefaultValue: []uint{},
		ParseFunc: func(s string) []uint {
			var arr []uint
			if err := json.Unmarshal([]byte(s), &arr); err != nil {
				return nil
			}
			return arr
		},
	}
)

type config struct {
	Token string
	Cache *struct {
		ExpiryHours uint
		MaxSize     uint
	}
	Port                      *uint
	TrustInsecureCertificates bool
	WhitelistedChannels       []uint
	Prefix                    string
}

type configEntry[K any] struct {
	Key          string
	DefaultValue K
	ParseFunc    func(string) K
}

func Read() (*config, error) {
	token := environmentValueOrDefault(tokenEntry)
	cacheExpiryHours := environmentValueOrDefault(cacheExpiryHoursEntry)
	cacheMaxSize := environmentValueOrDefault(cacheMaxSizeEntry)
	port := environmentValueOrDefault(portEntry)
	trustInsecureCertificates := environmentValueOrDefault(trustInsecureCertificatesEntry)
	whitelistedChannels := environmentValueOrDefault(whitelistedChannelsEntry)
	prefix := environmentValueOrDefault(prefixEntry)

	if token == nil {
		return nil, ErrTokenNotSupplied
	}

	return &config{
		Token: *token,
		Cache: &struct {
			ExpiryHours uint
			MaxSize     uint
		}{
			ExpiryHours: cacheExpiryHours,
			MaxSize:     cacheMaxSize,
		},
		Port:                      port,
		TrustInsecureCertificates: trustInsecureCertificates,
		WhitelistedChannels:       whitelistedChannels,
		Prefix:                    prefix,
	}, nil
}

func environmentValueOrDefault[K any](e configEntry[K]) K {
	v := os.Getenv(e.Key)

	if v == "" {
		return e.DefaultValue
	}

	return e.ParseFunc(v)
}

func parseUint(s string) uint {
	v, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		panic(err)
	}

	return uint(v)
}
