package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func parseDirectiveArguments(l *LexerWrapper) (args []ast.Argument, err error) {
	if err = l.Skip(gogqllexer.ParenL); err != nil {
		return nil, err
	}

	for {
		arg := ast.Argument{}

		if arg.Name, err = l.ReadNameValue(); err != nil {
			return nil, err
		}
		if l.Skip(gogqllexer.Colon) != nil {
			return nil, err
		}
		if err = l.PeekAndMustBe(
			[]gogqllexer.Kind{gogqllexer.Int, gogqllexer.Float, gogqllexer.Name, gogqllexer.String, gogqllexer.BlockString},
			func(t gogqllexer.Token, advanceLexer func()) error {
				defer advanceLexer()
				arg.Value = t.Value
				return nil
			},
		); err != nil {
			return nil, err
		}

		args = append(args, arg)

		if l.SkipIf(gogqllexer.ParenR) {
			break
		}
	}

	return args, err
}

func parseDirectives(l *LexerWrapper) (directives []ast.Directive, err error) {
	for {
		var d ast.Directive

		d, err = parseDirective(l)
		if err != nil {
			return nil, err
		}

		directives = append(directives, d)

		if !l.CheckKind(gogqllexer.At) {
			break
		}
	}

	return directives, nil
}

func parseDirective(l *LexerWrapper) (d ast.Directive, err error) {
	if err = l.Skip(gogqllexer.At); err != nil {
		return d, err
	}

	if d.Name, err = l.ReadNameValue(); err != nil {
		return d, err
	}

	if l.CheckKind(gogqllexer.ParenL) {
		if d.Arguments, err = parseDirectiveArguments(l); err != nil {
			return d, err
		}
	}

	return d, err
}

func ParsDirectiveLocation(v string) ast.DirectiveLocation {
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

func parseDirectiveLocations(l *LexerWrapper) (locs []ast.DirectiveLocation, err error) {
	l.SkipIf(gogqllexer.Pipe)

	for {
		var locationValue string
		if locationValue, err = l.ReadNameValue(); err != nil {
			return nil, err
		}
		loc := ParsDirectiveLocation(locationValue)
		if loc == ast.DirectiveLocationUnknown {
			return nil, fmt.Errorf("unexpected token %+v", locationValue)
		}

		locs = append(locs, loc)

		if l.SkipIf(gogqllexer.Pipe) {
			continue
		}
		break
	}

	return locs, nil
}

// ParseDirectiveDefinition parses a directive definition
//
// Reference: https://spec.graphql.org/October2021/#sec-Type-System.Directives
func ParseDirectiveDefinition(l *LexerWrapper, description string) (def *ast.DirectiveDefinition, err error) {
	def = &ast.DirectiveDefinition{
		Description: description,
	}

	if err = l.SkipKeyword("directive"); err != nil {
		return nil, err
	}
	if err = l.Skip(gogqllexer.At); err != nil {
		return nil, err
	}
	if def.Name, err = l.ReadNameValue(); err != nil {
		return nil, err
	}

	if l.CheckKind(gogqllexer.ParenL) {
		if def.ArgumentsDefinition, err = ParseArgumentsDefinition(l); err != nil {
			return nil, err
		}
	}

	if l.SkipKeywordIf("repeatable") {
		def.IsRepeatable = true
	}

	if err = l.SkipKeyword("on"); err != nil {
		return nil, err
	}

	def.DirectiveLocations, err = parseDirectiveLocations(l)
	if err != nil {
		return nil, err
	}

	return def, nil
}
