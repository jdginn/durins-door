package explorer

import (
	"debug/dwarf"
	"fmt"
	"strings"

	"github.com/jdginn/durins-door/client"
	"github.com/jdginn/durins-door/explorer/plat"
	"github.com/jdginn/durins-door/parser"
)

// Struct that mediates DWARF parsing as well as reading and writing
type Explorer struct {
	readerFile string
	reader     *dwarf.Reader
	client     client.Client
	ctx        explorerCtx
}

func NewExplorer() *Explorer {
	return &Explorer{
		ctx: explorerCtx{},
	}
}

func (c Explorer) jumpTo() {
  if len(c.ctx.levels) > 0 {
    c.reader.Seek(c.ctx.levels[len(c.ctx.levels)-1].Offset)
  }
}

// Creates a reader within this explorer, reading the specified file
func (e *Explorer) CreateReaderFromFile(fname string) error {
	e.readerFile = fname
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

// Returns the name of the DWARF file this explorer is reading
func (e *Explorer) GetReaderFilename() string {
	return e.readerFile
}

// Sets this explorer's client
func (e *Explorer) SetClient(c client.Client) error {
	e.client = c
	return nil
}

func (e *Explorer) GetTypeDefProxy(name string) (*parser.TypeDefProxy, error) {
	if e.reader == nil {
		return nil, fmt.Errorf("Cannot get TypeDef proxies without setting a reader. Create a reader using CreateReaderFromFile().")
	}
	levels := strings.Split(name, "/")
	entry, _, err := parser.GetEntry(e.reader, levels[0])
	if err != nil {
		return nil, err
	}
	p, err := parser.NewTypeDefProxy(e.reader, entry)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (e *Explorer) ShowAllChildren() ([]string, error) {
	if e.reader == nil {
		return nil, fmt.Errorf("Cannot List CUs without setting a reader. Create a reader using CreateReaderFromFile().")
	}
	entries, err := parser.GetChildren(e.reader, func(entry *dwarf.Entry) bool {
		return (entry.Tag == dwarf.TagVariable || entry.Tag == dwarf.TagCompileUnit)
	})
	if err != nil {
		return []string{}, err
	}
	ret := make([]string, len(entries), len(entries))
	for i, e := range entries {
		ret[i] = e.Val(dwarf.AttrName).(string)
	}
	return ret, nil
}

func (e *Explorer) StepIntoChild(childName string) {
	entry, _, err := parser.GetEntry(e.reader, childName)
	if err != nil {
		panic(err)
	}
	e.ctx.Push(entry)
}

func (e *Explorer) Up() {
	e.ctx.Pop()
	e.jumpTo()
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
