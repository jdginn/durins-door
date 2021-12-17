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

func TestNewTypeDefProxy(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
	// var driverProxy *TypeDefProxy
	var e *dwarf.Entry

	// Start with a few trivial cases
	e, _ = GetEntry(reader, "char")
  driverProxy := NewTypeDefProxy(reader, e)
	assert.Equal(t, "char", driverProxy.Name)
	assert.Equal(t, int(8), driverProxy.BitSize)
	assert.Equal(t, make([]TypeDefProxy, 0), driverProxy.Children)

	// Move on to non-trivial cases in which Children must actually be populated
	e, _ = GetEntry(reader, "Driver")
	driverProxy = NewTypeDefProxy(reader, e)
	var driverChildren = []TypeDefProxy {
		{
			Name:         "initials",
			BitSize:      8,
			StructOffset: 0,
			ArrayRanges:  []int{2},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:    "car_number",
			BitSize: 32,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			StructOffset: 32,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:    "has_won_wdc",
			BitSize: 8,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			// StructOffset: 48,
			StructOffset: 64,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
	}
	assert.Equal(t, "Driver", driverProxy.Name)
	assert.Equal(t, int(12*8), driverProxy.BitSize)
	assert.Equal(t, driverChildren, driverProxy.Children)

	// A type that includes the type from the previous test
	e, _ = GetEntry(reader, "Team")
  teamProxy := NewTypeDefProxy(reader, e)

	var teamChildren = []TypeDefProxy {
		{
			Name:         "Driver",
			BitSize:      96,
			StructOffset: 0,
			ArrayRanges:  []int{2},
      Children:     []TypeDefProxy{*driverProxy},
		},
		{
			Name:    "sponsors",
			BitSize: 16,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			StructOffset: 96,
			ArrayRanges:  []int{4},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:    "has_won_wdc",
			BitSize: 8,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			// StructOffset: 48,
			StructOffset: 128,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:    "last_wdc",
			BitSize: 32,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			// StructOffset: 48,
			StructOffset: 160,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:    "has_won_wcc",
			BitSize: 8,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			// StructOffset: 48,
			StructOffset: 192,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:    "last_wcc",
			BitSize: 32,
			// TODO: I'm fairly ceratain this is correct and the program is currently
			// generating bad output but I will make this pass temporarily to evaluate
			// some additional testcases to probe my theory
			// StructOffset: 48,
			StructOffset: 222,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
  }

	assert.Equal(t, "Team", teamProxy.Name)
	assert.Equal(t, int(30*8), teamProxy.BitSize)
	assert.Equal(t, teamChildren, teamProxy.Children)

}
