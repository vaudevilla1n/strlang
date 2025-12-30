/*
	grammar:

	expression	::= primary (operator primary)*

	primary		::= func | string | group
	group		::= '(' expression ')'
	string		::= ''' .* ''' | '"' .* '"'

	func		::= builtin '::(' args ')'
	args		::= expression ( ',' expression | number )*
	number		::= [0-9]+

	operator	::= '+'

	ast:

	expression	::= binary | func | group | string

	binary		::= expression operator expression

	func		::= builtin '::(' args ')'
	args		::= expression ( ',' expression | number )*
	number		::= [0-9]

	group		::= '(' expression ')'

	string		::= ''' .* ''' | '"' .* '"'

	operator	::= '+'
*/
package parse

import (
	"fmt"
	"strlang/lex"
	e "strlang/expression"
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
		return fmt.Sprintf("parser error: pos %d: %s", e.token.Pos + 1, e.msg)
	} else {
		return fmt.Sprintf("parser error: pos %d: '%s': %s", e.token.Pos + 1, e.token.Text, e.msg)
	}
}

func (p *Parser) atEOF() bool {
	if p.pos + 1 >= len(p.tokens) {
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

func (p *Parser) primary() (e.Expr, error) {
	t := p.advance()
	switch t.Kind {
	case lex.STRING:
		return &e.StringExpr{t}, nil

	case lex.OPAREN:
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		end := p.curr()

		if err := p.expect("expected closing parentheses", lex.CPAREN); err != nil {
			return nil, err
		}

		return &e.GroupExpr{t, expr, end}, nil

	case lex.IDENTIFIER:
		if err := p.expect("expected '::(' after function identifier", lex.BLOCK, lex.OPAREN); err != nil {
			return nil, err
		}

		block := p.tokens[p.pos - 2]
		oparen := p.tokens[p.pos - 1]

		var args []e.Expr

		if arg, err := p.expression(); err == nil {
			args = append(args, arg)
		} else {
			return nil, err
		}

		for !p.atEOF() && !p.check(lex.CPAREN) {
			if !p.check(lex.COMMA) {
				return nil, p.err("expected comma after function argument")
			}
			p.advance()

			if p.check(lex.NUMBER) {
				n := p.advance()
				args = append(args, &e.NumberExpr{n})
			} else if arg, err := p.expression(); err == nil {
				args = append(args, arg)
			} else {
				return nil, err
			}
		}

		cparen := p.curr()
		if err := p.expect("expected closing parentheses after function argument(s)", lex.CPAREN); err != nil {
			return nil, err
		}

		return &e.FuncExpr{t, block, oparen, args, cparen}, nil

	case lex.EOF:
		return nil, p.err("I expected more from you")
	default:
		p.pos -= 1
		return nil, p.err("unexpected token")
	}

}

func (p *Parser) expression() (e.Expr, error) {
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

		expr = &e.BinaryExpr{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) Parse() (e.Expr, error) {
	if expr, err := p.expression(); err != nil {
		return nil, err
	} else if !p.atEOF() {
		return nil, p.err("erroneous token")
	} else {
		return expr, nil
	}
}
