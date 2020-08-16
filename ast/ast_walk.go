package ast

import (
	"fmt"
)

// A Visitor implements a Visit method, which is invoked for each Expression
// encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of Expression with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(expr Expression) (w Visitor)
}

// A VisitReplacer is an extension of visitor that allows
// the visitor to replace the node currently being visited
// by calling the supplied replacer function
type VisitReplacer interface {
	VisitReplace(expr Expression, replacer func(Expression)) VisitReplacer
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(expr); Expression must not be nil. If the visitor w returned by
// v.Visit(expr) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of Expression, followed by a call of
// w.Visit(nil).
// if the Visitor v is also a VisitReplacer, it will be "upgraded"
// and `VisitReplace(expr, ...)` will be called in place of .Visit(expr)
func Walk(v Visitor, expr Expression) {
	visitReplacer, upgrade := v.(VisitReplacer)
	if upgrade {
		walkReplacer(visitReplacer, expr, nil, 0)
	} else {
		walkVisitor(v, expr)
	}
}

func walkVisitor(v Visitor, expr Expression) {
	if v = v.Visit(expr); v == nil {
		return
	}

	switch expr := expr.(type) {
	case *ActionExpr:
		walkVisitor(v, expr.Expr)
	case *AndCodeExpr:
		// Nothing to do
	case *AndExpr:
		walkVisitor(v, expr.Expr)
	case *AnyMatcher:
		// Nothing to do
	case *CharClassMatcher:
		// Nothing to do
	case *ChoiceExpr:
		for _, e := range expr.Alternatives {
			walkVisitor(v, e)
		}
	case *Grammar:
		for _, e := range expr.Rules {
			walkVisitor(v, e)
		}
	case *LabeledExpr:
		walkVisitor(v, expr.Expr)
	case *LitMatcher:
		// Nothing to do
	case *NotCodeExpr:
		// Nothing to do
	case *NotExpr:
		walkVisitor(v, expr.Expr)
	case *OneOrMoreExpr:
		walkVisitor(v, expr.Expr)
	case *Rule:
		walkVisitor(v, expr.Expr)
	case *RuleRefExpr:
		// Nothing to do
	case *SeqExpr:
		for _, e := range expr.Exprs {
			walkVisitor(v, e)
		}
	case *StateCodeExpr:
		// Nothing to do
	case *ZeroOrMoreExpr:
		walkVisitor(v, expr.Expr)
	case *ZeroOrOneExpr:
		walkVisitor(v, expr.Expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
}

func walkReplacer(v VisitReplacer, expr, parent0 Expression, index int) {
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

	if v = v.VisitReplace(expr, replacer); v == nil {
		return
	}

	switch expr := expr.(type) {
	case *ActionExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	case *AndCodeExpr:
		// Nothing to do
	case *AndExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	case *AnyMatcher:
		// Nothing to do
	case *CharClassMatcher:
		// Nothing to do
	case *ChoiceExpr:
		for i, e := range expr.Alternatives {
			walkReplacer(v, e, expr, i)
		}
	case *Grammar:
		for i, e := range expr.Rules {
			walkReplacer(v, e, expr, i)
		}
	case *LabeledExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	case *LitMatcher:
		// Nothing to do
	case *NotCodeExpr:
		// Nothing to do
	case *NotExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	case *OneOrMoreExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	case *Rule:
		walkReplacer(v, expr.Expr, expr, 0)
	case *RuleRefExpr:
		// Nothing to do
	case *SeqExpr:
		for i, e := range expr.Exprs {
			walkReplacer(v, e, expr, i)
		}
	case *StateCodeExpr:
		// Nothing to do
	case *ZeroOrMoreExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	case *ZeroOrOneExpr:
		walkReplacer(v, expr.Expr, expr, 0)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}

}

type inspector func(Expression) bool

func (f inspector) Visit(expr Expression) Visitor {
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
