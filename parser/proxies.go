package parser

import (
	"debug/dwarf"
	"fmt"
	// "strings"

	"github.com/jdginn/durins-door/client"
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
	Name         string
	BitSize      int
	StructOffset int
	ArrayRanges  []int
	Children     []TypeDefProxy
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
		Name:         name,
		BitSize:      0,
		StructOffset: 0,
		ArrayRanges:  arrayRanges,
		Children:     make([]TypeDefProxy, 0),
	}

	// The offset into the struct is defined by the member, not its type
	if hasAttr(e, dwarf.AttrDataMemberLoc) {
		proxy.StructOffset = int(e.Val(dwarf.AttrDataMemberLoc).(int64)) * 8
	}

	// TODO: this probably needs an else case where we compute size from walking
	// the typedef, which we will do anyway.
	if hasAttr(typeEntry, dwarf.AttrByteSize) || hasAttr(typeEntry, dwarf.AttrBitSize) {
		var bitSize int
		bitSize, err = GetBitSize(typeEntry)
		proxy.BitSize = bitSize
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
			proxy.Children = append(proxy.Children, *childProxy)
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

func (p *TypeDefProxy) string() string {
	// TODO: for now, we don't print children
	var str string = fmt.Sprintf("Typedef %s\n  BitSize: %d\n  ArrayRanges %v\n  Children %#v\n", p.Name, p.BitSize, p.ArrayRanges, p.Children)
	return str
}

func (p *TypeDefProxy) GoString() string {
	return p.string()
}

func (p *TypeDefProxy) GetChild(childName string) (*TypeDefProxy, error) {
	for _, c := range p.Children {
		if c.Name == childName {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("Could not find child %s for %s", childName, p.GoString())
}

// Represents a variable and facilitates interacting with that
// variable. For our purposes, a variable is an object of known Type
// located at a known address, whose value we can find by reading
// a known size from that address. Alternatively, we can set the value
// of this variables, or its members if it is a struct, and then write
// those values back to memory.
//
// Data is stored internally as bytes and parse into fields on demand.
//
// Writing data to memory is handled elsewhere; this proxy instructs
// a client which addresses to read and provides a writeable stream
// of bytes to allow the client to write the variable back to memory.
type VariableProxy struct {
	Name    string
	Type    TypeDefProxy
	Address int
	value   []byte
	client  client.Client
}

// Construct a new VariableProxy for a variable known to the DWARF
// debug info.
//
// To create a variable from scratch , use *some other method*
func NewVariableProxy(reader *dwarf.Reader, entry *dwarf.Entry) (*VariableProxy, error) {
	typeDefProxy, err := NewTypeDefProxy(reader, entry)
	if err != nil {
		return nil, err
	}
	loc, err := GetLocation(entry)
	proxy := &VariableProxy{
		Name:    entry.Val(dwarf.AttrName).(string),
		Type:    *typeDefProxy,
		Address: ParseLocation(loc),
		value:   []byte{},
		client:  nil,
	}
	return proxy, err
}

func (p *VariableProxy) Init(reader *dwarf.Reader, entry *dwarf.Entry) error {
	typeDefProxy, err := NewTypeDefProxy(reader, entry)
	if err != nil {
		return err
	}
	loc, err := GetLocation(entry)
	p.Name = entry.Val(dwarf.AttrName).(string)
	p.Type = *typeDefProxy
	p.Address = ParseLocation(loc)
	return err
}

// TODO: change the child hierarchy to use ordered maps not slices for lookup speed?
func (p *VariableProxy) GetChild(childName string) (*TypeDefProxy, error) {
	var err error
	typeDef := p.Type
	for _, child := range typeDef.Children {
		if child.Name == childName {
			return &child, err
		}
	}
	return nil, fmt.Errorf("Could not find child %s for %s", childName, p.GoString())
}

func (p *VariableProxy) string() string {
	var str string = fmt.Sprintf("Variable %s at address %x of type:\n%v", p.Name, p.Address, p.Type.string())
	return str
}

func (p *VariableProxy) GoString() string {
	return p.string()
}

// Set the value of this entire variable
//
// In the case of a multi-field struct, this is most useful for
// initializing our proxy of a variable by having the client read
// the entire variable out of memory. Once the proxy is poplulated,
// we can access fields as required.
func (p *VariableProxy) Set(value []byte) error {
	var err error = nil
	if len(value)*8 > p.Type.BitSize {
		err = fmt.Errorf("Attempted to set value size %d bits, larger than type with size %d bits", len(value)*8, p.Type.BitSize)
	}
	p.value = value
	return err
}

// Set the value of a single field within this variable
// Typically this will be used for a struct or class
//
// NOTE: at present, fields must be byte-aligned
func (p *VariableProxy) SetField(field string, value int) error {
	// TODO: what if the field is not byte-aligned?
	fieldEntry, err := p.GetChild(field)
	if err != nil {
		return err
	}
	startIndex := fieldEntry.StructOffset / 8
	n := fieldEntry.BitSize / 8
	// TODO: surely there is a mroe elegant way
	if fieldEntry.BitSize%8 != 0 {
		n += 1
	}
	for i := 0; i < n; i++ {
		p.value[startIndex+i] = byte(value >> ((n - i - 1) * 8) & 0xff)
	}
	return err
}

// Return the value of the entire variable
//
// In the case of a multi-field struct, this is most useful to
// enable the client to write the entire variable back to memory.
func (p *VariableProxy) Get() ([]byte, error) {
	return p.value, nil
}

// // Return the value of a single field within this variable
// //
// // Take a path to the desired field through arbitrary levels
// // in the struct hierarchy. The hierarchy is formatted as it
// // would be in C. Struct members are delimited by dots and
// // and array indices are delimted by brackets. For example:
// //
// // myProxy.GetField("thisMember[3].thatMember.theOtherMember[2]")
// func (p *VariableProxy) GetField(field string) (int, error) {
//   // Look for members
//   members := strings.Split(field, ".")
//   for _, m := range members {
//     split := strings.SplitAfterN(m, "[", 2)
//     child, err := p.GetChild(split[0])
//     if err != nil {
//       return 0, err
//     }
//     if len(split) == 2 {
//       // We have an array to index into
//       // TODO: fix
//       index := strings.SplitAfterN(split[1], "]", 2)[0]
//       child = child[index]
//     // TODO: multidemensional arrays
//     } else if len(split) > 2 {
//       return 0, fmt.Errorf("Parsed too many array indices out of path")
//     }
//   }
// }

// Return the value of a single field within this variable
// Typically this will be used for a struct or class
//
// NOTE: at present, fields must be byte-aligned
func (p *VariableProxy) GetField(field string) (int, error) {
	fieldEntry, err := p.GetChild(field)
	if err != nil {
		return 0, err
	}
	if p.value == nil {
		err = fmt.Errorf("Proxy has no internal data to get")
	}
	startByte := fieldEntry.StructOffset / 8
	byteLen := fieldEntry.BitSize / 8
	endByte := startByte + byteLen - 1
	if len(p.value) < (startByte + byteLen) {
		err = fmt.Errorf("Internal data len %d bytes is smaller than the requested field %s at bytes %d:%d", len(p.value), fieldEntry.Name, startByte, endByte)
		return 0, err
	}
	valInt := 0
	for i := 0; i < byteLen; i++ {
		b := p.value[startByte+i]
		shiftAmt := (byteLen - i - 1) * 8
		valInt = valInt + int(b)<<shiftAmt
	}
	return valInt, err
}

func (p *VariableProxy) SetClient(c client.Client) {
	p.client = c
}

func (p *VariableProxy) Read() error {
	if p.client == nil {
		return fmt.Errorf("Cannot read proxy %s: no client is set!", p.string())
	}
	// TODO: what if this isn't byte-aligned?
	data, err := p.client.Read(p.Address, p.Type.BitSize/8)
	if err != nil {
		return err
	}
	p.Set(data)
	return nil
}

func (p *VariableProxy) Write() error {
	if p.client == nil {
		return fmt.Errorf("Cannot write proxy %s: no client is set!", p.string())
	}
	return p.client.Write(p.Address, p.value)
}
