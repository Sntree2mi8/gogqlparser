package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseInterfaceTypeDefinition(t *testing.T) {
	type args struct {
		l           *LexerWrapper
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantD   *ast.InterfaceTypeDefinition
		wantErr bool
	}{
		{
			name: "implements interface",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
interface User implements Node {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.InterfaceTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Interfaces:  []string{"Node"},
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
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
			name: "with directive",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
interface User @deprecated(reason: "this is deprecated") {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.InterfaceTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Directives: []ast.Directive{
					{
						Name: "deprecated",
						Arguments: []ast.Argument{
							{
								Name:  "reason",
								Value: "\"this is deprecated\"",
							},
						},
					},
				},
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
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
			name: "with directive, implements interface",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(
						strings.NewReader(`
interface User implements Node @deprecated(reason: "this is deprecated") {
	id: ID!
	name: String!
}
`,
						),
					),
				),
				description: "this is description",
			},
			wantD: &ast.InterfaceTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Interfaces:  []string{"Node"},
				Directives: []ast.Directive{
					{
						Name: "deprecated",
						Arguments: []ast.Argument{
							{
								Name:  "reason",
								Value: "\"this is deprecated\"",
							},
						},
					},
				},
				FieldDefinitions: []*ast.FieldDefinition{
					{
						Name: "id",
						Type: ast.Type{
							NamedType: "ID",
							NotNull:   true,
						},
					},
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
			gotD, err := ParseInterfaceTypeDefinition(tt.args.l, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterfaceTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("ParseInterfaceTypeDefinition() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}
