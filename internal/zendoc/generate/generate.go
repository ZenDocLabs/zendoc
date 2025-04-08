package generate

import (
	"fmt"
	"os"

	"github.com/dterbah/zendoc/config"
	"github.com/dterbah/zendoc/core/export"
	"github.com/dterbah/zendoc/core/parser"
	"github.com/dterbah/zendoc/internal"
	"github.com/fatih/color"
)

/*
@description Generate the documentation in a JSON format, or in a web app
@param outputFormat string - Either "json" or "web"
@author Dorian TERBAH
@return error - An error if the generation has failed
*/
func GenerateDoc(outputFormat string) error {
	var docExporter export.DocExporter

	projectConfig, err := config.GetConfiguration()
	if err != nil {
		return fmt.Errorf("error when reading the zendoc configuration : %s", err)
	}

	if outputFormat == internal.JSON_EXPORT_TYPE {
		docExporter = export.JSONExporter{}
	} else {
		docExporter = export.WebExporter{}
	}

	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	docParser := parser.DocParser{
		FileValidators: createFilevalidators(*projectConfig),
	}
	projectDoc, err := docParser.ParseDocForDir(cwd)
	if err != nil {
		color.Red("error when parse your project %s", err)
		return err
	}

	return docExporter.Export(*projectDoc)
}

func createFilevalidators(configuration config.Config) []parser.DocParserFileValidator {
	validators := []parser.DocParserFileValidator{}

	validators = append(validators, IsTestFileWrapper(configuration.DocConfig.IncludeTests))

	return validators
}
