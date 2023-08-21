package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#sec-Interfaces
func ParseInterfaceTypeDefinition(l *LexerWrapper, description string) (d *ast.InterfaceTypeDefinition, err error) {
	d = &ast.InterfaceTypeDefinition{}

	d.Description = description

	if err = l.SkipKeyword("interface"); err != nil {
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

	return d, err
}

// ParseInterfaceTypeExtension parse interface type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Interface-Extensions
func ParseInterfaceTypeExtension(l *LexerWrapper) (def *ast.InterfaceTypeExtension, err error) {
	def = &ast.InterfaceTypeExtension{}

	if err = l.SkipKeyword("interface"); err != nil {
		return nil, err
	}

	if def.Name, err = l.ReadNameValue(); err != nil {
		return nil, err
	}

	if l.CheckKeyword("implements") {
		if def.ImplementInterfaces, err = parseImplementsInterfaces(l); err != nil {
			return nil, err
		}
	}

	if l.CheckKind(gogqllexer.At) {
		if def.Directives, err = parseDirectives(l); err != nil {
			return nil, err
		}
	}

	if l.CheckKind(gogqllexer.BraceL) {
		if def.FieldsDefinition, err = parseFieldsDefinition(l); err != nil {
			return nil, err
		}
	}

	return def, nil
}
