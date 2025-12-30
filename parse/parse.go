/*
	grammar:

	expression	::= primary (operator primary)*
	primary		::= func | literal | group
	func		::= builtin '::(' (expression ',')* expression ')'
	group		::= '(' expression ')'
	literal		::= ''' .* ''' | '"' .* '"'
	operator	::= '+'

	ast:

	expression	::= binary | func | group | literal
	binary		::= expression operator expression
	func		::= builtin '::(' (expression ',')* expression ')'
	group		::= '(' expression ')'
	literal		::= ''' .* ''' | '"' .* '"'
	operator	::= '+'
*/
package parse

import (
	"fmt"
	"strlang/lex"
)

type Parser struct {
	tokens []lex.Token
	pos	int
}

func NewParser(tokens []lex.Token) Parser {
	return Parser{tokens, 0}
}

type ParseError struct {
	token lex.Token
	msg string
}

func (e *ParseError) Error() string {
	if e.token.Kind == lex.EOF {
		return fmt.Sprintf("parser error: pos %d: %s", e.token.Pos, e.msg)
	} else {
		return fmt.Sprintf("parser error: pos %d: '%s': %s", e.token.Pos, e.token.Text, e.msg)
	}
}

func (p *Parser) atEOF() bool {
	if p.pos >= len(p.tokens) {
		return true
	} else {
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

func (p *Parser) err(msg string) *ParseError {
	return &ParseError{p.curr(), msg}
}

func (p *Parser) expect(msg string, kinds ...lex.TokenKind) *ParseError {
	start := p.pos
	for _, kind := range kinds {
		if !p.check(kind) {
			p.pos = start
			return p.err(msg)
		}
		
		p.advance()
	}

	return nil
}

func (p *Parser) primary() (Expr, error) {
	t := p.advance()
	switch (t.Kind) {
	case lex.LITERAL:
		return &literalExpr{t}, nil

	case lex.OPAREN:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		end := p.curr()

		if err := p.expect("expected closing parentheses", lex.CPAREN); err != nil {
			return nil, err
		}

		return &groupExpr{t, expr, end}, nil

	case lex.IDENTIFIER:
		if err := p.expect("expected '::(' after function identifier", lex.BLOCK, lex.OPAREN); err != nil {
			return nil, err
		}

		block := p.tokens[p.pos - 2]
		oparen := p.tokens[p.pos - 1]

		var args []Expr
		for !p.atEOF() {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			args = append(args, expr)

			if p.check(lex.COMMA) {
				p.advance()
				continue
			}

			if p.check(lex.CPAREN) {
				break
			}

			return nil, p.err("no comma or closing parentheses after function argument(s)")
		}

		cparen := p.curr()
		if err := p.expect("no comma or closing parentheses after function argument(s)", lex.CPAREN); err != nil {
			return nil, err
		}

		return &funcExpr{t, block, oparen, args, cparen}, nil

	case lex.EOF:
		return nil, p.err("I expected more from you")
	default:
		return nil, p.err("unexpected token")
	}

}

func (p *Parser) expression() (Expr, error) {
// love you
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.check(lex.PLUS) {
		op := p.advance()
		right, err := p.primary()
		if err != nil {
			return nil, err
		}

		expr = &binaryExpr{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) Parse() (Expr, error) {
	if expr, err := p.expression(); err != nil {
		return nil, err
	} else if !p.atEOF() {
		return nil, p.err("erroneous token")
	} else {
		return expr, nil
	}
}
