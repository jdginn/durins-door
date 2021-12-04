package main

import (
  "github.com/jdginn/dwarf-experiments/parser"
	"fmt"
	"os"
)

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

	entryReader := parser.GetReader(filename)

	for {
		entry, _ := entryReader.Next()
		if entry == nil {
			fmt.Println("Encountered a nil entry")
			break
		}
		parser.PrintDieInfo(entryReader, entry)
	}
}
