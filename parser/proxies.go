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
		Name:    e.Val(dwarf.AttrName).(string),
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
	Name         string
	BitSize      int
	StructOffset int
	// TODO: Is it safe for us to put this dwarf TypeDefProxy
	// in an Outward-facing struct? Probably not.
	// DwarfOffset dwarf.Offset
	// TODO: is this the best name for this?
	ArrayRanges []int
	Children    []TypeDefProxy
}

func (p *TypeDefProxy) string() string {
	// TODO: for now, we don't print children
	var str string = fmt.Sprintf("Typedef %s\n  BitSize: %d\n  ArrayRanges %v\n  Children %#v\n", p.Name, p.BitSize, p.ArrayRanges, p.Children)
	return str
}

func (p *TypeDefProxy) GoString() string {
	return p.string()
}

func NewTypeDefProxy(reader *dwarf.Reader, e *dwarf.Entry) *TypeDefProxy {
	typeEntry, _ := GetTypeEntry(reader, e)
	proxy := &TypeDefProxy{
		Name:         e.Val(dwarf.AttrName).(string),
		BitSize:      0,
		StructOffset: 0,
		// DwarfOffset:  e.Offset,
		ArrayRanges:  []int{0},
		Children:     make([]TypeDefProxy, 0),
	}

	// TODO: this probably needs an else case where we compute size from walking
	// the typedef, which we will do anyway.
	if hasAttr(typeEntry, dwarf.AttrByteSize) || hasAttr(typeEntry, dwarf.AttrBitSize) {
		proxy.BitSize = GetBitSize(typeEntry)
	}

	// Need to handle traversing *through* array entries to get to the underlying
	// typedefs. This is hard so for now we just ignore them (terrible hack)
	if typeEntry.Tag == dwarf.TagArrayType {
		typeEntry, _ = GetTypeEntry(reader, typeEntry)
		fmt.Println("Jumping through TagArrayType to get to:")
		PrintEntryInfo(typeEntry)
	}

	fmt.Println("Parsing typedef for:")
	PrintEntryInfo(typeEntry)
	if typeEntry.Children {
		for {
			child, err := reader.Next()
			fmt.Println("Next child:")
			PrintEntryInfo(child)
			if err != nil {
				fmt.Println("Error iterating children; **this error handling needs to be improved!**")
			}

			// When we've finished iterating over members, we are done with the meaningful
			// children of this typedef. We are also finished if we reach the end of the DWARF
			// section during this iteration.
			if child == nil {
				fmt.Println("Bailing from populating children because we saw the final entry in the DWARF")
				break
			}
			if child.Tag == 0 {
				fmt.Println("Bailing from populating children because we found a null entry")
				break
			}

			// Note that constructing proxies for all children makes this constructor
			// itself recursive.
			childProxy := NewTypeDefProxy(reader, child)
			fmt.Printf("%#v\n", childProxy)
			// TODO: is this the right way to do this in go?
			proxy.Children = append(proxy.Children, *childProxy)
			// How do we appropriately parse this stuff without having to jump around a bunch in the reader?
			// ^ Is jumping around in the reader expensive?
			reader.Seek(child.Offset)
			reader.Next()
		}
	}
	return proxy
}

type VariableProxy struct {
	Type    TypeDefProxy
	Address int
	Value   int
}
