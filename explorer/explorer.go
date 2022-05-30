package explorer

import (
	"debug/dwarf"
	"fmt"
	// "log"
	"strings"

	"github.com/jdginn/durins-door/client"
	"github.com/jdginn/durins-door/explorer/plat"
	"github.com/jdginn/durins-door/parser"
)

// Struct that mediates DWARF parsing as well as reading and writing
type Explorer struct {
	DwarfFile string
	reader    *dwarf.Reader
	client    client.Client
	ctx       stack
}

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

func NewExplorerFromFile(fname string) *Explorer {
	e := NewExplorer()
	err := e.CreateReaderFromFile(fname)
	if err != nil {
		panic(err)
	}
	return e
}

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

func (e *Explorer) Back() error {
	e.ctx.Pop()
	if e.ctx.CurrMode() == modeEntry {
		e.reader.Seek(e.ctx.CurrEntry().Offset)
	}
	return nil
}

func (e *Explorer) GetType() error {
	switch e.ctx.CurrProxy().(type) {
	// If we are already looking at a typeDef, there is nothing to do
	case parser.TypeDefProxy:
	case parser.VariableProxy:
		e.ctx.Push(modeProxy, nil, e.ctx.CurrProxy().(parser.VariableProxy).Type)
	}
	return nil
}

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

// func (e *Explorer) GetTypeDefProxy() (*parser.TypeDefProxy, error) {
// 	if e.reader == nil {
// 		return nil, fmt.Errorf("Cannot get TypeDef proxies without setting a reader. Create a reader using CreateReaderFromFile().")
// 	}
// 	entry := e.ctx.CurrEntry()
// 	p, err := parser.NewTypeDefProxy(e.reader, entry)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return p, nil
// }

func (e *Explorer) Up() bool {
	c, ok := e.ctx.Pop()
	if ok {
		e.reader.Seek(c.entry.Offset)
	}
	return ok
}

func (e *Explorer) Info() string {
	entry := e.ctx.CurrEntry()
	return parser.FormatEntryInfo(entry)

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

// Returns a VariableProxy to work with
func (e *Explorer) GetVariableProxy(name string) (*parser.VariableProxy, error) {
	if e.reader == nil {
		return nil, fmt.Errorf("Cannot get Variable proxies without setting a reader. Create a reader using CreateReaderFromFile().")
	}
	levels := strings.Split(name, "/")
	entry, cu, err := parser.GetEntry(e.reader, levels[0])
	offset := int64(cu.AttrField(dwarf.AttrLowpc).Val.(uint64))
	e.client.SetOffset(offset)
	if err != nil {
		return nil, err
	}
	p, err := parser.NewVariableProxy(e.reader, entry)
	if err != nil {
		return nil, err
	}
	p.SetClient(e.client)
	return p, nil
}
