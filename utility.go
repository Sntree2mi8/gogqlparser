package gogqlparser

import (
	"fmt"
	"github.com/Sntree2mi8/gogqllexer"
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

type mustBeCallback func(t gogqllexer.Token) (ok bool)

func mustBe(l *lexerWrapper, callbacks ...mustBeCallback) error {
	for _, callback := range callbacks {
		t := l.NextToken()
		if !callback(t) {
			return fmt.Errorf("unexpected token %v", t)
		}
	}
	return nil
}

type maybeCallback func(t gogqllexer.Token) (ok bool)

func maybe(l *lexerWrapper, callbacks ...maybeCallback) error {
	for _, callback := range callbacks {
		t := l.PeekToken()
		if callback(t) {
			l.NextToken()
		}
	}
	return nil
}
