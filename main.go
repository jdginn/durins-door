package main

import (
	"bufio"
	"fmt"
	"os"
  "strings"
)

import "github.com/jdginn/dwarf-experiments/parser"

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

	entryReader := parser.GetReader(filename)
	// Start by printing the first DIE we find
	entry, _ := entryReader.Next()
	parser.PrintDieInfo(entry)

	r := bufio.NewReader(os.Stdin)

	// Parse input from the user
	for {
		fmt.Printf("(Enter command) > ")
		command, _ := r.ReadString('\n')
    command = strings.TrimSpace(command)
		switch command {
    case "h":
    {
        fmt.Println("Supported commands are:")
        fmt.Println("  h: display this message")
        fmt.Println("  n: iterate to next DIE in the current context")
        fmt.Println("  t: display this DIE's type DIE")
      }
		case "n":
			{
				if entry == nil {
					fmt.Println("Encountered a nil entry")
					break
				}
				entry, _ := entryReader.Next()
				parser.PrintDieInfo(entry)
			}
    case "t":
      typeDie := parser.GetTypeDie(entryReader, entry) 
      parser.PrintDieInfo(typeDie)
		}

	}
}
