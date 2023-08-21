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
		l           *LexerWrapper
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.InputObjectTypeDefinition
		wantErr bool
	}{
		{
			name: "simple input object",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
input User {
	name: String!
}
`,
						),
					),
				),
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
input User @deprecated @danger {
	name: String!
}
`,
						),
					),
				),
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
input User {
	"name is deprecated"
	name: String! @deprecated
}
`,
						),
					),
				),
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
			gotDef, err := ParseInputObjectTypeDefinition(tt.args.l, tt.args.description)
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
	type args struct {
		l *LexerWrapper
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.InputObjectTypeExtension
		wantErr bool
	}{
		{
			name: "simple input object extension",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
input Restaurant {
	name: String!
}
`,
						),
					),
				),
			},
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
input Restaurant @input_directive
`,
						),
					),
				),
			},
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
extend input Restaurant @input_directive {
    name: String!
}
`,
						),
					),
				),
			},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDef, err := ParseInputObjectTypeExtension(tt.args.l)
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
