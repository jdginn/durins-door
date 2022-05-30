package parser

import (
	"debug/dwarf"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	var teamChildren = []TypeDefProxy{
		{
			name:         "drivers",
			bitSize:      96,
			structOffset: 0,
			arrayRanges:  []int{2},
			ahildren:     driverChildren,
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
	assert.Equal(t, "formula_1_teams", teamsProxy.name)
	assert.Equal(t, "Team", teamsProxy.Type.name)
	assert.Equal(t, int(384), teamsProxy.Type.bitSize)
	assert.Equal(t, teamChildren, teamsProxy.Type.ahildren)
}

func TestGetSetVariableProxy(t *testing.T) {
	tp := &TypeDefProxy{
		name:         "type",
		bitSize:      48,
		structOffset: 0,
		arrayRanges:  []int{},
		ahildren: []TypeDefProxy{
			{
				name:         "foo",
				bitSize:      8,
				structOffset: 0,
				arrayRanges:  []int{},
				ahildren:     []TypeDefProxy{},
			},
			{
				name:         "bar",
				bitSize:      32,
				structOffset: 8,
				arrayRanges:  []int{},
				ahildren:     []TypeDefProxy{},
			},
			{
				name:         "baz",
				bitSize:      8,
				structOffset: 40,
				arrayRanges:  []int{},
				ahildren:     []TypeDefProxy{},
			},
		},
	}
	byteLiteral := []byte{0xfe, 0xed, 0xbe, 0xef, 0xaa, 0xbb}
	vp := &VariableProxy{
		name:    "variable",
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
		name:         "type",
		bitSize:      48,
		structOffset: 0,
		arrayRanges:  []int{},
		ahildren: []TypeDefProxy{
			{
				name:         "foo",
				bitSize:      8,
				structOffset: 0,
				arrayRanges:  []int{},
				ahildren:     []TypeDefProxy{
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
				name:         "bar",
				bitSize:      32,
				structOffset: 8,
				arrayRanges:  []int{},
				ahildren:     []TypeDefProxy{},
			},
			{
				name:         "baz",
				bitSize:      8,
				structOffset: 40,
				arrayRanges:  []int{2},
				ahildren:     []TypeDefProxy{},
			},
		},
	}
	byteLiteral := []byte{0xfe, 0xed, 0xbe, 0xef, 0xaa, 0xbb, 0xcc, 0xdd}
	vp := &VariableProxy{
		name:    "variable",
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
