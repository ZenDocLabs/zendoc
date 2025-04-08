package zendoc

import (
	"fmt"
	"os"
	"path"

	"github.com/dterbah/zendoc/config"
)

func InitZenDoc() error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Erreur:", err)
		return err
	}

	projectName := path.Base(cwd)

	defaultConfiguration := config.Config{
		ProjectConfig: config.ProjectConfig{
			Name:        projectName,
			Description: "Description of your project",
			Version:     "1.0",
		},
		DocConfig: config.DocConfig{
			IncludePrivate: false,
			IncludeTests:   false,
		},
	}

	return config.SaveConfiguration(defaultConfiguration)
}
