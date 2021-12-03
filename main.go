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
    if entry == nil {
      fmt.Println("Encountered a nil entry")
      break
    }
    // fmt.Println("Found a struct")
    fmt.Printf("Found a %s\n", entry.Tag)
    name := ""
    var byte_size int64
    fmt.Printf("  Attributes:\n")
    for _, field := range entry.Field {
      fmt.Printf("      %s\n", field.Attr)
      if field.Attr == dwarf.AttrName {
        name = field.Val.(string)
        fmt.Printf("  DW_AT_name: %s\n", name)
      }
      if field.Attr == dwarf.AttrByteSize {
        byte_size = field.Val.(int64)
        fmt.Printf("  DW_AT_byte_size: %d\n", byte_size)
      }
    }
	}
}
