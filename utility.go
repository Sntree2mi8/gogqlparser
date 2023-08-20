package gogqlparser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"slices"
)

type lexerWrapper struct {
	lexer     *gogqllexer.Lexer
	keepToken *gogqllexer.Token
}

func (l *lexerWrapper) NextToken() gogqllexer.Token {
	if l.keepToken != nil {
		t := *l.keepToken
		l.keepToken = nil
		return t
	}

	return l.lexer.NextToken()
}

func (l *lexerWrapper) PeekToken() gogqllexer.Token {
	if l.keepToken == nil {
		t := l.lexer.NextToken()
		l.keepToken = &t
	}

	return *l.keepToken
}

func (l *lexerWrapper) PeekAndMustBe(kinds []gogqllexer.Kind, callback func(t gogqllexer.Token, advanceLexer func()) error) error {
	t := l.PeekToken()
	if slices.Contains(kinds, t.Kind) {
		return callback(t, func() { l.NextToken() })
	}
	return fmt.Errorf("unexpected token %v", t)
}

func (l *lexerWrapper) Skip(kind gogqllexer.Kind) error {
	t := l.NextToken()
	if t.Kind != kind {
		return fmt.Errorf("unexpected token %v", t)
	}
	return nil
}

func (l *lexerWrapper) SkipIf(kind gogqllexer.Kind) (skip bool) {
	defer func() {
		if skip {
			l.NextToken()
		}
	}()
	t := l.PeekToken()
	return t.Kind == kind
}

func (l *lexerWrapper) SkipKeyword(keyword string) error {
	t := l.NextToken()
	if t.Kind != gogqllexer.Name || t.Value != keyword {
		return fmt.Errorf("unexpected token %v", t)
	}
	return nil
}
