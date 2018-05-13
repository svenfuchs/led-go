package led

import (
	keys "github.com/svenfuchs/led-go/keys"
)

// NewLed creates a line editor that resembles most of Linenoise's functionality
func NewLed(prompt string, t ...Iterm) *Prompt {
	p := NewPrompt(prompt, t...)
	p.Handle(keys.Chars, func(p *Prompt, k keys.Key) { p.Insert(k.Chars) })
	p.Handle(keys.CtrlA, func(p *Prompt, k keys.Key) { p.Return() })
	p.Handle(keys.CtrlB, func(p *Prompt, k keys.Key) { p.Left() })
	p.Handle(keys.CtrlC, func(p *Prompt, k keys.Key) { p.Discard() })
	p.Handle(keys.CtrlD, func(p *Prompt, k keys.Key) { p.Delete() })
	p.Handle(keys.CtrlE, func(p *Prompt, k keys.Key) { p.End() })
	p.Handle(keys.CtrlF, func(p *Prompt, k keys.Key) { p.Right() })
	p.Handle(keys.CtrlK, func(p *Prompt, k keys.Key) { p.DeleteFromCursor() })
	p.Handle(keys.CtrlT, func(p *Prompt, k keys.Key) { p.Transpose() })
	p.Handle(keys.CtrlU, func(p *Prompt, k keys.Key) { p.Reset() })
	p.Handle(keys.CtrlW, func(p *Prompt, k keys.Key) { p.BackWord() })
	p.Handle(keys.Enter, func(p *Prompt, k keys.Key) { p.Newline() })
	p.Handle(keys.Backspace, func(p *Prompt, k keys.Key) { p.Back() })
	p.Handle(keys.Delete, func(p *Prompt, k keys.Key) { p.Delete() })
	p.Handle(keys.Left, func(p *Prompt, k keys.Key) { p.Left() })
	p.Handle(keys.Right, func(p *Prompt, k keys.Key) { p.Right() })
	return p
}
