/*
	expression	::= binary | func | group | literal
	binary		::= expression operator expression
	func		::= builtin argument_count expression+
	group		::= '(' expression ')'
	literal		::= string literal
	operator	::= '+'
*/
package parse

import (
	"fmt"
	"strings"
	"strlang/lex"
)

type NodeKind int

const (
	EOF 		NodeKind = iota 
	ERROR 
	FUNC
	BINARY
	GROUP
	LITERAL
)

var nodeKindName = map[NodeKind]string {
	EOF:		"EOF",
	ERROR:		"ERROR",
	FUNC:		"FUNC",
	BINARY:		"BINARY",
	GROUP:		"GROUP",
	LITERAL:	"LITERAL",
}

type Parser struct {
	tokens []lex.Token
	pos	int
}

func NewParser(tokens []lex.Token) Parser {
	return Parser{tokens, 0}
}

type Node struct {
	Kind NodeKind
	tokens []lex.Token
}

func (n *Node) String() string {
	var tstrs []string
	for _, token := range n.tokens {
		tstrs = append(tstrs, token.String())
	}

	return fmt.Sprintf("%s: { %s }", nodeKindName[n.Kind], strings.Join(tstrs, ", "))
}

func (p *Parser) curr() lex.Token {
	len := len(p.tokens)
	if p.pos >= len {
		return p.tokens[len - 1]
	} else {
		t := p.tokens[p.pos]
		return t;
	}
}

func (p *Parser) next() lex.Token {
	t := p.curr()
	if (t.Kind != lex.EOF) {
		p.pos += 1
	}
	return t;
}

func (p *Parser) expect(kinds ...lex.TokenKind) bool {
	for _, kind := range kinds {
		if p.curr().Kind != kind {
			return false
		}
		
		p.next()
	}

	return true
}

func (p *Parser) Next() Node {
	start := p.pos
	switch (p.next().Kind) {
	case lex.EOF:
		return Node{EOF, p.tokens[start: p.pos]}

	case lex.LITERAL:
		return Node{LITERAL, p.tokens[start: p.pos]}

	case lex.OPAREN:
		for p.curr().Kind != lex.EOF && p.curr().Kind != lex.CPAREN {
			p.next()
		}

		var kind NodeKind
		if !p.expect(lex.CPAREN) {
			kind = ERROR
		} else {
			kind = GROUP
		}

		return Node{kind, p.tokens[start:p.pos]}

	case lex.IDENTIFIER:
		if !p.expect(lex.BLOCK, lex.OPAREN) {
			return Node{ERROR, p.tokens[start:p.pos]}
		}

		for p.curr().Kind != lex.CPAREN {
			p.next()
		}

		if !p.expect(lex.CPAREN) {
			return Node{ERROR, p.tokens[start:p.pos]}
		}

		return Node{FUNC, p.tokens[start:p.pos]}

	default:
		return Node{ERROR, p.tokens[start:p.pos]}
	}

}
