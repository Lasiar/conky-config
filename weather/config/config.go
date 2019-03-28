package config

import (
	"github.com/go-yaml/yaml"
	"log"
	"os"
	"os/user"
	"path"
	"sync"
)

type Config struct {
	ApiKey string
	CityID int
}

var (
	_once   sync.Once
	_config *Config
)

func loadConfig(confgPath string) *Config {
	configFile, err := os.Open(confgPath)
	if err != nil {
		log.Panicf("Error open config file [path]: %v: [error]: %v", confgPath, err)
	}

	rawConfig := struct {
		Weather struct {
			ApiKey string `yaml:"api_key"`
			CityID int    `yaml:"city_id"`
		}
	}{}

	if err := yaml.NewDecoder(configFile).Decode(&rawConfig); err != nil {
		log.Panicf("Error decode config: %v", err)
	}

	return &Config{rawConfig.Weather.ApiKey, rawConfig.Weather.CityID}
}

func GetConfig() *Config {
	_once.Do(func() {
		u, err := user.Current()
		if err != nil {
			log.Panicf("err read user: %v", err)
		}

		pathConfig := path.Join(u.HomeDir, ".conky/secrets.yml")

		_config = new(Config)
		_config = loadConfig(pathConfig)
	})
	return _config
}
