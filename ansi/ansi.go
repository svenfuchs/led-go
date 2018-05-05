package ansi

import (
	"bytes"
	"fmt"
	"regexp"
)

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

const (
	Rgt int = iota
	Lft
)

type Code struct {
	Code  int
	Chars []byte
	Hint  string
}

func (c Code) Str() string {
	return string(c.Chars)
}

var ansi = map[int]Code{
	Clear:      {Clear, []byte("\x1b[0K"), "<clear>"},
	ClearLine:  {ClearLine, []byte("\x1b[2K"), "<clear-line>"},
	ShowCursor: {ShowCursor, []byte("\x1b[?25h"), "<show-crsr>"},
	HideCursor: {HideCursor, []byte("\x1b[?25l"), "<hide-crsr>"},
	Newline:    {Newline, []byte("\n"), "<nl>"},
	Del:        {Del, []byte("\x7F"), "<del>"},
	Cr:         {Cr, []byte("\r"), "<cr>"},
	Red:        {Red, []byte("\x1b[0;31m"), "<red>"},
	Green:      {Green, []byte("\x1b[0;32m"), "<green>"},
	Reset:      {Reset, []byte("\x1b[0m"), "<reset>"},
}

func Ansi(c int) []byte {
	return ansi[c].Chars
}

func Colored(c int, b []byte) []byte {
	return concat(ansi[c].Chars, b, ansi[Reset].Chars)
}

var regCrsr = regexp.MustCompile("\x1b\\[([0-9]+)(C|D)")
var open = []byte("\x1b[")
var dir = map[int][]byte{
	Rgt: []byte("C"),
	Lft: []byte("D"),
}

func SetCursor(pos int) []byte {
	p := []byte(fmt.Sprintf("%d", pos))
	return concat(dup(ansi[Cr].Chars), open, p, dir[Rgt])
}

func MoveCursor(pos int, d int) []byte {
	p := []byte(fmt.Sprintf("%d", pos))
	return concat(open, p, dir[d])
}

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
