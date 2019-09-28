package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
)

type _package struct{
	pkg *types.Package
	fset *token.FileSet
	info *types.Info
	files []*ast.File
	qualifier func(*types.Package) string
}

func newPackage(path string, fset *token.FileSet, p *ast.Package) (*_package, error) {
	conf := &types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
		Implicits: make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes: make(map[ast.Node]*types.Scope),
	}
	files := make([]*ast.File, 0, len(p.Files))
	for _, f := range p.Files {
		files = append(files, f)
	}
	pkg, err := conf.Check(path, fset, files, info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to type-check")
	}
	qualifier := func(other *types.Package) string {
		if pkg == other {
			return ""
		}
		return other.Path()
	}
	return &_package{pkg, fset, info, files, qualifier}, nil
}

func searchPackage(pkg *ast.Package, path string, fset *token.FileSet, queryers []queryerInfo) ([]queryInfo, error) {
	p, err := newPackage(path, fset, pkg)
	if err != nil {
		return nil, err
	}

	queryersMap := make(map[string]int)
	for _, qi := range queryers {
		queryersMap[qi.FullName] = qi.QueryPos
	}

	end2fn := make(map[token.Pos]*types.Func)
	for ident, obj := range p.info.Uses {
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
	for _, f := range p.files {
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
