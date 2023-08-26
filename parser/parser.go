package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"slices"
	"strings"
)

type parser struct {
	lexer *gogqllexer.Lexer

	keepToken *gogqllexer.Token
}

func (p *parser) NextToken() gogqllexer.Token {
	if p.keepToken != nil {
		t := *p.keepToken
		p.keepToken = nil
		return t
	}

	return p.lexer.NextToken()
}

func (p *parser) PeekToken() gogqllexer.Token {
	if p.keepToken == nil {
		t := p.lexer.NextToken()
		p.keepToken = &t
	}

	return *p.keepToken
}

func (p *parser) PeekAndMustBe(kinds []gogqllexer.Kind, callback func(t gogqllexer.Token, advanceLexer func()) error) error {
	t := p.PeekToken()
	if slices.Contains(kinds, t.Kind) {
		return callback(t, func() { p.NextToken() })
	}
	return fmt.Errorf("unexpected token %v", t)
}

func (p *parser) Skip(kind gogqllexer.Kind) error {
	t := p.NextToken()
	if t.Kind != kind {
		return fmt.Errorf("unexpected token %v", t)
	}
	return nil
}

func (p *parser) SkipIf(kind gogqllexer.Kind) (skip bool) {
	defer func() {
		if skip {
			p.NextToken()
		}
	}()
	t := p.PeekToken()
	return t.Kind == kind
}

func (p *parser) SkipKeywordIf(keyword string) (skip bool) {
	defer func() {
		if skip {
			p.NextToken()
		}
	}()
	t := p.PeekToken()
	return t.Kind == gogqllexer.Name && t.Value == keyword
}

func (p *parser) CheckKind(kind gogqllexer.Kind) bool {
	return p.PeekToken().Kind == kind
}

func (p *parser) CheckKeyword(keyword string) bool {
	t := p.PeekToken()
	return t.Kind == gogqllexer.Name && t.Value == keyword
}

func (p *parser) ReadNameValue() (string, error) {
	t := p.NextToken()
	if t.Kind != gogqllexer.Name {
		return "", fmt.Errorf("unexpected token %v", t)
	}
	return t.Value, nil
}

func (p *parser) SkipKeyword(keyword string) error {
	t := p.NextToken()
	if t.Kind != gogqllexer.Name || t.Value != keyword {
		return fmt.Errorf("unexpected token %v", t)
	}
	return nil
}

func (p *parser) ReadDescription() (description string, ok bool) {
	t := p.PeekToken()
	if t.Kind == gogqllexer.String || t.Kind == gogqllexer.BlockString {
		p.NextToken()
		return t.Value, ok
	}

	return "", false
}

func ParseTypeSystemExtensionDocument(src *ast.Source) (doc *ast.TypeSystemExtensionDocument, err error) {
	doc = &ast.TypeSystemExtensionDocument{
		TypeDefinitions: []ast.TypeDefinition{},
	}
	p := &parser{
		lexer: gogqllexer.New(strings.NewReader(src.Body)),
	}

	for {
		description, _ := p.ReadDescription()

		t := p.PeekToken()
		if t.Kind == gogqllexer.EOF {
			break
		}
		if t.Kind != gogqllexer.Name {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}

		switch t.Value {
		case "type":
			def, err := p.ParseObjectTypeDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.TypeDefinitions = append(doc.TypeDefinitions, def)
		case "interface":
			def, err := p.ParseInterfaceTypeDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.TypeDefinitions = append(doc.TypeDefinitions, def)
		case "union":
			def, err := p.ParseUnionTypeDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.TypeDefinitions = append(doc.TypeDefinitions, def)
		case "enum":
			def, err := p.ParseEnumTypeDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.TypeDefinitions = append(doc.TypeDefinitions, def)
		case "input":
			def, err := p.ParseInputObjectTypeDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.TypeDefinitions = append(doc.TypeDefinitions, def)
		case "scalar":
			def, err := p.ParseScalarTypeDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.TypeDefinitions = append(doc.TypeDefinitions, def)
		case "directive":
			directiveDefinition, err := p.ParseDirectiveDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.DirectiveDefinitions = append(doc.DirectiveDefinitions, *directiveDefinition)
		case "schema":
			schemaDef, err := p.ParseSchemaDefinition(description)
			if err != nil {
				return nil, err
			}
			doc.SchemaDefinitions = append(doc.SchemaDefinitions, *schemaDef)
		case "extend":
			if err = p.SkipKeyword("extend"); err != nil {
				return nil, err
			}

			keyword, err := p.ReadNameValue()
			if err != nil {
				return nil, err
			}
			switch keyword {
			case "type":
				def, err := p.ParseObjectTypeExtension()
				if err != nil {
					return nil, err
				}
				doc.TypeSystemExtensions = append(doc.TypeSystemExtensions, def)
			case "interface":
				def, err := p.ParseInterfaceTypeExtension()
				if err != nil {
					return nil, err
				}
				doc.TypeSystemExtensions = append(doc.TypeSystemExtensions, def)
			case "union":
				def, err := p.ParseUnionTypeExtension()
				if err != nil {
					return nil, err
				}
				doc.TypeSystemExtensions = append(doc.TypeSystemExtensions, def)
			case "enum":
				def, err := p.ParseEnumTypeExtension()
				if err != nil {
					return nil, err
				}
				doc.TypeSystemExtensions = append(doc.TypeSystemExtensions, def)
			case "input":
				def, err := p.ParseInputObjectTypeExtension()
				if err != nil {
					return nil, err
				}
				doc.TypeSystemExtensions = append(doc.TypeSystemExtensions, def)
			case "scalar":
				def, err := p.ParseScalarTypeExtension()
				if err != nil {
					return nil, err
				}
				doc.TypeSystemExtensions = append(doc.TypeSystemExtensions, def)
			case "schema":
				def, err := p.ParseSchemaExtension()
				if err != nil {
					return nil, err
				}
				doc.SchemaExtensions = append(doc.SchemaExtensions, *def)
			default:
				return nil, fmt.Errorf("unexpected token %+v", t.Value)
			}
		default:
			return nil, fmt.Errorf("unexpected token %+v", t.Value)
		}
	}

	return doc, nil
}
