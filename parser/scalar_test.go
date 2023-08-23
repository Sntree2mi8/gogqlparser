package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseScalarTypeExtension(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		wantDef *ast.ScalarTypeExtension
		wantErr bool
	}{
		{
			name: "simple scalar type extension",
			schema: `
scalar Int @max(n: 100)
`,
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
			schema: `
scalar Int

type Example {
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
			gotDef, err := p.ParseScalarTypeExtension()
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
