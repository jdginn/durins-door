package main

import (
	"debug/dwarf"
	// "debug/elf"
	"debug/macho"
	"fmt"
	// "io/ioutil"
	"log"
	"os"
)

func parse_location(location []uint8) int64 {
  // Ignore the first entry in the slice
  // --> This somehow communicates a format?
  // Build the last slice from right to left
  location = location[1:]
  // for i := 0; i < len(location)/2; i++ {
  //   j := len(location) - i - 1
  //   location[i], location[j] = location[j], location[i]
  // }
  fmt.Printf("Flipped slice: %v\n", location)
  var location_as_int int64
  location_as_int = 0
  for i := 0; i < len(location); i++ {
    location_as_int += int64(location[i]) << (8 * i)
  }
  return location_as_int
}

func follow_typedef(entry *dwarf.Entry) {
}

// JDG TODO: make sure I'm using the right DT_AT names here
func print_die_info(entry *dwarf.Entry) {
  fmt.Printf("Found a %s\n", entry.Tag)
  // fmt.Printf("  Attributes:\n")
  for _, field := range entry.Field {
    // fmt.Printf("      %s\n", field.Attr)
    if field.Attr == dwarf.AttrName {
      name := field.Val.(string)
      fmt.Printf("  DW_AT_name: %s\n", name)
    }
    if field.Attr == dwarf.AttrByteSize {
      byte_size := field.Val.(int64)
      fmt.Printf("  DW_AT_byte_size: %d\n", byte_size)
    }
    if field.Attr == dwarf.AttrLocation{
      location := field.Val.([]uint8)
      fmt.Printf("  DW_AT_location: %x\n", parse_location(location))
    }
    if field.Attr == dwarf.AttrDataMemberLoc{
      location := field.Val
      fmt.Printf("  DW_AT_data_member_location: %x\n", location)
    }
    if field.Attr == dwarf.AttrCompDir {
      comp_dir := field.Val
      fmt.Printf("  DW_AT_comp_dir: %s\n", comp_dir)
    }
    if field.Attr == dwarf.AttrType{
      type_die := field.Val
      fmt.Printf("  DW_AT_type_die ... TODO: offset?: %v\n", type_die)
    }
  }
}

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
    print_die_info(entry)
	}
}
