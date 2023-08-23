package parser

import (
	"github.com/Sntree2mi8/gogqllexer"
	"github.com/Sntree2mi8/gogqlparser/ast"
	"reflect"
	"strings"
	"testing"
)

func TestParseUnionTypeDefinition(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name    string
		schema  string
		args    args
		wantDef *ast.UnionTypeDefinition
		wantErr bool
	}{
		{
			name: "simple union",
			schema: `
union User = SuperUser | NormalUser
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.UnionTypeDefinition{
				Description: "this is description",
				Name:        "User",
				MemberTypes: []ast.Type{
					{
						NamedType: "SuperUser",
					},
					{
						NamedType: "NormalUser",
					},
				},
			},
		},
		{
			name: "with directives",
			schema: `
union User @deprecated = SuperUser | NormalUser
`,
			args: args{
				description: "this is description",
			},
			wantDef: &ast.UnionTypeDefinition{
				Description: "this is description",
				Name:        "User",
				Directives: []ast.Directive{
					{
						Name: "deprecated",
					},
				},
				MemberTypes: []ast.Type{
					{
						NamedType: "SuperUser",
					},
					{
						NamedType: "NormalUser",
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
			gotDef, err := p.ParseUnionTypeDefinition(tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUnionTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseUnionTypeDefinition() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}

func TestParseUnionTypeExtension(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		wantDef *ast.UnionTypeExtension
		wantErr bool
	}{
		{
			name: "simple union type extension",
			schema: `
union Restaurant = ItalianRestaurant
`,
			wantDef: &ast.UnionTypeExtension{
				Name: "Restaurant",
				MemberTypes: []ast.Type{
					{
						NamedType: "ItalianRestaurant",
					},
				},
			},
		},
		{
			name: "only directive",
			schema: `
union Restaurant @union_directive
`,
			wantDef: &ast.UnionTypeExtension{
				Name: "Restaurant",
				Directives: []ast.Directive{
					{
						Name: "union_directive",
					},
				},
			},
		},
		{
			name: "extend but do nothing",
			schema: `
union Restaurant
`,
			wantDef: &ast.UnionTypeExtension{
				Name: "Restaurant",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				lexer: gogqllexer.New(strings.NewReader(tt.schema)),
			}
			gotDef, err := p.ParseUnionTypeExtension()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUnionTypeExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDef, tt.wantDef) {
				t.Errorf("ParseUnionTypeExtension() gotDef = %v, want %v", gotDef, tt.wantDef)
			}
		})
	}
}
