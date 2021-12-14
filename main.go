package main

import (
	"debug/dwarf"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/golang-collections/collections/stack"
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
		{Text: "back", Description: "Move context back to the previous DIE we were targeting"},
		{Text: "list_all_attributes", Description: "List each attribute we find under this entry"},
	}
	return prompt.FilterHasSuffix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

	entryReader, _ := parser.GetReader(filename)
	// Start by printing the first DIE we find
	entry, _ := entryReader.Next()
	parser.PrintEntryInfo(entry)

	// Stack onto which we will push entries when we switch
	// contexts to allow navigating "back" to a previous context
	var entryStack stack.Stack

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
			parser.PrintEntryInfo(entry)
		case "next":
			{
				if entry == nil {
					fmt.Println("Encountered a nil entry")
					break
				}
				entry, _ = entryReader.Next()
				parser.PrintEntryInfo(entry)
			}
		case "type":
			entryStack.Push(entry.Offset)
			entry = parser.GetTypeDie(entryReader, entry)
			parser.PrintEntryInfo(entry)
		case "back":
			if entryStack.Len() == 0 {
				fmt.Println("No context to go backwards to")
			} else {
				// Restore context to the previous DIE we were viewing
				entryOffset := entryStack.Pop().(dwarf.Offset)
				entryReader.Seek(entryOffset)
				entry, _ = entryReader.Next()
				parser.PrintEntryInfo(entry)
			}
    case "list_all_attributes":
      parser.ListAllAttributes(entry)
		}
	}
}
