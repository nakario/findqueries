package findqueries

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

type querierInfo struct{
	FullName string `json:"full_name"`
	QueryPos int    `json:"query_pos"`
}

func defaultQuerierInfoPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "querier", "queriers.json")
}

func loadQuerierInfo(path string) ([]querierInfo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read queriers data")
	}
	qis := make([]querierInfo, 0)
	if err := json.Unmarshal(data, &qis); err != nil {
		return nil, errors.Wrap(err, "failed to read queriers data")
	}
	return qis, nil
}
