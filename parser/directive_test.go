package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseDirectiveDefinition(t *testing.T) {
	type args struct {
		l           *LexerWrapper
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.DirectiveDefinition
		wantErr bool
	}{
		{
			name: "simple directive",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing on SCHEMA
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
			},
		},
		{
			name: "with description",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing on SCHEMA
`,
						),
					),
				),
				description: "this is description",
			},
			wantDef: &ast.DirectiveDefinition{
				Description: "this is description",
				Name:        "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
			},
		},
		{
			name: "with argumentsDefinition",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing(arg1: String!) on SCHEMA
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
				ArgumentsDefinition: []ast.InputValueDefinition{
					{
						Name: "arg1",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
		{
			name: "with repeatable",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing repeatable on SCHEMA
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
				IsRepeatable: true,
			},
		},
		{
			name: "with all optional items",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing(arg1: String!) repeatable on SCHEMA
`,
						),
					),
				),
				description: "this is description",
			},
			wantDef: &ast.DirectiveDefinition{
				Description:  "this is description",
				Name:         "testing",
				IsRepeatable: true,
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
				ArgumentsDefinition: []ast.InputValueDefinition{
					{
						Name: "arg1",
						Type: ast.Type{
							NamedType: "String",
							NotNull:   true,
						},
					},
				},
			},
		},
		{
			name: "single directive location",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing on SCHEMA
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
			},
		},
		{
			name: "single directive location with leading pipe",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing on | SCHEMA
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
				},
			},
		},
		{
			name: "multiple directive locations",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing on SCHEMA | SCALAR | OBJECT
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
					ast.DirectiveLocationScalar,
					ast.DirectiveLocationObject,
				},
			},
		},
		{
			name: "parse all type system directive locations",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
directive @testing on SCHEMA 
| SCALAR 
| OBJECT 
| FIELD_DEFINITION 
| ARGUMENT_DEFINITION 
| INTERFACE 
| UNION 
| ENUM 
| ENUM_VALUE 
| INPUT_OBJECT 
| INPUT_FIELD_DEFINITION
`,
						),
					),
				),
			},
			wantDef: &ast.DirectiveDefinition{
				Name: "testing",
				DirectiveLocations: []ast.DirectiveLocation{
					ast.DirectiveLocationSchema,
					ast.DirectiveLocationScalar,
					ast.DirectiveLocationObject,
					ast.DirectiveLocationFieldDefinition,
					ast.DirectiveLocationArgumentDefinition,
					ast.DirectiveLocationInterface,
					ast.DirectiveLocationUnion,
					ast.DirectiveLocationEnum,
					ast.DirectiveLocationEnumValue,
					ast.DirectiveLocationInputObject,
					ast.DirectiveLocationInputFieldDefinition,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDef, err := ParseDirectiveDefinition(tt.args.l, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectiveDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseDirectiveDefinition() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
