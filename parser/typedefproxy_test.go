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
	assert.Equal(t, "char", driverProxy.name)
	assert.Equal(t, int(8), driverProxy.bitSize)
	assert.Equal(t, make([]TypeDefProxy, 0), driverProxy.ahildren)

	// Move on to non-trivial cases in which Children must actually be populated
	e, _, err = GetEntry(reader, "Driver")
	assert.NoError(t, err)
	driverProxy, err = NewTypeDefProxy(reader, e)
	assert.NoError(t, err)
	// NOTE: clang chooses to pad bools out to 4 bytes despite the typical implementation
	// being only 1 byte
	var driverChildren = []TypeDefProxy{
		{
			name:         "initials",
			bitSize:      8,
			structOffset: 0,
			arrayRanges:  []int{2},
			ahildren:     make([]TypeDefProxy, 0),
		},
		{
			name:         "car_number",
			bitSize:      32,
			structOffset: 32,
			arrayRanges:  []int{0},
			ahildren:     make([]TypeDefProxy, 0),
		},
		{
			name:         "has_won_wdc",
			bitSize:      8,
			structOffset: 64,
			arrayRanges:  []int{0},
			ahildren:     make([]TypeDefProxy, 0),
		},
	}
	assert.Equal(t, "Driver", driverProxy.name)
	assert.Equal(t, int(12*8), driverProxy.bitSize)
	assert.Equal(t, driverChildren, driverProxy.ahildren)

	// A type that includes the type from the previous test
	e, _, err = GetEntry(reader, "Team")
	assert.NoError(t, err)
	teamProxy, err := NewTypeDefProxy(reader, e)
	assert.NoError(t, err)

	var teamChildren = []TypeDefProxy{
		{
			name:         "drivers",
			bitSize:      96,
			structOffset: 0,
			arrayRanges:  []int{2},
			ahildren:     driverProxy.ahildren,
		},
		{
			name:         "sponsors",
			bitSize:      16,
			structOffset: 192,
			arrayRanges:  []int{4},
			ahildren:     make([]TypeDefProxy, 0),
		},
		{
			name:         "has_won_wdc",
			bitSize:      8,
			structOffset: 256,
			arrayRanges:  []int{0},
			ahildren:     make([]TypeDefProxy, 0),
		},
		{
			name:         "last_wdc",
			bitSize:      32,
			structOffset: 288,
			arrayRanges:  []int{0},
			ahildren:     make([]TypeDefProxy, 0),
		},
		{
			name:         "has_won_wcc",
			bitSize:      8,
			structOffset: 320,
			arrayRanges:  []int{0},
			ahildren:     make([]TypeDefProxy, 0),
		},
		{
			name:         "last_wcc",
			bitSize:      32,
			structOffset: 352,
			arrayRanges:  []int{0},
			ahildren:     make([]TypeDefProxy, 0),
		},
	}

	assert.Equal(t, "Team", teamProxy.name)
	assert.Equal(t, int(384), teamProxy.bitSize)
	assert.Equal(t, teamChildren, teamProxy.ahildren)
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
	assert.Equal(t, "initials", initialsProxy.name)
	assert.Equal(t, int(8), initialsProxy.bitSize)
	assert.Equal(t, int(0), initialsProxy.structOffset)
	assert.Equal(t, []int{2}, initialsProxy.arrayRanges)

	// Navigate two levels down
	e, _, err = GetEntry(reader, "Team")
	assert.NoError(t, err)
	teamProxy, err := NewTypeDefProxy(reader, e)
	assert.NoError(t, err)
	driverProxy, err = teamProxy.GetChild("drivers")
	assert.NoError(t, err)
	assert.Equal(t, "drivers", driverProxy.name)
	assert.Equal(t, int(12*8), driverProxy.bitSize)
	assert.Equal(t, int(0), driverProxy.structOffset)
	assert.Equal(t, []int{2}, driverProxy.arrayRanges)
	initialsProxy, err = driverProxy.GetChild("initials")
	assert.NoError(t, err)
	assert.Equal(t, "initials", initialsProxy.name)
	assert.Equal(t, int(8), initialsProxy.bitSize)
	assert.Equal(t, int(0), initialsProxy.structOffset)
	assert.Equal(t, []int{2}, initialsProxy.arrayRanges)
}
