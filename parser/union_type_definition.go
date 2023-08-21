package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"log"
)

// https://spec.graphql.org/October2021/#UnionMemberTypes
func parseUnionMemberTypes(l *LexerWrapper) (memberTypes []ast.Type, err error) {
	if err = l.Skip(gogqllexer.Equal); err != nil {
		return nil, err
	}

	l.SkipIf(gogqllexer.Pipe)

	for {
		var mt ast.Type
		if mt.NamedType, err = l.ReadNameValue(); err != nil {
			log.Println("here")
			return nil, err
		}

		memberTypes = append(memberTypes, mt)

		if !l.SkipIf(gogqllexer.Pipe) {
			break
		}
	}

	return memberTypes, nil
}

// https://spec.graphql.org/October2021/#sec-Unions
func ParseUnionTypeDefinition(l *LexerWrapper, description string) (def *ast.UnionTypeDefinition, err error) {
	def = &ast.UnionTypeDefinition{
		Description: description,
	}

	if err = l.SkipKeyword("union"); err != nil {
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

	if err = l.PeekAndMustBe(
		[]gogqllexer.Kind{gogqllexer.Equal},
		func(t gogqllexer.Token, advanceLexer func()) error {
			if def.MemberTypes, err = parseUnionMemberTypes(l); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return def, nil
}

// ParseUnionTypeExtension parses union type extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphql.org/October2021/#sec-Union-Extensions
func ParseUnionTypeExtension(l *LexerWrapper) (def *ast.UnionTypeExtension, err error) {
	return def, err
}
