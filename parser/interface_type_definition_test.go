package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseInterfaceTypeDefinition(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name    string
		schema  string
		args    args
		wantD   *ast.InterfaceTypeDefinition
		wantErr bool
	}{
		{
			name: "implements interface",
			schema: `
interface User implements Node {
	id: ID!
	name: String!
}
`,
			args: args{
				description: "this is description",
			},
			wantD: &ast.InterfaceTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Interfaces:  []string{"Node"},
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
		{
			name: "with directive",
			schema: `
interface User @deprecated(reason: "this is deprecated") {
	id: ID!
	name: String!
}
`,
			args: args{
				description: "this is description",
			},
			wantD: &ast.InterfaceTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Directives: []ast.Directive{
					{
						Name: "deprecated",
						Arguments: []ast.Argument{
							{
								Name:  "reason",
								Value: "\"this is deprecated\"",
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
		{
			name: "with directive, implements interface",
			schema: `
interface User implements Node @deprecated(reason: "this is deprecated") {
	id: ID!
	name: String!
}
`,
			args: args{
				description: "this is description",
			},
			wantD: &ast.InterfaceTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Interfaces:  []string{"Node"},
				Directives: []ast.Directive{
					{
						Name: "deprecated",
						Arguments: []ast.Argument{
							{
								Name:  "reason",
								Value: "\"this is deprecated\"",
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
			p := &parser{
				lexer: gogqllexer.New(strings.NewReader(tt.schema)),
			}
			gotD, err := p.ParseInterfaceTypeDefinition(tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterfaceTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("ParseInterfaceTypeDefinition() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestParseInterfaceTypeExtension(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		wantDef *ast.InterfaceTypeExtension
		wantErr bool
	}{
		{
			name: "simple interface type extension",
			schema: `
interface RestaurantInterface {
    address: String!
}
`,
			wantDef: &ast.InterfaceTypeExtension{
				Name: "RestaurantInterface",
				FieldsDefinition: []*ast.FieldDefinition{
					{
						Name: "address",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
		{
			name: "implements other interface",
			schema: `
interface RestaurantInterface implements Store
`,
			wantDef: &ast.InterfaceTypeExtension{
				Name:                "RestaurantInterface",
				ImplementInterfaces: []string{"Store"},
			},
		},
		{
			name: "onply directive",
			schema: `
interface RestaurantInterface @interface_directive
`,
			wantDef: &ast.InterfaceTypeExtension{
				Name: "RestaurantInterface",
				Directives: []ast.Directive{
					{
						Name: "interface_directive",
					},
				},
			},
		},
		{
			name: "extend but do nothing",
			schema: `
interface RestaurantInterface
`,
			wantDef: &ast.InterfaceTypeExtension{
				Name: "RestaurantInterface",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				lexer: gogqllexer.New(strings.NewReader(tt.schema)),
			}
			gotDef, err := p.ParseInterfaceTypeExtension()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterfaceTypeExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseInterfaceTypeExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
