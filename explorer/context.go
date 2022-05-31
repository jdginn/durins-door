package explorer

import (
	"debug/dwarf"
	// "sync"

	"github.com/jdginn/durins-door/parser"
)

type mode int

const (
	modeCUs mode = iota
	modeEntry
	modeProxy
)

type ctxLevel struct {
	mode  mode
	entry *dwarf.Entry
	proxy parser.Proxy
}

type stack struct {
	// mux    sync.Mutex
	levels []ctxLevel
}

func NewStack() *stack {
	// s := stack{sync.Mutex{}, make([]ctxLevel, 0, 128)}
	s := &stack{make([]ctxLevel, 0, 128)}
	s.Push(modeCUs, nil, nil)
	return s
}

// Returns the current entry pointed to by this context
func (c *stack) CurrEntry() *dwarf.Entry {
	// c.mux.Lock()
	// defer c.mux.Unlock()
	return c.levels[len(c.levels)-1].entry
}

// Returns the current entry pointed to by this context
func (c *stack) CurrProxy() parser.Proxy {
	// c.mux.Lock()
	// defer c.mux.Unlock()
	return c.levels[len(c.levels)-1].proxy
}

func (c *stack) CurrMode() mode {
	if len(c.levels) < 1 {
		return -1
	}
	// c.mux.Lock()
	// defer c.mux.Unlock()
	return c.levels[len(c.levels)-1].mode
}

func (c *stack) Push(m mode, e *dwarf.Entry, p parser.Proxy) {
	// c.mux.Lock()
	// defer c.mux.Unlock()
	c.levels = append(c.levels, ctxLevel{m, e, p})
}

func (c *stack) Pop() (ctxLevel, bool) {
	// c.mux.Lock()
	// defer c.mux.Unlock()
	var level ctxLevel
	if len(c.levels) == 0 {
		return ctxLevel{}, false
	}
	level, c.levels = c.levels[len(c.levels)-1], c.levels[:len(c.levels)-1]
	return level, true
}
