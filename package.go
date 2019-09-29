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
	unresolved := make([]*ast.CallExpr, 0)
	var err error
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
				query, err = strconv.Unquote(q.Value)
				if err != nil {
					err = errors.Wrap(err, "failed to unquote query")
					break
				}
			default:
				// TODO
				unresolved = append(unresolved, ce)
				return true
			}
			queries = append(queries, queryInfo{query})
			return false
		})
	}
	if err != nil {
		return nil, err
	}

	if len(unresolved) > 0 {
		fmt.Fprintln(os.Stderr, "UNRESOLVED")
	}
	for _, ce := range unresolved {
		fmt.Fprintln(os.Stderr, pkg.Fset.Position(ce.Pos()), types.ExprString(ce))
	}

	return queries, nil
}
