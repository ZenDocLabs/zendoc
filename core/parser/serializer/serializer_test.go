package serializer

import (
	"testing"

	"github.com/dterbah/zendoc/internal/doc"
	"github.com/stretchr/testify/assert"
)

func TestSerializeToJSON_Success(t *testing.T) {
	projectDoc := doc.ProjectDoc{
		PackageDocs: map[string][]doc.FileDoc{
			"main": {
				{
					Docs: []doc.FuncDoc{
						{
							Name:        "MyFunction",
							Description: "Does something",
							Params: []doc.Param{
								{Name: "param1", Type: "string", Description: "First parameter"},
							},
							Return: &doc.Return{Type: "int", Description: "Return value"},
						},
					},
				},
			},
		},
	}

	result, err := SerializeToJSON(projectDoc)

	assert.NoError(t, err)
	assert.Contains(t, result, `"MyFunction"`)
	assert.Contains(t, result, `"param1"`)
	assert.Contains(t, result, `"Return value"`)
}
