package parser

import (
	"debug/dwarf"
	"debug/macho"
	"fmt"
)

func hasAttr(entry *dwarf.Entry, attr dwarf.Attr) bool {
  for _, field := range entry.Field {
    if field.Attr == attr {
      return true
    }
  }
  return false
}

func SetContextToThisEntry(reader *dwarf.Reader, entry *dwarf.Entry) {
  reader.Seek(entry.Offset)
}

func ParseLocation(location []uint8) int64 {
	// Ignore the first entry in the slice
	// --> This somehow communicates a format?
	// Build the last slice from right to left
	location = location[1:]
	var locationAsInt int64
	locationAsInt = 0
	for i := 0; i < len(location); i++ {
		locationAsInt += int64(location[i]) << (8 * i)
	}
	return locationAsInt
}

func GetTypeDie(reader *dwarf.Reader, entry *dwarf.Entry) *dwarf.Entry {
  if !hasAttr(entry, dwarf.AttrType) {
    fmt.Println("This DIE does not have a type DIE - returning it as-is")
    return entry
  }
	var typeDie *dwarf.Entry
	for _, field := range entry.Field {
		if field.Attr == dwarf.AttrType {
			typeDieOffset := field.Val.(dwarf.Offset)
			fmt.Printf("  DW_AT_type_die: %v\n", typeDieOffset)
			reader.Seek(typeDieOffset)
			typeDie, _ := reader.Next()
			return typeDie
		}
	}
	return typeDie
}

func ListAllAttributes(entry *dwarf.Entry) {
  fmt.Println("All fields in this entry:")
  for _, field := range entry.Field {
    fmt.Printf("  %v\n", field.Attr)
  }
}

// JDG TODO: make sure I'm using the right DT_AT names here
func PrintDieInfo(entry *dwarf.Entry) {
  fmt.Printf("Tag: %s\n", entry.Tag)
  fmt.Printf("  Children: %v\n", entry.Children)
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
			fmt.Printf("  DW_AT_location: %x\n", ParseLocation(location))
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
			fmt.Printf("  DW_AT_type_die at offset: %v\n", field.Val)
		}
	}
}

func GetReader(filename string) (*dwarf.Reader, error) {
	fmt.Println("Filename: ", filename)

	fmt.Println("Opening file")
	machoFile, err := macho.Open(filename)
	fmt.Println("Opened elfFile")

	dwarfData, err := machoFile.DWARF()
	fmt.Println("Collected dwarfData")

	entryReader := dwarfData.Reader()

	return entryReader, err
}
