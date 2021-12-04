package main

import (
	"bufio"
	"fmt"
	"os"
)

import "github.com/jdginn/dwarf-experiments/parser"

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

	entryReader := parser.GetReader(filename)
	// Start by printing the first DIE we find
	entry, _ := entryReader.Next()
	parser.PrintDieInfo(entryReader, entry)

	r := bufio.NewReader(os.Stdin)

	// Parse input from the user
	for {
		command, _ := r.ReadString('\n')
		switch {
		case command == "n":
			if entry == nil {
				fmt.Println("Encountered a nil entry")
				break
			}
			entry, _ := entryReader.Next()
			parser.PrintDieInfo(entryReader, entry)
		}

	}
}
