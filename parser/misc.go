package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

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

func parseInputValueDefinition(l *LexerWrapper, description string) (def ast.InputValueDefinition, err error) {
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

func ParseArgumentsDefinition(l *LexerWrapper) (defs []ast.InputValueDefinition, err error) {
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

func parseType(l *LexerWrapper) (t ast.Type, err error) {
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
