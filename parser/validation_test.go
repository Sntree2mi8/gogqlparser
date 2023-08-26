package parser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"log"
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
		{
			name: "The directive must not contain the use of a directive which references itself directly",
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
									Directives: []ast.Directive{
										{
											Name: "test",
										},
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationInputFieldDefinition,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "A directive definition must not contain the use of a directive which references itself indirectly by referencing a Type or Directive which transitively includes a reference to this directive.",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					TypeDefinitions: []ast.TypeDefinition{
						&ast.InputObjectTypeDefinition{
							Name: "testInputObject",
							Directives: []ast.Directive{
								{
									Name: "test",
								},
							},
						},
					},
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "test",
							ArgumentsDefinition: []ast.InputValueDefinition{
								{
									Name: "arg",
									Type: ast.Type{
										NamedType: "testInputObject",
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationInputObject,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "A directive definition must not contain the use of a directive which references itself indirectly by referencing a Type or Directive which transitively includes a reference to this directive.",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					TypeDefinitions: []ast.TypeDefinition{
						&ast.InputObjectTypeDefinition{
							Name: "testInputObject",
							InputFields: []ast.InputValueDefinition{
								{
									Name: "test",
									Type: ast.Type{
										NamedType: "String",
									},
									Directives: []ast.Directive{
										{
											Name: "directiveForInputFieldDefinition",
										},
									},
								},
							},
						},
					},
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "directiveForInputFieldDefinition",
							ArgumentsDefinition: []ast.InputValueDefinition{
								{
									Name: "arg",
									Type: ast.Type{
										NamedType: "testInputObject",
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationInputFieldDefinition,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Argument type must be scalar, enum, input object",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					TypeDefinitions: []ast.TypeDefinition{
						&ast.ObjectTypeDefinition{
							Name: "TestObject",
						},
					},
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "test",
							ArgumentsDefinition: []ast.InputValueDefinition{
								{
									Name: "__arg",
									Type: ast.Type{
										NamedType: "TestObject",
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationInputFieldDefinition,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Argument type must be scalar, enum, input object",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					TypeDefinitions: []ast.TypeDefinition{
						&ast.UnionTypeDefinition{
							Name: "TestUnion",
						},
					},
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "test",
							ArgumentsDefinition: []ast.InputValueDefinition{
								{
									Name: "__arg",
									Type: ast.Type{
										NamedType: "TestUnion",
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationInputFieldDefinition,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Argument type must be scalar, enum, input object",
			args: args{
				doc: &ast.TypeSystemExtensionDocument{
					TypeDefinitions: []ast.TypeDefinition{
						&ast.InterfaceTypeDefinition{
							Name: "TestInterface",
						},
					},
					DirectiveDefinitions: []ast.DirectiveDefinition{
						{
							Name: "test",
							ArgumentsDefinition: []ast.InputValueDefinition{
								{
									Name: "__arg",
									Type: ast.Type{
										NamedType: "TestInterface",
									},
								},
							},
							DirectiveLocations: []ast.DirectiveLocation{
								ast.DirectiveLocationInputFieldDefinition,
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
			v := newValidator(tt.args.doc)
			if err := v.validateDirectiveDefinitions(); (err != nil) != tt.wantErr {
				t.Errorf("validateDirectiveDefinitions() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				log.Println(err)
			}
		})
	}
}
