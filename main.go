package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	queryerInfoPath := defaultQueryerInfoPath()
	switch len(os.Args) {
	case 2:
		break
	case 3:
		var err error
		queryerInfoPath, err = filepath.Abs(os.Args[2])
		if err != nil {
			log.Fatalln("failed to get path of queryers.json:", err)
		}
	default:
		exe, err := os.Executable()
		if err != nil {
			log.Fatalln("failed to get executable name:", err)
		}
		fmt.Println("Usage:", exe, "path/to/dir", "[path/to/queryers.json]")
		return
	}
	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalln("failed to get absolute path of", os.Args[1], ":", err)
	}

	queryers, err := loadQueryerInfo(queryerInfoPath)
	if err != nil {
		log.Fatalln("failed to load queryers:", err)
	}

	queries, err := findQueries(dir, queryers)
	if err != nil {
		log.Fatalln("failed to find queries:", err)
	}

	b, err := json.Marshal(queries)
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, b); err != nil {
		log.Fatalln("failed to compact json:", err)
	}
	fmt.Println(buf)
}
