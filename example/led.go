package main

import (
	"bytes"
	. "github.com/svenfuchs/led-go"
	. "github.com/svenfuchs/led-go/keys"
	"io/ioutil"
	"os"
)

func main() {
	l := NewLed("travis $ ")
	l.Handle(Enter, func(p *Prompt, k Key) { enter(p) })
	l.Handle(Chars, func(p *Prompt, k Key) { chars(p, k) })
	l.Handle(Backspace, func(p *Prompt, k Key) { back(p, k) })
	l.Handle(Delete, func(p *Prompt, k Key) { delete(p, k) })
	l.Handle(Tab, func(p *Prompt, k Key) { tab(p, k) })
	l.Handle(ShiftTab, func(p *Prompt, k Key) { shiftTab(p, k) })
	l.Handle(Up, func(p *Prompt, k Key) { prev(p) })
	l.Handle(Down, func(p *Prompt, k Key) { next(p) })
	l.Run()
}

var cmds = [][]byte{
	[]byte("repo"),
	[]byte("repos"),
	[]byte("user"),
	[]byte("users"),
}

func enter(p *Prompt) {
	p.Pause()
	println("\n\rEntered: " + p.Str())
	historyAdd(p.Str())
	p.Resume()
	p.Reset()
}

func chars(p *Prompt, k Key) {
	p.Insert(k.Chars)
	suggest(p)
}

func back(p *Prompt, k Key) {
	p.Back()
	suggest(p)
}

func delete(p *Prompt, k Key) {
	p.Delete()
	suggest(p)
}

func tab(p *Prompt, k Key) {
	p.CompleteNext(cmds)
}

func shiftTab(p *Prompt, k Key) {
	p.CompletePrev(cmds)
}

func suggest(p *Prompt) {
	p.Suggest(cmds)
}

func prev(p *Prompt) {
	p.HistoryPrev(history())
}

func next(p *Prompt) {
	p.HistoryNext(history())
}

var filename = "/tmp/led.history"

func history() [][]byte {
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		return compact(bytes.Split(data, []byte("\n")))
	}
	return [][]byte{}
}

func historyAdd(line string) {
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(line + "\n"))
	f.Close()
}

func compact(lines [][]byte) [][]byte {
	r := [][]byte{}
	for _, l := range lines {
		if len(l) > 0 {
			r = append(r, l)
		}
	}
	return r
}
