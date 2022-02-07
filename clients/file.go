package file

import (
  "os"
	"fmt"
  
  "github.com/jdginn/dwarf-explore"
)

type ProxyWrapper struct {
  proxy *parser.VariableProxy
  rw *os.File
}

// TODO: this needs a lot of work. From the files, construct everything
// we need for the NewVariableProxy and then read its contents
// from binFile.
//
// Encapsulate anything having to do with dwarf.
// Still take filenames (for compatibility) but break out the handling
// of those so it's not so clunky and they don't need to appear
// in this package (still make sure they don't have to appear in clients
// of this package)
//
// Keep in mind, the next step here is building a server
func NewVariableProxy(dwarfFile parser.DebugFile, binFile *os.File, entryName string) (*ProxyWrapper, error) {
  dwarfReader, err := parser.GetReader(dwarfFile)
  entry, err := parser.GetEntry(dwarfReader, entryName)
  if err != nil {
    return nil, err
  }
  newProxy, err := parser.NewVariableProxy(dwarfReader, entry)
  p := &ProxyWrapper{
    proxy: newProxy,
    rw: binFile,
  }
  return p, err
}

func (p *ProxyWrapper) Write() error { 
  m := p.proxy.GetAccessMetadata()
  value, err := p.proxy.Get()
  if err != nil {
    return err
  }
  p.rw.WriteAt(value, int64(m.Address))
  return err
}

func (p *ProxyWrapper) Read() error { 
  m := p.proxy.GetAccessMetadata()
  value := make([]byte, m.Size)
  n, err := p.rw.ReadAt(value, int64(m.Address))
  if err != nil {
    return err
  }
  if n != m.Size {
    err = fmt.Errorf("Read the incorrect number of bytes\n Expected: %d bytes; Read %d", m.Size, n)
  }
  p.proxy.Set(value)
  return err
}
