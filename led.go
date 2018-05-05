package led

import (
	. "github.com/svenfuchs/led-go/keys"
)

func NewLed(prompt string, t ...Iterm) *Prompt {
	p := NewPrompt(prompt, t...)
	p.Handle(Chars, func(p *Prompt, k Key) { p.Insert(k.Chars) })
	p.Handle(CtrlA, func(p *Prompt, k Key) { p.Return() })
	p.Handle(CtrlB, func(p *Prompt, k Key) { p.Left() })
	p.Handle(CtrlC, func(p *Prompt, k Key) { p.Discard() })
	p.Handle(CtrlD, func(p *Prompt, k Key) { p.Delete() })
	p.Handle(CtrlE, func(p *Prompt, k Key) { p.End() })
	p.Handle(CtrlF, func(p *Prompt, k Key) { p.Right() })
	p.Handle(CtrlK, func(p *Prompt, k Key) { p.DeleteFromCursor() })
	p.Handle(CtrlT, func(p *Prompt, k Key) { p.Transpose() })
	p.Handle(CtrlU, func(p *Prompt, k Key) { p.Reset() })
	p.Handle(CtrlW, func(p *Prompt, k Key) { p.BackWord() })
	p.Handle(Enter, func(p *Prompt, k Key) { p.Newline() })
	p.Handle(Backspace, func(p *Prompt, k Key) { p.Back() })
	p.Handle(Delete, func(p *Prompt, k Key) { p.Delete() })
	p.Handle(Left, func(p *Prompt, k Key) { p.Left() })
	p.Handle(Right, func(p *Prompt, k Key) { p.Right() })
	return p
}
