package parser

import (
	"debug/dwarf"
	"fmt"
)

// Description of
type TypeEntryProxy struct {
	Name    string
	Offset  int
	BitSize int
}

func NewTypeEntryProxy(reader *dwarf.Reader, e *dwarf.Entry) *TypeEntryProxy {
	typeEntry, _ := GetTypeEntry(reader, e)
	proxy := &TypeEntryProxy{
		Name:    e.AttrField(dwarf.AttrName).Val.(string),
		Offset:  int(typeEntry.Offset),
		BitSize: GetBitSize(typeEntry),
	}
	return proxy
}

// Outward-facing representation of a typedef representing what a user
// may care about without any DWARF implementation information. This proxy
// represents the key information required to understand the layout of a particular type
// type and then read or create variables of this type idiomatically from within a target
// language (most immediately relevant is Go but this should be generic enough to apply
// to other languages through Go bindings).
//
// All relevant DWARF parsing is handled when this proxy is created and no intermediate
// DWARF data is included here. The proxy should be ready to hand off as-is to a user.
type TypeDefProxy struct {
	Name        string
	BitSize     int
	DwarfOffset int
	// TODO: is this the best name for this?
	ArrayRanges []int
	Children    []TypeDefProxy
}

func NewTypeDefProxy(reader *dwarf.Reader, e *dwarf.Entry) *TypeDefProxy {
	typeEntry, _ := GetTypeEntry(reader, e)
	proxy := &TypeDefProxy{
		Name:        e.AttrField(dwarf.AttrName).Val.(string),
		BitSize:     0,
		DwarfOffset: 0,
		ArrayRanges: []int{0},
		Children:    make([]TypeDefProxy, 0),
	}

	// TODO: this probably needs an else case where we compute size from walking
	// the typedef, which we will do anyway.
	if hasAttr(typeEntry, dwarf.AttrByteSize) || hasAttr(typeEntry, dwarf.AttrBitSize) {
		proxy.BitSize = GetBitSize(typeEntry)
	}

  fmt.Println("Parsing typedef for:")
	PrintEntryInfo(typeEntry)
	if typeEntry.Children {
		for {
			child, err := reader.Next()
			if err != nil {
				fmt.Println("Error iterating children; **this error handling needs to be improved!**")
			}

			// When we've finished iterating over members, we are done with the meaningful
			// children of this typedef. We are also finished if we reach the end of the DWARF
			// section during this iteration.
			if (child == nil) {
        fmt.Println("Bailing from populating children because we saw the final entry in the DWARF")
				break
			}
			if (child.Tag == 0) {
        fmt.Println("Bailing from populating children because we found a null entry")
				break
			}

			// Note that constructing proxies for all children makes this constructor
			// itself recursive.

      // TODO: the problem here is that we are resolving the type of this inside, which does not work
      // upon recursive calls
			childProxy := NewTypeDefProxy(reader, child)
      // TODO: is this the right way to do this in go?
			proxy.Children = append(proxy.Children, *childProxy)
		}
	}

	return proxy
}

type VariableProxy struct {
  Type TypeDefProxy
	Address  int
  Value int
}
