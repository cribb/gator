package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

const configFilename = ".gatorconfig.json"

func Read() (Config, error) {
	fmt.Println("Reading config...")
	configFilename, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Config file (%v) not found.\n", configFilename)
		return Config{}, err
	}

	configFileContents, err := os.ReadFile(configFilename)
	if err != nil {
		fmt.Printf("Error reading config file (%v) contents.\n", configFilename)
		return Config{}, err
	}
	// fmt.Printf("Config file contents: %v\n", configFileContents)

	var cfg Config
	if err := json.Unmarshal(configFileContents, &cfg); err != nil {
		fmt.Printf("err: %v\n", err)
		return Config{}, err
	}

	// fmt.Printf("Config file contents: %+v\n", config)

	return cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUsername = username
	return write(*cfg)
}

func write(config Config) error {
	// fmt.Printf("Current config:\n  %+v\n", config)

	configFilepath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Data marshaling error: %v\n", err)
		return err
	}

	if err := os.WriteFile(configFilepath, []byte(data), 0644); err != nil {
		fmt.Printf("Error writing file: %v.\n", err)
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := userHome + "/" + configFilename

	_, err = os.Stat(configFilePath)
	if err != nil {
		fmt.Printf("Config file (%v) not found.\n", configFilePath)
		return "", err

	}

	// fmt.Printf("User homedir: %v\n", userHome)
	// fmt.Printf("Config file found at: %v\n", configFilePath)

	return configFilePath, nil

}
