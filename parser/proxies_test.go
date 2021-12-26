package parser

import (
	"debug/dwarf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTypeDefProxy(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
	// var driverProxy *TypeDefProxy
	var e *dwarf.Entry
	var err error

	// Start with a few trivial cases
	e, err = GetEntry(reader, "char")
	assert.Equal(t, nil, err)
	driverProxy, err := NewTypeDefProxy(reader, e)
	assert.Equal(t, nil, err)
	assert.Equal(t, "char", driverProxy.Name)
	assert.Equal(t, int(8), driverProxy.BitSize)
	assert.Equal(t, make([]TypeDefProxy, 0), driverProxy.Children)

	// Move on to non-trivial cases in which Children must actually be populated
	e, err = GetEntry(reader, "Driver")
	assert.Equal(t, nil, err)
	driverProxy, err = NewTypeDefProxy(reader, e)
	assert.Equal(t, nil, err)
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
	e, err = GetEntry(reader, "Team")
	assert.Equal(t, nil, err)
	teamProxy, err := NewTypeDefProxy(reader, e)
	assert.Equal(t, nil, err)

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

func TestNewVariableProxy(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
	var teamsProxy *VariableProxy
	var e *dwarf.Entry
	var err error

	// // Start with a few trivial cases
	// e, err = GetEntry(reader, "char")
	//  assert.Equal(t, nil, err)
	//  teamsProxy, err := NewTypeDefProxy(reader, e)
	//  assert.Equal(t, nil, err)
	// assert.Equal(t, "char", teamsProxy.Name)
	// assert.Equal(t, int(8), teamsProxy.BitSize)
	// assert.Equal(t, make([]TypeDefProxy, 0), teamsProxy.Children)

	// Move on to non-trivial cases in which Children must actually be populated
	e, err = GetEntry(reader, "formula_1_teams")
	assert.Equal(t, nil, err)
	teamsProxy, err = NewVariableProxy(reader, e)
	assert.Equal(t, nil, err)
	// NOTE: clang chooses to pad bools out to 4 bytes despite the typical implementation
	// being only 1 byte
	// var driverChildren = []TypeDefProxy{
	// 	{
	// 		Name:         "initials",
	// 		BitSize:      8,
	// 		StructOffset: 0,
	// 		ArrayRanges:  []int{2},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// 	{
	// 		Name:         "car_number",
	// 		BitSize:      32,
	// 		StructOffset: 32,
	// 		ArrayRanges:  []int{0},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// 	{
	// 		Name:         "has_won_wdc",
	// 		BitSize:      8,
	// 		StructOffset: 64,
	// 		ArrayRanges:  []int{0},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// }
	assert.Equal(t, "formula_1_teams", teamsProxy.Name)
	assert.Equal(t, int(384), teamsProxy.Type.BitSize)
	// assert.Equal(t, driverChildren, teamsProxy.Children)

	// A type that includes the type from the previous test
	e, err = GetEntry(reader, "formula_1_teams")
	assert.Equal(t, nil, err)
	teamProxy, err := NewVariableProxy(reader, e)
	assert.Equal(t, nil, err)

	// var teamChildren = []TypeDefProxy{
	// 	{
	// 		Name:         "drivers",
	// 		BitSize:      96,
	// 		StructOffset: 0,
	// 		ArrayRanges:  []int{2},
	// 		Children:     driverChildren,
	// 	},
	// 	{
	// 		Name:         "sponsors",
	// 		BitSize:      16,
	// 		StructOffset: 192,
	// 		ArrayRanges:  []int{4},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// 	{
	// 		Name:         "has_won_wdc",
	// 		BitSize:      8,
	// 		StructOffset: 256,
	// 		ArrayRanges:  []int{0},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// 	{
	// 		Name:         "last_wdc",
	// 		BitSize:      32,
	// 		StructOffset: 288,
	// 		ArrayRanges:  []int{0},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// 	{
	// 		Name:         "has_won_wcc",
	// 		BitSize:      8,
	// 		StructOffset: 320,
	// 		ArrayRanges:  []int{0},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// 	{
	// 		Name:         "last_wcc",
	// 		BitSize:      32,
	// 		StructOffset: 352,
	// 		ArrayRanges:  []int{0},
	// 		Children:     make([]TypeDefProxy, 0),
	// 	},
	// }

	assert.Equal(t, "formula_1_teams", teamProxy.Name)
	assert.Equal(t, "Team", teamProxy.Type.Name)
	assert.Equal(t, int(384), teamProxy.Type.BitSize)
	assert.Equal(t, uint64(0x100003f50), teamProxy.Address)
	// assert.Equal(t, teamChildren, teamProxy.Children)
}

func TestGetSetVariableProxy(t *testing.T) {
	tp := &TypeDefProxy{
		Name:         "type",
		BitSize:      48,
		StructOffset: 0,
		ArrayRanges:  []int{},
		Children: []TypeDefProxy{
			{
				Name:        "foo",
				BitSize:     8,
				StructOffset: 0,
				ArrayRanges: []int{},
				Children:    []TypeDefProxy{},
			},
			{
				Name:        "bar",
				BitSize:     32,
				StructOffset: 8,
				ArrayRanges: []int{},
				Children:    []TypeDefProxy{},
			},
			{
				Name:        "baz",
				BitSize:     8,
				StructOffset: 40,
				ArrayRanges: []int{},
				Children:    []TypeDefProxy{},
			},
		},
	}
  byteLiteral := []byte{0xfe, 0xed, 0xbe, 0xef, 0xaa, 0xbb, 0xcc}
	vp := &VariableProxy{
		Name: "variable",
    Type: *tp,
    Address: 0xfeedbeef,
    value: byteLiteral,
	}
 
  val, _ := vp.Get()
  assert.Equal(t, byteLiteral, val)
  foo, _ := vp.GetField("foo")
  assert.Equal(t, 0xfe, foo)
  bar, _ := vp.GetField("bar")
  assert.Equal(t, 0xedbeefaabb, bar)
  baz, _ := vp.GetField("baz")
  assert.Equal(t, 0xcc, baz)
}
