/*
	grammar

	expression	::= term
	term		::= primary '+' primary
	primary		::= func | string | group
	func		::= builtin '::(' expression (',' expression)* ')'
	string		::= '"' text '"' | ''' text '''
	group		::= '(' expression ')'
	builtin		::= any builtin function

	ast

	expression	::= binary | primary
	binary		::= expression (operator expression)*
	primary		::= func | group | literal
	func		::= builtin expression+
	group		::= '(' expression ')'
	literal		::= string literal
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

	for _, token := range tokens {
		fmt.Println(token.String())
	}

	parser := parse.NewParser(tokens)
	for n := parser.Next(); n.Kind != parse.EOF && n.Kind != parse.ERROR; n = parser.Next() {
		fmt.Println(n.String())
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
