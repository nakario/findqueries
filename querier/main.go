package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"runtime"
	"sort"
)

type querierInfo struct{
	FullName string `json:"full_name"`
	QueryPos int    `json:"query_pos"`
}

func main() {
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Join(filepath.Dir(filename), "queriers.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}
	conf := &types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Uses: make(map[*ast.Ident]types.Object),
	}
	_, err = conf.Check("github.com/nakario/findqueries/test", fset, []*ast.File{f}, info)
	if err != nil {
		panic(err)
	}
	qis := make([]querierInfo, 0)
	for _, obj := range info.Uses {
		if fn, ok := obj.(*types.Func); ok {
			fullName := fn.FullName()
			sig := fn.Type().(*types.Signature)
			params := sig.Params()
			pos := -1
			for i := 0; i < params.Len(); i++ {
				if params.At(i).Name() == "query" {
					pos = i
					break
				}
			}
			if pos == -1 {
				panic("argument \"query\" not found")
			}
			qis = append(qis, querierInfo{fullName, pos})
		}
	}
	sort.Slice(qis, func(i, j int) bool {
		return qis[i].FullName < qis[j].FullName
	})
	data, err := json.Marshal(qis)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, data); err != nil {
		panic(err)
	}
	fmt.Println(buf)
}
