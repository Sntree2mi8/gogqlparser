package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#FieldsDefinition
func parseFieldsDefinition(l *LexerWrapper) (d []*ast.FieldDefinition, err error) {
	if err = l.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		var fieldDefinition *ast.FieldDefinition
		fieldDefinition, err = parseFieldDefinition(l)
		if err != nil {
			return nil, err
		}

		d = append(d, fieldDefinition)

		if l.SkipIf(gogqllexer.BraceR) {
			break
		}
	}

	return d, nil
}

// https://spec.graphql.org/October2021/#FieldDefinition
func parseFieldDefinition(l *LexerWrapper) (d *ast.FieldDefinition, err error) {
	d = &ast.FieldDefinition{}

	if err = l.PeekAndMayBe(
		[]gogqllexer.Kind{gogqllexer.String, gogqllexer.BlockString},
		func(t gogqllexer.Token, advanceLexer func()) error {
			defer advanceLexer()

			d.Description = t.Value
			return nil
		},
	); err != nil {
		return nil, err
	}

	if d.Name, err = l.ReadNameValue(); err != nil {
		return nil, err
	}

	if l.CheckKind(gogqllexer.ParenL) {
		if d.ArgumentDefinition, err = ParseArgumentsDefinition(l); err != nil {
			return nil, err
		}
	}

	if err = l.Skip(gogqllexer.Colon); err != nil {
		return nil, err
	}

	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Name},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if d.Type, err = parseType(l); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	if l.CheckKind(gogqllexer.At) {
		if d.Directives, err = parseDirectives(l); err != nil {
			return nil, err
		}
	}

	return d, err
}
