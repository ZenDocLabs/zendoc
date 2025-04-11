package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dterbah/zendoc/internal/doc"
	"github.com/dterbah/zendoc/internal/export/app"
	"github.com/dterbah/zendoc/internal/export/helper"
	"github.com/dterbah/zendoc/internal/system"
	"github.com/fatih/color"
)

const TEMPLATE_GIT_LINK = "https://github.com/ZenDocLabs/zendoc-ui-template"

/*
@description Struct that implements the DocExporter interface and exports the documentation in a web-friendly format.
@author Dorian TERBAH
@field DocExporter DocExporter - Embedded base exporter providing common exporting behavior.
*/
type WebExporter struct {
	DocExporter
	GitLink     string
	AppName     string
	MainBranch  string
	DocPath     string
	Version     string
	Description string
	FileSystem  system.FileSystem
	CmdRunner   system.CommandRunner
}

/*
@description Export the project documentation as JSON to stdout
@param projectDoc doc.ProjectDoc - The documentation to export
@return error - An error if the export fails
@example WebExporter{}.Export(projectDoc)
*/
func (webExport WebExporter) Export(projectDoc doc.ProjectDoc) error {
	b, err := json.Marshal(projectDoc)
	if err != nil {
		return fmt.Errorf("error when exporting the documentation in JSON: %w", err)
	}

	currentPath, _ := os.Getwd()
	docPath := filepath.Join(currentPath, webExport.DocPath, webExport.AppName)

	if err := webExport.ensureTemplate(docPath); err != nil {
		return err
	}

	if err := webExport.updateAppConfig(docPath, webExport.Version, webExport.Description); err != nil {
		return err
	}

	color.Green("App config file updated !")

	if err := webExport.writeDocumentationFile(docPath, b); err != nil {
		return err
	}

	if err := webExport.writeEnvFile(docPath, webExport.GitLink, webExport.AppName, webExport.MainBranch); err != nil {
		return err
	}

	color.Green("Env file updated !")

	color.Green("Documentation v%s saved!", webExport.Version)
	return nil
}

// ensureTemplate checks if the template exists, and installs it if not
func (webExport WebExporter) ensureTemplate(docPath string) error {
	if helper.IsFileExist(docPath) {
		return nil
	}

	color.HiYellow("No documentation found, installing template ...")
	return webExport.installWebTemplate(filepath.Dir(docPath), webExport.AppName)
}

func (webExport WebExporter) updateAppConfig(docPath, version, description string) error {
	appPath := filepath.Join(docPath, "src", "assets", "app.json")
	if err := app.UpdateAppConfig(appPath, version, description); err != nil {
		return fmt.Errorf("error when updating the versions of your documentation: %w", err)
	}
	color.Green("Version file updated!")
	return nil
}

// writeDocumentationFile saves the doc content as a JSON file
func (webExport WebExporter) writeDocumentationFile(docPath string, content []byte) error {
	docFile := filepath.Join(docPath, "src", "assets", fmt.Sprintf("doc-%s.json", webExport.Version))
	if err := webExport.FileSystem.WriteFile(docFile, content, 0644); err != nil {
		return fmt.Errorf("error when saving your project documentation: %w", err)
	}
	return nil
}

func (webExport WebExporter) writeEnvFile(docPath, gitLink, appName, mainBranch string) error {
	envFile := filepath.Join(docPath, ".env")
	fileContent := fmt.Sprintf("VITE_GIT_LINK=%s\nVITE_APP_NAME=%s\nVITE_MAIN_BRANCH=%s", gitLink, appName, mainBranch)

	if err := webExport.FileSystem.WriteFile(envFile, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("error when saving your project documentation: %w", err)
	}
	return nil
}

// installWebTemplate clones the template repo and installs its dependencies
func (webExport WebExporter) installWebTemplate(docPath string, appName string) error {
	err := webExport.FileSystem.MkdirAll(docPath, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = webExport.CmdRunner.Execute(docPath, "git", "clone", TEMPLATE_GIT_LINK)

	if err != nil {
		return err
	}

	err = webExport.FileSystem.Rename(filepath.Join(docPath, "zendoc-ui-template"), filepath.Join(docPath, appName))
	if err != nil {
		return err
	}

	_, err = webExport.CmdRunner.Execute(filepath.Join(docPath, appName), "npm", "i")

	if err != nil {
		return err
	}

	return nil
}
