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
		l           *LexerWrapper
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.EnumTypeDefinition
		wantErr bool
	}{
		{
			name: "simple enum",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
enum UserKind {
	ADMIN
	NORMAL
}
`,
						),
					),
				),
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
enum UserKind @deprecated {
	ADMIN @deprecated
	NORMAL
}
`,
						),
					),
				),
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
enum UserKind {
	"this is admin" ADMIN
	"""
	this is normal
	"""
	NORMAL
}
`,
						),
					),
				),
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
			gotDef, err := ParseEnumTypeDefinition(tt.args.l, tt.args.description)
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
	type args struct {
		l *LexerWrapper
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.EnumTypeExtension
		wantErr bool
	}{
		{
			name: "simple enum extension",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
enum RestaurantKind {
        CHINESE
}
`,
						),
					),
				),
			},
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
enum RestaurantKind @enum_directive
`,
						),
					),
				),
			},
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
enum RestaurantKind
`,
						),
					),
				),
			},
			wantDef: &ast.EnumTypeExtension{
				Name: "RestaurantKind",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDef, err := ParseEnumExtensionDefinition(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnumExtensionDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseEnumExtensionDefinition() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
