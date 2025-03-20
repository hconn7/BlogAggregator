package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL string `json:"db_url"`
	User  string `json:"current_user_name"`
}

func ReadDB() (Config, error) {
	filePath, _ := getConfigFilePath()

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func WriteFile(cfg *Config) error {
	filePath, _ := getConfigFilePath()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.User = userName
	WriteFile(cfg)
	return nil
}
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	jsonFilePath := filepath.Join(homeDir, configFileName)
	return jsonFilePath, nil
}
