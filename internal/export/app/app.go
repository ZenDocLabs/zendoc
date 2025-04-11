package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dterbah/zendoc/internal/export/helper"
)

type AppConfig struct {
	Versions    []string `json:"versions"`
	Description string   `json:"description"`
}

func UpdateAppConfig(appPath, currentVersion, description string) error {
	if !helper.IsFileExist(appPath) {
		config := AppConfig{
			Versions:    []string{currentVersion},
			Description: description,
		}

		return saveAppConfig(appPath, config)
	}

	data, err := os.ReadFile(appPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("version file does not exist: %w", err)
		}
		return fmt.Errorf("error reading version file: %w", err)
	}

	var config AppConfig

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	updateVersion := true
	for _, v := range config.Versions {
		if v == currentVersion {
			updateVersion = true
			break
		}
	}

	if updateVersion {
		config.Versions = append(config.Versions, currentVersion)
	}

	return saveAppConfig(appPath, config)
}

func saveAppConfig(appPath string, config AppConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(appPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
