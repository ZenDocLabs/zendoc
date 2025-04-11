package export

import (
	"encoding/json"
	"fmt"

	"github.com/dterbah/zendoc/internal/doc"
)

/*
@description Struct that implements the DocExporter interface and exports the documentation in a web-friendly format.
@author Dorian TERBAH
@field DocExporter DocExporter - Embedded base exporter providing common exporting behavior.
*/
type WebExporter struct {
	DocExporter
}

/*
@description Export the project documentation as JSON to stdout
@param projectDoc doc.ProjectDoc - The documentation to export
@return error - An error if the export fails
@example WebExporter{}.Export(projectDoc)
@author Dorian TERBAH
*/
func (jsonExport WebExporter) Export(projectDoc doc.ProjectDoc) error {
	b, err := json.Marshal(projectDoc)
	if err != nil {
		return fmt.Errorf("error when exporting the documentation in JSON")
	}

	fmt.Print(b)

	return nil
}
