package expression

import (
	"fmt"
	"strconv"
	"strings"
	"strlang/lex"
)

type ExprKind int

const (
	BINARY	ExprKind = iota
	FUNC
	GROUP
	STRING
	NUMBER
)

var exprKindName = map[ExprKind]string {
	BINARY:		"BINARY",
	FUNC:		"FUNC",
	GROUP:		"GROUP",
	STRING:		"STRING",
	NUMBER:		"NUMBER",
}

type Expr interface {
	Kind() ExprKind
	Pos() (int, int)
	String() string
	Eval() (string, error)
}

type EvalError struct {
	msg string
}
func (e *EvalError) Error() string {
	return fmt.Sprintf("evaluation error: %s", e.msg)
}

type BinaryExpr struct {
	Left Expr
	Op lex.Token
	Right Expr
}
func (b *BinaryExpr) Kind() ExprKind {
	return BINARY
}
func (b *BinaryExpr) Pos() (int, int) {
	pos, _ := b.Left.Pos()
	_, end := b.Right.Pos()
	return pos, end
}
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Op.Text, b.Right.String())
}
func (b *BinaryExpr) Eval() (string, error) {
	left, err := b.Left.Eval()
	if err != nil {
		return "", err
	}

	right, err := b.Right.Eval()
	if err != nil {
		return "", err
	}

	switch b.Op.Kind {
	case lex.PLUS:
		return left + right, nil
	default:
		panic("unreachable")
	}
}

type builtinFunc func(...string)(string, error)

func expect_argument_count(func_name string, expected int, args []string) error {
	if len(args) != expected {
		var plural string
		if expected > 1 {
			plural = "s"
		} else {
			plural = ""
		}

		err := fmt.Sprintf("%s: expected %d argument%s, got %d", func_name, expected, plural, len(args))
		return &EvalError{err}
	} else {
		return nil
	}
}

func __reverse(args ...string) (string, error) {
	if err := expect_argument_count("reverse", 1, args); err != nil {
		return "", err
	}

	s := []rune(args[0])
	for i, j := 0, len(s) - 1; i < j; i, j = i + 1, j - 1 {
		s[i], s[j] = s[j], s[i]
	}

	return string(s), nil
}

func __substr(args ...string) (string, error) {
	switch len(args) {
	case 2:
		s := []rune(args[0])

		start, err := strconv.Atoi(args[1])
		if err != nil {
			return "", &EvalError{"substr: second argument is not an integer"}
		}

		len := len(s)
		if start < 0 {
			start += len
		}

		if 0 > start || start >= len {
			return "", &EvalError{"substr: invalid range"}
		}

		return string(s[start:]), nil
	case 3:
		s := []rune(args[0])

		start, err := strconv.Atoi(args[1])
		if err != nil {
			return "", &EvalError{"substr: second argument is not an integer"}
		}

		end, err := strconv.Atoi(args[2])
		if err != nil {
			return "", &EvalError{"substr: third argument is not an integer"}
		}

		len := len(s)
		if start < 0 {
			start += len
		}
		if end < 0 {
			end += len
		}

		if 0 > start || start >= len || 0 > end || end > len || start >= end {
			return "", &EvalError{"substr: invalid range"}
		}
	}

	err := fmt.Sprintf("substr: expected either 2 or 3 arguments, got %d", len(args))
	return "", &EvalError{err}
}

func __repeat(args ...string) (string, error) {
	if err := expect_argument_count("repeat", 2, args); err != nil {
		return "", err
	}

	count, err := strconv.Atoi(args[1])
	if err != nil || count < 0 {
		return "", &EvalError{"repeat: second argument is not a valid integer"}
	}

	var s string
	for i := 0; i < count; i++ {
		s += args[0]
	}

	return s, nil
}

var builtinFuncMap = map[string]builtinFunc {
	"reverse":	__reverse,
	"substr":	__substr,
	"repeat":	__repeat,
}

/*
	builtin is a terrible name for the function identifier
	but im a dickhead so im not changing it
*/
type FuncExpr struct {
	Builtin lex.Token
	Block lex.Token
	Oparen lex.Token
	Args []Expr
	Cparen lex.Token
}
func (f *FuncExpr) Kind() ExprKind {
	return FUNC
}
func (f *FuncExpr) Pos() (int, int) {
	return f.Builtin.Pos, f.Cparen.Pos
}
func (f *FuncExpr) String() string {
	var args []string
	for _, arg := range f.Args {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("%s::(%s)", f.Builtin.Text, strings.Join(args, ", "))
}
func (f *FuncExpr) Eval() (string, error) {
	builtin_function := builtinFuncMap[f.Builtin.Text]
	if builtin_function == nil {
		return "", &EvalError{"\"" + f.Builtin.Text + "\"" + ": no such function"}
	}
	
	var args []string
	for _, arg := range f.Args {
		if s, err := arg.Eval(); err == nil {
			args = append(args, s)
		} else {
			return "", err
		}
	}

	if ret, err := builtin_function(args...); err == nil {
		return ret, nil
	} else {
		return "", err
	}
}

type GroupExpr struct {
	Oparen lex.Token
	Expression Expr
	Cparen lex.Token
}
func (g *GroupExpr) Kind() ExprKind {
	return GROUP
}
func (g *GroupExpr) Pos() (int, int) {
	return g.Oparen.Pos, g.Cparen.Pos
}
func (g *GroupExpr) String() string {
	return fmt.Sprintf("(%s)", g.Expression.String())
}
func (g *GroupExpr) Eval() (string, error) {
	return g.Expression.Eval()
}

type StringExpr struct {
	Str lex.Token
}
func (s *StringExpr) Kind() ExprKind {
	return STRING
}
func (s *StringExpr) Pos() (int, int) {
	return s.Str.Pos, s.Str.Pos
}
func (s *StringExpr) String() string {
	return s.Str.Text
}
func (s *StringExpr) Eval() (string, error) {
	return s.Str.Text[1:len(s.Str.Text) - 1], nil
}

type NumberExpr struct {
	Num lex.Token
}
func (n *NumberExpr) Kind() ExprKind {
	return NUMBER
}
func (n *NumberExpr) Pos() (int, int) {
	return n.Num.Pos, n.Num.Pos
}
func (n *NumberExpr) String() string {
	return n.Num.Text
}
func (n *NumberExpr) Eval() (string, error) {
	return n.Num.Text, nil
}
