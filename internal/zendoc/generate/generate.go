package generate

import (
	"fmt"
	"os"
	"regexp"

	"github.com/dterbah/zendoc/config"
	"github.com/dterbah/zendoc/core/export"
	"github.com/dterbah/zendoc/core/parser"
	"github.com/dterbah/zendoc/internal"
	"github.com/fatih/color"
)

func wrapFileValidator(condition bool, wrappedFunc parser.DocParserFileValidator) parser.DocParserFileValidator {
	return func(filePath string) bool {
		res := wrappedFunc(filePath)
		if res {
			return condition
		}

		return true
	}
}

func wrapFunctionValidator(condition bool, wrappedFunc parser.DocParserFunctionValidator) parser.DocParserFunctionValidator {
	return func(name string) bool {
		res := wrappedFunc(name)
		if res {
			return condition
		}
		return true
	}
}

func wrapRegexFileValidator(regex []string) parser.DocParserFileValidator {
	return func(filePath string) bool {
		for _, reg := range regex {
			exp := regexp.MustCompile(reg)
			// skip invalid regex
			if exp != nil {
				return !exp.MatchString(filePath)
			}
		}
		return true
	}
}

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
		docExporter = export.WebExporter{
			GitLink:    projectConfig.ProjectConfig.GitLink,
			AppName:    projectConfig.ProjectConfig.Name,
			MainBranch: projectConfig.ProjectConfig.MainBranch,
			DocPath:    projectConfig.ProjectConfig.DocPath,
			Version:    projectConfig.ProjectConfig.Version,
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	docParser := parser.DocParser{
		FileValidators:     createFilevalidators(*projectConfig),
		FunctionValidators: createFunctionsValidators(*projectConfig),
	}
	projectDoc, err := docParser.ParseDocForDir(cwd, "")
	if err != nil {
		color.Red("error when parse your project %s", err)
		return err
	}

	return docExporter.Export(*projectDoc)
}

func createFilevalidators(configuration config.Config) []parser.DocParserFileValidator {
	validators := []parser.DocParserFileValidator{}

	validators = append(validators, wrapFileValidator(configuration.DocConfig.IncludeTests, IsTestFile),
		wrapFileValidator(configuration.DocConfig.IncludeMain, IsMainFile),
		wrapRegexFileValidator(configuration.DocConfig.ExcludeFiles),
	)

	return validators
}

func createFunctionsValidators(configuration config.Config) []parser.DocParserFunctionValidator {
	validators := []parser.DocParserFunctionValidator{}

	validators = append(validators, wrapFunctionValidator(configuration.DocConfig.IncludePrivate, IsPrivateFunction))

	return validators
}
