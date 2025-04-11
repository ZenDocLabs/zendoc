package export

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dterbah/zendoc/core/export/helper"
	"github.com/dterbah/zendoc/core/export/version"
	"github.com/dterbah/zendoc/internal/doc"
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
	GitLink    string
	AppName    string
	MainBranch string
	DocPath    string
	Version    string
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

	if err := webExport.updateVersionFile(docPath); err != nil {
		return err
	}

	if err := webExport.writeDocumentationFile(docPath, b); err != nil {
		return err
	}

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

// updateVersionFile updates the version.json file with the current version
func (webExport WebExporter) updateVersionFile(docPath string) error {
	versionPath := filepath.Join(docPath, "src", "assets", "versions.json")
	if err := version.UpdateVersions(versionPath, webExport.Version); err != nil {
		return fmt.Errorf("error when updating the versions of your documentation: %w", err)
	}
	color.Green("Version file updated!")
	return nil
}

// writeDocumentationFile saves the doc content as a JSON file
func (webExport WebExporter) writeDocumentationFile(docPath string, content []byte) error {
	docFile := filepath.Join(docPath, "src", "assets", fmt.Sprintf("doc-%s.json", webExport.Version))
	if err := os.WriteFile(docFile, content, 0644); err != nil {
		return fmt.Errorf("error when saving your project documentation: %w", err)
	}
	return nil
}

// installWebTemplate clones the template repo and installs its dependencies
func (webExport WebExporter) installWebTemplate(docPath string, appName string) error {
	err := os.MkdirAll(docPath, os.ModePerm)
	if err != nil {
		return err
	}

	gitCommand := exec.Command("git", "clone", TEMPLATE_GIT_LINK)
	gitCommand.Dir = docPath

	_, err = gitCommand.Output()
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(docPath, "zendoc-ui-template"), filepath.Join(docPath, appName))
	if err != nil {
		return err
	}

	npmInstallCommand := exec.Command("npm", "i")
	npmInstallCommand.Dir = filepath.Join(docPath, appName)

	_, err = npmInstallCommand.Output()
	if err != nil {
		return err
	}

	return nil
}
