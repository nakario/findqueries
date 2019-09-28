package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		exe, err := os.Executable()
		if err != nil {
			log.Fatalln("failed toget executable name:", err)
		}
		fmt.Println("Usage:", exe, "path/to/dir")
		return
	}
	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalln("failed to get absolute path of", os.Args[1], ":", err)
	}

	findQueries(dir)
}