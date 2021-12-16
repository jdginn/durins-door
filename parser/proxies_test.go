package parser

import (
	"debug/dwarf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTypeEntryProxy(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
	var p *TypeEntryProxy
	var e *dwarf.Entry

	e, _ = GetEntry(reader, "Driver")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, p.Name, "Driver")
	assert.Equal(t, p.BitSize, int(12*8))

	e, _ = GetEntry(reader, "char")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, p.Name, "char")
	assert.Equal(t, p.BitSize, int(8))

	// Make sure we can get the same proxy twice
	e, _ = GetEntry(reader, "Driver")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, p.Name, "Driver")
	assert.Equal(t, p.BitSize, int(12*8))

	e, _ = GetEntry(reader, "char")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, p.Name, "char")
	assert.Equal(t, p.BitSize, int(8))
}

func TestNewTypeDefProxy(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
	var p *TypeDefProxy
	var e *dwarf.Entry

	// Start with a few trivial cases
	e, _ = GetEntry(reader, "Driver")
	p = NewTypeDefProxy(reader, e)
	assert.Equal(t, p.Name, "Driver")
	assert.Equal(t, p.BitSize, int(12*8))
	assert.Equal(t, p.Children, make(map[string]TypeEntryProxy))

	e, _ = GetEntry(reader, "char")
	p = NewTypeDefProxy(reader, e)
	assert.Equal(t, p.Name, "char")
	assert.Equal(t, p.BitSize, int(8))
	assert.Equal(t, p.Children, make(map[string]TypeEntryProxy))

	// Move on to non-trivial cases in which Children must actually be populated
	e, _ = GetEntry(reader, "Driver")
	p = NewTypeDefProxy(reader, e)
	var driverChildren = map[string]TypeDefProxy{
		"initials": {
			Name:        "initials",
			BitSize:     8,
			DwarfOffset: 0,
			ArrayRanges: []int{2},
			Children:    make(map[string]TypeEntryProxy),
		},
		"car_number": {
			Name:        "car_number",
			BitSize:     32,
			DwarfOffset: 16,
			ArrayRanges: []int{0},
			Children:    make(map[string]TypeEntryProxy),
		},
		"has_won_wdc": {
			Name:        "has_won_wdc",
			BitSize:     8,
			DwarfOffset: 48,
			ArrayRanges: []int{0},
			Children:    make(map[string]TypeEntryProxy),
		},
	}
	assert.Equal(t, p.Name, "Driver")
	assert.Equal(t, p.BitSize, int(12*8))
	assert.Equal(t, p.Children, driverChildren)
}
