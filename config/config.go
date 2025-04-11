package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dterbah/zendoc/internal/system"
)

const ZENDOC_CONFIG_FILE = ".zendoc.config.json"

type ProjectConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	GitLink     string `json:"gitLink"`
	MainBranch  string `json:"mainBranch"`
	DocPath     string `json:"docPath"`
}

type DocConfig struct {
	IncludePrivate bool     `json:"includePrivate"`
	IncludeTests   bool     `json:"includeTests"`
	IncludeMain    bool     `json:"includeMain"`
	ExcludeFiles   []string `json:"excludeFiles"`
}

type Config struct {
	ProjectConfig ProjectConfig `json:"projectConfig"`
	DocConfig     DocConfig     `json:"docConfig"`
}

/*
@description Load the ZenDoc configuration from the configuration file
@return (*Config, error) - A pointer to the loaded configuration and an error if loading fails
@example GetConfiguration() => &Config{...}, nil
@author Dorian TERBAH
*/
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

/*
@description Save the given ZenDoc configuration to the configuration file
@param config Config - The configuration to save
@return error - An error if the saving process fails, otherwise nil
@example SaveConfiguration(myConfig)
@author Dorian TERBAH
*/

func SaveConfiguration(config Config, fs system.FileSystem) error {
	content, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = fs.WriteFile(ZENDOC_CONFIG_FILE, content, 0644)

	if err != nil {
		return fmt.Errorf("error when saving the config of ZenDoc %s", err)
	}

	return nil
}
