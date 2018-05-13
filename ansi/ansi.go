package ansi

import (
	"bytes"
	"fmt"
	"regexp"
)

// Known ansi codes
const (
	Clear int = iota
	ClearLine
	Cursor
	ShowCursor
	HideCursor
	Del
	Cr
	Newline
	Red
	Green
	Reset
)

// Directions
const (
	Rgt int = iota
	Lft
)

// Code represents an ansi code
type Code struct {
	Code  int
	Chars []byte
	Hint  string
}

// Str returns the ansi code's chars as a string
func (c Code) Str() string {
	return string(c.Chars)
}

var ansi = map[int]Code{
	Clear:      {Clear, []byte("\x1b[0K"), "<clear>"},
	ClearLine:  {ClearLine, []byte("\x1b[2K"), "<clear-line>"},
	ShowCursor: {ShowCursor, []byte("\x1b[?25h"), "<show-crsr>"},
	HideCursor: {HideCursor, []byte("\x1b[?25l"), "<hide-crsr>"},
	Newline:    {Newline, []byte("\n"), "<nl>"},
	Cr:         {Cr, []byte("\r"), "<cr>"},
	Del:        {Del, []byte("\x7F"), "<del>"},
	Red:        {Red, []byte("\x1b[0;31m"), "<red>"},
	Green:      {Green, []byte("\x1b[0;32m"), "<green>"},
	Reset:      {Reset, []byte("\x1b[0m"), "<reset>"},
}

// Ansi returns the chars for a given ansi code
func Ansi(c int) []byte {
	return ansi[c].Chars
}

// Colored returns the given chars wrapped in color codes
func Colored(c int, b []byte) []byte {
	return concat(ansi[c].Chars, b, ansi[Reset].Chars)
}

var regCrsr = regexp.MustCompile("\x1b\\[([0-9]+)(C|D)")
var open = []byte("\x1b[")
var dir = map[int][]byte{
	Rgt: []byte("C"),
	Lft: []byte("D"),
}

// SetCursor returns ansi codes for setting the cursor to the given horizontal
// position.
func SetCursor(pos int) []byte {
	p := []byte(fmt.Sprintf("%d", pos))
	return concat(dup(ansi[Cr].Chars), open, p, dir[Rgt])
}

// MoveCursor returns ansi codes for moving the cursor by the given number of
// chars in the given direction.
func MoveCursor(pos int, d int) []byte {
	p := []byte(fmt.Sprintf("%d", pos))
	return concat(open, p, dir[d])
}

// Deansi replaces ansi codes in the given byte array with hints. This is
// useful for testing.
func Deansi(b []byte) []byte {
	for _, k := range ansi {
		b = bytes.Replace(b, k.Chars, []byte(k.Hint), 99)
	}
	for m := regCrsr.FindSubmatch(b); len(m) > 0; m = regCrsr.FindSubmatch(b) {
		var name string
		if bytes.Equal(m[2], dir[Rgt]) {
			name = "rgt"
		} else {
			name = "lft"
		}
		tag := tag(name, m[1])
		b = bytes.Replace(b, m[0], tag, 1)
	}
	return b
}

func tag(name string, body []byte) []byte {
	return concat([]byte("<"+name+"-"), body, []byte(">"))
}

func concat(a []byte, b ...[]byte) []byte {
	r := append(a, b[0]...)
	if len(b) > 1 {
		r = concat(r, b[1:]...)
	}
	return r
}

func dup(b []byte) []byte {
	return append([]byte{}, b...)
}
