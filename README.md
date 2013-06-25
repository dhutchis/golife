glife
=====
Side project in development.  This is a Google Go implementation of the [Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life) using message passing via the structured grid parallel programming paradigm.  Yes, there is a really easy [Matlab matrix shifting implementation](http://www.mathworks.com/moler/exm/chapters/life.pdf) of the Game of Life, but I want to give this a shot for learning and comparison.

Interfaces with the [plaintext *.cells* file format](http://conwaylife.com/wiki/Plaintext).

## Installation ##

	go get github.com/denine99/glife/glife


## Misc ##
Here's some output from a Blinker:

    O O . . 
	O . . . 
	. . . O 
	. . O O 
	
	O O . . 
	O O . . 
	. . O O 
	. . O O 
	
	O O . . 
	O . . . 
	. . . O 
	. . O O 
	
	O O . . 
	O O . . 
	. . O O 
	. . O O 

	O O . . 
	O . . . 
	. . . O 
	. . O O 

Next todo: write some tests using a string to load a Field, and a string to compare the expected state of a Field after running it N times.