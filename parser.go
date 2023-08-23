package gogqlparser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"github.com/Sntree2mi8/gogqlparser/parser"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) ParseTypeSystem(sources []*ast.Source) (*ast.TypeSystemExtensionDocument, error) {
	typeSystemDocs := make([]*ast.TypeSystemExtensionDocument, len(sources))
	// TODO: parallelize
	for i, src := range sources {
		doc, err := p.parseTypeSystemDocument(src)
		if err != nil {
			return nil, err
		}
		typeSystemDocs[i] = doc
	}
	return MergeTypeSystemDocument(typeSystemDocs), nil
}

func MergeTypeSystemDocument(documents []*ast.TypeSystemExtensionDocument) *ast.TypeSystemExtensionDocument {
	return nil
}

func (p *Parser) parseTypeSystemDocument(src *ast.Source) (*ast.TypeSystemExtensionDocument, error) {
	return parser.ParseTypeSystemExtensionDocument(src)
}
