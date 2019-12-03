package findqueries

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/inspector"
)

type ExprResolver interface {
	ResolveFrom(pos token.Pos) ast.Expr
}

type exprResolver struct {
	pos2expr map[token.Pos]ast.Expr
}

func (er exprResolver) ResolveFrom(pos token.Pos) ast.Expr {
	return er.pos2expr[pos]
}

func NewExprResolver(inspctr *inspector.Inspector) ExprResolver {
	pos2expr := make(map[token.Pos]ast.Expr)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.GoStmt)(nil),
		(*ast.DeferStmt)(nil),
	}
	inspctr.Preorder(nodeFilter, func(n ast.Node) {
		switch expr := n.(type) {
		case *ast.CallExpr:
			pos2expr[expr.Lparen] = expr
		case *ast.GoStmt:
			pos2expr[expr.Go] = expr.Call
		case *ast.DeferStmt:
			pos2expr[expr.Defer] = expr.Call
		}
	})
	return exprResolver{pos2expr: pos2expr}
}
