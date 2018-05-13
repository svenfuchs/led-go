package main

import (
	"bytes"
	. "github.com/svenfuchs/led-go"
	. "github.com/svenfuchs/led-go/keys"
	"io/ioutil"
	"os"
)

func main() {
	r := NewReadline("travis $ ")
	r.Handle(Enter, func(e *Ed, k Key) { enter(e) })
	r.Handle(Chars, func(e *Ed, k Key) { chars(e, k) })
	r.Handle(Backspace, func(e *Ed, k Key) { back(e, k) })
	r.Handle(Delete, func(e *Ed, k Key) { delete(e, k) })
	r.Handle(Tab, func(e *Ed, k Key) { tab(e, k) })
	r.Handle(ShiftTab, func(e *Ed, k Key) { shiftTab(e, k) })
	r.Handle(Up, func(e *Ed, k Key) { prev(e) })
	r.Handle(Down, func(e *Ed, k Key) { next(e) })
	r.Run()
}

var cmds = [][]byte{
	[]byte("repo"),
	[]byte("repos"),
	[]byte("user"),
	[]byte("users"),
}

func enter(e *Ed) {
	e.Pause()
	println("\n\rEntered: " + e.Str())
	historyAdd(e.Str())
	e.Resume()
	e.Reset()
}

func chars(e *Ed, k Key) {
	e.Insert(k.Chars)
	suggest(e)
}

func back(e *Ed, k Key) {
	e.Back()
	suggest(e)
}

func delete(e *Ed, k Key) {
	e.Delete()
	suggest(e)
}

func tab(e *Ed, k Key) {
	e.CompleteNext(cmds)
}

func shiftTab(e *Ed, k Key) {
	e.CompletePrev(cmds)
}

func suggest(e *Ed) {
	e.Suggest(cmds)
}

func prev(e *Ed) {
	e.HistoryPrev(history())
}

func next(e *Ed) {
	e.HistoryNext(history())
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
