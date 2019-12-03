package findqueries

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

type builderInfo struct {
	FullName string `json:"full_name"`
	ArgIndex int    `json:"arg_index"`
	RetIndex int    `json:"ret_index"`
}

func defaultBuilderInfoPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "querier", "builders.json")
}

func loadBuilderInfo(path string) ([]builderInfo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read builders data")
	}
	bis := make([]builderInfo, 0)
	if err := json.Unmarshal(data, &bis); err != nil {
		return nil, errors.Wrap(err, "failed to read builders data")
	}
	return bis, nil
}
