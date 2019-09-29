package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func searchPackage(pkg *packages.Package, path string, fset *token.FileSet, queryers []queryerInfo) ([]queryInfo, error) {
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
	unresolved := make([]*ast.CallExpr, 0)
	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			ce, ok := n.(*ast.CallExpr)
			if !ok {
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
				query = q.Value
			default:
				// TODO
				unresolved = append(unresolved, ce)
				return true
			}
			queries = append(queries, queryInfo{query})
			return false
		})
	}

	fmt.Println("UNRESOLVED")
	for end, fn := range end2fn {
		fmt.Println(fset.Position(end), fn)
	}
	for _, ce := range unresolved {
		fmt.Println(fset.Position(ce.Pos()), types.ExprString(ce))
	}

	return queries, nil
}
