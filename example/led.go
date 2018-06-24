package main

import (
	"bytes"
	e "github.com/svenfuchs/led-go"
	"io/ioutil"
	"os"
)

func main() {
	r := e.NewReadline("travis $ ")
	r.Handle(e.Enter, func(e *e.Ed, k e.Key) { enter(e) })
	r.Handle(e.Chars, func(e *e.Ed, k e.Key) { chars(e, k) })
	r.Handle(e.Backspace, func(e *e.Ed, k e.Key) { back(e, k) })
	r.Handle(e.Delete, func(e *e.Ed, k e.Key) { delete(e, k) })
	r.Handle(e.Tab, func(e *e.Ed, k e.Key) { tab(e, k) })
	r.Handle(e.ShiftTab, func(e *e.Ed, k e.Key) { shiftTab(e, k) })
	r.Handle(e.Up, func(e *e.Ed, k e.Key) { prev(e) })
	r.Handle(e.Down, func(e *e.Ed, k e.Key) { next(e) })
	r.Run()
}

var cmds = [][]byte{
	[]byte("repo"),
	[]byte("repos"),
	[]byte("user"),
	[]byte("users"),
}

func enter(e *e.Ed) {
	e.Pause()
	println("\n\rEntered: " + e.Str())
	historyAdd(e.Str())
	e.Resume()
	e.Reset()
}

func chars(e *e.Ed, k e.Key) {
	e.Insert(k.Chars)
	suggest(e)
}

func back(e *e.Ed, k e.Key) {
	e.Back()
	suggest(e)
}

func delete(e *e.Ed, k e.Key) {
	e.Delete()
	suggest(e)
}

func tab(e *e.Ed, k e.Key) {
	e.CompleteNext(cmds)
}

func shiftTab(e *e.Ed, k e.Key) {
	e.CompletePrev(cmds)
}

func suggest(e *e.Ed) {
	// e.Suggest(cmds)
}

func prev(e *e.Ed) {
	e.HistoryPrev(history())
}

func next(e *e.Ed) {
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
