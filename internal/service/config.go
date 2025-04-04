package service

import (
	"os"
	"strings"
)

type Config struct {
	props map[string]string
}

func (c *Config) Get(key string) (string, bool) {
	value, ok := c.props[key]
	return value, ok
}

func ReadConfig(filePath string) (*Config, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	keys := strings.Split(string(dat), "\n")
	config := Config{map[string]string{}}
	for line := range keys {
		pair := strings.Split(keys[line], "=")
		if len(pair) == 2 {
			config.props[pair[0]] = pair[1]
		}
	}
	return &config, nil
}
