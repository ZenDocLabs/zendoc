package export

import "github.com/dterbah/zendoc/internal/doc"

/*
@description Export a project doc
@author Dorian TERBAH
*/
type DocExporter interface {
	/*
		@description Export documentation
		@param projectDoc doc.ProjectDoc - The associated doc
		@author Dorian TERBAH
		@return error - If there is any problem during the export
	*/
	Export(projectDoc doc.ProjectDoc) error
}
