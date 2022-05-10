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

func TestNewVariableProxy(t *testing.T) {
	reader, _ := getReaderFromFile(testcaseFilename)
	var teamsProxy *VariableProxy
	var e *dwarf.Entry
	var err error

	// Move on to non-trivial cases in which Children must actually be populated
	e, _, err = GetEntry(reader, "formula_1_teams")
	assert.NoError(t, err)
	teamsProxy, err = NewVariableProxy(reader, e)
	assert.NoError(t, err)
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
	assert.Equal(t, int(0x100008010), teamsProxy.Address)
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
	assert.NoError(t, err)
	assert.Equal(t, byteLiteral, val)
	foo, err := vp.GetField("foo")
	assert.NoError(t, err)
	assert.Equal(t, int(0xfe), foo)
	bar, err := vp.GetField("bar")
	assert.NoError(t, err)
	assert.Equal(t, int(0xedbeefaa), bar)
	baz, err := vp.GetField("baz")
	assert.NoError(t, err)
	assert.Equal(t, int(0xbb), baz)

	// Should fail because this data cannot fit in this variable's type
	err = vp.Set([]byte{0x22, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
	assert.Error(t, err)

	err = vp.Set([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66})

	assert.NoError(t, err)
	foo, err = vp.GetField("foo")
	assert.NoError(t, err)
	assert.Equal(t, int(0x11), foo)
	bar, err = vp.GetField("bar")
	assert.NoError(t, err)
	assert.Equal(t, int(0x22334455), bar)
	baz, err = vp.GetField("baz")
	assert.Equal(t, int(0x66), baz)
	assert.NoError(t, err)

	err = vp.SetField("foo", int(0xff))
	assert.NoError(t, err)
	err = vp.SetField("bar", int(0x00c0ffee))
	assert.NoError(t, err)
	err = vp.SetField("baz", int(0x00))
	assert.NoError(t, err)

	val, err = vp.Get()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xff, 0x00, 0xc0, 0xff, 0xee, 0x00}, val)

	foo, err = vp.GetField("foo")
	assert.NoError(t, err)
	assert.Equal(t, int(0xff), foo)
	bar, err = vp.GetField("bar")
	assert.NoError(t, err)
	assert.Equal(t, int(0xc0ffee), bar)
	baz, err = vp.GetField("baz")
	assert.NoError(t, err)
	assert.Equal(t, int(0x00), baz)
	assert.NoError(t, err)
}

func TestGetSetMultilevelVariableProxy(t *testing.T) {
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
				Children:     []TypeDefProxy{
					// {
					// 	Name:         "jonah",
					// 	BitSize:      32,
					// 	StructOffset: 8,
					// 	ArrayRanges:  []int{},
					// 	Children:     []TypeDefProxy{},
					// },
					// {
					// 	Name:         "noah",
					// 	BitSize:      8,
					// 	StructOffset: 40,
					// 	ArrayRanges:  []int{2},
					// 	Children:     []TypeDefProxy{},
					// },
				},
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
				ArrayRanges:  []int{2},
				Children:     []TypeDefProxy{},
			},
		},
	}
	byteLiteral := []byte{0xfe, 0xed, 0xbe, 0xef, 0xaa, 0xbb, 0xcc, 0xdd}
	vp := &VariableProxy{
		Name:    "variable",
		Type:    *tp,
		Address: 0xfeedbeef,
		value:   byteLiteral,
	}

	val, err := vp.Get()
	assert.NoError(t, err)
	assert.Equal(t, byteLiteral, val)
	foo, err := vp.GetField("foo")
	assert.NoError(t, err)
	assert.Equal(t, int(0xfe), foo)
	bar, err := vp.GetField("bar")
	assert.NoError(t, err)
	assert.Equal(t, int(0xedbeefaa), bar)
	baz, err := vp.GetField("baz")
	assert.Equal(t, int(0xbb), baz)
	assert.NoError(t, err)

	// Should fail because this data cannot fit in this variable's type
	err = vp.Set([]byte{0x22, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
	assert.Error(t, err)

	err = vp.Set([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66})

	assert.NoError(t, err)
	foo, err = vp.GetField("foo")
	assert.NoError(t, err)
	assert.Equal(t, int(0x11), foo)
	bar, err = vp.GetField("bar")
	assert.NoError(t, err)
	assert.Equal(t, int(0x22334455), bar)
	baz, err = vp.GetField("baz")
	assert.Equal(t, int(0x66), baz)
	assert.NoError(t, err)

	err = vp.SetField("foo", int(0xff))
	assert.NoError(t, err)
	err = vp.SetField("bar", int(0x00c0ffee))
	assert.NoError(t, err)
	err = vp.SetField("baz", int(0x00))
	assert.NoError(t, err)

	val, err = vp.Get()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xff, 0x00, 0xc0, 0xff, 0xee, 0x00}, val)

	foo, err = vp.GetField("foo")
	assert.NoError(t, err)
	assert.Equal(t, int(0xff), foo)
	bar, err = vp.GetField("bar")
	assert.NoError(t, err)
	assert.Equal(t, int(0xc0ffee), bar)
	baz, err = vp.GetField("baz")
	assert.Equal(t, int(0x00), baz)
	assert.NoError(t, err)
}

func TestGetAccessMetadata(t *testing.T) {}
