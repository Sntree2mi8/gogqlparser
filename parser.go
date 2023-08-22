package gogqlparser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"github.com/Sntree2mi8/gogqlparser/parser"
	"strings"
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
	d := &ast.TypeSystemExtensionDocument{
		SchemaDefinitions:    []ast.SchemaDefinition{},
		TypeDefinitions:      []ast.TypeDefinition{},
		DirectiveDefinitions: []ast.DirectiveDefinition{},
	}
	l := parser.NewLexerWrapper(gogqllexer.New(strings.NewReader(src.Body)))

	for {
		description, _ := l.ReadDescription()

		t := l.PeekToken()
		if t.Kind == gogqllexer.EOF {
			break
		}
		if t.Kind != gogqllexer.Name {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}

		switch t.Value {
		case "type":
			def, err := parser.ParseObjectTypeDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions = append(d.TypeDefinitions, def)
		case "interface":
			def, err := parser.ParseInterfaceTypeDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions = append(d.TypeDefinitions, def)
		case "union":
			def, err := parser.ParseUnionTypeDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions = append(d.TypeDefinitions, def)
		case "enum":
			def, err := parser.ParseEnumTypeDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions = append(d.TypeDefinitions, def)
		case "input":
			def, err := parser.ParseInputObjectTypeDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions = append(d.TypeDefinitions, def)
		case "directive":
			directiveDefinition, err := parser.ParseDirectiveDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.DirectiveDefinitions = append(d.DirectiveDefinitions, *directiveDefinition)
		case "schema":
			schemaDef, err := parser.ParseSchemaDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.SchemaDefinitions = append(d.SchemaDefinitions, *schemaDef)
		default:
			return nil, fmt.Errorf("unexpected token %+v", t.Value)
		}
	}

	return d, nil
}
