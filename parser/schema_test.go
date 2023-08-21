package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseSchemaExtension(t *testing.T) {
	type args struct {
		l *LexerWrapper
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.SchemaExtension
		wantErr bool
	}{
		{
			name: "simple schema extension",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
schema {
    query: Query
}
`,
						),
					),
				),
			},
			wantDef: &ast.SchemaExtension{
				Query: &ast.RootOperationTypeDefinition{Type: "Query"},
			},
		},
		{
			name: "with directive",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
schema @schema_directive {
    query: Query
}
`,
						),
					),
				),
			},
			wantDef: &ast.SchemaExtension{
				Directives: []ast.Directive{
					{
						Name: "schema_directive",
					},
				},
				Query: &ast.RootOperationTypeDefinition{Type: "Query"},
			},
		},
		{
			name: "directive only",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
schema @schema_directive
`,
						),
					),
				),
			},
			wantDef: &ast.SchemaExtension{
				Directives: []ast.Directive{
					{
						Name: "schema_directive",
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
schema
`,
						),
					),
				),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDef, err := ParseSchemaExtension(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSchemaExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseSchemaExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
