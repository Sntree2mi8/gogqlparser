package ast

type DirectiveDefinition struct{}

type Directives []Directive

type Directive struct {
	Name      string
	Arguments Arguments
}

type Arguments []Argument

type Argument struct {
	Name  string
	Value string
}
