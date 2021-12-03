package main

// import "debug/dwarf"
import (
	"debug/dwarf"
	// "debug/elf"
	"debug/macho"
	"fmt"
	// "io/ioutil"
	"log"
	"os"
)

func main() {
	filename := os.Args[1]
	fmt.Println("Filename: ", filename)

	fmt.Println("Opening file")
	machoFile, err := macho.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Opened elfFile")

	dwarfData, err := machoFile.DWARF()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Collected dwarfData")

	entryReader := dwarfData.Reader()

	for {
		entry, _ := entryReader.Next()
    // fmt.Println("Entry:")
    // fmt.Println(entry)
    // fmt.Println("")
    if entry == nil {
      break
    }
		if entry != nil && entry.Tag == dwarf.TagBaseType {
			fmt.Println("Found a struct")
			for _, field := range entry.Field {
        fmt.Printf("Found an attribute: %s\n", field.Attr)
        if field.Attr == dwarf.AttrName {
          fmt.Println(field.Val.(string))
        }
			}
		}
	}
}
