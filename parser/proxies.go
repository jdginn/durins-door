package parser

import (
	"debug/dwarf"
	// "fmt"
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
	Name     string
	BitSize  int
  DwarfOffset int
  // TODO: is this the best name for this?
  ArrayRanges []int
	Children map[string]TypeEntryProxy
}

func NewTypeDefProxy(reader *dwarf.Reader, e *dwarf.Entry) *TypeDefProxy {
	typeEntry, _ := GetTypeEntry(reader, e)
	proxy := &TypeDefProxy{
		Name:     e.AttrField(dwarf.AttrName).Val.(string),
		BitSize:  0,
    DwarfOffset: 0,
    ArrayRanges: []int{0},
		Children: make(map[string]TypeEntryProxy),
	}

	// TODO: this probably needs an else case where we compute size from walking
  // the typedef, which we will do anyway.
	if hasAttr(typeEntry, dwarf.AttrByteSize) || hasAttr(typeEntry, dwarf.AttrBitSize) {
		proxy.BitSize = GetBitSize(typeEntry)
	}

	return proxy
}

// TODO: this should probably be a new thing called TypeDefProxy, representing the full TypeDef
// // Traverse the entire type hierarchy underneath this
// // entry to populate the `children` maps at each level
// func (t *TypedieProxy) Populate() *TypedieProxy {
//   //TODO
//   return t
// }

// func (t *TypedieProxy) Flatten() *TypedieProxy {
//   //TODO: what does this return?
//   return t
// }

type VariableProxy struct {
	Name     string
	Address  int
	Size     int
	entry    *dwarf.Entry
	children map[string]VariableProxy
}
