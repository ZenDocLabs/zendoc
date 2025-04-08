package export

import (
	"encoding/json"
	"fmt"

	"github.com/dterbah/zendoc/internal/doc"
)

func ExportToJSON(projectDoc doc.ProjectDoc) (string, error) {
	b, err := json.Marshal(projectDoc)
	if err != nil {
		return "", fmt.Errorf("error when exporting the documentation in JSON")
	}

	return string(b), nil
}
