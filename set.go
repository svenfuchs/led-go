package led

import (
	"bytes"
)

var blank = []byte{}

// NewSet returns a list of strings that are used for completion, history, and
// suggestions.
func NewSet(strs [][]byte, str []byte, mode int) *Set {
	if mode == Comp {
		str = lastWord(str)
	} else {
		str = firstWord(str)
	}
	return &Set{strs: strs, str: str, curr: -1}
}

// Set represents a list of strings that are used for completion, history, and
// suggestions.
type Set struct {
	curr int
	strs [][]byte
	str  []byte
}

// Next returns the next string from the list.
func (c *Set) Next() []byte {
	strs := c.matches()
	if len(strs) == 0 {
		return blank
	}
	c.curr = c.curr + 1
	if c.curr >= len(strs) {
		c.curr = 0
	}
	return strs[c.curr]
}

// Prev returns the previous string from the list.
func (c *Set) Prev() []byte {
	strs := c.matches()
	if len(strs) == 0 {
		return blank
	}
	c.curr = c.curr - 1
	if c.curr < 0 {
		c.curr = len(strs) - 1
	}
	return strs[c.curr]
}

func (c *Set) eq(other *Set) bool {
	lft, rgt := c.strs, other.strs
	if len(lft) != len(rgt) {
		return false
	}
	for i, v := range lft {
		if !bytes.Equal(v, rgt[i]) {
			return false
		}
	}
	return true
}

func (c *Set) matches() [][]byte {
	if len(c.str) == 0 {
		return c.strs
	}
	m := [][]byte{}
	for _, s := range c.strs {
		if bytes.HasPrefix(s, c.str) {
			m = append(m, s)
		}
	}
	return m
}
