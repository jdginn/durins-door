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
	var location_as_int int64
	location_as_int = 0
	for i := 0; i < len(location); i++ {
		location_as_int += int64(location[i]) << (8 * i)
	}
	return location_as_int
}

// JDG TODO: make sure I'm using the right DT_AT names here
func print_die_info(reader *dwarf.Reader, entry *dwarf.Entry) {
	fmt.Printf("Found a %s\n", entry.Tag)
	for _, field := range entry.Field {
		if field.Attr == dwarf.AttrName {
			name := field.Val.(string)
			fmt.Printf("  DW_AT_name: %s\n", name)
		}
		if field.Attr == dwarf.AttrByteSize {
			byte_size := field.Val.(int64)
			fmt.Printf("  DW_AT_byte_size: %d\n", byte_size)
		}
		if field.Attr == dwarf.AttrLocation {
			location := field.Val.([]uint8)
			fmt.Printf("  DW_AT_location: %x\n", parse_location(location))
		}
		if field.Attr == dwarf.AttrDataMemberLoc {
			location := field.Val
			fmt.Printf("  DW_AT_data_member_location: %x\n", location)
		}
		if field.Attr == dwarf.AttrCompDir {
			comp_dir := field.Val
			fmt.Printf("  DW_AT_comp_dir: %s\n", comp_dir)
		}
		if field.Attr == dwarf.AttrType {
      curr_offset := entry.Offset
			type_die_offset := field.Val.(dwarf.Offset)
      fmt.Printf("  DW_AT_type_die: %v\n", type_die_offset)
      fmt.Printf("    curr_offset: %v\n", curr_offset)
      reader.Seek(type_die_offset)
      type_die, _ := reader.Next()
      fmt.Printf("  DW_AT_type_die:\n")
      print_die_info(reader, type_die)
      fmt.Println("")
      // Restore us to the offset we were reading before we jumped to follow the type
      fmt.Printf("    Restoring offset to %v\n", curr_offset)
      reader.Seek(curr_offset)
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
		print_die_info(entryReader, entry)
	}
}
