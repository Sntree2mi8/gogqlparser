package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseInputObjectTypeDefinition(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name    string
		schema  string
		args    args
		wantDef *ast.InputObjectTypeDefinition
		wantErr bool
	}{
		{
			name: "simple input object",
			schema: `
input User {
	name: String!
}
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.InputObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				InputFields: []ast.InputValueDefinition{
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
						RawDefaultValue: "",
					},
				},
			},
		},
		{
			name: "with directive",
			schema: `
input User @deprecated @danger {
	name: String!
}
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.InputObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Directives: []ast.Directive{
					{
						Name: "deprecated",
					},
					{
						Name: "danger",
					},
				},
				InputFields: []ast.InputValueDefinition{
					{
						Name: "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
						RawDefaultValue: "",
					},
				},
			},
		},
		{
			name: "with field directive",
			schema: `
input User {
	"name is deprecated"
	name: String! @deprecated
}
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.InputObjectTypeDefinition{
				Description: "this is description",
				Name:        "User",
				InputFields: []ast.InputValueDefinition{
					{
						Description: "\"name is deprecated\"",
						Name:        "name",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
						RawDefaultValue: "",
						Directives: []ast.Directive{
							{
								Name: "deprecated",
							},
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
			gotDef, err := p.ParseInputObjectTypeDefinition(tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInputObjectTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseInputObjectTypeDefinition() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}

func TestParseInputObjectTypeExtension(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		wantDef *ast.InputObjectTypeExtension
		wantErr bool
	}{
		{
			name: "simple input object extension",
			schema: `
input Restaurant {
	name: String!
}
`,
			wantDef: &ast.InputObjectTypeExtension{
				Name: "Restaurant",
				InputsFieldDefinition: []ast.InputValueDefinition{
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
			schema: `
input Restaurant @input_directive
`,
			wantDef: &ast.InputObjectTypeExtension{
				Name: "Restaurant",
				Directives: []ast.Directive{
					{
						Name: "input_directive",
					},
				},
			},
		},
		{
			name: "directive and field",
			schema: `
input Restaurant @input_directive {
    name: String!
}
`,
			wantDef: &ast.InputObjectTypeExtension{
				Name: "Restaurant",
				Directives: []ast.Directive{
					{
						Name: "input_directive",
					},
				},
				InputsFieldDefinition: []ast.InputValueDefinition{
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
			name: "extend input object type needs at least one field or directive",
			schema: `
input Restaurant 

type OtherType {
	name: String!
}
`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				lexer: gogqllexer.New(strings.NewReader(tt.schema)),
			}
			gotDef, err := p.ParseInputObjectTypeExtension()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInputObjectTypeExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseInputObjectTypeExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
