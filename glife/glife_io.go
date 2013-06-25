// glife_io.go contains functions to create/save/load Field structures.

package glife

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// CreateAllDead creates a new Field of given dimensions
func CreateAllDead(nr, nc int) Field {

	f := new(cellarr)
	f.ca = make([][]cell, nr)
	// allocate the 2D array all at once for efficiency, since it's seldom resized
	cells := make([]cell, nr*nc)
	for r := range f.ca {
		f.ca[r], cells = cells[:nc], cells[nc:]
	}
	return f
}

// ReadFieldFrom loads a Field from a Writer
func ReadFieldFrom(re io.Reader) Field {
	ca := new(cellarr)
	ca = ca.ReadFrom(re).(*cellarr)
	fmt.Println("ca:\n", ca)
	return ca
}

// writes the [][]cell to a Writer and returns the number of bytes written,
// and any error that may have occurred
func (f cellarr) WriteTo(wr io.Writer) (n int, err error) {
	w := bufio.NewWriter(wr)
	var p int
	for r := range f.ca {
		for c := range f.ca[r] {
			p, err = w.WriteString(f.ca[r][c].String())
			n += p
			if err != nil {
				return n, err
			}
		}
		p, err = w.WriteRune('\n')
		n += p
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// function ReadFieldFrom reads from a Writer completely and forms a Field from
//   the contents.  '.'==dead cell, 'O'==alive cell, '!' starts a comment line.
// Then, it copies the Field into the current object, first clearing all contents
//   and then copying in the new Field
func (f *cellarr) ReadFrom(re io.Reader) Field {
	// 1. Read in Field from re (the [][]cell may have unequal lengths)
	maxcols := 0 // keep track of largest column
	ftmp := make([][]cell, 0)
	scanner := bufio.NewScanner(re)
	for scanner.Scan() {
		text := scanner.Text()
		// ignore the lines beginnning with '!'
		if strings.HasPrefix(text, "!") {
			continue
		}
		// process a row
		if len(text) > maxcols {
			maxcols = len(text)
		}
		row := make([]cell, len(text))
		for i, t := range text {
			if t == 'O' {
				row[i].alive = true
			}
		}
		ftmp = append(ftmp, row)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println("ftmp:\n", cellarr{ftmp})

	// 2. Resize if necessary
	if len(f.ca) == 0 { // special case
		f = CreateAllDead(len(ftmp), maxcols).(*cellarr)
	}
	if maxcols > len(f.ca[0]) {
		//coldiff :=
		for r := range f.ca {
			f.ca[r] = append(f.ca[r], make([]cell, maxcols-len(f.ca[0]))...)
		}
	}
	for diff := len(ftmp) - len(f.ca); diff > 0; diff-- {
		f.ca = append(f.ca, make([]cell, len(f.ca[0])))
	}

	// 3. Copy over ftmp to f and clear anything outside
	zeroRow := make([]cell, len(f.ca[0]))
	var r int
	for r = 0; r < len(ftmp); r++ {
		// copy ftmp to f
		copy(f.ca[r][:len(ftmp[r])], ftmp[r])
		// zero remainder of row of f
		copy(f.ca[r][len(ftmp[r]):], zeroRow[:len(ftmp[r])])
	}
	for ; r < len(f.ca); r++ {
		// zero remaining rows of f past ftmp
		copy(f.ca[r], zeroRow)
	}
	fmt.Println("f:\n", f)
	return f
}
