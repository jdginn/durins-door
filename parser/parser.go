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
	machoFile, err := macho.Open(filename)
	dwarfData, err := machoFile.DWARF()
	entryReader := dwarfData.Reader()
	return entryReader, err
}

// Iterate once through the remaining entries looking for
// an entry by name
func getEntryByNameFromRemaining(reader *dwarf.Reader, name string) (*dwarf.Entry, error) {
	for {
		e, err := reader.Next()
		if err != nil {
			return nil, err
		}
		if e == nil {
			return e, err
		}
		// TODO: there may be an optimization to skip children in some cases?
		if e.AttrField(dwarf.AttrName) == nil {
			continue
		}
		if e.Val(dwarf.AttrName) == name {
			return e, err
		}
	}
}

// Search for an entry matching a requested name
func GetEntry(reader *dwarf.Reader, name string) (*dwarf.Entry, error) {
	e, err := getEntryByNameFromRemaining(reader, name)
	// If we don't find the entry by the time we reach the end of the DWARF
	// section, we need to start searching again from the beginning. We avoid
	// always seeking back to the beginning because in most cases, the entry
	// we are looking for is more likely to come after the most recent
	// entry.
	if e == nil {
		reader.Seek(0)
		e, err = getEntryByNameFromRemaining(reader, name)
	}
	if e == nil {
		err = errors.New(fmt.Sprintf("Could not find entry %v", name))
	}
	return e, err
}

// Find the size of the type defined by this entry, in bits
func GetBitSize(entry *dwarf.Entry) (int, error) {
	if hasAttr(entry, dwarf.AttrBitSize) {
		return entry.Val(dwarf.AttrBitSize).(int), nil
	} else if hasAttr(entry, dwarf.AttrByteSize) {
		return int(entry.Val(dwarf.AttrByteSize).(int64) * 8), nil
	} else {
		return 0, errors.New(fmt.Sprintf("Could not get bit size of entry:\n%v", FormatEntryInfo(entry)))
	}
}

// Return a slice with en entry for the range of each array dimension
//
// Scalar types will have range []int{0}. The length of the return defines
// the dimension of the array.
func GetArrayRanges(reader *dwarf.Reader, entry *dwarf.Entry) ([]int, error) {
	_, err := GetTypeEntry(reader, entry)
	ranges := make([]int, 0)
	// typeEntry, err := GetTypeEntry(reader, entry)
	// var err error
	for {
		// fmt.Println("Stepping through a subrange:")
		subrange, _ := reader.Next()
		// fmt.Println(FormatEntryInfo(subrange))

		// When we've finished iterating over members, we are done with the meaningful
		// children of this typedef. We are also finished if we reach the end of the DWARF
		// section during this iteration.
		if subrange == nil || subrange.Tag == 0 {
			break
		}

		if hasAttr(subrange, dwarf.AttrCount) {
			ranges = append(ranges, int(subrange.Val(dwarf.AttrCount).(int64)))
		}
	}
	return ranges, err
}

// Format key information about this entry as a string; strive to be easily readable.
func FormatEntryInfo(entry *dwarf.Entry) string {
	if entry == nil {
		fmt.Println("ERROR: nil entry passed")
	}
	// JDG TODO: make sure I'm using the right DW_AT names here
	var str string
	str = fmt.Sprintf("Tag: %s\n", entry.Tag)
	str += fmt.Sprintf("  Children: %v\n", entry.Children)
	str += fmt.Sprintf("  Offset: %v\n", entry.Offset)
	for _, field := range entry.Field {
		// str += fmt.Sprintf("  %s: %v", field.Val)
		if field.Attr == dwarf.AttrName {
			name := field.Val.(string)
			str += fmt.Sprintf("  DW_AT_name: %s\n", name)
		}
		if field.Attr == dwarf.AttrByteSize {
			byte_size := field.Val.(int64)
			str += fmt.Sprintf("  DW_AT_byte_size: %d\n", byte_size)
		}
		if field.Attr == dwarf.AttrLocation {
			location, err := GetLocation(entry)
			str += fmt.Sprintf("  DW_AT_location: %x\n", ParseLocation(location))
		}
		if field.Attr == dwarf.AttrDataMemberLoc {
			location := field.Val
			str += fmt.Sprintf("  DW_AT_data_member_location: %x\n", location)
		}
		if field.Attr == dwarf.AttrCompDir {
			comp_dir := field.Val
			str += fmt.Sprintf("  DW_AT_comp_dir: %s\n", comp_dir)
		}
		if field.Attr == dwarf.AttrType {
			str += fmt.Sprintf("  DW_AT_type_die at offset: %v\n", field.Val)
		}
		if field.Attr == dwarf.AttrDataLocation {
			str += fmt.Sprintf("  DW_AT_data_location: %v\n", field.Val)
		}
		if field.Attr == dwarf.AttrDataMemberLoc {
			str += fmt.Sprintf("  DW_AT_data_member_loc: %v\n", field.Val)
		}
		if field.Attr == dwarf.AttrDataBitOffset {
			str += fmt.Sprintf("  DW_AT_data_bit_offset: %v\n", field.Val)
		}
		if field.Attr == dwarf.AttrCount {
			str += fmt.Sprintf("  DW_AT_count: %d\n", field.Val)
		}
	}
	return str
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

func GetLocation(entry *dwarf.Entry) ([]uint8, error) {
  var err error = nil
  loc := entry.Val(dwarf.AttrDataLocation)
  if loc == nil {
		err = errors.New(fmt.Sprintf("Could not find data location for %v", entry))
    return nil, err
  }
	return entry.Val(dwarf.AttrDataLocation).([]uint8), nil
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
		// fmt.Printf("Entry %v does not have a type entry - returning it as-is\n", entry.Val(dwarf.AttrName))
		return entry, nil
	}
	var typeDie *dwarf.Entry
	for _, field := range entry.Field {
		if field.Attr == dwarf.AttrType {
			typeDieOffset := field.Val.(dwarf.Offset)
			reader.Seek(typeDieOffset)
			typeDie, _ := reader.Next()
			return typeDie, nil
			// if typeDie.Tag == dwarf.TagTypedef {
			//   typeDie, _ = GetTypeEntry(reader, typeDie)
			// }
			// break
		}
	}
	return typeDie, nil
}
