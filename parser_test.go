package gogqlparser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
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

func TestParser_parseTypeSystemDocument_ParseSchemaDefinition(t *testing.T) {
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

func TestParser_parseTypeSystemDocument_ParseDirectiveDefinition(t *testing.T) {
	testdataDir := "./testdata/parseTypeSystemDocument/ParseDirectiveDefinition"

	tests := []struct {
		name       string
		schemaPath string
		want       *ast.TypeSystemExtensionDocument
		wantErr    bool
	}{
		{
			name:       "parse schema definition",
			schemaPath: path.Join(testdataDir, "parseDirectiveDefinition.graphql"),
			want: &ast.TypeSystemExtensionDocument{
				SchemaDefinitions: []ast.SchemaDefinition{},
				TypeDefinitions:   []ast.TypeDefinition{},
				DirectiveDefinitions: []ast.DirectiveDefinition{
					{
						Description:  "",
						Name:         "directive",
						IsRepeatable: false,
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
					{
						Description:  "",
						Name:         "repeatableDirective",
						IsRepeatable: true,
						DirectiveLocations: []ast.DirectiveLocation{
							ast.DirectiveLocationSchema,
						},
					},
					{
						Description: "",
						Name:        "argumentDirective",
						ArgumentsDefinition: []ast.InputValueDefinition{
							{
								Name: "argStr",
								Type: ast.Type{
									NamedType: "String",
									NotNull:   true,
								},
							},
							{
								Name: "argStrNullable",
								Type: ast.Type{
									NamedType: "String",
									NotNull:   false,
								},
							},
							{
								Name: "argStrArray",
								Type: ast.Type{
									ListType: &ast.Type{
										NamedType: "String",
										NotNull:   true,
									},
									NotNull: true,
								},
							},
							{
								Name: "argStrNullableArray",
								Type: ast.Type{
									ListType: &ast.Type{
										NamedType: "String",
										NotNull:   false,
									},
									NotNull: true,
								},
							},
							{
								Name: "argStrNullableArrayNullable",
								Type: ast.Type{
									ListType: &ast.Type{
										NamedType: "String",
										NotNull:   false,
									},
									NotNull: false,
								},
							},
							{
								Name: "argStrNullableDefault",
								Type: ast.Type{
									NamedType: "String",
									NotNull:   false,
								},
								RawDefaultValue: `"default"`,
							},
							{
								Name: "argInt",
								Type: ast.Type{
									NamedType: "Int",
									NotNull:   true,
								},
							},
							{
								Name: "argIntNullable",
								Type: ast.Type{
									NamedType: "Int",
									NotNull:   false,
								},
							},
							{
								Name: "argIntNullableDefault",
								Type: ast.Type{
									NamedType: "Int",
									NotNull:   false,
								},
								RawDefaultValue: "0",
							},
							{
								Name: "argFloat",
								Type: ast.Type{
									NamedType: "Float",
									NotNull:   true,
								},
							},
							{
								Name: "argFloatNullable",
								Type: ast.Type{
									NamedType: "Float",
									NotNull:   false,
								},
							},
							{
								Name: "argFloatNullableDefault",
								Type: ast.Type{
									NamedType: "Float",
									NotNull:   false,
								},
								RawDefaultValue: "0.0",
							},
							{
								Name: "argBool",
								Type: ast.Type{
									NamedType: "Boolean",
									NotNull:   true,
								},
							},
							{
								Name: "argBoolNullable",
								Type: ast.Type{
									NamedType: "Boolean",
									NotNull:   false,
								},
							},
							{
								Name: "argBoolNullableDefault",
								Type: ast.Type{
									NamedType: "Boolean",
									NotNull:   false,
								},
								RawDefaultValue: "false",
							},
						},
						IsRepeatable: false,
						DirectiveLocations: []ast.DirectiveLocation{
							ast.DirectiveLocationSchema,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := os.ReadFile(tt.schemaPath)
			if err != nil {
				t.Fatal(err)
			}

			p := &Parser{}
			got, err := p.parseTypeSystemDocument(&ast.Source{
				Name: tt.schemaPath,
				Body: string(d),
			})
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
