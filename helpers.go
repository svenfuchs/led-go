package led

import (
	"bytes"
)

func min(i ...int) int {
	min := i[0]
	for _, num := range i[0:] {
		if num < min {
			min = num
		}
	}
	return min
}

func trimSpace(b []byte) []byte {
	for bytes.HasSuffix(b, space) {
		bytes.TrimSuffix(b, space)
	}
	return b
}

func trimLastWord(b []byte) []byte {
	a := bytes.Split(b, space)
	return bytes.Join(a[:len(a)-1], space)
}

func firstWord(b []byte) []byte {
	w := bytes.Split(b, space)
	if len(w) > 0 {
		return w[0]
	}
	return []byte{}
}

func lastWord(b []byte, skipSpace ...bool) []byte {
	w := bytes.Split(b, space)
	if len(w) == 0 {
		return []byte{}
	} else if len(skipSpace) > 0 && skipSpace[0] {
		b = []byte{}
		for i := len(w) - 1; i >= 0; i-- {
			if len(w[i]) == 0 {
				b = append(b, space...)
			} else {
				return append(w[i], b...)
			}
		}
		return b
	} else {
		return w[len(w)-1]
	}
}

func hasTailingSpace(b []byte) bool {
	return len(b) > 0 && b[len(b)-1] == ' '
}

func insert(a []byte, b []byte, i int) []byte {
	return append(a[:i], append(b, a[i:]...)...)
}

func delete(a []byte, i int, c int) []byte {
	return append(a[:i], a[i+c:]...)
}

func swap(b []byte, lft int, rgt int) []byte {
	c := b[lft]
	b[lft] = b[rgt]
	b[rgt] = c
	return b
}

func concat(a []byte, b ...[]byte) []byte {
	r := append(a, b[0]...)
	if len(b) > 1 {
		r = concat(r, b[1:]...)
	}
	return r
}

func indexOf(strs [][]byte, str []byte) int {
	for i, s := range strs {
		if bytes.Equal(s, str) {
			return i
		}
	}
	return -1
}
