package explorer

import (
	"debug/dwarf"
	"fmt"
	// "log"

	"github.com/jdginn/durins-door/client"
	"github.com/jdginn/durins-door/explorer/plat"
	"github.com/jdginn/durins-door/parser"
)

// Struct that mediates DWARF parsing as well as reading and writing
type Explorer struct {
	DwarfFile string
	reader    *dwarf.Reader
	client    client.Client
	ctx       *stack
}

// Returns a new explorer struct with sane defaults
func NewExplorer() *Explorer {
	return &Explorer{
		ctx: NewStack(),
	}
}

// Creates a reader within this explorer, reading the specified file
func (e *Explorer) CreateReaderFromFile(fname string) error {
	e.DwarfFile = fname
	fh, err := plat.GetReaderFromFile(fname)
	if err != nil {
		return err
	}
	reader, err := parser.GetReader(fh)
	if err != nil {
		return err
	}
	e.reader = reader
	return nil
}

// Returns a new explorer with reader set to the specified file
func NewExplorerFromFile(fname string) *Explorer {
	e := NewExplorer()
	err := e.CreateReaderFromFile(fname)
	if err != nil {
		panic(err)
	}
	return e
}

// Returns a slice containing the names of each child of this Entry
func (e *Explorer) listEntryChildren() []string {
	entries, err := parser.GetChildren(e.reader, func(entry *dwarf.Entry) bool {
		return (entry.Tag == dwarf.TagVariable || entry.Tag == dwarf.TagCompileUnit)
	})
	if err != nil {
		return []string{}
	}
	ret := make([]string, len(entries), len(entries))
	for i, e := range entries {
		ret[i] = e.Val(dwarf.AttrName).(string)
	}
	return ret
}

// Returns a slice containing the names of each child of the current item
func (e *Explorer) ListChildren() []string {
	switch e.ctx.CurrMode() {
	case modeCUs:
		s, _ := e.ListCUs()
		return s
	case modeEntry:
		return e.listEntryChildren()
	case modeProxy:
		return e.ctx.CurrProxy().ListChildren()
	default:
		return []string{"No children to display"}
	}
}

// Returns a string representation of the curent mode
func (e *Explorer) CurrMode() string {
	switch e.ctx.CurrMode() {
	case modeCUs:
		return "modeCUs"
	case modeEntry:
		return "modeEntry"
	case modeProxy:
		return "modeProxy"
	}
	return "bad mode"
}

// Returns the name of the current item
func (e *Explorer) CurrName() string {
	switch e.ctx.CurrMode() {
	case modeCUs:
		return "all CUs"
	case modeEntry:
		return e.ctx.CurrEntry().Val(dwarf.AttrName).(string)
	case modeProxy:
		return e.ctx.CurrProxy().Name()
	default:
		return "bad mode"
	}
}

// Moves the context to the specified child of the current item
//
// Child is specified by name
func (e *Explorer) StepIntoChild(childName string) error {
	switch e.ctx.CurrMode() {
	case modeCUs:
		entry, _, err := parser.GetEntry(e.reader, childName)
		if err != nil {
			return err
		}
		e.ctx.Push(modeEntry, entry, nil)
		return nil
	case modeEntry:
		entry, _, err := parser.GetEntry(e.reader, childName)
		if err != nil {
			return err
		}
		p, err := e.getProxy(entry)
		if err != nil {
			return err
		}
		e.ctx.Push(modeProxy, nil, p)
		return nil
	case modeProxy:
		p, err := e.ctx.CurrProxy().GetChild(childName)
		if err != nil {
			return err
		}
		e.ctx.Push(modeProxy, nil, p)
		return nil
	default:
		return fmt.Errorf("Invalid mode")
	}
}

// Moves the context to the previous item
func (e *Explorer) Back() error {
	e.ctx.Pop()
	if e.ctx.CurrMode() == modeEntry {
		e.reader.Seek(e.ctx.CurrEntry().Offset)
	}
	return nil
}

// Moves the context to the applicable TypeDef proxy
//
// Creates a TypeDefProxy from the current item if it is either an
// entry or a VariableProxy. No action if the current item is already a
// TypeDefProxy.
func (e *Explorer) GetType() error {
	switch e.ctx.CurrProxy().(type) {
	// If we are already looking at a typeDef, there is nothing to do
	case parser.TypeDefProxy:
	case parser.VariableProxy:
		e.ctx.Push(modeProxy, nil, e.ctx.CurrProxy().(parser.VariableProxy).Type)
	}
	return nil
}

// Creates the proxy corresponding to the passed entry
func (e *Explorer) getProxy(entry *dwarf.Entry) (parser.Proxy, error) {
	switch entry.Tag {
	case dwarf.TagVariable:
		return parser.NewVariableProxy(e.reader, entry)
	case dwarf.TagTypedef:
		return parser.NewTypeDefProxy(e.reader, entry)
	default:
		return nil, fmt.Errorf("Invalid tag %s for entry %s", entry.Tag.String(), parser.FormatEntryInfo(entry))
	}
}

func (e *Explorer) Up() bool {
	panic("explorer.Up() not implemented yet")
}

// Returns a string representing key info about the current entry, if there is one
func (e *Explorer) Info() string {
	switch e.ctx.CurrMode() {
	case modeCUs, modeEntry:
		entry := e.ctx.CurrEntry()
		return parser.FormatEntryInfo(entry)
	default:
		return ""
	}
}

// Returns a list of all CUs in this file
func (e *Explorer) ListCUs() ([]string, error) {
	if e.reader == nil {
		return nil, fmt.Errorf("Cannot List CUs without setting a reader. Create a reader using CreateReaderFromFile().")
	}
	CUs, err := parser.GetCUs(e.reader)
	if err != nil {
		return []string{}, err
	}
	ret := make([]string, len(CUs), len(CUs))
	for i, cu := range CUs {
		ret[i] = cu.Val(dwarf.AttrName).(string)
	}
	return ret, nil
}
