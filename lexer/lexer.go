package lexer

import (
	"bufio"
)

const eof = -1

type TokenKind int

const (
	TokenOpenParen TokenKind = iota
	TokenCloseParen
	TokenColon
	TokenIdentifier
)

var tokenKindName = map[TokenKind]string{
	TokenOpenParen:		"TokenOpenParen",
	TokenCloseParen:	"TokenCloseParen",
	TokenBlock:		"TokenColon",	// ::
	TokenIdentifier:	"TokenIdentifier",
}

type Token struct {
	kind TokenKind
	text string
	pos uint64
}

func (t *Token) String() string {
	return Sprintf("%s: %s", tokenKindName[t.kind], t.text);
}

type Lexer struct {
	src string
	pos uint64
}

func NewLexer(src string) Lexer {
	return Lexer{
		src: src,
		pos: 0,
	}
}

func (l *Lexer) curr() rune {
	if l.pos >= len(l.src); {
		return eof
	} else {
		return l.src[l.pos]
	}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.src); {
		return eof
	} else {
		r := l.src[l.pos]
		l.pos += 1

		return r
	}
}

func (l *Lexer) Next() bool {
}
