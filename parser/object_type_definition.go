package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

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
		if d.Interfaces, err = parseImplementsInterfaces(l); err != nil {
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
			if d.FieldDefinitions, err = parseFieldsDefinition(l); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return d, nil
}

// ParseObjectTypeExtension parse object type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Object-Extensions
func ParseObjectTypeExtension(l *LexerWrapper) (def *ast.ObjectTypeExtension, err error) {
	def = &ast.ObjectTypeExtension{}

	if err = l.SkipKeyword("type"); err != nil {
		return nil, err
	}

	if def.Name, err = l.ReadNameValue(); err != nil {
		return nil, err
	}

	var canOmitFields bool
	if l.CheckKeyword("implements") {
		if def.ImplementInterfaces, err = parseImplementsInterfaces(l); err != nil {
			return nil, err
		}

		canOmitFields = true
	}

	if l.CheckKind(gogqllexer.At) {
		if def.Directives, err = parseDirectives(l); err != nil {
			return nil, err
		}

		canOmitFields = true
	}

	if l.CheckKind(gogqllexer.BraceL) {
		if def.FieldsDefinition, err = parseFieldsDefinition(l); err != nil {
			return nil, err
		}
	} else if !canOmitFields {
		return nil, fmt.Errorf("unexpected token. expected interface implementation or directive or fields definition")
	}

	return def, nil
}
