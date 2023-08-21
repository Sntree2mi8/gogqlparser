package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func parseInputValueDefinition(l *LexerWrapper, description string) (def ast.InputValueDefinition, err error) {
	def.Description = description
	if def.Name, err = l.ReadNameValue(); err != nil {
		return def, err
	}

	if err = l.Skip(gogqllexer.Colon); err != nil {
		return def, err
	}

	def.Type, err = parseType(l)
	if err != nil {
		return def, err
	}

	if l.SkipIf(gogqllexer.Equal) {
		if err = l.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Int, gogqllexer.Float, gogqllexer.String, gogqllexer.BlockString, gogqllexer.Name},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()

				def.RawDefaultValue = t.Value
				return nil
			},
		); err != nil {
			return def, err
		}
	}

	if l.CheckKind(gogqllexer.At) {
		if def.Directives, err = parseDirectives(l); err != nil {
			return def, err
		}
	}

	return def, nil
}

// https://spec.graphql.org/October2021/#InputFieldsDefinition
func parseInputFieldsDefinition(l *LexerWrapper) (defs []ast.InputValueDefinition, err error) {
	if err = l.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		description, _ := l.ReadDescription()
		def, err := parseInputValueDefinition(l, description)
		if err != nil {
			return nil, err
		}

		defs = append(defs, def)

		if l.SkipIf(gogqllexer.BraceR) {
			break
		}
	}

	return defs, nil
}

// https://spec.graphql.org/October2021/#sec-Input-Objects
func ParseInputObjectTypeDefinition(l *LexerWrapper, description string) (def *ast.InputObjectTypeDefinition, err error) {
	def = &ast.InputObjectTypeDefinition{
		Description: description,
	}

	if err := l.SkipKeyword("input"); err != nil {
		return nil, err
	}

	if def.Name, err = l.ReadNameValue(); err != nil {
		return nil, err
	}

	if l.CheckKind(gogqllexer.At) {
		if def.Directives, err = parseDirectives(l); err != nil {
			return nil, err
		}
	}

	if def.InputFields, err = parseInputFieldsDefinition(l); err != nil {
		return nil, err
	}

	return def, nil
}

// ParseInputObjectTypeExtension parse input object type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Input-Object-Extensions
func ParseInputObjectTypeExtension(l *LexerWrapper) (def *ast.InputObjectTypeExtension, err error) {
	return def, nil
}
