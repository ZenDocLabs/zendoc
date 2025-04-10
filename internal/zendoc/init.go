package zendoc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dterbah/zendoc/config"
)

func InitZenDoc() error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Erreur:", err)
		return err
	}

	projectName := filepath.Base(cwd)

	defaultConfiguration := config.Config{
		ProjectConfig: config.ProjectConfig{
			Name:        projectName,
			Description: "Description of your project",
			Version:     "1.0",
		},
		DocConfig: config.DocConfig{
			IncludePrivate: false,
			IncludeTests:   false,
			IncludeMain:    false,
			ExcludeFiles:   []string{},
		},
	}

	return config.SaveConfiguration(defaultConfiguration)
}
