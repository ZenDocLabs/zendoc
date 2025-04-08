package export

import "github.com/dterbah/zendoc/internal/doc"

type DocExporter interface {
	Export(projectDoc doc.ProjectDoc) error
}
