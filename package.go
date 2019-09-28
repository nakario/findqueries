package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"

	"github.com/pkg/errors"
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

func searchPackage(pkg *ast.Package, path string, fset *token.FileSet) ([]queryInfo, error) {
	p, err := newPackage(path, fset, pkg)
	if err != nil {
		return nil, err
	}
	queries := make([]queryInfo, 0)

	for _, f := range p.files {
		ast.Inspect(f, func(n ast.Node) bool {
			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			t := types.TypeString(p.info.TypeOf(ce), p.qualifier)
			fmt.Println(t)
			return false
		})
	}

	return queries, nil
}
