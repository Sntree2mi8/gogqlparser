package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseEnumTypeDefinition(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name    string
		schema  string
		args    args
		wantDef *ast.EnumTypeDefinition
		wantErr bool
	}{
		{
			name: "simple enum",
			schema: `
enum UserKind {
	ADMIN
	NORMAL
}
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.EnumTypeDefinition{
				Description: "this is description",
				Name:        "UserKind",
				EnumValue: []ast.EnumValueDefinition{
					{
						Value: ast.EnumValue{
							Value: "ADMIN",
						},
					},
					{
						Value: ast.EnumValue{
							Value: "NORMAL",
						},
					},
				},
			},
		},
		{
			name: "with directives",
			schema: `
enum UserKind @deprecated {
	ADMIN @deprecated
	NORMAL
}
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.EnumTypeDefinition{
				Description: "this is description",
				Name:        "UserKind",
				Directives: []ast.Directive{
					{
						Name: "deprecated",
					},
				},
				EnumValue: []ast.EnumValueDefinition{
					{
						Directives: []ast.Directive{
							{
								Name: "deprecated",
							},
						},
						Value: ast.EnumValue{
							Value: "ADMIN",
						},
					},
					{
						Value: ast.EnumValue{
							Value: "NORMAL",
						},
					},
				},
			},
		},
		{
			name: "with enum value description",
			schema: `
enum UserKind {
	"this is admin" ADMIN
	"""
	this is normal
	"""
	NORMAL
}
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.EnumTypeDefinition{
				Description: "this is description",
				Name:        "UserKind",
				EnumValue: []ast.EnumValueDefinition{
					{
						Description: "\"this is admin\"",
						Value: ast.EnumValue{
							Value: "ADMIN",
						},
					},
					{
						Description: "\"\"\"\n\tthis is normal\n\t\"\"\"",
						Value: ast.EnumValue{
							Value: "NORMAL",
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
			gotDef, err := p.ParseEnumTypeDefinition(tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnumTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseEnumTypeDefinition() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}

// NOTION:
// "extend" keyword is assumed to be consumed before this function is called
func TestParseEnumExtensionDefinition(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		wantDef *ast.EnumTypeExtension
		wantErr bool
	}{
		{
			name: "simple enum extension",
			schema: `
enum RestaurantKind {
        CHINESE
}
`,
			wantDef: &ast.EnumTypeExtension{
				Name: "RestaurantKind",
				EnumValue: []ast.EnumValueDefinition{
					{
						Value: ast.EnumValue{
							Value: "CHINESE",
						},
					},
				},
			},
		},
		{
			name: "only directive",
			schema: `
enum RestaurantKind @enum_directive
`,
			wantDef: &ast.EnumTypeExtension{
				Name: "RestaurantKind",
				Directives: []ast.Directive{
					{
						Name: "enum_directive",
					},
				},
			},
		},
		{
			name: "extend but do nothing",
			schema: `
enum RestaurantKind
`,
			wantDef: &ast.EnumTypeExtension{
				Name: "RestaurantKind",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				lexer: gogqllexer.New(strings.NewReader(tt.schema)),
			}
			gotDef, err := p.ParseEnumTypeExtension()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnumTypeExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseEnumTypeExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
