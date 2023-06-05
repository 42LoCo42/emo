package shared

import (
	"fmt"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token  []byte
	Server string
	Socket string
}

var config = &Config{}

func ConfigPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", Wrap(err, "could not get user config dir")
	}

	return path.Join(base, "emo.yml"), nil
}

func LoadConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return Wrap(err, "could not get config path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Print("Creating example config at ", path)

			config = &Config{
				Token:  []byte{},
				Server: "https://example.com",
				Socket: fmt.Sprintf("/run/user/%d/emo.socket", os.Getuid()),
			}

			if err := SaveConfig(); err != nil {
				return Wrap(err, "could not create example config")
			}

			return LoadConfig()
		}

		return Wrap(err, "could not read config")
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return Wrap(err, "could not load config")
	}

	return nil
}

func GetConfig() *Config {
	return config
}

func SaveConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return Wrap(err, "could not get config path")
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return Wrap(err, "could not marshal config")
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return Wrap(err, "could not save config")
	}

	return nil
}
