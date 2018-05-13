package keys

import (
	"bytes"
)

type reader interface {
	Read(b []byte) (int, error)
}

// Key represents a key
type Key struct {
	Code  int
	Chars []byte
	Name  string
}

// Str returns the key's chars as a string
func (k Key) Str() string {
	return string(k.Chars)
}

// Known keys
const (
	Chars int = iota
	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlH
	Tab
	CtrlK
	CtrlL
	Enter
	CtrlN
	CtrlP
	CtrlT
	CtrlU
	CtrlW
	Esc
	Backspace
	Delete
	ShiftTab
	Up
	Down
	Right
	Left
)

// Keys defines known keys
var Keys = map[int]Key{
	CtrlA:     {CtrlA, []byte{0x1}, "Ctrl-A"},
	CtrlB:     {CtrlB, []byte{0x2}, "Ctrl-B"},
	CtrlC:     {CtrlC, []byte{0x3}, "Ctrl-C"},
	CtrlD:     {CtrlD, []byte{0x4}, "Ctrl-D"},
	CtrlE:     {CtrlE, []byte{0x5}, "Ctrl-E"},
	CtrlF:     {CtrlF, []byte{0x6}, "Ctrl-F"},
	CtrlH:     {CtrlH, []byte{0x8}, "Ctrl-H"},
	Tab:       {Tab, []byte{0x9}, "Tab"},
	CtrlK:     {CtrlK, []byte{0x0b}, "Ctrl-K"},
	CtrlL:     {CtrlL, []byte{0x0c}, "Ctrl-L"},
	Enter:     {Enter, []byte{0x0d}, "Enter"},
	CtrlN:     {CtrlN, []byte{0x0e}, "Ctrl-N"},
	CtrlP:     {CtrlP, []byte{0x10}, "Ctrl-P"},
	CtrlT:     {CtrlT, []byte{0x14}, "Ctrl-T"},
	CtrlU:     {CtrlU, []byte{0x15}, "Ctrl-U"},
	CtrlW:     {CtrlW, []byte{0x17}, "Ctrl-W"},
	Esc:       {Esc, []byte{0x1b}, "Esc"},
	Backspace: {Backspace, []byte{0x7f}, "Backspace"},
	Delete:    {Delete, []byte{0x1b, 0x5b, 0x33, 0x7E}, "Delete"}, // \x1b[3~
	ShiftTab:  {ShiftTab, []byte{0x1b, 0x5b, 0x5a}, "Shift-Tab"},
	Up:        {Up, []byte{0x1b, 0x5b, 0x41}, "Up"},
	Down:      {Down, []byte{0x1b, 0x5b, 0x42}, "Down"},
	Right:     {Right, []byte{0x1b, 0x5b, 0x43}, "Right"},
	Left:      {Left, []byte{0x1b, 0x5b, 0x44}, "Left"},
}

// Read returns a channel for reading keys. Terminates on ctrl-d.
func Read(tty reader) chan Key {
	bytes := make([]byte, 10)
	keys := make(chan Key, 1)

	go func() {
	loop:
		for {
			i, _ := tty.Read(bytes)
			if i == 0 {
				continue
			}
			key := find(bytes[:i])
			if key.Code == CtrlD {
				break loop
			}
			keys <- key
		}
		close(keys)
	}()

	return keys
}

func find(b []byte) Key {
	for _, k := range Keys {
		if bytes.Equal(k.Chars, b) {
			return k
		}
	}
	return Key{Code: Chars, Chars: b}
}
