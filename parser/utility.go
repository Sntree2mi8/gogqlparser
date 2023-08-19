package parser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
	"slices"
)

type LexerWrapper struct {
	lexer     *gogqllexer.Lexer
	keepToken *gogqllexer.Token
}

func NewLexerWrapper(lexer *gogqllexer.Lexer) *LexerWrapper {
	return &LexerWrapper{
		lexer: lexer,
	}
}

func (l *LexerWrapper) NextToken() gogqllexer.Token {
	if l.keepToken != nil {
		t := *l.keepToken
		l.keepToken = nil
		return t
	}

	return l.lexer.NextToken()
}

func (l *LexerWrapper) PeekToken() gogqllexer.Token {
	if l.keepToken == nil {
		t := l.lexer.NextToken()
		l.keepToken = &t
	}

	return *l.keepToken
}

func (l *LexerWrapper) PeekAndMayBe(kinds []gogqllexer.Kind, callback func(t gogqllexer.Token, advanceLexer func()) error) error {
	t := l.PeekToken()
	if slices.Contains(kinds, t.Kind) {
		return callback(t, func() { l.NextToken() })
	}
	return nil
}

func (l *LexerWrapper) PeekAndMustBe(kinds []gogqllexer.Kind, callback func(t gogqllexer.Token, advanceLexer func()) error) error {
	t := l.PeekToken()
	if slices.Contains(kinds, t.Kind) {
		return callback(t, func() { l.NextToken() })
	}
	return fmt.Errorf("unexpected token %v", t)
}

func (l *LexerWrapper) Skip(kind gogqllexer.Kind) error {
	t := l.NextToken()
	if t.Kind != kind {
		return fmt.Errorf("unexpected token %v", t)
	}
	return nil
}

func (l *LexerWrapper) SkipIf(kind gogqllexer.Kind) (skip bool) {
	defer func() {
		if skip {
			l.NextToken()
		}
	}()
	t := l.PeekToken()
	return t.Kind == kind
}

func (l *LexerWrapper) SkipKeyword(keyword string) error {
	t := l.NextToken()
	if t.Kind != gogqllexer.Name || t.Value != keyword {
		return fmt.Errorf("unexpected token %v", t)
	}
	return nil
}
