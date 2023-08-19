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
		TypeDefinitions:      map[string]ast.TypeDefinition{},
		DirectiveDefinitions: []ast.DirectiveDefinition{},
	}
	l := parser.NewLexerWrapper(gogqllexer.New(strings.NewReader(src.Body)))

ParseSystemDocumentLoop:
	for {
		var description string
		if err := l.PeekAndMayBe(
			[]gogqllexer.Kind{gogqllexer.String, gogqllexer.BlockString},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()

				description = t.Value
				return nil
			},
		); err != nil {
			return nil, err
		}

		t := l.PeekToken()
		if t.Kind == gogqllexer.EOF {
			break ParseSystemDocumentLoop
		}
		if t.Kind != gogqllexer.Name {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}

		switch t.Value {
		case "type":
			typeObjectDefinition, err := parser.ParseObjectTypeDefinition(l, description)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions[typeObjectDefinition.Name] = typeObjectDefinition
		case "interface":
			typeInterfaceDefinition, err := parser.ParseTypeInterfaceDefinition(l)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions[typeInterfaceDefinition.Name] = typeInterfaceDefinition

		case "directive":
			l.NextToken()
			t = l.NextToken()
			if t.Kind != gogqllexer.At {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}
			t = l.NextToken()
			if t.Kind != gogqllexer.Name {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}

			directiveName := t.Value
			isRepeatable := false

			t = l.PeekToken()
			var err error
			var inputValueDefinitions []ast.InputValueDefinition
			if t.Kind == gogqllexer.ParenL {
				inputValueDefinitions, err = parser.ParseArgumentsDefinition(l)
				if err != nil {
					return nil, err
				}
			}

			t = l.PeekToken()
			if t.Kind != gogqllexer.Name {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}

			if t.Value == "repeatable" {
				isRepeatable = true
				l.NextToken()
			}

			t = l.NextToken()
			if t.Kind != gogqllexer.Name {
				return nil, fmt.Errorf("unexpected token %+v", t)
			} else if t.Value != "on" {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}

			t = l.NextToken()
			if t.Kind != gogqllexer.Name {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}
			dLocations := make([]ast.DirectiveLocation, 0)
			dl := parser.ParsDirectiveLocation(t.Value)
			if dl == ast.DirectiveLocationUnknown {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}
			dLocations = append(dLocations, dl)

			for {
				t = l.PeekToken()
				if t.Kind != gogqllexer.Pipe {
					break
				} else {
					l.NextToken()
				}

				t = l.NextToken()
				if t.Kind != gogqllexer.Name {
					return nil, fmt.Errorf("unexpected token %+v", t)
				}
				dl := parser.ParsDirectiveLocation(t.Value)
				if dl == ast.DirectiveLocationUnknown {
					return nil, fmt.Errorf("unexpected token %+v", t)
				}
				dLocations = append(dLocations, dl)
			}

			d.DirectiveDefinitions = append(d.DirectiveDefinitions, ast.DirectiveDefinition{
				Description:         description,
				Name:                directiveName,
				ArgumentsDefinition: inputValueDefinitions,
				IsRepeatable:        isRepeatable,
				DirectiveLocations:  dLocations,
			})

		case "schema":
			if err := l.SkipKeyword("schema"); err != nil {
				return nil, err
			}
			// TODO: schemaだった場合にこれからのトークンになくてはならない並び順がある
			// TODO: directiveを一旦飛ばしているのであとで実装する
			if err := l.Skip(gogqllexer.BraceL); err != nil {
				return nil, err
			}

			schemaDef := ast.SchemaDefinition{
				Description: description,
				Directives:  []ast.Directive{},
			}

		ParseRootOperationLoop:
			for {
				t = l.NextToken()
				switch t.Kind {
				case gogqllexer.BraceR:
					break ParseRootOperationLoop
				case gogqllexer.Name:
					switch t.Value {
					case "query":
						if err := l.Skip(gogqllexer.Colon); err != nil {
							return nil, err
						}
						if err := l.PeekAndMustBe(
							[]gogqllexer.Kind{gogqllexer.Name},
							func(t gogqllexer.Token, advanceLexer func()) error {
								defer advanceLexer()
								schemaDef.Query = &ast.RootOperationTypeDefinition{Type: t.Value}
								return nil
							},
						); err != nil {
							return nil, err
						}
					case "mutation":
						if err := l.Skip(gogqllexer.Colon); err != nil {
							return nil, err
						}
						if err := l.PeekAndMustBe(
							[]gogqllexer.Kind{gogqllexer.Name},
							func(t gogqllexer.Token, advanceLexer func()) error {
								defer advanceLexer()
								schemaDef.Mutation = &ast.RootOperationTypeDefinition{Type: t.Value}
								return nil
							},
						); err != nil {
							return nil, err
						}
					case "subscription":
						if err := l.Skip(gogqllexer.Colon); err != nil {
							return nil, err
						}
						if err := l.PeekAndMustBe(
							[]gogqllexer.Kind{gogqllexer.Name},
							func(t gogqllexer.Token, advanceLexer func()) error {
								defer advanceLexer()
								schemaDef.Subscription = &ast.RootOperationTypeDefinition{Type: t.Value}
								return nil
							},
						); err != nil {
							return nil, err
						}
					default:
						return nil, fmt.Errorf("unexpected token %+v", t)
					}
				default:
					return nil, fmt.Errorf("unexpected token %+v", t)
				}
			}

			d.SchemaDefinitions = append(d.SchemaDefinitions, schemaDef)
		default:
			return nil, fmt.Errorf("unexpected token %+v", t.Value)
		}
	}

	return d, nil
}
