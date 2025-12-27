package main

/*
	grammar

	expression	::= term
	term		::= primary '+' primary
	primary		::= func | string | group
	func		::= builtin '::(' expression (',' expression)* ')'
	string		::= text
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

func run(line []byte) {
	lex := lexer.NewLexer(line);
	for lex.Next() {
		t := lex.Token()
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
