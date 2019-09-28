package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

type queryerInfo struct{
	FullName string `json:"full_name"`
	QueryPos int    `json:"query_pos"`
}

func defaultQueryerInfoPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "queryer", "queryers.json")
}

func loadQueryerInfo(path string) ([]queryerInfo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read queryers data")
	}
	qis := make([]queryerInfo, 0)
	if err := json.Unmarshal(data, &qis); err != nil {
		return nil, errors.Wrap(err, "failed to read queryers data")
	}
	return qis, nil
}
