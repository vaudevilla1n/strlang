/*
	expression	::= binary | primary
	binary		::= expression (operator expression)*
	primary		::= func | literal
	func		::= builtin expression+
	literal		::= func | group | literal
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

type LiteralKind int

const (
	GROUP	LiteralKind = iota 
	LITERAL
)

var literalKindName = map[LiteralKind]string {
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

type Literal struct {
	Kind LiteralKind
	tokens []lex.Token
}

func (p *Literal) String() string {
	var tstrs []string
	for _, token := range p.tokens {
		tstrs = append(tstrs, token.String())
	}

	return fmt.Sprintf("%s: { %s }", literalKindName[p.Kind], strings.Join(tstrs, ", "))
}

type Function struct {
	builtin Literal
	arguments []Literal
}

func (p *parser) atEOF() bool {
	if p.pos >= len(p.tokens) {
		return true
	] else {
		return false
	}
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

func (p *Parser) advance() lex.Token {
	t := p.curr()
	if (t.Kind != lex.EOF) {
		p.pos += 1
	}
	return t;
}

func (p *Parser) check(kinds ...lex.TokenKind) bool {
	for _, kind := range kinds {
		if p.curr().Kind == kind {
			return true
		}
	}

	return false
}

func (p *Parser) expect(kinds ...lex.TokenKind) bool {
	for _, kind := range kinds {
		if !p.check(kind) {
			return false
		}
		
		p.advance()
	}

	return true
}

func (p *Parser) literal() (Literal, string) {
	start := p.pos

	switch (p.advance().Kind) {
	case lex.LITERAL:
		return Literal{LITERAL, p.tokens[start: p.pos]}, nil

	case lex.OPAREN:
		for p.curr().Kind != lex.EOF && p.curr().Kind != lex.CPAREN {
			p.advance()
		}

		var kind LiteralKind
		if !p.expect(lex.CPAREN) {
			kind = ERROR
		} else {
			kind = GROUP
		}

		return Literal{kind, p.tokens[start:p.pos]}, nil

	case lex.IDENTIFIER:
		if !p.expect(lex.BLOCK, lex.OPAREN) {
			return Literal{ERROR, p.tokens[start:p.pos]}
		}

		for p.curr().Kind != lex.CPAREN {
			p.advance()
		}

		if !p.expect(lex.CPAREN) {
			return Literal{ERROR, p.tokens[start:p.pos]}
		}

		return Literal{FUNC, p.tokens[start:p.pos]}, nil

	case lex.EOF:
		return nil, "expected token(s)"
	default:
		return nil, "unexpected token"
	}

}

func (p *Parser) function() Function, string {
	if (!p.expect(lex.IDENTIFIER))
		return nil, "expected function name"
	if (!p.expect(lex.BLOCK, lex.OPAREN))
		return nil, "expected '::(' after function name"

	var args []Literal
	
	for !p.atEOF() && !p.check(lex.EOF) {
		literal, err := p.literal()
		if err != nil {
			return nil, err
		}

		args = append(args, literal)
	}

	if (!p.expect(lex.CPAREN))
		return nil, "expected closing parenthese"

	return Function{builtin, args}
}

type Primary interface {
	Literal | Function
}

type Binary struct {
	left Primary
	op lex.Token
	right Primary
}

func (p *Parser) primary() (Primary, string) {
	switch p.curr().Kind {
	case lex.IDENTIFER:
		return p.function()
	default:
		return p.literal()
	}
}

func (p *Parser) binary() Binary {
	var b Binary

	b.left, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.check(lex.PLUS) {
		op = p.advance()

		right, err := p.primary()
		if err != nil {
			return nil, err
		}

		b = Binary{b.left, op, right}
	}
// love you
}
