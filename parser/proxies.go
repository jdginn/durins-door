package parser

import (
	"debug/dwarf"
	"fmt"
)

// Relatively toothless sibling of TypeDefProxy
type TypeEntryProxy struct {
	Name    string
	Offset  int
	BitSize int
}

func NewTypeEntryProxy(reader *dwarf.Reader, e *dwarf.Entry) (*TypeEntryProxy, error) {
	typeEntry, _ := GetTypeEntry(reader, e)
	bitSize, err := GetBitSize(typeEntry)
	if err != nil {
		return &TypeEntryProxy{}, err
	}
	proxy := &TypeEntryProxy{
		Name:    e.Val(dwarf.AttrName).(string),
		Offset:  int(typeEntry.Offset),
		BitSize: bitSize,
	}
	return proxy, err
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
	ArrayRanges  []int
	Children     []TypeDefProxy
}

func NewTypeDefProxy(reader *dwarf.Reader, e *dwarf.Entry) (*TypeDefProxy, error) {
  var arrayRanges = []int{0}
	var err error = nil
  fmt.Println("NewTypeDefProxy")
	fmt.Println(FormatEntryInfo(e))
	typeEntry, _ := GetTypeEntry(reader, e)

  if typeEntry.Tag == dwarf.TagArrayType{
    arrayRanges, err = GetArrayRanges(reader, typeEntry)
    typeEntry, err = GetTypeEntry(reader, typeEntry)
  }

  fmt.Println("Type entry:")
	fmt.Println(FormatEntryInfo(typeEntry))

  // TODO: handle the situation where we have no AttrName
	proxy := &TypeDefProxy{
		Name:         e.Val(dwarf.AttrName).(string),
		BitSize:      0,
		StructOffset: 0,
		ArrayRanges:  arrayRanges,
		Children:     make([]TypeDefProxy, 0),
	}

	// The offset into the struct is defined by the member, not its type
	if hasAttr(e, dwarf.AttrDataMemberLoc) {
		proxy.StructOffset = int(e.Val(dwarf.AttrDataMemberLoc).(int64)) * 8
	}

	// Need to handle traversing through array entries to get to the underlying typedefs.
	if typeEntry.Tag == dwarf.TagArrayType {
		ranges, _ := GetArrayRanges(reader, e)
		proxy.ArrayRanges = ranges
		// Having resolved the array information the real type is behind the ArrayType Entry
		// This entry describes the array
		typeEntry, _ = GetTypeEntry(reader, typeEntry)
		// This entry describes the type of object the array is made of
		typeEntry, _ = GetTypeEntry(reader, typeEntry)
	}

	// TODO: this probably needs an else case where we compute size from walking
	// the typedef, which we will do anyway.
	if hasAttr(typeEntry, dwarf.AttrByteSize) || hasAttr(typeEntry, dwarf.AttrBitSize) {
		var bitSize int
		bitSize, err = GetBitSize(typeEntry)
		proxy.BitSize = bitSize
	}

	// fmt.Println("Parsing type entry:")
	// fmt.Println(FormatEntryInfo(typeEntry))
	if typeEntry.Children {
		for {
			child, err := reader.Next()
			if err != nil {
				fmt.Println("Error iterating children; **this error handling needs to be improved!**")
			}

			// When we've finished iterating over members, we are done with the meaningful
			// children of this typedef. We are also finished if we reach the end of the DWARF
			// section during this iteration.
			if child == nil {
				// fmt.Println("Bailing from populating children because we saw the final entry in the DWARF")
				break
			}
			if child.Tag == 0 {
				fmt.Println("Bailing from populating children because we found a null entry")
				break
			}

			// Note that constructing proxies for all children makes this constructor
			// itself recursive.
			childProxy, err := NewTypeDefProxy(reader, child)
			// TODO: is this the right way to do this in go?
			proxy.Children = append(proxy.Children, *childProxy)
			// How do we appropriately parse this stuff without having to jump around a bunch in the reader?
			// ^ Is jumping around in the reader expensive?
			reader.Seek(child.Offset)
			reader.Next()
		}
	}
	return proxy, err
}

func (p *TypeDefProxy) string() string {
	// TODO: for now, we don't print children
	var str string = fmt.Sprintf("Typedef %s\n  BitSize: %d\n  ArrayRanges %v\n  Children %#v\n", p.Name, p.BitSize, p.ArrayRanges, p.Children)
	return str
}

func (p *TypeDefProxy) GoString() string {
	return p.string()
}

type VariableProxy struct {
	// TODO: split this out from the notion of having children?
	// i.e. define type information element-wise rather than wholesale
	Name     string
	Type     TypeDefProxy
	Address  int64
	Value    int64
	Children []VariableProxy
}

// Construct a new VariableProxy
func NewVariableProxy(reader *dwarf.Reader, entry *dwarf.Entry) (*VariableProxy, error) {
	typeDefProxy, err := NewTypeDefProxy(reader, entry)
  loc, err := GetLocation(entry)
	proxy := &VariableProxy{
		Name:    entry.Val(dwarf.AttrName).(string),
		Type:    *typeDefProxy,
		Address: ParseLocation(loc),
		Value:   0,
	}
	return proxy, err
}

func NewVariableProxyFromTypedef(typeDef TypeDefProxy) *VariableProxy {
	proxy := &VariableProxy{
		Type:    typeDef,
		Address: 0,
		Value:   0,
	}
	return proxy
}
