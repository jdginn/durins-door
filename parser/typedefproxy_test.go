package parser

import (
	"debug/dwarf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTypeDefProxy(t *testing.T) {
	reader, _ := getReaderFromFile(testcaseFilename)
	// var driverProxy *TypeDefProxy
	var e *dwarf.Entry
	var err error

	// Start with a few trivial cases
	e, _, err = GetEntry(reader, "char")
	assert.Equal(t, nil, err)
	driverProxy, err := NewTypeDefProxy(reader, e)
	assert.NoError(t, err)
	assert.Equal(t, "char", driverProxy.Name)
	assert.Equal(t, int(8), driverProxy.BitSize)
	assert.Equal(t, make([]TypeDefProxy, 0), driverProxy.Children)

	// Move on to non-trivial cases in which Children must actually be populated
	e, _, err = GetEntry(reader, "Driver")
	assert.NoError(t, err)
	driverProxy, err = NewTypeDefProxy(reader, e)
	assert.NoError(t, err)
	// NOTE: clang chooses to pad bools out to 4 bytes despite the typical implementation
	// being only 1 byte
	var driverChildren = []TypeDefProxy{
		{
			Name:         "initials",
			BitSize:      8,
			StructOffset: 0,
			ArrayRanges:  []int{2},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:         "car_number",
			BitSize:      32,
			StructOffset: 32,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:         "has_won_wdc",
			BitSize:      8,
			StructOffset: 64,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
	}
	assert.Equal(t, "Driver", driverProxy.Name)
	assert.Equal(t, int(12*8), driverProxy.BitSize)
	assert.Equal(t, driverChildren, driverProxy.Children)

	// A type that includes the type from the previous test
	e, _, err = GetEntry(reader, "Team")
	assert.NoError(t, err)
	teamProxy, err := NewTypeDefProxy(reader, e)
	assert.NoError(t, err)

	var teamChildren = []TypeDefProxy{
		{
			Name:         "drivers",
			BitSize:      96,
			StructOffset: 0,
			ArrayRanges:  []int{2},
			Children:     driverProxy.Children,
		},
		{
			Name:         "sponsors",
			BitSize:      16,
			StructOffset: 192,
			ArrayRanges:  []int{4},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:         "has_won_wdc",
			BitSize:      8,
			StructOffset: 256,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:         "last_wdc",
			BitSize:      32,
			StructOffset: 288,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:         "has_won_wcc",
			BitSize:      8,
			StructOffset: 320,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
		{
			Name:         "last_wcc",
			BitSize:      32,
			StructOffset: 352,
			ArrayRanges:  []int{0},
			Children:     make([]TypeDefProxy, 0),
		},
	}

	assert.Equal(t, "Team", teamProxy.Name)
	assert.Equal(t, int(384), teamProxy.BitSize)
	assert.Equal(t, teamChildren, teamProxy.Children)
}

func TestTypeDefProxyGetChild(t *testing.T) {
	reader, _ := getReaderFromFile(testcaseFilename)

	// Navigate one level down
	e, _, err := GetEntry(reader, "Driver")
	assert.NoError(t, err)
	driverProxy, err := NewTypeDefProxy(reader, e)
	assert.NoError(t, err)
	initialsProxy, err := driverProxy.GetChild("initials")
	assert.NoError(t, err)
	assert.Equal(t, "initials", initialsProxy.Name)
	assert.Equal(t, int(8), initialsProxy.BitSize)
	assert.Equal(t, int(0), initialsProxy.StructOffset)
	assert.Equal(t, []int{2}, initialsProxy.ArrayRanges)

	// Navigate two levels down
	e, _, err = GetEntry(reader, "Team")
	assert.NoError(t, err)
	teamProxy, err := NewTypeDefProxy(reader, e)
	assert.NoError(t, err)
	driverProxy, err = teamProxy.GetChild("drivers")
	assert.NoError(t, err)
	assert.Equal(t, "drivers", driverProxy.Name)
	assert.Equal(t, int(12*8), driverProxy.BitSize)
	assert.Equal(t, int(0), driverProxy.StructOffset)
	assert.Equal(t, []int{2}, driverProxy.ArrayRanges)
	initialsProxy, err = driverProxy.GetChild("initials")
	assert.NoError(t, err)
	assert.Equal(t, "initials", initialsProxy.Name)
	assert.Equal(t, int(8), initialsProxy.BitSize)
	assert.Equal(t, int(0), initialsProxy.StructOffset)
	assert.Equal(t, []int{2}, initialsProxy.ArrayRanges)
}
