package ast

type ValueKind int

const (
	ValueKindVariable ValueKind = iota
	ValueKindInt
	ValueKindFloat
	ValueKindString
	ValueKindBoolean
	ValueKindNull
	ValueKindEnum
	ValueKindList
	ValueKindObject
)

type Value interface {
	ValueKind() ValueKind
}

type IntValue struct {
	Value int64
}

func (i IntValue) ValueKind() ValueKind {
	return ValueKindInt
}

type FloatValue struct {
	Value float64
}

func (f FloatValue) ValueKind() ValueKind {
	return ValueKindFloat
}

type StringValue struct {
	Value string
}

func (s StringValue) ValueKind() ValueKind {
	return ValueKindString
}

type BooleanValue struct {
	Value bool
}

func (b BooleanValue) ValueKind() ValueKind {
	return ValueKindBoolean
}

type NullValue struct{}

func (n NullValue) ValueKind() ValueKind {
	return ValueKindNull
}

type EnumValue struct {
	Value string
}

func (e EnumValue) ValueKind() ValueKind {
	return ValueKindEnum
}

type ListValue struct {
	Values []Value
}

func (l ListValue) ValueKind() ValueKind {
	return ValueKindList
}

type ObjectValue struct {
	Fields []ObjectField
}

func (o ObjectValue) ValueKind() ValueKind {
	return ValueKindObject
}

type ObjectField struct {
	Name  string
	Value Value
}
