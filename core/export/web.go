package export

import (
	"encoding/json"
	"fmt"

	"github.com/dterbah/zendoc/internal/doc"
)

type WebExporter struct {
	DocExporter
}

func (jsonExport WebExporter) Export(projectDoc doc.ProjectDoc) error {
	b, err := json.Marshal(projectDoc)
	if err != nil {
		return fmt.Errorf("error when exporting the documentation in JSON")
	}

	fmt.Print(b)

	return nil
}
