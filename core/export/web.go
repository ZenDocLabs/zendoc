package export

import (
	"encoding/json"
	"fmt"

	"github.com/dterbah/zendoc/internal/doc"
)

type WebExporter struct {
	DocExporter
}

/*
@description Export the project documentation as JSON to stdout
@param projectDoc doc.ProjectDoc - The documentation to export
@return error - An error if the export fails
@example WebExporter{}.Export(projectDoc)
*/
func (jsonExport WebExporter) Export(projectDoc doc.ProjectDoc) error {
	b, err := json.Marshal(projectDoc)
	if err != nil {
		return fmt.Errorf("error when exporting the documentation in JSON")
	}

	fmt.Print(b)

	return nil
}
