// Package glife implements the game of life via a structured mesh parallel paradigm.
package glife

import (
	"bytes"
)

// The zeroed cell has no neighbors, is dead and is at round 0
type cell struct {
	neighborCh []chan bool
	alive      bool
	round      int
}

func (c cell) String() string {
	if c.alive {
		return "O"
	} else {
		return "."
	}
}

// The zeroed cellarr has all zeroed cells
type cellarr [][]cell

// could be implemented as a [][]cell
type Field interface {
	SetAlive(x, y int, alive bool)
	Alive(x, y int) bool
	String() string
}

func CreateAllDead(nr, nc int) Field {
	f := make([][]cell, nr)
	for r := 0; r < nr; r++ {
		f[r] = make([]cell, nc)
	}
	return cellarr(f)
}

func (f cellarr) SetAlive(x, y int, alive bool) {
	f[x][y].alive = alive
}
func (f cellarr) Alive(x, y int) bool {
	return f[x][y].alive
}

// String returns the cellarr in a printable matrix representation
func (f cellarr) String() string {
	var buf bytes.Buffer
	// this could be abstracted into a FOLDL call
	for r := 0; r < len(f); r++ {
		for c := 0; c < len(f[r]); c++ {
			buf.WriteString(f[r][c].String())
			buf.WriteByte(' ')
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
