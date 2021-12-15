package parser

import (
	"debug/dwarf"
	"debug/macho"
  "errors"
	"fmt"
)

// Return a dwarf.Reader object for a macho file
// TODO: support ELF in addition to macho...
func GetReader(filename string) (*dwarf.Reader, error) {
	fmt.Println("Opening ", filename)
	machoFile, err := macho.Open(filename)

	dwarfData, err := machoFile.DWARF()
	fmt.Println("Collected dwarfData")

	entryReader := dwarfData.Reader()

	return entryReader, err
}

// Search for an entry matching a requested name
func GetEntry(reader *dwarf.Reader, name string) (*dwarf.Entry, error) {
  fmt.Printf("Locating %s\n", name)
  for {
    e, err := reader.Next()
    if err != nil {
      return nil, err
    }
    if e == nil {
      return nil, errors.New("entry could not be found")
    }
    if e.AttrField(dwarf.AttrName) == nil {
      continue
    }
    if e.AttrField(dwarf.AttrName).Val == name {
      return e, err
    }
  }
}

func GetBitSize(entry *dwarf.Entry) int {
  if hasAttr(entry, dwarf.AttrBitSize) {
    return entry.AttrField(dwarf.AttrBitSize).Val.(int)
  } else {
    return int(entry.AttrField(dwarf.AttrByteSize).Val.(int64) * 8)
  }
}

// Display key information about this entry; strive to be easily readable.
func PrintEntryInfo(entry *dwarf.Entry) {
  // JDG TODO: make sure I'm using the right DW_AT names here
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

// Print each attribute for this entry.
func ListAllAttributes(entry *dwarf.Entry) {
	fmt.Println("All fields in this entry:")
	for _, field := range entry.Field {
		fmt.Printf("  %v\n", field.Attr)
	}
}

// Return whether this entry contains a requested attribute
func hasAttr(entry *dwarf.Entry, attr dwarf.Attr) bool {
	for _, field := range entry.Field {
		if field.Attr == attr {
			return true
		}
	}
	return false
}

// Translate a DW_AT_locationn attribute into an address
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

// Return the entry defining the type for a given entry. Returns self if
// no such entry can be found. Leaves the reader at the new entry.
func GetTypeEntry(reader *dwarf.Reader, entry *dwarf.Entry) (*dwarf.Entry, error) {
	if !hasAttr(entry, dwarf.AttrType) {
		fmt.Println("This entry does not have a type entry - returning it as-is")
		return entry, nil
	}
	var typeDie *dwarf.Entry
	for _, field := range entry.Field {
		if field.Attr == dwarf.AttrType {
			typeDieOffset := field.Val.(dwarf.Offset)
			fmt.Printf("  DW_AT_type_die: %v\n", typeDieOffset)
			reader.Seek(typeDieOffset)
			typeDie, _ := reader.Next()
			return typeDie, nil
		}
	}
	return typeDie, nil
}
