package led

import (
	"bytes"
	"github.com/svenfuchs/led-go/ansi"
	"github.com/svenfuchs/led-go/keys"
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

// NewPrompt creates a line editor (Prompt) without any handlers attached. See
// NewLed for a line editor that resembles (most of) Linenoise's functionality.
func NewPrompt(prompt string, t ...Iterm) *Prompt {
	return &Prompt{
		term:      StartTerm(t...),
		handlers:  map[int]func(*Prompt, keys.Key){},
		Prompt:    []byte(prompt),
		pos:       0,
		Chars:     []byte{},
		Suggested: []byte{},
	}
}

// Prompt represents the line editor
type Prompt struct {
	term      *Term
	handlers  map[int]func(*Prompt, keys.Key)
	Prompt    []byte
	pos       int
	Chars     []byte
	Suggested []byte
	set       *Set
}

// Handle attaches a handler for a key
func (p *Prompt) Handle(key int, handler func(*Prompt, keys.Key)) {
	p.handlers[key] = handler
}

// Run runs the editor
func (p *Prompt) Run() {
	p.refresh()
	for k := range p.term.Read() {
		if p.handlers[k.Code] != nil {
			p.handlers[k.Code](p, k)
		}
	}
	p.Stop()
}

// Pause pauses the editor, should be used before outputting text to the
// terminal, e.g. in an Enter handler.
func (p *Prompt) Pause() {
	p.term.Pause()
}

// Resume resumes the editor, and places the terminal in raw mode.
func (p *Prompt) Resume() {
	p.term.Resume()
}

// Left moves the cursor one char to the left.
func (p *Prompt) Left() {
	if p.pos > 0 {
		p.MoveCursor(1, Back)
	}
}

// Right moves the cursor one char to the right.
func (p *Prompt) Right() {
	if p.pos < len(p.Chars) {
		p.MoveCursor(1, Forw)
	}
}

// Append appends the given chars at the end of the line.
func (p *Prompt) Append(b []byte) {
	p.Chars = append(p.Chars, b...)
	p.pos = len(p.Chars) - 1
	p.Write(b)
	p.update()
}

// Insert inserts the given chars at the current cursor position.
func (p *Prompt) Insert(b []byte) {
	p.Chars = insert(p.Chars, b, p.pos)
	p.Write(b)
	p.pos += len(b)
	p.update()
}

// Reject rejects the given chars by printing them at the current
// cursor position in red, and removing them after 100 milliseconds.
func (p *Prompt) Reject(chars []byte) {
	p.term.Write(ansi.Colored(ansi.Red, chars))
	time.Sleep(100 * time.Millisecond)
	p.SetCursor()
	p.clear()
}

// Back removes one char before the current cursor position.
func (p *Prompt) Back() {
	if p.pos == 0 {
		return
	}

	p.MoveCursor(1, Back)
	p.Chars = delete(p.Chars, p.pos, 1)
	p.Del()
	p.update()
}

// BackWord removes one word before the current cursor position.
func (p *Prompt) BackWord() {
	if p.pos == 0 {
		return
	}

	w := lastWord(p.Chars[:p.pos], true)
	p.MoveCursor(len(w), Back)
	p.Chars = delete(p.Chars, p.pos, len(w))
	p.Del(len(w))
	p.update()
}

// Delete removes one char after the current cursor position.
func (p *Prompt) Delete() {
	if p.pos == len(p.Chars) {
		return
	}

	p.Chars = delete(p.Chars, p.pos, 1)
	p.Del()
	p.update()
}

// DeleteFromCursor deletes all chars from the current cursor position to the
// end of the line.
func (p *Prompt) DeleteFromCursor() {
	if p.pos == len(p.Chars) {
		return
	}

	p.Chars = delete(p.Chars, p.pos, len(p.Chars)-p.pos)
	p.clear()
	p.update()
}

// Transpose transposes the char before the cursor position with the one on the
// cursor position, and moves the cursor one char to the right, if possible. If
// the cursor is at the beginning of the line it transposes the current char
// with the next one, and moves the cursor two chars to the right. This
// resembles Zshell's behaviour.
func (p *Prompt) Transpose() {
	if len(p.Chars) < 2 {
		return
	}

	offset := 0
	if p.pos == len(p.Chars) {
		offset = 2
	} else if p.pos > 0 {
		offset = 1
	}

	p.SetCursor(p.pos - offset)
	p.Chars = swap(p.Chars, p.pos, p.pos+1)
	p.Write(p.Chars[p.pos:])
	p.SetCursor(min(p.pos+2, len(p.Chars)))
}

// Discard discards the given input by moving to the next line and starting
// over with an empty prompt. This resembles the behaviour Ctrl-C in Bash or
// Zshell.
func (p *Prompt) Discard() {
	p.Newline()
	p.Reset()
}

// CompleteNext displays the next completion from the given slice.
func (p *Prompt) CompleteNext(strs [][]byte) {
	p.Complete(strs, Forw)
}

// CompletePrev displays the previous completion from the given slice.
func (p *Prompt) CompletePrev(strs [][]byte) {
	p.Complete(strs, Back)
}

// Complete displays the previous or next completion from the given slice
// depending on the given direction.
func (p *Prompt) Complete(strs [][]byte, dir int) {
	p.cycle(strs, Comp, dir)
}

// HistoryNext displays the next line from the given slice.
func (p *Prompt) HistoryNext(strs [][]byte) {
	p.History(strs, Forw)
}

// HistoryPrev displays the previous line from the given slice.
func (p *Prompt) HistoryPrev(strs [][]byte) {
	p.History(strs, Back)
}

// History displays the previous or next line from the given slice
// depending on the given direction.
func (p *Prompt) History(strs [][]byte, dir int) {
	p.cycle(strs, Hist, dir)
}

func (p *Prompt) cycle(strs [][]byte, mode int, dir int) {
	c := NewSet(strs, p.Chars, mode)
	if p.set == nil || !p.set.eq(c) {
		p.set = c
	}

	b := []byte{}
	if mode == Comp {
		b = trimLastWord(p.Chars)
		if len(b) > 0 {
			b = concat(b, space)
		}
	}

	if dir == Back {
		b = concat(b, p.set.Prev())
	} else {
		b = concat(b, p.set.Next())
	}

	p.Set(b)
}

// Suggest appends the first matching suggestion from the given slice in green
// after the current cursor position.
func (p *Prompt) Suggest(strs [][]byte) {
	c := NewSet(strs, lastWord(p.Chars), Sugg)
	w := lastWord(p.Chars)
	s := bytes.TrimPrefix(c.Next(), w)

	p.Suggested = s
	p.clear()
	p.Write(concat(p.Chars[p.pos:], ansi.Colored(ansi.Green, p.Suggested)))
	p.SetCursor()
}

// Newline writes a newline char to the terminal.
func (p *Prompt) Newline() {
	p.term.Newline()
}

// Return moves the cursor to the beginning of the line.
func (p *Prompt) Return() {
	p.SetCursor(0)
}

// End moves the cursor to the end of the line.
func (p *Prompt) End() {
	p.SetCursor(len(p.Chars))
}

// Set sets the content of the editor to the given line.
func (p *Prompt) Set(b []byte) {
	p.SetCursor(0)
	p.clear()
	p.Chars = b
	p.Write(b)
	p.pos += len(b)
	p.update()
}

// Reset resets the editor and starts over with an empty line.
func (p *Prompt) Reset() {
	p.reset()
	p.refresh()
}

// Write writes the given chars to the terminal.
func (p *Prompt) Write(b []byte) {
	p.term.Write(b)
}

// SetCursor sets the cursor to the given position, defaults to the current
// position.
func (p *Prompt) SetCursor(pos ...int) {
	if len(pos) > 0 {
		p.pos = pos[0]
	}
	p.term.SetCursor(p.pos + len(p.Prompt))
}

// MoveCursor moves the cursor by the given number of chars in the given
// direction.
func (p *Prompt) MoveCursor(i int, dir int) {
	if dir == Forw {
		if p.pos == len(p.Chars)-1 {
			return
		}
		p.pos += i
	} else {
		if p.pos == 0 {
			return
		}
		p.pos -= i
	}
	p.term.MoveCursor(i, dir)
}

// Del writes a delete char (`\x7F`) to the terminal.
func (p *Prompt) Del(i ...int) {
	p.term.Del(i...)
}

// Stop stops the terminal.
func (p *Prompt) Stop() {
	p.Newline()
	p.term.ClearLine()
	p.term.Stop()
}

// Str returns the current line.
func (p *Prompt) Str() string {
	return string(p.Chars)
}

func (p *Prompt) reset() {
	p.pos = 0
	p.Chars = []byte{}
	p.Suggested = []byte{}
	p.set = nil
}

func (p *Prompt) refresh() {
	p.update()
	p.clearLine()
	if len(p.Chars) > 0 {
		p.Write(concat(p.Chars, p.Suggested))
	}
	p.SetCursor()
}

func (p *Prompt) update() {
	if len(p.Chars) == 0 {
		p.reset()
	}

	if hasTailingSpace(p.Chars) {
		p.set = nil
	}
}

func (p *Prompt) clearLine() {
	p.term.ClearLine()
	p.Write(p.Prompt)
}

func (p *Prompt) clear() {
	p.term.Clear()
}
