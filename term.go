package led

import (
	"bytes"
	"github.com/pkg/term"
	"github.com/svenfuchs/led-go/ansi"
	"github.com/svenfuchs/led-go/keys"
)

// StartTerm starts a terminal.
func StartTerm(ts ...Iterm) *Term {
	t := Term{}
	if len(ts) > 0 {
		t.tty = ts[0]
	} else {
		t.tty = &termWrap{}
	}
	t.tty.Start()
	return &t
}

// Term represents a terminal
type Term struct {
	tty Iterm
	pos int
}

// Read returns a channel for reading keys from the terminal. See keys.Read.
func (t *Term) Read() chan keys.Key {
	return keys.Read(t.tty)
}

// Write writes the given chars to the terminal.
func (t *Term) Write(b []byte) {
	t.tty.Write(b)
}

// Del writes the given number of delete chars (`\x7F`) to the terminal
// (defaults to 1), deleting the char after the cursor position.
func (t *Term) Del(i ...int) {
	if len(i) == 0 {
		i = []int{1}
	}
	t.tty.Write(bytes.Repeat(chars(ansi.Del), i[0]))
}

// Newline writes a newline char to the terminal.
func (t *Term) Newline() {
	t.Write(chars(ansi.Newline))
}

// Return writes a carriage return char to the terminal.
func (t *Term) Return() {
	t.Write(chars(ansi.Cr))
}

// ClearLine clears the current line.
func (t *Term) ClearLine() {
	t.Return()
	t.Clear()
}

// Clear clears from the current cursor position to the end of the line.
func (t *Term) Clear() {
	t.Write(chars(ansi.Clear))
}

// ShowCursor shows the cursor.
func (t *Term) ShowCursor() {
	t.Write(chars(ansi.ShowCursor))
}

// HideCursor hides the cursor.
func (t *Term) HideCursor() {
	t.Write(chars(ansi.HideCursor))
}

// SetCursor moves the cursor to the given horizontal position.
func (t *Term) SetCursor(pos int) {
	t.Write(ansi.SetCursor(pos))
}

// MoveCursor moves the cursor by the given number of chars in the given
// direction.
func (t *Term) MoveCursor(i int, dir int) {
	t.Write(ansi.MoveCursor(i, dir))
}

// Pause pauses the terminal, restoring the previous mode and settings.
func (t *Term) Pause() {
	t.tty.Restore()
}

// Resume resumes the terminal, setting the terminal in raw mode.
func (t *Term) Resume() {
	t.tty.RawMode()
}

// Stop stops the terminal, restoring the previous mode and settings, and
// closing the tty.
func (t *Term) Stop() {
	t.tty.Restore()
	t.tty.Close()
}

// Iterm represents a subset of the tty implemented in github.com/pkg/term.
type Iterm interface {
	Start()
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	Restore() error
	RawMode() error
	Close() error
}

type termWrap struct {
	tty *term.Term
}

func (t *termWrap) Start() {
	tty, _ := term.Open("/dev/tty")
	t.tty = tty
	t.RawMode()
}

func (t *termWrap) Read(b []byte) (int, error) {
	return t.tty.Read(b)
}

func (t *termWrap) Write(b []byte) (int, error) {
	return t.tty.Write(b)
}

func (t *termWrap) Restore() error {
	return t.tty.Restore()
}

func (t *termWrap) Close() error {
	return t.tty.Close()
}

func (t *termWrap) RawMode() error {
	return term.RawMode(t.tty)
}

func chars(c int) []byte {
	return ansi.Ansi(c)
}
