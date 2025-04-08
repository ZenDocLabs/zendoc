package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const ZENDOC_CONFIG_FILE = ".zendoc.config.json"

type ProjectConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type DocConfig struct {
	IncludePrivate bool `json:"includePrivate"`
	IncludeTests   bool `json:"includeTests"`
}

type Config struct {
	ProjectConfig ProjectConfig `json:"projectConfig"`
	DocConfig     DocConfig     `json:"docConfig"`
}

func GetConfiguration() (*Config, error) {
	config := Config{}

	fileBytes, err := os.ReadFile(ZENDOC_CONFIG_FILE)
	if err != nil {
		return nil, fmt.Errorf("error when reading the config file %s", err)
	}

	err = json.Unmarshal(fileBytes, &config)

	if err != nil {
		return nil, fmt.Errorf("error when loading the config file %s", err)
	}

	return &config, nil
}

func SaveConfiguration(config Config) error {
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(ZENDOC_CONFIG_FILE, content, 0644)

	if err != nil {
		return fmt.Errorf("error when saving the config of ZenDoc %s", err)
	}

	return nil
}
