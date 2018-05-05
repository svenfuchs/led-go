package led

import (
	"bytes"
	"github.com/svenfuchs/led-go/ansi"
	"github.com/svenfuchs/led-go/keys"
	"time"
)

const (
	Forw int = iota
	Back
	Hist
	Comp
	Sugg
)

var space = []byte{' '}

func NewPrompt(prompt string, t ...Iterm) *Prompt {
	return &Prompt{
		term:      StartTerm(t...),
		handlers:  map[int]func(*Prompt, keys.Key){},
		Prompt:    []byte(prompt),
		pos:       0,
		Chars:     []byte{},
		suggested: []byte{},
	}
}

type Prompt struct {
	term      *Term
	handlers  map[int]func(*Prompt, keys.Key)
	Prompt    []byte
	pos       int
	Chars     []byte
	suggested []byte
	set       *Set
}

func (p *Prompt) Handle(key int, handler func(*Prompt, keys.Key)) {
	p.handlers[key] = handler
}

func (p *Prompt) Run() {
	p.Refresh()
	for k := range p.term.Read() {
		if p.handlers[k.Code] != nil {
			p.handlers[k.Code](p, k)
		}
	}
	p.Stop()
}

func (p *Prompt) Pause() {
	p.term.Pause()
}

func (p *Prompt) Resume() {
	p.term.Resume()
}

func (p *Prompt) Left() {
	if p.pos > 0 {
		p.MoveCursor(1, Back)
	}
}

func (p *Prompt) Right() {
	if p.pos < len(p.Chars) {
		p.MoveCursor(1, Forw)
	}
}

func (p *Prompt) Append(b []byte) {
	p.Chars = append(p.Chars, b...)
	p.pos = len(p.Chars) - 1
	p.Write(b)
	p.Update()
}

func (p *Prompt) Insert(b []byte) {
	p.Chars = insert(p.Chars, b, p.pos)
	p.Write(b)
	p.pos += len(b)
	p.Update()
}

func (p *Prompt) Reject(chars []byte) {
	p.term.Write(ansi.Colored(ansi.Red, chars))
	time.Sleep(250 * time.Millisecond)
	p.SetCursor()
	p.Clear()
}

func (p *Prompt) Back() {
	if p.pos == 0 {
		return
	}

	p.MoveCursor(1, Back)
	p.Chars = delete(p.Chars, p.pos, 1)
	p.Del()
	p.Update()
}

func (p *Prompt) BackWord() {
	if p.pos == 0 {
		return
	}

	w := lastWord(p.Chars[:p.pos], true)
	p.MoveCursor(len(w), Back)
	p.Chars = delete(p.Chars, p.pos, len(w))
	p.Del(len(w))
	p.Update()
}

func (p *Prompt) Delete() {
	if p.pos == len(p.Chars) {
		return
	}

	p.Chars = delete(p.Chars, p.pos, 1)
	p.Del()
	p.Update()
}

func (p *Prompt) DeleteFromCursor() {
	if p.pos == len(p.Chars) {
		return
	}

	p.Chars = delete(p.Chars, p.pos, len(p.Chars)-p.pos)
	p.Clear()
	p.Update()
}

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

func (p *Prompt) Discard() {
	p.Newline()
	p.Reset()
}

func (p *Prompt) CompleteNext(strs [][]byte) {
	p.Complete(strs, Forw)
}

func (p *Prompt) CompletePrev(strs [][]byte) {
	p.Complete(strs, Back)
}

func (p *Prompt) Complete(strs [][]byte, dir int) {
	p.Cycle(strs, Comp, dir)
}

func (p *Prompt) HistoryNext(strs [][]byte) {
	p.History(strs, Forw)
}

func (p *Prompt) HistoryPrev(strs [][]byte) {
	p.History(strs, Back)
}

func (p *Prompt) History(strs [][]byte, dir int) {
	p.Cycle(strs, Hist, dir)
}

func (p *Prompt) Cycle(strs [][]byte, mode int, dir int) {
	c := NewSet(strs, p.Chars, mode)
	if p.set == nil || !p.set.Eq(c) {
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

func (p *Prompt) Suggest(strs [][]byte) {
	c := NewSet(strs, lastWord(p.Chars), Sugg)
	w := lastWord(p.Chars)
	s := bytes.TrimPrefix(c.Next(), w)

	p.suggested = s
	p.Clear()
	p.Write(concat(p.Chars[p.pos:], ansi.Colored(ansi.Green, p.suggested)))
	p.SetCursor()
}

func (p *Prompt) Suggesting(b []byte) []byte {
	if len(p.suggested) == 0 {
		return b
	}
	return concat(b, p.suggested)
}

func (p *Prompt) Suggested() []byte {
	return p.suggested
}

func (p *Prompt) Newline() {
	p.term.Newline()
}

func (p *Prompt) Return() {
	p.SetCursor(0)
}

func (p *Prompt) End() {
	p.SetCursor(len(p.Chars))
}

func (p *Prompt) Reset() {
	p.reset()
	p.Refresh()
}

func (p *Prompt) reset() {
	p.pos = 0
	p.Chars = []byte{}
	p.suggested = []byte{}
	p.set = nil
}

func (p *Prompt) Set(b []byte) {
	p.Chars = b
	p.SetCursor(0)
	p.Clear()
	p.Write(b)
	p.pos += len(b)
	p.Update()
}

func (p *Prompt) Write(b []byte) {
	p.term.Write(b)
}

func (p *Prompt) Refresh() {
	p.Update()
	p.clearLine()
	if len(p.Chars) > 0 {
		p.Write(p.Suggesting(p.Chars))
	}
	p.SetCursor()
}

func (p *Prompt) Update() {
	if len(p.Chars) == 0 {
		p.reset()
	}

	if hasTailingSpace(p.Chars) {
		p.set = nil
	}
}

func (p *Prompt) SetCursor(pos ...int) {
	if len(pos) > 0 {
		p.pos = pos[0]
	}
	p.term.SetCursor(p.pos + len(p.Prompt))
}

func (p *Prompt) MoveCursor(i int, dir int) {
	if dir == Forw {
		p.pos += i
	} else {
		p.pos -= i
	}
	p.term.MoveCursor(i, dir)
}

func (p *Prompt) clearLine() {
	p.term.ClearLine()
	p.Write(p.Prompt)
}

func (p *Prompt) Clear() {
	p.term.Clear()
}

func (p *Prompt) Del(i ...int) {
	p.term.Del(i...)
}

func (p *Prompt) Stop() {
	p.Newline()
	p.term.ClearLine()
	p.term.Stop()
}

func (p *Prompt) Str() string {
	return string(p.Chars)
}
