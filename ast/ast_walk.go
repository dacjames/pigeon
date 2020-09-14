package ast

import (
	"fmt"
)

// Backref holds a reference to the parent of the current expression being visited
// and providers a replacer convenience function for a visitee to replace itself
// in the parent expression. This allows visitee to replace themselves with an
// optimized expression.
type Backref struct {
	parent   Expression
	replacer func(Expression)
}

// A Visitor implements a Visit method, which is invoked for each Expression
// encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of Expression with the visitor w, followed by a call of w.Visit(nil).
// Passes a Backref on each visit, proving a reference back to the parent
// expression so visitees can change their identity while being visited.
type Visitor interface {
	Visit(expr Expression, br Backref) (w Visitor)
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(expr); Expression must not be nil. If the visitor w returned by
// v.Visit(expr) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of Expression, followed by a call of
// w.Visit(nil).
func Walk(v Visitor, expr Expression) {
	walk0(v, expr, nil, 0)
}

func walk0(v Visitor, expr, parent0 Expression, index int) {
	var replacer func(Expression)

	switch parent := parent0.(type) {
	case nil:
		replacer = func(expr Expression) {}
	case *ActionExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	case *AndExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	case *ChoiceExpr:
		replacer = func(expr Expression) {
			parent.Alternatives[index] = expr
		}
	case *Grammar:
		replacer = func(expr Expression) {
			parent.Rules[index] = expr.(*Rule)
		}
	case *LabeledExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	case *NotExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	case *OneOrMoreExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	case *Rule:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	case *SeqExpr:
		replacer = func(expr Expression) {
			parent.Exprs[index] = expr
		}
	case *ZeroOrMoreExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}

	case *ZeroOrOneExpr:
		replacer = func(expr Expression) {
			parent.Expr = expr
		}
	}

	if v = v.Visit(expr, Backref{
		parent:   parent0,
		replacer: replacer,
	}); v == nil {
		return
	}

	switch expr := expr.(type) {
	case *ActionExpr:
		walk0(v, expr.Expr, expr, 0)
	case *AndCodeExpr:
		// Nothing to do
	case *AndExpr:
		walk0(v, expr.Expr, expr, 0)
	case *AnyMatcher:
		// Nothing to do
	case *CharClassMatcher:
		// Nothing to do
	case *ChoiceExpr:
		for i, e := range expr.Alternatives {
			walk0(v, e, expr, i)
		}
	case *Grammar:
		for i, e := range expr.Rules {
			walk0(v, e, expr, i)
		}
	case *LabeledExpr:
		walk0(v, expr.Expr, expr, 0)
	case *LitMatcher:
		// Nothing to do
	case *NotCodeExpr:
		// Nothing to do
	case *NotExpr:
		walk0(v, expr.Expr, expr, 0)
	case *OneOrMoreExpr:
		walk0(v, expr.Expr, expr, 0)
	case *Rule:
		walk0(v, expr.Expr, expr, 0)
	case *RuleRefExpr:
		// Nothing to do
	case *SeqExpr:
		for i, e := range expr.Exprs {
			walk0(v, e, expr, i)
		}
	case *StateCodeExpr:
		// Nothing to do
	case *ZeroOrMoreExpr:
		walk0(v, expr.Expr, expr, 0)
	case *ZeroOrOneExpr:
		walk0(v, expr.Expr, expr, 0)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}

}

type inspector func(Expression) bool

func (f inspector) Visit(expr Expression, br Backref) Visitor {
	if f(expr) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(expr); expr must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of expr, followed by a
// call of f(nil).
func Inspect(expr Expression, f func(Expression) bool) {
	Walk(inspector(f), expr)
}
