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
				TypeDefinitions:      map[string]ast.TypeDefinition{},
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
				TypeDefinitions:   map[string]ast.TypeDefinition{},
				DirectiveDefinitions: []ast.DirectiveDefinition{
					{
						Name: "directive",
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
						Name:         "repeatableDirective",
						IsRepeatable: true,
						DirectiveLocations: []ast.DirectiveLocation{
							ast.DirectiveLocationSchema,
						},
					},
					{
						Description: "",
						Name:        "sampleDirective",
						ArgumentsDefinition: []ast.InputValueDefinition{
							{
								Name: "text",
								Type: ast.Type{
									NamedType: "String",
								},
							},
						},
						DirectiveLocations: []ast.DirectiveLocation{
							ast.DirectiveLocationInputFieldDefinition,
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
								Name: "argStrWithDirective",
								Type: ast.Type{
									NamedType: "String",
									NotNull:   true,
								},
								Directives: []ast.Directive{
									{
										Name: "sampleDirective",
										Arguments: []ast.Argument{
											{
												Name:  "text",
												Value: `"text"`,
											},
										},
									},
								},
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

func TestParser_parseTypeSystemDocument_parseRootOperationTypesSchema(t *testing.T) {
	parseRootOperationTypesSchema := `
schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}

type Query {
    example: String
}

type Mutation {
    example: String
}

type Subscription {
    example: String
}
`

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
			name: "parse root operation types schema",
			args: args{
				src: &ast.Source{
					Name: "parseRootOperationTypesSchema.graphql",
					Body: parseRootOperationTypesSchema,
				},
			},
			want: &ast.TypeSystemExtensionDocument{
				SchemaDefinitions: []ast.SchemaDefinition{
					{
						Directives: []ast.Directive{},
						RootOperationTypeDefinitions: []ast.RootOperationTypeDefinition{
							{
								OperationType: ast.OperationTypeQuery,
								Type:          "Query",
							},
							{
								OperationType: ast.OperationTypeMutation,
								Type:          "Mutation",
							},
							{
								OperationType: ast.OperationTypeSubscription,
								Type:          "Subscription",
							},
						},
					},
				},
				TypeDefinitions: map[string]ast.TypeDefinition{
					"Query": &ast.ObjectTypeDefinition{
						Name: "Query",
					},
					"Mutation": &ast.ObjectTypeDefinition{
						Name: "Mutation",
					},
					"Subscription": &ast.ObjectTypeDefinition{
						Name: "Subscription",
					},
				},
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
