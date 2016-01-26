package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Cell - Cells are either living or dead
type Cell bool

// Alive - Convenience function for determining whether a Cell is currently alive or dead.
func (c Cell) Alive() bool {
	return bool(c)
}

// String - Converts a Cell to a string for output
func (c Cell) String() string {
	if c {
		return "X"
	}
	return " "
}

// Board - the Game of Life game board!
type Board [][]Cell

// NewBoard - initializes a new Board pointer with all Cells dead.
func NewBoard(length, width int) Board {
	board := make([][]Cell, 0, length)
	for i := 0; i < length; i++ {
		row := make([]Cell, width, width)
		board = append(board, row)
	}
	return board
}

// CheckCell - Checks whether a a cell is alive or dead (performs bounds checking)
func (b Board) CheckCell(row int, col int) bool {

	if row < 0 || col < 0 || row > len(b)-1 || col > len(b[0])-1 {
		return false
	}

	return b[row][col].Alive()
}

// Equals - Does this Board have a sequence of cells that exactly match this other sequence?
func (b Board) Equals(b2 Board) bool {
	return BoardEqual(b, b2)
}

// BoardEqual - Determines whether to boards are equal
func BoardEqual(b Board, b2 Board) bool {

	if len(b) != len(b2) {
		return false
	}

	for rowi, row := range b {
		if len(row) != len(b2[rowi]) {
			return false
		}
		for coli, value := range row {
			if value != b2[rowi][coli] {
				return false
			}
		}
	}
	return true
}

// Game - A Game of Life Game!
type Game struct {
	// Length - length of the board
	Length int
	// Width - width of the board
	Width int
	// Board - the board
	Board
}

// NewGame - initialize a new Game pointer
func NewGame(length, width int) *Game {
	game := &Game{
		Length: length,
		Width:  width,
		Board:  NewBoard(length, width),
	}
	return game
}

// Tick - Advance the game by one turn
func (g *Game) Tick() {
	newBoard := NewBoard(g.Length, g.Width)

	for rowi, row := range g.Board {
		for coli, cell := range row {
			liveNeighbors := 0
			// top
			if g.Board.CheckCell(rowi-1, coli) {
				liveNeighbors++
			}

			// top right
			if g.Board.CheckCell(rowi-1, coli+1) {
				liveNeighbors++
			}

			// right
			if g.Board.CheckCell(rowi, coli+1) {
				liveNeighbors++
			}

			// bottom right
			if g.Board.CheckCell(rowi+1, coli+1) {
				liveNeighbors++
			}

			// bottom
			if g.Board.CheckCell(rowi+1, coli) {
				liveNeighbors++
			}

			// bottom left
			if g.Board.CheckCell(rowi+1, coli-1) {
				liveNeighbors++
			}

			// left
			if g.Board.CheckCell(rowi, coli-1) {
				liveNeighbors++
			}

			// top left
			if g.Board.CheckCell(rowi-1, coli-1) {
				liveNeighbors++
			}

			newBoard[rowi][coli] = cell

			if cell.Alive() {
				if liveNeighbors < 2 || liveNeighbors > 3 {
					newBoard[rowi][coli] = false
				}
			} else {
				if liveNeighbors == 3 {
					newBoard[rowi][coli] = true
				}
			}
		}
	}
	g.Board = newBoard
}

// PrintBoard - prints the cells
func (b Board) PrintBoard() {
	for _, row := range b {
		for _, col := range row {
			fmt.Printf("%v ", col)
		}
		fmt.Println("")
	}

	for _ = range b[0] {
		fmt.Print("- ")
	}
	fmt.Println("")
}

// SeedBoard - randomly seeds the board with dead and living cells
func (b Board) SeedBoard() {
	for i, row := range b {
		for y := range row {
			b[i][y] = RandomCell()
		}
	}
}

// RandomCell - Randomly sets a cell to living or dead
func RandomCell() Cell {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := r.Int()
	if randomInt%4 == 0 {
		return true
	}

	return false
}

// CheckIfSequence - Checks if there is a sequence
func CheckIfSequence(boards []Board, b Board) bool {
	for i := len(boards) - 1; i >= 0; i-- {
		if b.Equals(boards[i]) {
			fmt.Printf("Sequence detected %v iterations back.\n", len(boards)-i)
			return true
		}
	}
	return false
}

func exitError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "Game of Life"
	app.Usage = "play the game of life"

	app.Action = func(c *cli.Context) {
		if len(c.Args()) <= 0 {
			exitError("Usage: %v <length> <width> <num_rounds>", app.Name)
		}

		length, err := strconv.Atoi(c.Args()[0])
		if err != nil {
			exitError(err.Error())
		}
		width, err := strconv.Atoi(c.Args()[1])
		if err != nil {
			exitError(err.Error())
		}

		rounds, err := strconv.Atoi(c.Args()[2])
		if err != nil {
			exitError(err.Error())
		}

		game := NewGame(length, width)

		game.SeedBoard()
		game.PrintBoard()

		var previousBoards []Board
		b := game.Board

		previousBoards = append(previousBoards, b)

		for i := 0; i < rounds; i++ {
			game.Tick()

			game.PrintBoard()
			b = game.Board

			if CheckIfSequence(previousBoards, b) {
				break
			}

			previousBoards = append(previousBoards, b)

		}

	}

	app.Run(os.Args)

}
