package main

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

	expression	::= binary | func | group | literal
	binary		::= expression operator expression
	func		::= builtin argument_count expression+
	group		::= '(' expression ')'
	literal		::= string literal
	operator	::= '+'
*/

import (
	"os"
	"fmt"
	"bufio"
	"strlang/lex"
)

func run(line string) {
	lexer := lex.NewLexer(line);
	
	var t lex.Token
	for t = lexer.Next(); t.Kind != lex.EOF && t.Kind != lex.ILLEGAL; t = lexer.Next() {
		fmt.Println(t.String())
	}

	if t.Kind == lex.ILLEGAL {
		fmt.Fprintf(os.Stderr, "lexer error: pos %d: %s\n", t.Pos, t.Text)
	} else {
		fmt.Println(t.String())
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
