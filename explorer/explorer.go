package explorer

import (
	"debug/dwarf"
	"strings"

	"github.com/jdginn/dwarf-explore/client"
	"github.com/jdginn/dwarf-explore/explorer/plat"
	"github.com/jdginn/dwarf-explore/parser"
)

// Struct that mediates DWARF parsing as well as reading and writing
type Explorer struct {
	reader *dwarf.Reader
	client client.Client
}

func NewExplorer() *Explorer {
	return &Explorer{}
}

func (e *Explorer) GetReaderFromFile(fname string) error {
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

func (e *Explorer) SetClient(c client.Client) error {
	e.client = c
	return nil
}

func (e *Explorer) GetTypeDefProxy(name string) (*parser.TypeDefProxy, error) {
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

// Returns a VariableProxy to work with
func (e *Explorer) GetVariableProxy(name string) (*parser.VariableProxy, error) {
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
