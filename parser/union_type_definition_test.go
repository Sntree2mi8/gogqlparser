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
		l           *LexerWrapper
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.UnionTypeDefinition
		wantErr bool
	}{
		{
			name: "simple union",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(strings.NewReader(`
union User = SuperUser | NormalUser
`),
					),
				),
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(strings.NewReader(`
union User @deprecated = SuperUser | NormalUser
`),
					),
				),
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
			gotDef, err := ParseUnionTypeDefinition(tt.args.l, tt.args.description)
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
	type args struct {
		l *LexerWrapper
	}
	tests := []struct {
		name    string
		args    args
		wantDef *ast.UnionTypeExtension
		wantErr bool
	}{
		{
			name: "simple union type extension",
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(strings.NewReader(`
union Restaurant = ItalianRestaurant
`),
					),
				),
			},
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(strings.NewReader(`
union Restaurant @union_directive
`),
					),
				),
			},
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
			args: args{
				l: NewLexerWrapper(
					gogqllexer.New(strings.NewReader(`
union Restaurant
`),
					),
				),
			},
			wantDef: &ast.UnionTypeExtension{
				Name: "Restaurant",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDef, err := ParseUnionTypeExtension(tt.args.l)
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
