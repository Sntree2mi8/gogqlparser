package validator

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func ValidateTypeSystemExtensionDocument(doc *ast.TypeSystemExtensionDocument) (*ast.Schema, error) {
	doc = doc.Merge(BultinTypeSystemExtensionDocument)

	v, err := newValidator(doc)
	if err != nil {
		return nil, err
	}

	if err := v.validateDirectiveDefinitions(); err != nil {
		return nil, err
	}
	// TODO: TypeDefinitionsのvalidation
	// TODO: RootOperationTypeDefinitionsのvalidation

	// TODO: 最終的にSchemaにして返す
	return &ast.Schema{
		Description:  "",
		Query:        nil,
		Mutation:     nil,
		Subscription: nil,
		Types:        nil,
		Directives:   nil,
	}, nil
}

func (v *validator) validateTypeExtensions() error {
	return nil
}

type validator struct {
	doc *ast.TypeSystemExtensionDocument

	// temp
	typeDefs      map[string]ast.TypeDefinition
	directiveDefs map[string]ast.DirectiveDefinition
}

func newValidator(doc *ast.TypeSystemExtensionDocument) (*validator, error) {
	// 同じTypeが定義できない系はここでチェックしている
	typeDefs := make(map[string]ast.TypeDefinition, len(doc.TypeDefinitions))
	for _, def := range doc.TypeDefinitions {
		if _, ok := typeDefs[def.TypeName()]; ok {
			return nil, fmt.Errorf("duplicate type definition: %s", def.TypeName())
		}
		typeDefs[def.TypeName()] = def
	}

	directiveDefs := make(map[string]ast.DirectiveDefinition, len(doc.DirectiveDefinitions))
	for _, def := range doc.DirectiveDefinitions {
		if _, ok := directiveDefs[def.Name]; ok {
			return nil, fmt.Errorf("duplicate directive definition: %s", def.Name)
		}
		directiveDefs[def.Name] = def
	}

	return &validator{
		doc:           doc,
		typeDefs:      typeDefs,
		directiveDefs: directiveDefs,
	}, nil
}
