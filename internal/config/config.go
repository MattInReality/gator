package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName string = ".gatorconfig.json"

var path string = getFullConfigFilePath(configFileName)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(currentUser string) error {
	c.CurrentUserName = currentUser
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	var config Config
	configFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	bytes, _ := io.ReadAll(configFile)
	json.Unmarshal(bytes, &config)

	return config, nil
}

func getFullConfigFilePath(configFileName string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s", home, configFileName)
}

func write(cfg Config) error {
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer configFile.Close()
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = configFile.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
