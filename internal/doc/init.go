package doc

import (
	"os"
	"path/filepath"

	"github.com/dterbah/zendoc/config"
	"github.com/dterbah/zendoc/internal/system"
)

func InitZenDoc() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectName := filepath.Base(cwd)

	defaultConfiguration := config.Config{
		ProjectConfig: config.ProjectConfig{
			Name:        projectName,
			Description: "Description of your project",
			Version:     "1.0",
			GitLink:     "",
		},
		DocConfig: config.DocConfig{
			IncludePrivate: false,
			IncludeTests:   false,
			IncludeMain:    false,
			ExcludeFiles:   []string{},
		},
	}

	return config.SaveConfiguration(defaultConfiguration, system.OSFileSystem{})
}
