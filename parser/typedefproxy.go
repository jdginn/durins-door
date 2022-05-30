package parser

import (
	"debug/dwarf"
	"fmt"
	// "strings"
)

// A TypedefProxy is an outward-facing representation of a typedef representing what a user
// may care about without any DWARF implementation information. This proxy
// represents the key information required to understand the layout of a particular type
// and then read or create variables of this type idiomatically from within a target
// language (most immediately relevant is Go but this should be generic enough to
// apply to other languages through Go bindings or a socket server using json or rpc).
//
// All relevant DWARF parsing is handled when this proxy is created and no intermediate
// DWARF data is included here. The proxy should be ready to hand off as-is to a user.
type TypeDefProxy struct {
	name         string
	bitSize      int
	structOffset int
	arrayRanges  []int
	ahildren     []TypeDefProxy
}

// Construct a new TypeDefProxy
func NewTypeDefProxy(reader *dwarf.Reader, e *dwarf.Entry) (*TypeDefProxy, error) {
	var arrayRanges = []int{0}
	var name string
	var err error
	typeEntry, err := GetTypeEntry(reader, e)

	// Need to handle traversing through array entries to get to the underlying typedefs.
	if typeEntry.Tag == dwarf.TagArrayType {
		arrayRanges, err = GetArrayRanges(reader, e)
		if err != nil {
			return nil, err
		}
		// Having resolved the array information the real type is behind the ArrayType Entry
		// This entry describes the array
		typeEntry, err = GetTypeEntry(reader, typeEntry)
		if err != nil {
			return nil, err
		}
	}

	if typeEntry.Tag == dwarf.TagConstType {
		typeEntry, err = GetTypeEntry(reader, typeEntry)
		if err != nil {
			return nil, err
		}
		name = typeEntry.Val(dwarf.AttrName).(string)
		typeEntry, err = GetTypeEntry(reader, typeEntry)
		if err != nil {
			return nil, err
		}
	} else {
		name = e.Val(dwarf.AttrName).(string)
	}

	// Arrays and Consts may still have a typedef entry behind them. We need to step
	// through it to find the underlying struct or base type
	if typeEntry.Tag == dwarf.TagTypedef {
		typeEntry, err = GetTypeEntry(reader, typeEntry)
	}

	// TODO: handle the situation where we have no AttrName
	proxy := &TypeDefProxy{
		name:         name,
		bitSize:      0,
		structOffset: 0,
		arrayRanges:  arrayRanges,
		ahildren:     make([]TypeDefProxy, 0),
	}

	// The offset into the struct is defined by the member, not its type
	if HasAttr(e, dwarf.AttrDataMemberLoc) {
		proxy.structOffset = int(e.Val(dwarf.AttrDataMemberLoc).(int64)) * 8
	}

	// TODO: this probably needs an else case where we compute size from walking
	// the typedef, which we will do anyway.
	if HasAttr(typeEntry, dwarf.AttrByteSize) || HasAttr(typeEntry, dwarf.AttrBitSize) {
		var bitSize int
		bitSize, err = GetBitSize(typeEntry)
		proxy.bitSize = bitSize
	}

	// TODO: split this into its own method
	if typeEntry.Children {
		for {
			child, err := reader.Next()
			if err != nil {
				panic("Error iterating children; **this error handling needs to be improved!**")
			}

			// When we've finished iterating over members, we are done with the meaningful
			// children of this typedef. We are also finished if we reach the end of the DWARF
			// section during this iteration.
			if child == nil {
				// fmt.Println("Bailing from populating children because we saw the final entry in the DWARF")
				break
			}
			if child.Tag == 0 {
				break
			}

			// Note that constructing proxies for all children makes this constructor
			// itself recursive.
			childProxy, err := NewTypeDefProxy(reader, child)
			if err != nil {
				panic(err)
			}
			// TODO: is this the right way to do this in go?
			proxy.ahildren = append(proxy.ahildren, *childProxy)
			// How do we appropriately parse this stuff without having to jump around a bunch in the reader?
			// ^ Is jumping around in the reader expensive?
			reader.Seek(child.Offset)
			_, err = reader.Next()
			if err != nil {
				panic(err)
			}
		}
	}
	return proxy, err
}

func (p TypeDefProxy) Name() string {
  return p.name
}

func (p *TypeDefProxy) string() string {
	// TODO: for now, we don't print children
	var str string = fmt.Sprintf("Typedef %s\n  BitSize: %d\n  ArrayRanges %v\n  Children %#v\n", p.name, p.bitSize, p.arrayRanges, p.ahildren)
	return str
}

func (p *TypeDefProxy) GoString() string {
	return p.string()
}

// Retuns a slice of strings containing the name of each member of this TypeDef
func (p TypeDefProxy) ListChildren() []string {
	names := make([]string, len(p.ahildren))
	for i, c := range p.ahildren {
		names[i] = c.name
	}
	return names
}

// Returns the TypeDefProxy for a member of this TypeDef by name
func (p TypeDefProxy) GetChild(childName string) (*TypeDefProxy, error) {
	for _, c := range p.ahildren {
		if c.name == childName {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("Could not find child %s for %s", childName, p.GoString())
}
