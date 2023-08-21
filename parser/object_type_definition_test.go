package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseObjectTypeDefinition(t *testing.T) {
	type args struct {
		l           *LexerWrapper
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantD   *ast.ObjectTypeDefinition
		wantErr bool
	}{
		{
			name: "object implements interface",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type User implements Node {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.ObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
				Interfaces: []string{"Node"},
			},
		},
		{
			name: "object implements interface",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type User implements & Node {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.ObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
				Interfaces: []string{"Node"},
			},
		},
		{
			name: "object implements interfaces",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type User implements & Node & UserInterface {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.ObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
				Interfaces: []string{"Node", "UserInterface"},
			},
		},
		{
			name: "object with directive",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type User @role(role: "admin") {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.ObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Directives: []ast.Directive{
					{
						Name: "role",
						Arguments: []ast.Argument{
							{
								Name:  "role",
								Value: "\"admin\"",
							},
						},
					},
				},
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := ParseObjectTypeDefinition(tt.args.l, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseObjectTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("ParseObjectTypeDefinition() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestParseObjectTypeExtension(t *testing.T) {
	type args struct {
		l *LexerWrapper
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.ObjectTypeExtension
		wantErr bool
	}{
		{
			name: "simple object type extension",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type Restaurant {
    name: String!
}
`,
						),
					),
				),
			},
			wantDef: &ast.ObjectTypeExtension{
				Name: "Restaurant",
				FieldsDefinition: []*ast.FieldDefinition{
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
		{
			name: "directive and field",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type Restaurant @object_directive {
    name: String!
}
`,
						),
					),
				),
			},
			wantDef: &ast.ObjectTypeExtension{
				Name: "Restaurant",
				Directives: []ast.Directive{
					{
						Name: "object_directive",
					},
				},
				FieldsDefinition: []*ast.FieldDefinition{
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
		{
			name: "implement and directive and field",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type Restaurant implements Store @object_directive {
    name: String!
}
`,
						),
					),
				),
			},
			wantDef: &ast.ObjectTypeExtension{
				Name: "Restaurant",
				ImplementInterfaces: []string{
					"Store",
				},
				Directives: []ast.Directive{
					{
						Name: "object_directive",
					},
				},
				FieldsDefinition: []*ast.FieldDefinition{
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
		{
			name: "only directive",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type Restaurant @object_directive
`,
						),
					),
				),
			},
			wantDef: &ast.ObjectTypeExtension{
				Name: "Restaurant",
				Directives: []ast.Directive{
					{
						Name: "object_directive",
					},
				},
			},
		},
		{
			name: "only implement",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
type Restaurant implements Store
`,
						),
					),
				),
			},
			wantDef: &ast.ObjectTypeExtension{
				Name: "Restaurant",
				ImplementInterfaces: []string{
					"Store",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDef, err := ParseObjectTypeExtension(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseObjectTypeExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseObjectTypeExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
