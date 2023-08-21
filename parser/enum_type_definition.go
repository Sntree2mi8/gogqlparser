package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

// https://spec.graphql.org/October2021/#EnumValuesDefinition
func parseEnumValuesDefinition(l *LexerWrapper) (enumValuesDef []ast.EnumValueDefinition, err error) {
	if err = l.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		var enumValueDef ast.EnumValueDefinition
		enumValueDef.Description, _ = l.ReadDescription()

		if enumValueDef.Value.Value, err = l.ReadNameValue(); err != nil {
			return nil, err
		}

		if l.CheckKind(gogqllexer.At) {
			if enumValueDef.Directives, err = parseDirectives(l); err != nil {
				return nil, err
			}
		}

		enumValuesDef = append(enumValuesDef, enumValueDef)

		if l.SkipIf(gogqllexer.BraceR) {
			break
		}
	}
	return enumValuesDef, nil
}

// https://spec.graphql.org/October2021/#sec-Enums
func ParseEnumTypeDefinition(l *LexerWrapper, description string) (def *ast.EnumTypeDefinition, err error) {
	def = &ast.EnumTypeDefinition{
		Description: description,
	}

	if err = l.SkipKeyword("enum"); err != nil {
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

	if def.EnumValue, err = parseEnumValuesDefinition(l); err != nil {
		return nil, err
	}

	return def, nil
}

// https://spec.graphql.org/October2021/#sec-Enum-Extensions
// NOTION: consume "extend" keyword before call this function.
func ParseEnumExtensionDefinition(l *LexerWrapper) (def *ast.EnumTypeExtension, err error) {
	def = &ast.EnumTypeExtension{}

	if err := l.SkipKeyword("enum"); err != nil {
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

	if l.CheckKind(gogqllexer.BraceL) {
		if def.EnumValue, err = parseEnumValuesDefinition(l); err != nil {
			return nil, err
		}
	}

	return def, nil
}
