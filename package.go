package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func searchPackage(pkg *packages.Package, queryers []queryerInfo) ([]queryInfo, error) {
	queryersMap := make(map[string]int)
	for _, qi := range queryers {
		queryersMap[qi.FullName] = qi.QueryPos
	}

	end2fn := make(map[token.Pos]*types.Func)
	for ident, obj := range pkg.TypesInfo.Uses {
		fn, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		if _, ok := queryersMap[fn.FullName()]; ok {
			end2fn[ident.End()] = fn
		}
	}


	queries := make([]queryInfo, 0)
	unresolved := make([]queryInfo, 0)
	var err error
	for _, f := range pkg.Syntax {
		filescope, _ := pkg.TypesInfo.Scopes[f]
		scope2fn := make(map[*types.Scope]*types.Func)
		for _, obj := range pkg.TypesInfo.Defs {
			fn, ok := obj.(*types.Func)
			if !ok {
				continue
			}
			if fn.Scope().Parent() == filescope {
				scope2fn[fn.Scope()] = fn
			}
		}
		ast.Inspect(f, func(n ast.Node) bool {
			var qi queryInfo
			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			qi.Expr = types.ExprString(ce)
			qi.Pos = pkg.Fset.Position(ce.Pos()).String()
			scope := pkg.Types.Scope().Innermost(ce.Pos())
			if scope == nil {
				err = errors.New("unexpected call out of all scope: " + qi.Pos + " " + qi.Expr)
				return true
			}
			for scope.Parent() != nil {
				if fn, ok := scope2fn[scope]; ok {
					qi.Caller = fn.FullName()
					break
				}
				scope = scope.Parent()
			}
			if qi.Caller == "" {
				err = errors.New("unexpected call out of all top-level functions: " + qi.Pos + " " + qi.Expr)
				return true
			}
			end := astutil.Unparen(ce.Fun).End()
			fn, ok := end2fn[end]
			if !ok {
				return true
			}
			delete(end2fn, end)
			pos := queryersMap[fn.FullName()]
			queryExpr := ce.Args[pos]
			query := ""
			switch q := queryExpr.(type) {
			case *ast.BasicLit:
				query, err = strconv.Unquote(q.Value)
				if err != nil {
					err = errors.Wrap(err, "failed to unquote query")
					break
				}
			default:
				// TODO
			}
			qi.Query = query
			if qi.Query == "" || qi.Caller == "" {
				unresolved = append(unresolved, qi)
				return true
			}
			queries = append(queries, qi)
			return false
		})
	}
	if err != nil {
		return nil, err
	}

	if len(unresolved) > 0 {
		fmt.Fprintln(os.Stderr, "UNRESOLVED")
	}
	for _, qi := range unresolved {
		fmt.Fprintln(os.Stderr, qi.Pos, qi.Expr)
	}

	return queries, nil
}
