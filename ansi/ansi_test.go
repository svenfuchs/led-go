package ansi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeansi(t *testing.T) {
	assert.Equal(t, "<clear>", deansi(Ansi(Clear)))
	assert.Equal(t, "<clear-line>", deansi(Ansi(ClearLine)))
	assert.Equal(t, "<show-crsr>", deansi(Ansi(ShowCursor)))
	assert.Equal(t, "<hide-crsr>", deansi(Ansi(HideCursor)))
	assert.Equal(t, "<nl>", deansi(Ansi(Newline)))
	assert.Equal(t, "<cr>", deansi(Ansi(Cr)))
	assert.Equal(t, "<del>", deansi(Ansi(Del)))
	assert.Equal(t, "<red>", deansi(Ansi(Red)))
	assert.Equal(t, "<green>", deansi(Ansi(Green)))
	assert.Equal(t, "<reset>", deansi(Ansi(Reset)))
	assert.Equal(t, "<cr><rgt-5>", deansi(SetCursor(5)))
	assert.Equal(t, "<lft-5>", deansi(MoveCursor(5, Lft)))
	assert.Equal(t, "<rgt-5>", deansi(MoveCursor(5, Rgt)))
}

func deansi(b []byte) string {
	return string(Deansi(b))
}
