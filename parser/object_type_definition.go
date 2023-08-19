package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func parseInterfaceImplementations(l *LexerWrapper) (interfaces []string, err error) {
	if err = l.SkipKeyword("implements"); err != nil {
		return nil, err
	}

	// implements at least one interface
	l.SkipIf(gogqllexer.Amp)
	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Name},
		func(t gogqllexer.Token, advanceLexer func()) error {
			defer advanceLexer()

			interfaces = append(interfaces, t.Value)
			return nil
		},
	); err != nil {
		return nil, err
	}

	// read more interfaces
	for {
		if skip := l.SkipIf(gogqllexer.Amp); !skip {
			break
		}

		if err = l.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Name},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()

				interfaces = append(interfaces, t.Value)
				return nil
			},
		); err != nil {
			return nil, err
		}
	}

	return interfaces, nil
}

// https://spec.graphql.org/October2021/#sec-Objects
func ParseObjectTypeDefinition(l *LexerWrapper, description string) (d *ast.ObjectTypeDefinition, err error) {
	d = &ast.ObjectTypeDefinition{
		Description: description,
	}

	if err = l.SkipKeyword("type"); err != nil {
		return nil, err
	}

	if d.Name, err = l.ReadNameValue(); err != nil {
		return nil, err
	}

	if l.CheckKeyword("implements") {
		if d.Interfaces, err = parseInterfaceImplementations(l); err != nil {
			return nil, err
		}
	}

	if l.CheckKind(gogqllexer.At) {
		if d.Directives, err = parseDirectives(l); err != nil {
			return nil, err
		}
	}

	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.BraceL},
		func(t gogqllexer.Token, advanceLexer func()) error {
			l.NextToken()

			for {
				t = l.PeekToken()
				if t.Kind == gogqllexer.BraceR {
					l.NextToken()
					break
				}
				if t.Kind == gogqllexer.EOF {
					return fmt.Errorf("unexpected token %+v", t)
				}

				fieldDefinition, err := parseFieldDefinition(l)
				if err != nil {
					return err
				}

				d.FieldDefinitions = append(d.FieldDefinitions, fieldDefinition)
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return d, nil
}
