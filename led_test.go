package led

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestEd(t *testing.T) {
	prompt, term := setup()
	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
	})
}

// Chars

func TestKeys(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo ")
	receive(term, "bar")

	assert.Equal(t, "foo bar", prompt.Str())
	assert.Equal(t, 7, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo ",
		"bar",
	})
}

// Left

func TestLeftAtEnd(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo")
	prompt.Left()
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo",
		"<lft-1>",
	})
}

func TestLeftAtStart(t *testing.T) {
	prompt, term := setup()
	prompt.Left()
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
	})
}

// Right

func TestRightAtEnd(t *testing.T) {
	prompt, term := setup()
	prompt.Right()
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
	})
}

func TestRightAtStart(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo")
	prompt.Return()
	prompt.Right()
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo",
		"<cr><rgt-4>",
		"<rgt-1>",
	})
}

// Set

func TestSet(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo")
	prompt.Set([]byte("bar"))

	assert.Equal(t, "bar", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo",
		"<cr><rgt-4><clear>bar",
	})
}

// Insert

func TestInsert(t *testing.T) {
	prompt, term := setup()
	prompt.Insert([]byte("foo "))
	prompt.Insert([]byte("bar"))

	assert.Equal(t, "foo bar", prompt.Str())
	assert.Equal(t, 7, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo ",
		"bar",
	})
}

// Reject

func TestReject(t *testing.T) {
	prompt, term := setup()
	prompt.Reject([]byte("foo"))

	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"<red>foo<reset>",
		"<cr><rgt-4>",
		"<clear>",
	})
}

// Back

func TestBackAtEnd(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	assert.Equal(t, "bar", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)

	prompt.Back()
	prompt.Back()

	assert.Equal(t, "b", prompt.Str())
	assert.Equal(t, 1, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<lft-1><del>",
		"<lft-1><del>",
	})
}

func TestBackInMiddle(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Left()
	prompt.Back()

	assert.Equal(t, "br", prompt.Str())
	assert.Equal(t, 1, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<lft-1>",
		"<lft-1><del>",
	})
}

func TestBackAtStart(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Return()
	prompt.Back()

	assert.Equal(t, "bar", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<cr><rgt-4>",
	})
}

// BackWord

func TestBackWordAtStart(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Return()
	prompt.BackWord()

	assert.Equal(t, "bar", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<cr><rgt-4>",
	})
}

func TestBackWordAtEnd(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.BackWord()

	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<lft-3><del><del><del>",
	})
}

func TestBackWordAtEndSpace(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar  ")
	prompt.BackWord()

	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar  ",
		"<lft-5><del><del><del><del><del>",
	})
}

func TestBackWordInWord(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar  foo")
	prompt.Left()
	prompt.BackWord()
	prompt.BackWord()

	assert.Equal(t, "o", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar  foo",
		"<lft-1>",
		"<lft-2><del><del>",
		"<lft-5><del><del><del><del><del>",
	})
}

// Delete

func TestDeleteAtEnd(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Delete()

	assert.Equal(t, "bar", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
	})
}

func TestDeleteInMiddle(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Left()
	prompt.Left()
	prompt.Delete()

	assert.Equal(t, "br", prompt.Str())
	assert.Equal(t, 1, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<lft-1>",
		"<lft-1>",
		"<del>",
	})
}

func TestDeleteAtStart(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Return()
	prompt.Delete()

	assert.Equal(t, "ar", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<cr><rgt-4>",
		"<del>",
	})
}

// DeleteFromCursor

func TestDeleteFromCursorAtEnd(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo")
	prompt.DeleteFromCursor()

	assert.Equal(t, "foo", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo",
	})
}

func TestDeleteFromCursorInMiddle(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo")
	prompt.Left()
	prompt.Left()
	prompt.DeleteFromCursor()

	assert.Equal(t, "f", prompt.Str())
	assert.Equal(t, 1, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo",
		"<lft-1>",
		"<lft-1>",
		"<clear>",
	})
}

func TestDeleteFromCursorAtStart(t *testing.T) {
	prompt, term := setup()
	receive(term, "foo")
	prompt.Return()
	prompt.DeleteFromCursor()

	assert.Equal(t, "", prompt.Str())
	assert.Equal(t, 0, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"foo",
		"<cr><rgt-4>",
		"<clear>",
	})
}

// Transpose

func TestTransposeAtStart(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Return()
	prompt.Transpose()

	assert.Equal(t, "abr", prompt.Str())
	assert.Equal(t, 2, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<cr><rgt-4>",
		"<cr><rgt-4>abr<cr><rgt-6>",
	})
}

func TestTransposeInMiddle1(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Left()
	prompt.Left()
	prompt.Transpose()

	assert.Equal(t, "abr", prompt.Str())
	assert.Equal(t, 2, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<lft-1>",
		"<lft-1>",
		"<cr><rgt-4>abr<cr><rgt-6>",
	})
}

func TestTransposeInMiddle2(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Left()
	prompt.Transpose()

	assert.Equal(t, "bra", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<lft-1>",
		"<cr><rgt-5>ra<cr><rgt-7>",
	})
}

func TestTransposeAtEnd(t *testing.T) {
	prompt, term := setup()
	receive(term, "bar")
	prompt.Transpose()

	assert.Equal(t, "bra", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"bar",
		"<cr><rgt-5>ra<cr><rgt-7>",
	})
}

// History

func TestHistoryNextEmpty(t *testing.T) {
	h := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("baz"),
	}
	prompt, term := setup()
	prompt.HistoryNext(h)

	assert.Equal(t, "foo", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"<cr><rgt-4><clear>foo",
	})
}

func TestHistoryPrevEmpty(t *testing.T) {
	h := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("baz"),
	}
	prompt, term := setup()
	prompt.HistoryPrev(h)

	assert.Equal(t, "baz", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"<cr><rgt-4><clear>baz",
	})
}

// Complete

func TestCompleteEmpty(t *testing.T) {
	c := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
	}
	prompt, term := setup()
	prompt.CompleteNext(c)
	prompt.CompleteNext(c)

	assert.Equal(t, "bar", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"<cr><rgt-4><clear>foo",
		"<cr><rgt-4><clear>bar",
	})
}

func TestCompleteWithChar(t *testing.T) {
	c := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("baz"),
	}
	prompt, term := setup()
	receive(term, "b")
	prompt.CompleteNext(c)
	prompt.CompleteNext(c)

	assert.Equal(t, "baz", prompt.Str())
	assert.Equal(t, 3, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"b",
		"<cr><rgt-4><clear>bar",
		"<cr><rgt-4><clear>baz",
	})
}

// History & Completion

func TestHistoryAndCompletion(t *testing.T) {
	h := [][]byte{
		[]byte("foo 1"),
		[]byte("foo 2"),
		[]byte("bar 1"),
		[]byte("bar 2"),
	}
	c := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("baz"),
	}
	prompt, _ := setup()

	prompt.CompleteNext(c)
	assert.Equal(t, "foo", prompt.Str())

	prompt.CompleteNext(c)
	assert.Equal(t, "bar", prompt.Str())

	prompt.HistoryNext(h)
	assert.Equal(t, "bar 1", prompt.Str())

	prompt.HistoryNext(h)
	assert.Equal(t, "bar 2", prompt.Str())

	prompt.HistoryNext(c)
	assert.Equal(t, "bar", prompt.Str())
}

// Suggest

func TestSuggest(t *testing.T) {
	c := []byte("foo")
	prompt, term := setup()
	receive(term, "b")
	prompt.Suggest(c)

	assert.Equal(t, "b", prompt.Str())
	assert.Equal(t, 1, prompt.Pos)
	assertOut(t, term, []string{
		"<cr><clear>t ~ <cr><rgt-4>",
		"b<clear>",
		"<green>ar<reset><cr><rgt-5>",
	})
}

func assertOut(t *testing.T, term *testTerm, strs []string) {
	out := string(Deansi([]byte(term.out)))
	assert.Equal(t, strings.Join(strs, ""), out)
}

func setup() (*Ed, *testTerm) {
	term := newTestTerm()
	prompt := NewReadline("t ~ ", term)
	go prompt.Run()
	time.Sleep(1 * time.Millisecond)
	return prompt, term
}

func receive(t *testTerm, str string) {
	t.keys <- str
	time.Sleep(1 * time.Millisecond)
}

func reset(term *testTerm) {
	term.out = ""
}

func p(s string) {
	println("\"" + s + "\"")
}

func key(k int) string {
	return Keys[k].Str()
}

func newTestTerm() *testTerm {
	k := make(chan string, 1)
	t := testTerm{keys: k}
	return &t
}

type testTerm struct {
	keys chan string
	out  string
}

func (t *testTerm) Start() {
}

func (t *testTerm) Read(b []byte) (int, error) {
	a := <-t.keys
	copy(b, a)
	return len(a), nil
}

func (t *testTerm) Write(b []byte) (int, error) {
	// b = Deansi(b)
	t.out = t.out + string(b)
	return 1, nil
}

func (t *testTerm) Restore() error {
	return nil
}

func (t *testTerm) Close() error {
	close(t.keys)
	return nil
}

func (t *testTerm) RawMode() error {
	return nil
}
