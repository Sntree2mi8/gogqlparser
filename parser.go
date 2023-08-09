package gogqlparser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
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
	l := gogqllexer.New(strings.NewReader(src.Body))

ParseSystemDocumentLoop:
	for {
		t := l.NextToken()
		if t.Kind == gogqllexer.EOF {
			break ParseSystemDocumentLoop
		}
		if t.Kind == gogqllexer.Comment {
			// commentが入ると色々だるいので一旦無視してコメント以外をちゃんとparseできることをまずは達成する
			continue ParseSystemDocumentLoop
		}

		// 一番最初にdescriptionが来る可能性がある
		// ここでくる文字列はdescription
		var description string
		if t.Kind == gogqllexer.String || t.Kind == gogqllexer.BlockString {
			description = t.Value
			t = l.NextToken()
			if t.Kind == gogqllexer.EOF {
				break
			}
			if t.Kind == gogqllexer.Comment {
				// commentが入ると色々だるいので一旦無視してコメント以外をちゃんとparseできることをまずは達成する
				continue
			}
		}

		// CHECK: lexerにもpeekが絶対に必要になるのか
		if t.Kind != gogqllexer.Name {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}

		switch t.Value {
		case "schema":
			// TODO: schemaだった場合にこれからのトークンになくてはならない並び順がある
			// TODO: directiveを一旦飛ばしているのであとで実装する
			t = l.NextToken()
			if t.Kind != gogqllexer.BraceL {
				return nil, fmt.Errorf("unexpected token %+v", t)
			}

			rootOperationMap := make(map[ast.RootOperationTypeKind]string)
		ParseRootOperationLoop:
			for {
				t = l.NextToken()
				if t.Kind == gogqllexer.Comment {
					continue ParseRootOperationLoop
				}
				switch t.Kind {
				case gogqllexer.BraceR:
					break ParseRootOperationLoop
				case gogqllexer.Name:
					switch t.Value {
					case "query":
						if _, ok := rootOperationMap[ast.OperationTypeQuery]; ok {
							return nil, fmt.Errorf("duplicate root operation type query")
						}
						t = l.NextToken()
						if t.Kind != gogqllexer.Colon {
							return nil, fmt.Errorf("unexpected token %+v", t)
						}
						t = l.NextToken()
						if t.Kind != gogqllexer.Name {
							return nil, fmt.Errorf("unexpected token %+v", t)
						}
						rootOperationMap[ast.OperationTypeQuery] = t.Value
					case "mutation":
						if _, ok := rootOperationMap[ast.OperationTypeMutation]; ok {
							return nil, fmt.Errorf("duplicate root operation type mutation")
						}
						t = l.NextToken()
						if t.Kind != gogqllexer.Colon {
							return nil, fmt.Errorf("unexpected token %+v", t)
						}
						t = l.NextToken()
						if t.Kind != gogqllexer.Name {
							return nil, fmt.Errorf("unexpected token %+v", t)
						}
						rootOperationMap[ast.OperationTypeMutation] = t.Value
					case "subscription":
						if _, ok := rootOperationMap[ast.OperationTypeSubscription]; ok {
							return nil, fmt.Errorf("duplicate root operation type subscription")
						}
						t = l.NextToken()
						if t.Kind != gogqllexer.Colon {
							return nil, fmt.Errorf("unexpected token %+v", t)
						}
						t = l.NextToken()
						if t.Kind != gogqllexer.Name {
							return nil, fmt.Errorf("unexpected token %+v", t)
						}
						rootOperationMap[ast.OperationTypeSubscription] = t.Value
					default:
						return nil, fmt.Errorf("unexpected token %+v", t)
					}
				default:
					return nil, fmt.Errorf("unexpected token %+v", t)
				}
			}
			schemaDef := ast.SchemaDefinition{
				Description:                  description,
				Directives:                   ast.Directives{},
				RootOperationTypeDefinitions: []ast.RootOperationTypeDefinition{},
			}
			for k, v := range rootOperationMap {
				schemaDef.RootOperationTypeDefinitions = append(
					schemaDef.RootOperationTypeDefinitions,
					ast.RootOperationTypeDefinition{
						OperationType: k,
						Type:          v,
					},
				)
			}
			d.SchemaDefinitions = append(d.SchemaDefinitions, schemaDef)
		default:
			return nil, fmt.Errorf("unexpected token %+v", t.Value)
		}
	}

	return d, nil
}
