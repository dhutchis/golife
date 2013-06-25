// Package glife implements the game of life via a structured mesh parallel paradigm.
package glife

import (
	"bytes"
	//"fmt"
	"io"
)

// The zeroed cell has no neighbors, is dead and is at round 0
type cell struct {
	myCh       <-chan bool   // for receiving only
	neighborCh []chan<- bool // for sending only
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
type cellarr struct {
	ca [][]cell
	// store doneCh<- here?
}

/* could be implemented as a [][]cell
Other methods of interest:
	CreateAllDead(nr, nc int) Field
	ReadFieldFrom(re io.Reader) Field
*/
type Field interface {
	SetAlive(x, y int, alive bool)
	Alive(x, y int) bool
	String() string
	WriteTo(io.Writer) (n int, err error)
	ReadFrom(re io.Reader) Field
	// resize

	Run(numRounds int)
}

func (f *cellarr) SetAlive(x, y int, alive bool) {
	f.ca[x][y].alive = alive
}
func (f cellarr) Alive(x, y int) bool {
	return f.ca[x][y].alive
}

// String returns the cellarr in a printable matrix representation
func (f cellarr) String() string {
	var buf bytes.Buffer
	// this could be abstracted into a FOLDL call
	for r := 0; r < len(f.ca); r++ {
		for c := 0; c < len(f.ca[r]); c++ {
			buf.WriteString(f.ca[r][c].String())
			buf.WriteByte(' ')
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

// Run will simulate a number of rounds of execution on the given Field
//   Current implementation is blocking (will block until all cells finish)
func (f *cellarr) Run(numRounds int) {
	// Setup communication channels, reset round count to zero
	f.setupCells()

	doneCh := make(chan bool)
	for r := range f.ca {
		for c := range f.ca[r] {
			// spawn goroutine on each cell
			// 1. Send out state to neighbors
			// 2. Receive state from neighbors
			// 3. Update state
			// 4. Round++.  Go back to 1 until we did numRounds.
			// End. Send true on done channel.
			go f.ca[r][c].cellRun(numRounds, doneCh)
		}
	}
	// Receive on done channels.
	for i := 0; i < len(f.ca)*len(f.ca[0]); i++ {
		<-doneCh
	}

}

func (c *cell) cellRun(numRounds int, doneCh chan<- bool) {
	for round := 0; round < numRounds; round++ {
		// 1. Send out state to neighbors
		for _, ch := range c.neighborCh {
			ch <- c.alive
		}

		// 2. Receive state from neighbors
		numNeighborsAlive := 0
		for i := 0; i < len(c.neighborCh); i++ {
			if <-c.myCh {
				numNeighborsAlive++
			}
		}

		// 3. Update state
		if c.alive {
			if numNeighborsAlive < 2 || numNeighborsAlive > 3 {
				c.alive = false
			} // else remain alive
		} else if numNeighborsAlive == 3 {
			c.alive = true
		} // else remain dead

		// 4. Round++.  Go back to 1 until we did numRounds.
		round++
	}
	// End. Send true on done channel.
	doneCh <- true
}

// What should the buffer size of neighbor channels be? Is deadlock possible?

// Advanced ideas: have cells send their state to a displayer goroutine that
//   updates a display for the user when all cells finish the 5th, 10th, ... round
// Option for Field to wrap from left to right and top to bottom and vice versa.

// ?: should this be a pointer method? I think so...
func (f *cellarr) setupCells() {
	for r := range f.ca {
		for c := range f.ca[r] {
			f.ca[r][c].round = 0
			f.ca[r][c].neighborCh = make([]chan<- bool, 0, 8) // len 0, cap 8
		}
	}
	for r := range f.ca {
		for c := range f.ca[r] {
			ch := make(chan bool, 8) // TODO: specify buffer size
			f.ca[r][c].myCh = ch

			rowlen, collen := len(f.ca), len(f.ca[0])
			if r != 0 && c != 0 {
				f.ca[r-1][c-1].neighborCh = append(f.ca[r-1][c-1].neighborCh, ch)
			}
			if r != 0 {
				f.ca[r-1][c].neighborCh = append(f.ca[r-1][c].neighborCh, ch)
			}
			if r != 0 && c != collen-1 {
				f.ca[r-1][c+1].neighborCh = append(f.ca[r-1][c+1].neighborCh, ch)
			}
			if c != 0 {
				f.ca[r][c-1].neighborCh = append(f.ca[r][c-1].neighborCh, ch)
			}
			if c != collen-1 {
				f.ca[r][c+1].neighborCh = append(f.ca[r][c+1].neighborCh, ch)
			}
			if r != rowlen-1 && c != 0 {
				f.ca[r+1][c-1].neighborCh = append(f.ca[r+1][c-1].neighborCh, ch)
			}
			if r != rowlen-1 {
				f.ca[r+1][c].neighborCh = append(f.ca[r+1][c].neighborCh, ch)
			}
			if r != rowlen-1 && c != collen-1 {
				f.ca[r+1][c+1].neighborCh = append(f.ca[r+1][c+1].neighborCh, ch)
			}

		}
	}

}
