package main

import (
	"debug/dwarf"
  "debug/macho"
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
		{Text: "skip_children", Description: "Skip all children of this entry"},
		{Text: "type_die", Description: "Move context to this DIE's type DIE"},
		{Text: "type", Description: "Read the type corresponding to this entry if possible"},
		{Text: "back", Description: "Move context back to the previous DIE we were targeting"},
		{Text: "list_all_attributes", Description: "List each attribute we find under this entry"},
	}
	return prompt.FilterHasSuffix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

  file, err := macho.Open(filename)
  data, err := file.DWARF()
	reader := data.Reader()
  if err != nil {
    fmt.Printf("Failed with error %v", err)
  }
    
	// Start by printing the first DIE we find
	entry, _ := reader.Next()
	parser.PrintEntryInfo(entry)

	// Stack onto which we will push entries when we switch
	// contexts to allow navigating "back" to a previous context
	var entryStack stack.Stack

	// Parse input from the user
	for {
		command := prompt.Input("> ", completer)
		switch command {
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
				entry, _ = reader.Next()
				parser.PrintEntryInfo(entry)
			}
    case "skip_children":
      reader.SkipChildren()
		case "type_die":
			entryStack.Push(entry.Offset)
			entry = parser.GetTypeEntry(reader, entry)
			parser.PrintEntryInfo(entry)
		case "type":
      if entry.Tag == dwarf.TagTypedef {
        typ, _ := data.Type(entry.Offset)
        fmt.Printf("Type: %v\n", typ)
        fmt.Printf("  Size: %v\n", typ.Size())
      } else if entry.AttrField(dwarf.AttrType).Val != nil {
        typ, _ := data.Type(entry.AttrField(dwarf.AttrType).Val.(dwarf.Offset))
        fmt.Printf("Type: %v\n", typ)
        fmt.Printf("  Size: %v\n", typ.Size())
      } else {
        fmt.Println("Cannot get a type from this tag")
      }
		case "back":
			if entryStack.Len() == 0 {
				fmt.Println("No context to go backwards to")
			} else {
				// Restore context to the previous DIE we were viewing
				entryOffset := entryStack.Pop().(dwarf.Offset)
				reader.Seek(entryOffset)
				entry, _ = reader.Next()
				parser.PrintEntryInfo(entry)
			}
    case "list_all_attributes":
      parser.ListAllAttributes(entry)
		}
	}
}
