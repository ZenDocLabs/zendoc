package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/dterbah/zendoc/internal/doc"
)

func SerializeToJSON(doc doc.ProjectDoc) (string, error) {
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error when exporting the documentation in JSON")
	}

	return string(b), nil
}
