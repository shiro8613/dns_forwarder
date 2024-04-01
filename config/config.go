package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const path = "./config.yml"
const ipv4Regex = "^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])(\\:([0-35565]{2}))?"

var conf Config

func Load() error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, &conf); err != nil {
		return err
	}

	if err := splitV4V6(&conf); err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return conf
}
