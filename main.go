package main

import (
	"fmt"
  "github.com/c-bata/go-prompt"
	"os"
)

import "github.com/jdginn/dwarf-experiments/parser"

func completer(d prompt.Document) []prompt.Suggest { 
  s := []prompt.Suggest{
    {Text: "help", Description: "View help documentation"},
    {Text: "quit", Description: "Display current DIE's type DIE"},
    {Text: "print", Description: "Display current DIE"},
    {Text: "next", Description: "Advance to the next DIE in the current context"},
    {Text: "type_die", Description: "Move context to this DIE's type DIE"},
    // WIP, unimplemented
    // JDG TODO: need to store a stack of entry pointers
    // JDG TODO: how do we reset the reader to the appropriate context for an entry?
    {Text: "back", Description: "Move context back to the previous DIE we were targeting"},
  }
  return prompt.FilterHasSuffix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

	entryReader := parser.GetReader(filename)
	// Start by printing the first DIE we find
	entry, _ := entryReader.Next()
	parser.PrintDieInfo(entry)

	// Parse input from the user
	for {
    command := prompt.Input("> ", completer)
		switch command {
    case "help":
    {
        fmt.Println("Supported commands are:")
        fmt.Println("  help: display this message")
        fmt.Println("  quit: quit")
        fmt.Println("  next: iterate to next DIE in the current context")
        fmt.Println("  type: display this DIE's type DIE")
      }
    case "quit":
      return
    case "print":
      parser.PrintDieInfo(entry)
		case "next":
			{
				if entry == nil {
					fmt.Println("Encountered a nil entry")
					break
				}
				entry, _ = entryReader.Next()
				parser.PrintDieInfo(entry)
			}
    case "type":
      typeDie := parser.GetTypeDie(entryReader, entry) 
      parser.PrintDieInfo(typeDie)
    case "back":
      fmt.Println("back is not yet implemented")
		}
	}
}
