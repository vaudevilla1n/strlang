package lex

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type TokenKind int

const (
	EOF	TokenKind = iota
	ILLEGAL
	PLUS
	COMMA
	BLOCK	// ::
	LITERAL
	IDENTIFIER
	OPAREN
	CPAREN
)

var tokenKindName = map[TokenKind]string{
	EOF:		"EOF",
	ILLEGAL:	"ILLEGAL",
	PLUS:		"PLUS",
	COMMA:		"COMMA",
	BLOCK:		"BLOCK",	// ::
	LITERAL:	"LITERAL",
	IDENTIFIER:	"IDENTIFIER",
	OPAREN:		"OPAREN",
	CPAREN:		"CPAREN",
}

type Token struct {
	Kind TokenKind
	Text string
	Pos int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s: %s", tokenKindName[t.Kind], t.Text)
}

type Lexer struct {
	src string
	pos int
}

func NewLexer(src string) Lexer {
	return Lexer{
		src: src,
		pos: 0,
	}
}

func (l *Lexer) curr() rune {
	if l.pos >= len(l.src) {
		return eof
	} else {
		r, _ := utf8.DecodeRuneInString(l.src[l.pos:])
		return r
	}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.src) {
		return eof
	} else {
		r, w := utf8.DecodeRuneInString(l.src[l.pos:])
		l.pos += w
	
		return r
	}
}

func (l *Lexer) identifier() Token {
	start := l.pos

	if !unicode.IsLetter(l.curr()) && l.curr() != '_' {
		err := fmt.Sprintf("erroneous character '%c'", l.src[start])
		return Token{ILLEGAL, err, start}
	}

	for unicode.IsLetter(l.curr()) || unicode.IsNumber(l.curr()) || l.curr() == '_' {
		l.next()
	}

	return Token{IDENTIFIER, l.src[start:l.pos], start}
}

func (l *Lexer) literal() Token {
	start := l.pos
	// consume opening quote
	quote := l.next()

	for l.curr() != quote && l.curr() != eof {
		l.next()
	}

	if l.curr() == eof {
		return Token{ILLEGAL, "missing closing quotes", start}
	}

	l.next()

	return Token{LITERAL, l.src[start:l.pos], start}
}

func (l *Lexer) Next() (Token) {
	for unicode.IsSpace(l.curr()) {
		l.next()
	}

	if l.curr() == eof {
		return Token{EOF, "EOF", l.pos}
	}

	start := l.pos
	switch (l.next()) {
	case '(':
		return Token{OPAREN, l.src[start:l.pos], start}
	case ')':
		return Token{CPAREN, l.src[start:l.pos], start}
	case ',':
		return Token{COMMA, l.src[start:l.pos], start}
	case '+':
		return Token{PLUS, l.src[start:l.pos], start}
	case ':':
		if l.next() == ':' {
			return Token{BLOCK, l.src[start:l.pos], start}
		} 
		return Token{ILLEGAL, "erroneous ':'", start}
	case '"':
		l.pos = start
		return l.literal()
	case '\'':
		l.pos = start
		return l.literal()
	default:
		l.pos = start
		return l.identifier()
	}
}
