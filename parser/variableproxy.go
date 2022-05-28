package parser

import (
	"debug/dwarf"
	"fmt"
	// "strings"

	"github.com/jdginn/durins-door/client"
)

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
