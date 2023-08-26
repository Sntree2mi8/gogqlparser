package validator

import (
	"fmt"
	"github.com/Sntree2mi8/gogqlparser/ast"
)

func ValidateTypeSystemExtensionDocument(doc *ast.TypeSystemExtensionDocument) error {
	doc = doc.Merge(BultinTypeSystemExtensionDocument)
	v, err := newValidator(doc)
	if err != nil {
		return err
	}

	return v.validateDirectiveDefinitions()
}

type validator struct {
	doc *ast.TypeSystemExtensionDocument

	typeDefs      map[string]ast.TypeDefinition
	directiveDefs map[string]ast.DirectiveDefinition
}

func newValidator(doc *ast.TypeSystemExtensionDocument) (*validator, error) {
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
