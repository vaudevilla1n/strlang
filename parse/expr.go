package parse

import (
	"fmt"
	"strings"
	"strlang/lex"
)

type ExprKind int

const (
	BINARY	ExprKind = iota
	FUNC
	GROUP
	LITERAL
)

var exprKindName = map[ExprKind]string {
	BINARY:		"BINARY",
	FUNC:		"FUNC",
	GROUP:		"GROUP",
	LITERAL:	"LITERAL",
}

type Expr interface {
	Kind() ExprKind
	Pos() (int, int)
	String() string
}

type binaryExpr struct {
	left Expr
	op lex.Token
	right Expr
}
func (b *binaryExpr) Kind() ExprKind {
	return BINARY
}
func (b *binaryExpr) Pos() (int, int) {
	pos, _ := b.left.Pos()
	_, end := b.right.Pos()
	return pos, end
}
func (b *binaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", b.left.String(), b.op.Text, b.right.String())
}

/*
	builtin is a terrible name for the function identifier
	but im a dickhead so im not changing it
*/
type funcExpr struct {
	builtin lex.Token
	block lex.Token
	oparen lex.Token
	args []Expr
	cparen lex.Token
}
func (f *funcExpr) Kind() ExprKind {
	return FUNC
}
func (f *funcExpr) Pos() (int, int) {
	return f.builtin.Pos, f.cparen.Pos
}
func (f *funcExpr) String() string {
	var args []string
	for _, arg := range f.args {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("%s::(%s)", f.builtin.Text, strings.Join(args, ", "))
}

type groupExpr struct {
	oparen lex.Token
	expr Expr
	cparen lex.Token
}
func (g *groupExpr) Kind() ExprKind {
	return GROUP
}
func (g *groupExpr) Pos() (int, int) {
	return g.oparen.Pos, g.cparen.Pos
}
func (g *groupExpr) String() string {
	return fmt.Sprintf("(%s)", g.expr.String())
}

type literalExpr struct {
	str lex.Token
}
func (l *literalExpr) Kind() ExprKind {
	return LITERAL
}
func (l *literalExpr) Pos() (int, int) {
	return l.str.Pos, l.str.Pos
}
func (l *literalExpr) String() string {
	return l.str.Text
}
