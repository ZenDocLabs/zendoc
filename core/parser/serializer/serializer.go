package serializer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dterbah/zendoc/internal/doc"
)

/*
@description Serialize a ProjectDoc into a pretty-printed JSON string
@param doc doc.ProjectDoc - The project documentation to serialize
@return (string, error) - The resulting JSON string and an error if serialization fails
@example SerializeToJSON(myDoc)
*/
func SerializeToJSON(doc doc.ProjectDoc) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(doc)

	if err != nil {
		return "", fmt.Errorf("error when exporting the documentation in JSON")
	}

	return buf.String(), nil
}
