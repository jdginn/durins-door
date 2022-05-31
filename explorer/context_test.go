package explorer

import (
	// "fmt"
	"debug/dwarf"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jdginn/durins-door/parser"
)

func TestStack(t *testing.T) {
	s := NewStack()
	assert.Equal(t, modeCUs, s.CurrMode())
	assert.Nil(t, s.CurrEntry())
	assert.Nil(t, s.CurrProxy())

	e := &dwarf.Entry{}
	s.Push(modeEntry, e, nil)
	assert.Equal(t, modeEntry, s.CurrMode())
	assert.Equal(t, e, s.CurrEntry())
	assert.Nil(t, s.CurrProxy())

	p := parser.VariableProxy{}
	s.Push(modeProxy, nil, p)
	assert.Equal(t, modeProxy, s.CurrMode())
	assert.Nil(t, s.CurrEntry())
	assert.Equal(t, p, s.CurrProxy())

	s.Pop()
	assert.Equal(t, modeEntry, s.CurrMode())
	assert.Equal(t, e, s.CurrEntry())
	assert.Nil(t, s.CurrProxy())
}
