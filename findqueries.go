// The findqueries package defines an Analyzer that find SQL queries
// and all functions calling such queries inside themselves.

package findqueries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `find SQL queries and functions calling them internally`

var Analyzer = &analysis.Analyzer{
	Name: "findqueries",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
	},
	ResultType: reflect.TypeOf((*Result)(nil)),
	FactTypes:  []analysis.Fact{},
}

// Diagnostic categories
const (
	QUERY = "query"
)

var (
	queriersInfoPath string // -queriers flag
	buildersInfoPath string // -builders flag
	silent           bool   // -silent flag
)

func init() {
	Analyzer.Flags.StringVar(&queriersInfoPath, "queriers", defaultQuerierInfoPath(), "path to queriers.json")
	Analyzer.Flags.StringVar(&buildersInfoPath, "builders", defaultBuilderInfoPath(), "path to builders.json")
	Analyzer.Flags.BoolVar(&silent, "silent", false, "stop emitting json")
}

func run(pass *analysis.Pass) (interface{}, error) {
	withMessage := func(err error) error {
		return errors.WithMessagef(err, "failed to analyze pass %s", pass.String())
	}

	queriers, err := loadQuerierInfo(queriersInfoPath)
	if err != nil {
		return nil, withMessage(err)
	}

	builders, err := loadBuilderInfo(buildersInfoPath)
	if err != nil {
		return nil, withMessage(err)
	}

	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	inspctr := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	er := NewExprResolver(inspctr)

	result, err := analyzePackage(ssa.Pkg, er, queriers, builders)
	if err != nil {
		return nil, withMessage(err)
	}

	for _, qi := range result.Queries {
		pass.Report(analysis.Diagnostic{
			Pos:      qi.sitePos,
			Category: QUERY,
			Message:  qi.Query,
		})
	}

	for _, qi := range result.unresolved {
		fmt.Fprintln(os.Stderr, qi.Pos, qi.Expr)
		fmt.Fprintln(os.Stderr, qi.err)
	}

	if !silent {
		b, err := json.Marshal(result)
		if err != nil {
			return nil, withMessage(errors.Wrap(err, "failed to jsonify result"))
		}
		buf := new(bytes.Buffer)
		if err := json.Compact(buf, b); err != nil {
			return nil, withMessage(errors.Wrap(err, "failed to compact json"))
		}
		fmt.Println(buf)
	}

	return result, nil
}
