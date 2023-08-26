package validator

import "github.com/Sntree2mi8/gogqlparser/ast"

var BultinTypeSystemExtensionDocument = &ast.TypeSystemExtensionDocument{
	SchemaDefinitions: nil,
	TypeDefinitions: []ast.TypeDefinition{
		&ast.ScalarTypeDefinition{
			Name: "Int",
		},
		&ast.ScalarTypeDefinition{
			Name: "Float",
		},
		&ast.ScalarTypeDefinition{
			Name: "String",
		},
		&ast.ScalarTypeDefinition{
			Name: "Boolean",
		},
		&ast.ScalarTypeDefinition{
			Name: "ID",
		},
	},
	DirectiveDefinitions: []ast.DirectiveDefinition{
		{
			Name: "include",
			ArgumentsDefinition: []ast.InputValueDefinition{
				{
					Name: "if",
					Type: ast.Type{
						NamedType: "Boolean",
						NotNull:   true,
					},
				},
			},
			DirectiveLocations: []ast.DirectiveLocation{
				ast.DirectiveLocationField,
				ast.DirectiveLocationFragmentSpread,
				ast.DirectiveLocationInlineFragment,
			},
		},
		{
			Name: "skip",
			ArgumentsDefinition: []ast.InputValueDefinition{
				{
					Name: "if",
					Type: ast.Type{
						NamedType: "Boolean",
						NotNull:   true,
					},
				},
			},
			DirectiveLocations: []ast.DirectiveLocation{
				ast.DirectiveLocationField,
				ast.DirectiveLocationFragmentSpread,
				ast.DirectiveLocationInlineFragment,
			},
		},
		{
			Name: "deprecated",
			ArgumentsDefinition: []ast.InputValueDefinition{
				{
					Name: "reason",
					Type: ast.Type{
						NamedType: "String",
						NotNull:   true,
					},
					RawDefaultValue: `"No longer supported"`,
				},
			},
			DirectiveLocations: []ast.DirectiveLocation{
				ast.DirectiveLocationField,
				ast.DirectiveLocationArgumentDefinition,
				ast.DirectiveLocationInputFieldDefinition,
				ast.DirectiveLocationEnumValue,
			},
		},
		{
			Name: "specifiedBy",
			ArgumentsDefinition: []ast.InputValueDefinition{
				{
					Name: "url",
					Type: ast.Type{
						NamedType: "String",
						NotNull:   true,
					},
				},
			},
			DirectiveLocations: []ast.DirectiveLocation{
				ast.DirectiveLocationScalar,
			},
		},
	},
}
