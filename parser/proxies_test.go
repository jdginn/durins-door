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
	assert.Equal(t, "Driver", p.Name)
	assert.Equal(t, int(12*8), p.BitSize)

	e, _ = GetEntry(reader, "char")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, "char", p.Name)
	assert.Equal(t, int(8), p.BitSize)

	// Make sure we can get the same proxy twice
	e, _ = GetEntry(reader, "Driver")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, "Driver", p.Name)
	assert.Equal(t, int(12*8), p.BitSize)

	e, _ = GetEntry(reader, "char")
	p = NewTypeEntryProxy(reader, e)
	assert.Equal(t, "char", p.Name)
	assert.Equal(t, int(8), p.BitSize)
}

// func TestNewTypeDefProxy(t *testing.T) {
// 	reader, _ := GetReader(testcaseFilename)
// 	var p *TypeDefProxy
// 	var e *dwarf.Entry

// 	// Start with a few trivial cases
// 	e, _ = GetEntry(reader, "char")
// 	p = NewTypeDefProxy(reader, e)
// 	assert.Equal(t, "char", p.Name)
// 	assert.Equal(t, int(8), p.BitSize)
//   assert.Equal(t, make([]TypeDefProxy, 0), p.Children)

// 	// Move on to non-trivial cases in which Children must actually be populated
// 	e, _ = GetEntry(reader, "Driver")
// 	p = NewTypeDefProxy(reader, e)
// 	var driverChildren = []TypeDefProxy{
//     {
// 			Name:        "initials",
// 			BitSize:     8,
// 			DwarfOffset: 0,
// 			ArrayRanges: []int{2},
//       Children:    make([]TypeDefProxy, 0),
// 		},
//     {
// 			Name:        "car_number",
// 			BitSize:     32,
// 			DwarfOffset: 16,
// 			ArrayRanges: []int{0},
//       Children:    make([]TypeDefProxy, 0),
// 		},
//     {
// 			Name:        "has_won_wdc",
// 			BitSize:     8,
// 			DwarfOffset: 48,
// 			ArrayRanges: []int{0},
//       Children:    make([]TypeDefProxy, 0),
// 		},
// 	}
// 	assert.Equal(t, "Driver", p.Name)
// 	assert.Equal(t, int(12*8), p.BitSize)
// 	assert.Equal(t, driverChildren, p.Children)
// }
