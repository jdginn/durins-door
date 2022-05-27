package explorer

import (
	"debug/dwarf"
)

type explorerCtx struct {
	tags   []dwarf.Tag
	levels []*dwarf.Entry
}

// Returns the current entry pointed to by this context
func (c *explorerCtx) CurrEntry() (*dwarf.Entry, bool) {
	if len(c.levels) <= 0 {
		return nil, false
	}
	return c.levels[len(c.levels)-1], true
}

func (c *explorerCtx) Push(e *dwarf.Entry) {
	c.levels = append(c.levels, e)
}

func (c *explorerCtx) Pop() (*dwarf.Entry, bool) {
	var entry *dwarf.Entry
	if len(c.levels) == 0 {
		return nil, false
	}
	entry, c.levels = c.levels[len(c.levels)-1], c.levels[:len(c.levels)-1]
	return entry, true
}
