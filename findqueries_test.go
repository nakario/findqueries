package findqueries_test

import (
	"path/filepath"
	"testing"

	"github.com/nakario/findqueries"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	findqueries.Analyzer.Flags.Set("silent", "true")
	findqueries.Analyzer.Flags.Set("queriers", filepath.Join(testdata, "queriers.json"))
	analysistest.Run(t, testdata, findqueries.Analyzer, "a")
}
