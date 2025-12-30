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
package main

import (
	"os"
	"fmt"
	"bufio"
	"strlang/lex"
	"strlang/parse"
)

func run(line string) {
	lexer := lex.NewLexer(line);
	
	var t lex.Token
	var tokens []lex.Token
	for {
		t = lexer.Next()

		if (t.Kind == lex.ILLEGAL) {
			break
		}

		tokens = append(tokens, t)

		if t.Kind == lex.EOF {
			break
		}
	}

	if t.Kind == lex.ILLEGAL {
		fmt.Fprintf(os.Stderr, "lexer error: pos %d: %s\n", t.Pos, t.Text)
		return
	} 

	parser := parse.NewParser(tokens)
	if expr, err := parser.Parse(); err == nil {
		if val, err := expr.Eval(); err == nil {
			fmt.Printf("\"%s\"\n", val)
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")

		if (!input.Scan()) {
			fmt.Print("\n")
			break
		}

		line := input.Text()
		
		if len(line) > 0 {
			run(line)
		}
	}
}
