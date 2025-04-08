package zendoc

import (
	"os"

	"github.com/dterbah/zendoc/core/export"
	"github.com/dterbah/zendoc/core/parser"
	"github.com/dterbah/zendoc/internal"
	"github.com/fatih/color"
)

func GenerateDoc(outputFormat string) error {
	var docExporter export.DocExporter

	if outputFormat == internal.JSON_EXPORT_TYPE {
		docExporter = export.JSONExporter{}
	} else {
		docExporter = export.WebExporter{}
	}

	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	projectDoc, err := parser.ParseDocForDir(cwd)
	if err != nil {
		color.Red("error when parse your project %s", err)
		return err
	}

	return docExporter.Export(*projectDoc)
}
