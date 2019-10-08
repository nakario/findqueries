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
	builderInfoPath := defaultBuilderInfoPath()
	switch len(os.Args) {
	case 2:
		break
	case 3:
		var err error
		queryerInfoPath, err = filepath.Abs(os.Args[2])
		if err != nil {
			log.Fatalln("failed to get path of queryers.json:", err)
		}
	case 4:
		var err error
		builderInfoPath, err = filepath.Abs(os.Args[3])
		if err != nil {
			log.Fatalln("failed to get path of builders.json:", err)
		}
	default:
		exe, err := os.Executable()
		if err != nil {
			log.Fatalln("failed to get executable name:", err)
		}
		fmt.Println("Usage:", exe, "path/to/dir", "[path/to/queryers.json, [path/to/builders.json]]")
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

	builders, err := loadBuilderInfo(builderInfoPath)
	if err != nil {
		log.Fatalln("failed to load builders:", err)
	}

	result, err := analyze(dir, queryers, builders)
	if err != nil {
		log.Fatalln("failed to find queries:", err)
	}

	b, err := json.Marshal(result)
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, b); err != nil {
		log.Fatalln("failed to compact json:", err)
	}
	fmt.Println(buf)
}
