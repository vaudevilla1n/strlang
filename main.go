package main

/*
	grammar

	expression	::= term
	term		::= primary '+' primary
	primary		::= func | string | group
	func		::= builtin '::(' expression (',' expression)* ')'
	string		::= string literal
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
	strlang/lexer"
	"strlang/quote"
)

func main() {
	quote.Quote();
	lexer.LaLa();
}
