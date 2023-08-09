package gogqlparser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

var testSchemaOnly = `
"""
A simple GraphQL schema which is well described.
"""
schema {
    query: Query
}
`

func TestParser_parseTypeSystemDocument(t *testing.T) {
	type args struct {
		src *ast.Source
	}
	tests := []struct {
		name    string
		args    args
		want    *ast.TypeSystemExtensionDocument
		wantErr bool
	}{
		{
			name: "parse schema definition",
			args: args{src: &ast.Source{
				Name: "testSchemaOnly.graphql",
				Body: testSchemaOnly,
			},
			},
			want: &ast.TypeSystemExtensionDocument{
				SchemaDefinitions: []ast.SchemaDefinition{
					{
						Description: `"""
A simple GraphQL schema which is well described.
"""`,
						Directives: []ast.Directive{},
						RootOperationTypeDefinitions: []ast.RootOperationTypeDefinition{
							{
								OperationType: ast.OperationTypeQuery,
								Type:          "Query",
							},
						},
					},
				},
				TypeDefinitions:      []ast.TypeDefinition{},
				DirectiveDefinitions: []ast.DirectiveDefinition{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, err := p.parseTypeSystemDocument(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTypeSystemDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// diffがわかりづらいので足した。機能的に必要としているわけじゃないのであとで消す。
			if df := cmp.Diff(got, tt.want); df != "" {
				t.Errorf("parseTypeSystemDocument() diff = %v", df)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTypeSystemDocument() got = %v, want %v", got, tt.want)
			}
		})
	}
}
