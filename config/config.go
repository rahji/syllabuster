package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Scale struct {
	Letter string  `yaml:"Letter"`
	Min    float64 `yaml:"Min"`
}

type Config struct {
	Scale       []Scale  `yaml:"scale"`
	Assignments []string `yaml:"assignments"`
}

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var (
	k      = koanf.New(".")
	parser = yaml.Parser()
)

func ReadConfig(fn string) (Config, error) {
	if err := k.Load(file.Provider(fn), parser); err != nil {
		return Config{}, fmt.Errorf("error loading config: %s", err)
	}

	var conf Config
	if err := k.Unmarshal("", &conf); err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config: %s", err)
	}

	if conf.Scale == nil {
		return Config{}, fmt.Errorf("no scale defined (check configuration file)")
	}

	return conf, nil
}
