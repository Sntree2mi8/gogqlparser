package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseScalarTypeExtension(t *testing.T) {
	type args struct {
		l *LexerWrapper
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.ScalarTypeExtension
		wantErr bool
	}{
		{
			name: "simple scalar type extension",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
scalar Int @max(n: 100)
`,
						),
					),
				),
			},
			wantDef: &ast.ScalarTypeExtension{
				Name: "Int",
				Directives: []ast.Directive{
					{
						Name: "max",
						Arguments: []ast.Argument{
							{
								Name:  "n",
								Value: "100",
							},
						},
					},
				},
			},
		},
		{
			name: "scalar type extension needs at least one directive",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(
							`
scalar Int

type Example {
    name: String!
}
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
			gotDef, err := ParseScalarTypeExtension(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseScalarTypeExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseScalarTypeExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
