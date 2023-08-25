package parser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"testing"
)

func Test_validateDirectiveDefinitions(t *testing.T) {
	type args struct {
		doc *ast.TypeSystemExtensionDocument
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "The directive must not have a name which begins with the characters \"__\" (two underscores).",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "__test",
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationField,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "The argument must not have a name which begins with the characters \"__\" (two underscores).",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "test",
							ArgumentsDefinition: []ast.InputValueDefinition{
								{
									Name: "__arg",
									Type: ast.Type{
										NamedType: "String",
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationField,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDirectiveDefinitions(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("validateDirectiveDefinitions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
