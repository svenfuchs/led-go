package led

import (
	"bytes"
	"github.com/pkg/term"
	. "github.com/svenfuchs/led-go/ansi"
	"github.com/svenfuchs/led-go/keys"
)

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

type Term struct {
	tty Iterm
	pos int
}

func (t *Term) Read() chan keys.Key {
	return keys.Read(t.tty)
}

func (t *Term) Write(b []byte) {
	t.tty.Write(b)
}

func (t *Term) Del(i ...int) {
	if len(i) == 0 {
		i = []int{1}
	}
	t.tty.Write(bytes.Repeat(Ansi(Del), i[0]))
}

func (t *Term) Newline() {
	t.Write(Ansi(Newline))
}

func (t *Term) Return() {
	t.Write(Ansi(Cr))
}

func (t *Term) ClearLine() {
	t.Return()
	t.Clear()
}

func (t *Term) Clear() {
	t.Write(Ansi(Clear))
}

func (t *Term) ShowCursor() {
	t.Write(Ansi(ShowCursor))
}

func (t *Term) HideCursor() {
	t.Write(Ansi(HideCursor))
}

func (t *Term) SetCursor(pos int) {
	t.Write(SetCursor(pos))
}

func (t *Term) MoveCursor(i int, dir int) {
	t.Write(MoveCursor(i, dir))
}

func (t *Term) Pause() {
	t.tty.Restore()
}

func (t *Term) Resume() {
	t.tty.RawMode()
}

func (t *Term) Stop() {
	t.tty.Restore()
	t.tty.Close()
}

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
