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
@author Dorian TERBAH
*/
func (webExport WebExporter) Export(projectDoc doc.ProjectDoc) error {
	b, err := json.Marshal(projectDoc)
	if err != nil {
		return fmt.Errorf("error when exporting the documentation in JSON")
	}

	// jsonDoc := string(b)

	currentPath, _ := os.Getwd()
	docPath := filepath.Join(currentPath, webExport.DocPath, webExport.AppName)

	if !helper.IsFileExist(docPath) {
		color.HiYellow("No documentation found, installing template ...")
		// need to clone the repo
		err := webExport.installWebTemplate(filepath.Join(currentPath, webExport.DocPath), webExport.AppName)
		if err != nil {
			return fmt.Errorf("error when installing the web template %s", err)
		}
	}

	// find the existing versions of the doc
	versionPath := filepath.Join(docPath, "src", "assets", "versions.json")
	err = version.UpdateVersions(versionPath, webExport.Version)
	if err != nil {
		return fmt.Errorf("error when updating the versions of your documentation : %s", err)
	}

	color.Green("Version file updated !")

	// finally, add the doc in the corresponding file
	documentationPath := filepath.Join(docPath, "src", "assets", fmt.Sprintf("doc-%s.json", webExport.Version))
	err = os.WriteFile(documentationPath, b, 0644)

	if err != nil {
		return fmt.Errorf("error when saving your project documentation : %s", err)
	}

	color.Green("Documentation v%s saved !", webExport.Version)

	return nil
}

func (webExport WebExporter) installWebTemplate(docPath string, appName string) error {
	// create the dir
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

	// rename the cloned repo
	err = os.Rename(filepath.Join(docPath, "zendoc-ui-template"), filepath.Join(docPath, appName))
	if err != nil {
		return err
	}

	// install packages inside
	npmInstallCommand := exec.Command("npm", "i")
	npmInstallCommand.Dir = filepath.Join(docPath, appName)

	_, err = npmInstallCommand.Output()
	if err != nil {
		return err
	}

	return nil
}
