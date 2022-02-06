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

func TestTypeDefProxyGetChild(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)

	// Navigate one level down
	e, err := GetEntry(reader, "Driver")
	assert.Nil(t, err)
	driverProxy, err := NewTypeDefProxy(reader, e)
	assert.Nil(t, err)
	initialsProxy, err := driverProxy.GetChild("initials")
	assert.Nil(t, err)
	assert.Equal(t, "initials", initialsProxy.Name)
	assert.Equal(t, int(8), initialsProxy.BitSize)
	assert.Equal(t, int(0), initialsProxy.StructOffset)
	assert.Equal(t, []int{2}, initialsProxy.ArrayRanges)

	// Navigate two levels down
	e, err = GetEntry(reader, "Team")
	teamProxy, err := NewTypeDefProxy(reader, e)
	assert.Nil(t, err)
	driverProxy, err = teamProxy.GetChild("drivers")
	assert.Nil(t, err)
	assert.Equal(t, "drivers", driverProxy.Name)
	assert.Equal(t, int(12*8), driverProxy.BitSize)
	assert.Equal(t, int(0), driverProxy.StructOffset)
	assert.Equal(t, []int{2}, driverProxy.ArrayRanges)
	initialsProxy, err = driverProxy.GetChild("initials")
	assert.Nil(t, err)
	assert.Equal(t, "initials", initialsProxy.Name)
	assert.Equal(t, int(8), initialsProxy.BitSize)
	assert.Equal(t, int(0), initialsProxy.StructOffset)
	assert.Equal(t, []int{2}, initialsProxy.ArrayRanges)
}

func TestNewVariableProxy(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
	var teamsProxy *VariableProxy
	var e *dwarf.Entry
	var err error

	// Move on to non-trivial cases in which Children must actually be populated
	e, err = GetEntry(reader, "formula_1_teams")
	assert.Equal(t, nil, err)
	teamsProxy, err = NewVariableProxy(reader, e)
	assert.Equal(t, nil, err)
	// First we confirm that this variable includes the same type we found in
	// TestNewTypeDefProxy
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
	var teamChildren = []TypeDefProxy{
		{
			Name:         "drivers",
			BitSize:      96,
			StructOffset: 0,
			ArrayRanges:  []int{2},
			Children:     driverChildren,
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
	assert.Equal(t, "formula_1_teams", teamsProxy.Name)
	assert.Equal(t, "Team", teamsProxy.Type.Name)
	assert.Equal(t, int(384), teamsProxy.Type.BitSize)
	assert.Equal(t, teamChildren, teamsProxy.Type.Children)
	assert.Equal(t, uint64(0x100003f50), teamsProxy.Address)
}

func TestGetSetVariableProxy(t *testing.T) {
	tp := &TypeDefProxy{
		Name:         "type",
		BitSize:      48,
		StructOffset: 0,
		ArrayRanges:  []int{},
		Children: []TypeDefProxy{
			{
				Name:         "foo",
				BitSize:      8,
				StructOffset: 0,
				ArrayRanges:  []int{},
				Children:     []TypeDefProxy{},
			},
			{
				Name:         "bar",
				BitSize:      32,
				StructOffset: 8,
				ArrayRanges:  []int{},
				Children:     []TypeDefProxy{},
			},
			{
				Name:         "baz",
				BitSize:      8,
				StructOffset: 40,
				ArrayRanges:  []int{},
				Children:     []TypeDefProxy{},
			},
		},
	}
	byteLiteral := []byte{0xfe, 0xed, 0xbe, 0xef, 0xaa, 0xbb}
	vp := &VariableProxy{
		Name:    "variable",
		Type:    *tp,
		Address: 0xfeedbeef,
		value:   byteLiteral,
	}

	val, err := vp.Get()
	assert.Equal(t, byteLiteral, val)
	foo, err := vp.GetField("foo")
	assert.Equal(t, uint64(0xfe), foo)
	bar, err := vp.GetField("bar")
	assert.Equal(t, uint64(0xedbeefaa), bar)
	baz, err := vp.GetField("baz")
	assert.Equal(t, uint64(0xbb), baz)
	assert.Nil(t, err)

	// Should fail because this data cannot fit in this variable's type
	err = vp.Set([]byte{0x22, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
	assert.Error(t, err)

	err = vp.Set([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66})

	assert.Nil(t, err)
	foo, err = vp.GetField("foo")
	assert.Equal(t, uint64(0x11), foo)
	bar, err = vp.GetField("bar")
	assert.Equal(t, uint64(0x22334455), bar)
	baz, err = vp.GetField("baz")
	assert.Equal(t, uint64(0x66), baz)
	assert.Nil(t, err)

	err = vp.SetField("foo", uint64(0xff))
	err = vp.SetField("bar", uint64(0x00c0ffee))
	err = vp.SetField("baz", uint64(0x00))
	assert.Nil(t, err)

	val, err = vp.Get()
	assert.Equal(t, []byte{0xff, 0x00, 0xc0, 0xff, 0xee, 0x00}, val)

	foo, err = vp.GetField("foo")
	assert.Equal(t, uint64(0xff), foo)
	bar, err = vp.GetField("bar")
	assert.Equal(t, uint64(0xc0ffee), bar)
	baz, err = vp.GetField("baz")
	assert.Equal(t, uint64(0x00), baz)
	assert.Nil(t, err)
}
