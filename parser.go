package gogqlparser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"strings"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) ParseTypeSystem(sources []*ast.Source) (*TypeSystemExtensionDocument, error) {
	typeSystemDocs := make([]*TypeSystemExtensionDocument, len(sources))
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

type TypeSystemExtensionDocument struct {
	// TypeSystemDefinition
	SchemaDefinitions    []ast.SchemaDefinition
	TypeDefinitions      []ast.TypeDefinition
	DirectiveDefinitions []ast.DirectiveDefinition
}

func MergeTypeSystemDocument(documents []*TypeSystemExtensionDocument) *TypeSystemExtensionDocument {
	return nil
}

func (p *Parser) parseTypeSystemDocument(src *ast.Source) (*TypeSystemExtensionDocument, error) {
	gogqllexer.New(strings.NewReader(src.Body))
	return nil, nil
}
