package parser

import (
	"debug/dwarf"
	"errors"
	"fmt"
)

type DebugFile interface {
	DWARF() (*dwarf.Data, error)
}

// Returns a dwarf.Reader object
func GetReader[T DebugFile](fh T) (*dwarf.Reader, error) {
	dwarfData, err := fh.DWARF()
	if err != nil {
		panic(err)
	}
	entryReader := dwarfData.Reader()
	return entryReader, err
}

// Return a slice of all CompileUnits in the DWARF
func GetCUs(r *dwarf.Reader) ([]*dwarf.Entry, error) {
  entries := make([]*dwarf.Entry, 0)
  for {
    entry, err := r.Next()
    // Since we just want CUs, we never want to see their children
    r.SkipChildren()
    if err != nil {
      return entries, err
    }
    if entry == nil { break }
    if entry.Tag == dwarf.TagCompileUnit {
      entries = append(entries, entry)
    }
  }
  return entries, nil
}

// Iterates once through the remaining entries looking for an entry by name
//
// The second argument returns true if the entry could be found
//
// TODO: this many return values seems like a bad idea
func getFromRemaining(r *dwarf.Reader, name string) (*dwarf.Entry, *dwarf.Entry, bool, error) {
	var lastCU *dwarf.Entry
	for {
		entry, err := r.Next()
		if err != nil {
			return nil, nil, false, err
		}
		if entry == nil {
			return nil, nil, false, err
		}
		if entry.Tag == dwarf.TagCompileUnit {
			lastCU = entry
		}
		// TODO: there may be an optimization to skip children in some cases?
		if entry.AttrField(dwarf.AttrName) == nil {
			continue
		}
		if entry.Val(dwarf.AttrName) == name {
			return entry, lastCU, true, nil
		}
	}
}

// Searches for an entry matching a requested name
//
// TODO: this many return values seems like a bad idea
func GetEntry(r *dwarf.Reader, name string) (*dwarf.Entry, *dwarf.Entry, error) {
	e, lastCU, ok, err := getFromRemaining(r, name)
	if err != nil {
		return nil, lastCU, err
	}
	// If we don't find the entry by the time we reach the end of the DWARF
	// section, we need to start searching again from the beginning. We avoid
	// always seeking back to the beginning because in most cases, the entry
	// we are looking for is more likely to come after the most recent
	// entry.
	if !ok {
		r.Seek(0)
		e, lastCU, ok, err = getFromRemaining(r, name)
	}
	if !ok {
		err = fmt.Errorf("Could not find entry %v", name)
	}
	return e, lastCU, err
}

// Finds the size of the type defined by this entry, in bits
func GetBitSize(entry *dwarf.Entry) (int, error) {
	if hasAttr(entry, dwarf.AttrBitSize) {
		return entry.Val(dwarf.AttrBitSize).(int), nil
	} else if hasAttr(entry, dwarf.AttrByteSize) {
		return int(entry.Val(dwarf.AttrByteSize).(int64) * 8), nil
	} else {
		return 0, errors.New(fmt.Sprintf("Could not get bit size of entry:\n%v", FormatEntryInfo(entry)))
	}
}

// Returns a slice with en entry for the range of each array dimension
//
// Scalar types will have range []int{0}. The length of the return defines
// the dimension of the array.
func GetArrayRanges(r *dwarf.Reader, entry *dwarf.Entry) ([]int, error) {
	_, err := GetTypeEntry(r, entry)
	ranges := make([]int, 0)
	// typeEntry, err := GetTypeEntry(reader, entry)
	// var err error
	for {
		// fmt.Println("Stepping through a subrange:")
		subrange, _ := r.Next()
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

// Formats key information about this entry as a string; strives to be easily readable.
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
			byte_size := field.Val
			str += fmt.Sprintf("  DW_AT_byte_size: %d\n", byte_size)
		}
		if field.Attr == dwarf.AttrLocation {
			location, _ := GetLocation(entry)
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

// Prints each attribute for this entry.
func ListAllAttributes(entry *dwarf.Entry) {
	fmt.Println("All fields in this entry:")
	for _, field := range entry.Field {
		fmt.Printf("  %v\n", field.Attr)
	}
}

// Return true if this entry contains the requested attribute
func hasAttr(entry *dwarf.Entry, attr dwarf.Attr) bool {
	for _, field := range entry.Field {
		if field.Attr == attr {
			return true
		}
	}
	return false
}

// Returns the location of an entry in memory
func GetLocation(entry *dwarf.Entry) ([]uint8, error) {
	var err error
	loc := entry.Val(dwarf.AttrLocation)
	if loc == nil {
		err = errors.New(fmt.Sprintf("Could not find data location for %v", entry))
		return nil, err
	}
	return entry.Val(dwarf.AttrLocation).([]uint8), nil
}

// Translates a DW_AT_locationn attribute into an address
func ParseLocation(location []uint8) int {
	if location == nil {
		panic("Cannot parse location for an empty slice!")
	}
	// Ignore the first entry in the slice
	// --> This somehow communicates a format?
	// Build the last slice from right to left
	location = location[1:]
	var locationAsInt int
	locationAsInt = 0
	for i := 0; i < len(location); i++ {
		locationAsInt += int(location[i]) << (8 * i)
	}
	return locationAsInt
}

// Returns the entry defining the type for a given entry. Returns self if
// no such entry can be found. Leaves the reader at the new entry.
func GetTypeEntry(reader *dwarf.Reader, entry *dwarf.Entry) (*dwarf.Entry, error) {

	var err error = nil
	if !hasAttr(entry, dwarf.AttrType) {
		// fmt.Printf("Entry %v does not have a type entry - returning it as-is\n", entry.Val(dwarf.AttrName))
		return entry, nil
	}
	var typeDie *dwarf.Entry
	for _, field := range entry.Field {
		if field.Attr == dwarf.AttrType {
			typeDieOffset := field.Val.(dwarf.Offset)
			reader.Seek(typeDieOffset)
			typeDie, err = reader.Next()
		}
	}
	return typeDie, err
}

// Repeatedly calls GetTypeEntry until arriving at an entry that truly describes
// the underlying type of this entry. Skip over array and other non-typedef
// entries
func ResolveTypeEntry(reader *dwarf.Reader, entry *dwarf.Entry) (*dwarf.Entry, error) {
	typeEntry, err := GetTypeEntry(reader, entry)
	switch typeEntry.Tag {
	case dwarf.TagArrayType:
		return ResolveTypeEntry(reader, typeEntry)
	default:
		return typeEntry, err
	}
}
