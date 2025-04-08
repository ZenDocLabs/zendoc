package export

import (
	"os"

	"github.com/dterbah/zendoc/core/parser/serializer"
	"github.com/dterbah/zendoc/internal/doc"
)

type JSONExporter struct {
	DocExporter
}

const EXPORT_FILE = "doc.json"

/*
@description Export the project documentation to a JSON file
@param projectDoc doc.ProjectDoc - The documentation to export
@return error - An error if writing to file fails
@example JSONExporter{}.Export(projectDoc)
*/
func (jsonExport JSONExporter) Export(projectDoc doc.ProjectDoc) error {
	projectDocJSON, err := serializer.SerializeToJSON(projectDoc)

	if err != nil {
		return nil
	}

	err = os.WriteFile(EXPORT_FILE, []byte(projectDocJSON), 0644)
	return err
}
