package parser

import (
  "debug/dwarf"
	"testing"
  "github.com/stretchr/testify/assert"
)

func TestNewTypeEntryProxy(t *testing.T) {
  reader, _ := GetReader(testcaseFilename)
  var p *TypeEntryProxy
  var e *dwarf.Entry
  e, _ = GetEntry(reader, "Driver")
  p = NewTypeEntryProxy(reader, e)
  assert.Equal(t, p.Name, "Driver")
  assert.Equal(t, p.BitSize, int(12 * 8))
  e, _ = GetEntry(reader, "char")
  p = NewTypeEntryProxy(reader, e)
  assert.Equal(t, p.Name, "char")
  assert.Equal(t, p.BitSize, int(8))
}
