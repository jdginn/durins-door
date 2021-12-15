package parser

import (
  "debug/dwarf"
  "fmt"
)

// Description of
type TypeEntryProxy struct {
  Name string
  Offset int
  BitSize int
  // entry *dwarf.Entry
  // children map[string]TypedefProxy
}

func NewTypeEntryProxy(reader *dwarf.Reader, e *dwarf.Entry) *TypeEntryProxy {
  fmt.Println(e.AttrField(dwarf.AttrName).Val)
  typeEntry, _ := GetTypeEntry(reader, e)
  proxy := &TypeEntryProxy{
    Name: e.AttrField(dwarf.AttrName).Val.(string),
    Offset: int(typeEntry.Offset),
    BitSize: GetBitSize(typeEntry),
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
  Name string
  Address int
  Size int
  entry *dwarf.Entry
  children map[string]VariableProxy
}
