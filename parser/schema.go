package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func (p *parser) ParseRootOperationTypeDefinitions() (defs map[string]*ast.RootOperationTypeDefinition, err error) {
	defs = make(map[string]*ast.RootOperationTypeDefinition)

	if err := p.Skip(gogqllexer.BraceL); err != nil {
		return nil, err
	}

	for {
		var rootOperationName string
		var rootOperationTypeName string

		if rootOperationName, err = p.ReadNameValue(); err != nil {
			return nil, err
		}

		switch rootOperationName {
		default:
			return nil, fmt.Errorf("unexpected root operation name %s", rootOperationName)
		case "query", "mutation", "subscription":
		}

		if err = p.Skip(gogqllexer.Colon); err != nil {
			return nil, err
		}

		if rootOperationTypeName, err = p.ReadNameValue(); err != nil {
			return nil, err
		}

		if _, ok := defs[rootOperationName]; ok {
			return nil, fmt.Errorf("duplicate root operation name %s", rootOperationName)
		}

		defs[rootOperationName] = &ast.RootOperationTypeDefinition{
			Type: rootOperationTypeName,
		}

		if p.SkipIf(gogqllexer.BraceR) {
			break
		}
	}

	if len(defs) == 0 {
		return nil, fmt.Errorf("schema definition must have at least one root operation type definition")
	}

	return defs, err
}

func (p *parser) ParseSchemaDefinition(description string) (def *ast.SchemaDefinition, err error) {
	def = &ast.SchemaDefinition{
		Description: description,
	}

	if err = p.SkipKeyword("schema"); err != nil {
		return nil, err
	}

	if p.CheckKind(gogqllexer.BraceL) {
		rootOperationTypeDefs, err := p.ParseRootOperationTypeDefinitions()
		if err != nil {
			return nil, err
		}
		def.Query = rootOperationTypeDefs["query"]
		def.Mutation = rootOperationTypeDefs["mutation"]
		def.Subscription = rootOperationTypeDefs["subscription"]
	} else {
		return nil, fmt.Errorf("schema definition must have at least one root operation type definition")
	}

	return def, err
}

// ParseSchemaExtension parse schema extension.
// "extend" keyword must be consumed before calling this function.
//
// Reference: https://spec.graphqp.org/October2021/#sec-Schema-Extension
func (p *parser) ParseSchemaExtension() (def *ast.SchemaExtension, err error) {
	def = &ast.SchemaExtension{}

	if err = p.SkipKeyword("schema"); err != nil {
		return nil, err
	}
	var canOmitRootOperationTypes bool
	if p.CheckKind(gogqllexer.At) {
		if def.Directives, err = p.parseDirectives(); err != nil {
			return nil, err
		}

		canOmitRootOperationTypes = true
	}

	if p.CheckKind(gogqllexer.BraceL) {
		rootOperationTypeDefs, err := p.ParseRootOperationTypeDefinitions()
		if err != nil {
			return nil, err
		}
		def.Query = rootOperationTypeDefs["query"]
		def.Mutation = rootOperationTypeDefs["mutation"]
		def.Subscription = rootOperationTypeDefs["subscription"]
	} else if !canOmitRootOperationTypes {
		return nil, fmt.Errorf("schema extension must have at least one root operation type definition or directive")
	}

	return def, err
}
