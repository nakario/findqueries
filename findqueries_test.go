package findqueries_test

import (
	"testing"

	"github.com/nakario/findqueries"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, findqueries.Analyzer, "a")
}
