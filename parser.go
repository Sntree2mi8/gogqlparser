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

// https://spec.graphql.org/October2021/#FieldDefinition
// TODO: parse arguments
func parseFieldDefinition(l *lexerWrapper) (d *ast.FieldDefinition, err error) {
	d = &ast.FieldDefinition{
		Directives: make([]ast.Directive, 0),
	}

	// description is optional
	if err = maybe(l, func(t gogqllexer.Token) bool {
		if t.Kind == gogqllexer.String || t.Kind == gogqllexer.BlockString {
			d.Description = t.Value
			return true
		}
		return false
	}); err != nil {
		return nil, err
	}

	if err = mustBe(
		l,
		// name is required
		func(t gogqllexer.Token) bool {
			if t.Kind == gogqllexer.Name {
				d.Name = t.Value
				return true
			}
			return false
		},
		// colon is required
		func(t gogqllexer.Token) bool {
			return t.Kind == gogqllexer.Colon
		},
	); err != nil {
		return nil, err
	}

	// type is required
	// parse type
	if d.Type, err = parseType(l); err != nil {
		return nil, err
	}

	// directives are optional
	for {
		t := l.PeekToken()
		if t.Kind == gogqllexer.At {
			directive, err := parseDirective(l)
			if err != nil {
				return nil, err
			}
			d.Directives = append(d.Directives, directive)
			continue
		}
		break
	}

	return d, err
}

// https://spec.graphql.org/October2021/#sec-Objects
// TODO: parse implements interface
// TODO: parse directives
func parseTypeObjectDefinition(l *lexerWrapper) (d *ast.ObjectTypeDefinition, err error) {
	d = &ast.ObjectTypeDefinition{
		FieldDefinitions: make([]*ast.FieldDefinition, 0),
	}

	var t gogqllexer.Token

	err = mustBe(
		l,
		// start with "type"
		func(t gogqllexer.Token) bool {
			return t.Kind == gogqllexer.Name && t.Value == "type"
		},
		// name of object is required
		func(t gogqllexer.Token) bool {
			if t.Kind == gogqllexer.Name {
				d.Name = t.Value
				return true
			}
			return false
		},
		// open brace is required
		func(t gogqllexer.Token) bool {
			return t.Kind == gogqllexer.BraceL
		},
	)
	if err != nil {
		return nil, err
	}

	for {
		t = l.PeekToken()
		if t.Kind == gogqllexer.BraceR {
			l.NextToken()
			break
		}
		if t.Kind == gogqllexer.EOF {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}

		fieldDefinition, err := parseFieldDefinition(l)
		if err != nil {
			return nil, err
		}

		d.FieldDefinitions = append(d.FieldDefinitions, fieldDefinition)
	}

	return d, nil
}

func (p *Parser) parseTypeSystemDocument(src *ast.Source) (*ast.TypeSystemExtensionDocument, error) {
	d := &ast.TypeSystemExtensionDocument{
		SchemaDefinitions:    []ast.SchemaDefinition{},
		TypeDefinitions:      map[string]ast.TypeDefinition{},
		DirectiveDefinitions: []ast.DirectiveDefinition{},
	}
	l := &lexerWrapper{lexer: gogqllexer.New(strings.NewReader(src.Body))}

ParseSystemDocumentLoop:
	for {
		t := l.PeekToken()
		if t.Kind == gogqllexer.EOF {
			break ParseSystemDocumentLoop
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
		}

		if t.Kind != gogqllexer.Name {
			return nil, fmt.Errorf("unexpected token %+v", t)
		}

		switch t.Value {
		case "type":
			typeObjectDefinition, err := parseTypeObjectDefinition(l)
			if err != nil {
				return nil, err
			}
			d.TypeDefinitions[typeObjectDefinition.Name] = typeObjectDefinition

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
				inputValueDefinitions, err = parseArgumentsDefinition(l)
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
			dl := parsDirectiveLocation(t.Value)
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
				dl := parsDirectiveLocation(t.Value)
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
			l.NextToken()
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
				Directives:                   []ast.Directive{},
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

func parsDirectiveLocation(v string) ast.DirectiveLocation {
	switch v {
	case "QUERY":
		return ast.DirectiveLocationQuery
	case "MUTATION":
		return ast.DirectiveLocationMutation
	case "SUBSCRIPTION":
		return ast.DirectiveLocationSubscription
	case "FIELD":
		return ast.DirectiveLocationField
	case "FRAGMENT_DEFINITION":
		return ast.DirectiveLocationFragmentDefinition
	case "FRAGMENT_SPREAD":
		return ast.DirectiveLocationFragmentSpread
	case "INLINE_FRAGMENT":
		return ast.DirectiveLocationInlineFragment
	case "VARIABLE_DEFINITION":
		return ast.DirectiveLocationVariableDefinition

	case "SCHEMA":
		return ast.DirectiveLocationSchema
	case "SCALAR":
		return ast.DirectiveLocationScalar
	case "OBJECT":
		return ast.DirectiveLocationObject
	case "FIELD_DEFINITION":
		return ast.DirectiveLocationFieldDefinition
	case "ARGUMENT_DEFINITION":
		return ast.DirectiveLocationArgumentDefinition
	case "INTERFACE":
		return ast.DirectiveLocationInterface
	case "UNION":
		return ast.DirectiveLocationUnion
	case "ENUM":
		return ast.DirectiveLocationEnum
	case "ENUM_VALUE":
		return ast.DirectiveLocationEnumValue
	case "INPUT_OBJECT":
		return ast.DirectiveLocationInputObject
	case "INPUT_FIELD_DEFINITION":
		return ast.DirectiveLocationInputFieldDefinition
	default:
		return ast.DirectiveLocationUnknown
	}
}

func parseInputValueDefinition(l *lexerWrapper, description string) (def ast.InputValueDefinition, err error) {
	var t gogqllexer.Token

	t = l.NextToken()
	if t.Kind != gogqllexer.Name {
		return def, fmt.Errorf("unexpected token %+v", t)
	}
	def.Name = t.Value

	if t = l.NextToken(); t.Kind != gogqllexer.Colon {
		return def, fmt.Errorf("unexpected token %+v", t)
	}

	argType, err := parseType(l)
	if err != nil {
		return def, err
	}
	def.Type = argType

	t = l.PeekToken()
	if t.Kind == gogqllexer.Equal {
		l.NextToken()
		t = l.NextToken()
		switch t.Kind {
		case gogqllexer.Int, gogqllexer.Float, gogqllexer.String, gogqllexer.BlockString, gogqllexer.Name:
			def.RawDefaultValue = t.Value
		default:
		}
	}

	directives := make([]ast.Directive, 0)
	for {
		t = l.PeekToken()
		if t.Kind != gogqllexer.At {
			break
		}
		d, err := parseDirective(l)
		if err != nil {
			return def, err
		}
		directives = append(directives, d)
	}
	if len(directives) > 0 {
		def.Directives = directives
	}
	def.Description = description

	return def, nil
}

func parseDirective(l *lexerWrapper) (d ast.Directive, err error) {
	t := l.NextToken()
	if t.Kind != gogqllexer.At {
		return d, fmt.Errorf("unexpected token %+v", t)
	}

	t = l.NextToken()
	if t.Kind != gogqllexer.Name {
		return d, fmt.Errorf("unexpected token %+v", t)
	}
	d.Name = t.Value

	args := make([]ast.Argument, 0)
	t = l.PeekToken()
	if t.Kind == gogqllexer.ParenL {
		arg := ast.Argument{}
		l.NextToken()
		t = l.NextToken()
		if t.Kind != gogqllexer.Name {
			return d, fmt.Errorf("unexpected token %+v", t)
		}
		arg.Name = t.Value

		t = l.NextToken()
		if t.Kind != gogqllexer.Colon {
			return d, fmt.Errorf("unexpected token %+v", t)
		}

		t = l.NextToken()
		switch t.Kind {
		default:
			return d, fmt.Errorf("unexpected token %+v", t)
		case gogqllexer.Int, gogqllexer.Float, gogqllexer.Name, gogqllexer.String, gogqllexer.BlockString:
			arg.Value = t.Value
		}

		t = l.NextToken()
		if t.Kind != gogqllexer.ParenR {
			return d, fmt.Errorf("unexpected token %+v", t)
		}

		args = append(args, arg)
		d.Arguments = args
	}

	return d, err
}

func parseArgumentsDefinition(l *lexerWrapper) (defs []ast.InputValueDefinition, err error) {
	t := l.NextToken()
	if t.Kind != gogqllexer.ParenL {
		return defs, err
	}

	for {
		inputValDescription := ""

		t = l.PeekToken()
		switch t.Kind {
		case gogqllexer.String, gogqllexer.BlockString:
			inputValDescription = t.Value
			l.NextToken()
		case gogqllexer.Name:
			ivd, err := parseInputValueDefinition(l, inputValDescription)
			if err != nil {
				return defs, err
			}
			defs = append(defs, ivd)
		case gogqllexer.ParenR:
			if len(defs) == 0 {
				return defs, fmt.Errorf("unexpected token %+v", t)
			}
			l.NextToken()
			return defs, nil
		default:
			return defs, fmt.Errorf("unexpected token %+v", t)
		}
	}
}

func parseType(l *lexerWrapper) (t ast.Type, err error) {
	var token gogqllexer.Token

	token = l.NextToken()
	switch token.Kind {
	case gogqllexer.Name:
		t.NamedType = token.Value
		token = l.PeekToken()
		if token.Kind == gogqllexer.Bang {
			t.NotNull = true
			l.NextToken()
		}
	case gogqllexer.BracketL:
		t.ListType = &ast.Type{}

		elmType, err := parseType(l)
		if err != nil {
			return t, err
		}
		t.ListType = &elmType

		token = l.NextToken()
		if token.Kind != gogqllexer.BracketR {
			return t, fmt.Errorf("unexpected token %+v", token)
		}

		token = l.PeekToken()
		if token.Kind == gogqllexer.Bang {
			t.NotNull = true
			l.NextToken()
		}
	default:
		return t, fmt.Errorf("unexpected token %+v", token)
	}

	return t, nil
}
