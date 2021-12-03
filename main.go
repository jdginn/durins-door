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
    // if entry == nil {
    //   fmt.Println("Encountered a nil entry")
    //   break
    // }
		if entry.Tag == dwarf.TagBaseType {
			fmt.Println("Found a struct")
      name := ""
      var byte_size int64
			for _, field := range entry.Field {
        if field.Attr == dwarf.AttrName {
          name = field.Val.(string)
          // fmt.Println(field.Val.(string))
        }
        if field.Attr == dwarf.AttrByteSize {
          byte_size = field.Val.(int64)
        }
			}
      fmt.Printf("Name: %s \n Size: %d \n\n", name, byte_size)
		}
	}
}
