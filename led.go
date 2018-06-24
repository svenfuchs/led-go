package led

import (
	"bytes"
	"time"
)

// Forw/Back - directions
// Hist/Comp/Sugg - modes for cycling through sets
const (
	Forw int = iota
	Back
	Hist
	Comp
	Sugg
)

var space = []byte{' '}

// NewReadline creates a line editor that resembles most of Linenoise's functionality
func NewReadline(led string, t ...Iterm) *Ed {
	e := NewEd(led, t...)
	e.Handle(Chars, func(e *Ed, k Key) { e.Insert(k.Chars) })
	e.Handle(CtrlA, func(e *Ed, k Key) { e.Return() })
	e.Handle(CtrlB, func(e *Ed, k Key) { e.Left() })
	e.Handle(CtrlC, func(e *Ed, k Key) { e.Discard() })
	e.Handle(CtrlD, func(e *Ed, k Key) { e.Delete() })
	e.Handle(CtrlE, func(e *Ed, k Key) { e.End() })
	e.Handle(CtrlF, func(e *Ed, k Key) { e.Right() })
	e.Handle(CtrlK, func(e *Ed, k Key) { e.DeleteFromCursor() })
	e.Handle(CtrlT, func(e *Ed, k Key) { e.Transpose() })
	e.Handle(CtrlU, func(e *Ed, k Key) { e.Reset() })
	e.Handle(CtrlW, func(e *Ed, k Key) { e.BackWord() })
	e.Handle(Enter, func(e *Ed, k Key) { e.Newline() })
	e.Handle(Backspace, func(e *Ed, k Key) { e.Back() })
	e.Handle(Delete, func(e *Ed, k Key) { e.Delete() })
	e.Handle(Left, func(e *Ed, k Key) { e.Left() })
	e.Handle(Right, func(e *Ed, k Key) { e.Right() })
	return e
}

// NewEd creates a line editor (Ed) without any handlers attached. See
// NewReadline for a line editor that resembles (most of) Linenoise's
// functionality.
func NewEd(led string, t ...Iterm) *Ed {
	return &Ed{
		term:      StartTerm(t...),
		handlers:  map[int]func(*Ed, Key){},
		Prompt:    []byte(led),
		Pos:       0,
		Chars:     []byte{},
		Suggested: []byte{},
	}
}

// Ed represents the line editor
type Ed struct {
	term      *Term
	handlers  map[int]func(*Ed, Key)
	Prompt    []byte
	Pos       int
	Chars     []byte
	Suggested []byte
	list      *List
}

// Handle attaches a handler for a key
func (e *Ed) Handle(key int, handler func(*Ed, Key)) {
	e.handlers[key] = handler
}

// Run runs the editor
func (e *Ed) Run() {
	e.Refresh()
	for k := range e.term.Read() {
		if e.handlers[k.Code] != nil {
			e.handlers[k.Code](e, k)
		}
	}
	e.Stop()
}

// Pause pauses the editor, should be used before outputting text to the
// terminal, e.g. in an Enter handler.
func (e *Ed) Pause() {
	e.term.Pause()
}

// Resume resumes the editor, and places the terminal in raw mode.
func (e *Ed) Resume() {
	e.term.Resume()
}

// Left moves the cursor one char to the left.
func (e *Ed) Left() {
	if e.Pos > 0 {
		e.MoveCursor(1, Back)
	}
}

// Right moves the cursor one char to the right.
func (e *Ed) Right() {
	if e.Pos < len(e.Chars) {
		e.MoveCursor(1, Forw)
	}
}

// Append appends the given chars at the end of the line.
func (e *Ed) Append(b []byte) {
	e.Chars = append(e.Chars, b...)
	e.Pos = len(e.Chars) - 1
	e.Write(b)
	e.update()
}

// Insert inserts the given chars at the current cursor position.
func (e *Ed) Insert(b []byte) {
	e.Chars = insert(e.Chars, b, e.Pos)
	e.Write(b)
	e.Pos += len(b)
	e.update()
}

// Reject rejects the given chars by printing them at the current
// cursor position in red, and removing them after 100 milliseconds.
func (e *Ed) Reject(chars []byte) {
	e.term.Write(Colored(Red, chars))
	time.Sleep(100 * time.Millisecond)
	e.SetCursor()
	e.clear()
}

// Back removes one char before the current cursor position.
func (e *Ed) Back() {
	if e.Pos == 0 {
		return
	}

	e.MoveCursor(1, Back)
	e.Chars = delete(e.Chars, e.Pos, 1)
	e.Del()
	e.update()
}

// BackWord removes one word before the current cursor position.
func (e *Ed) BackWord() {
	if e.Pos == 0 {
		return
	}

	w := lastWord(e.Chars[:e.Pos], true)
	e.MoveCursor(len(w), Back)
	e.Chars = delete(e.Chars, e.Pos, len(w))
	e.Del(len(w))
	e.update()
}

// Delete removes one char after the current cursor position.
func (e *Ed) Delete() {
	if e.Pos == len(e.Chars) {
		return
	}

	e.Chars = delete(e.Chars, e.Pos, 1)
	e.Del()
	e.update()
}

// DeleteFromCursor deletes all chars from the current cursor position to the
// end of the line.
func (e *Ed) DeleteFromCursor() {
	if e.Pos == len(e.Chars) {
		return
	}

	e.Chars = delete(e.Chars, e.Pos, len(e.Chars)-e.Pos)
	e.clear()
	e.update()
}

// Transpose transposes the char before the cursor position with the one on the
// cursor position, and moves the cursor one char to the right, if possible. If
// the cursor is at the beginning of the line it transposes the current char
// with the next one, and moves the cursor two chars to the right. This
// resembles Zshell's behaviour.
func (e *Ed) Transpose() {
	if len(e.Chars) < 2 {
		return
	}

	offset := 0
	if e.Pos == len(e.Chars) {
		offset = 2
	} else if e.Pos > 0 {
		offset = 1
	}

	e.SetCursor(e.Pos - offset)
	e.Chars = swap(e.Chars, e.Pos, e.Pos+1)
	e.Write(e.Chars[e.Pos:])
	e.SetCursor(min(e.Pos+2, len(e.Chars)))
}

// Discard discards the given input by moving to the next line and starting
// over with an empty prompt. This resembles the behaviour Ctrl-C in Bash or
// Zshele.
func (e *Ed) Discard() {
	e.Newline()
	e.Reset()
}

// CompleteNext displays the next completion from the given slice.
func (e *Ed) CompleteNext(strs [][]byte) {
	e.Complete(strs, Forw)
}

// CompletePrev displays the previous completion from the given slice.
func (e *Ed) CompletePrev(strs [][]byte) {
	e.Complete(strs, Back)
}

// Complete displays the previous or next completion from the given slice
// depending on the given direction.
func (e *Ed) Complete(strs [][]byte, dir int) {
	e.cycle(strs, Comp, dir)
}

// HistoryNext displays the next line from the given slice.
func (e *Ed) HistoryNext(strs [][]byte) {
	e.History(strs, Forw)
}

// HistoryPrev displays the previous line from the given slice.
func (e *Ed) HistoryPrev(strs [][]byte) {
	e.History(strs, Back)
}

// History displays the previous or next line from the given slice
// depending on the given direction.
func (e *Ed) History(strs [][]byte, dir int) {
	e.cycle(strs, Hist, dir)
}

func (e *Ed) cycle(strs [][]byte, mode int, dir int) {
	c := NewList(strs, e.Chars, mode)
	if e.list == nil || !e.list.eq(c) {
		e.list = c
	}

	b := []byte{}
	if mode == Comp {
		b = trimLastWord(e.Chars)
		if len(b) > 0 {
			b = concat(b, space)
		}
	}

	if dir == Back {
		b = concat(b, e.list.Prev())
	} else {
		b = concat(b, e.list.Next())
	}

	if bytes.Equal(b, e.Chars) && mode == Comp {
		e.list = nil
	} else {
		e.Set(b)
	}
}

// Suggest appends the first matching suggestion from the given slice in green
// after the current cursor position.
func (e *Ed) Suggest(str []byte) {
	if e.Pos == 0 {
		e.clearLine()
		return
	}

	// w := lastWord(e.Chars)
	s := bytes.TrimPrefix(str, e.Chars)

	e.Suggested = s
	if len(s) > 0 {
		e.clear()
		e.Write(concat(e.Chars[e.Pos:], Colored(Green, e.Suggested)))
		e.SetCursor()
	}
}

// Newline writes a newline char to the terminae.
func (e *Ed) Newline() {
	e.term.Newline()
}

// Return moves the cursor to the beginning of the line.
func (e *Ed) Return() {
	e.SetCursor(0)
}

// End moves the cursor to the end of the line.
func (e *Ed) End() {
	e.SetCursor(len(e.Chars))
}

// Set sets the content of the editor to the given line.
func (e *Ed) Set(b []byte) {
	e.SetCursor(0)
	e.clear()
	e.Chars = b
	e.Write(b)
	e.Pos += len(b)
	e.update()
}

// Reset resets the editor and starts over with an empty line.
func (e *Ed) Reset() {
	e.reset()
	e.Refresh()
}

// Write writes the given chars to the terminae.
func (e *Ed) Write(b []byte) {
	e.term.Write(b)
}

// SetCursor sets the cursor to the given position, defaults to the current
// position.
func (e *Ed) SetCursor(pos ...int) {
	if len(pos) > 0 {
		e.Pos = pos[0]
	}
	e.term.SetCursor(e.Pos + len(e.Prompt))
}

// MoveCursor moves the cursor by the given number of chars in the given
// direction.
func (e *Ed) MoveCursor(i int, dir int) {
	if dir == Forw {
		if e.Pos == len(e.Chars)-1 {
			return
		}
		e.Pos += i
	} else {
		if e.Pos == 0 {
			return
		}
		e.Pos -= i
	}
	e.term.MoveCursor(i, dir)
}

// Del writes a delete char (`\x7F`) to the terminae.
func (e *Ed) Del(i ...int) {
	e.term.Del(i...)
}

// Stop stops the terminae.
func (e *Ed) Stop() {
	e.Newline()
	e.term.ClearLine()
	e.term.Stop()
}

// Str returns the current line.
func (e *Ed) Str() string {
	return string(e.Chars)
}

func (e *Ed) reset() {
	e.Pos = 0
	e.Chars = []byte{}
	e.Suggested = []byte{}
	e.list = nil
}

func (e *Ed) Refresh() {
	e.update()
	e.clearLine()
	if len(e.Chars) > 0 {
		e.Write(concat(e.Chars, e.Suggested))
	}
	e.SetCursor()
}

func (e *Ed) update() {
	if len(e.Chars) == 0 {
		e.reset()
	}

	if hasTailingSpace(e.Chars) {
		e.list = nil
	}
}

func (e *Ed) clearLine() {
	e.term.ClearLine()
	e.Write(e.Prompt)
}

func (e *Ed) clear() {
	e.term.Clear()
}
