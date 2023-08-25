package parser

import (
	"github.com/Sntree2mi8/gogqlparser/ast"
	"testing"
)

func Test_validateDirectiveDefinitions_Valid(t *testing.T) {
	type args struct {
		dds []ast.DirectiveDefinition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid directive",
			args: args{
				dds: []ast.DirectiveDefinition{
					{
						Name: "test",
						DirectiveLocations: []ast.DirectiveLocation{
							ast.DirectiveLocationField,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDirectiveDefinitions(tt.args.dds); (err != nil) != tt.wantErr {
				t.Errorf("validateDirectiveDefinitions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateDirectiveDefinitions_Invalid(t *testing.T) {
	type args struct {
		dds []ast.DirectiveDefinition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "The directive must not have a name which begins with the characters \"__\" (two underscores).",
			args: args{
				dds: []ast.DirectiveDefinition{
					{
						Name: "__test",
						DirectiveLocations: []ast.DirectiveLocation{
							ast.DirectiveLocationField,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDirectiveDefinitions(tt.args.dds); (err != nil) != tt.wantErr {
				t.Errorf("validateDirectiveDefinitions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
