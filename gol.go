package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Cell bool

func (c Cell) Alive() bool {
	if c {
		return true
	}
	return false
}

func (c Cell) String() string {
	switch c {
	case true:
		return "X"
	default:
		return " "
	}
}

type Board [][]Cell

func NewBoard(length int, width int) Board {
	board := [][]Cell{}
	for i := 0; i < length; i++ {
		row := make([]Cell, width, width)
		board = append(board, row)
	}
	return board
}

func (b Board) CheckCell(row int, col int) bool {

	if row < 0 || col < 0 || row > len(b)-1 || col > len(b[0])-1 {
		return false
	}

	return b[row][col].Alive()
}

func (b Board) Equals(b2 Board) bool {
	return BoardEqual(b, b2)
}

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

type Game struct {
	Length int
	Width  int
	Board
}

func NewGame(length int, width int) *Game {
	game := &Game{
		Length: length,
		Width:  width,
		Board:  NewBoard(length, width),
	}
	return game
}

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

func (b Board) PrintBoard() {
	for _, row := range b {
		for _, col := range row {
			fmt.Printf("%v ", col)
		}
		fmt.Println("")
	}

	for _, _ = range b[0] {
		fmt.Print("- ")
	}
	fmt.Println("")
}

func (b Board) SeedBoard() {
	for i, row := range b {
		for y, _ := range row {
			b[i][y] = RandomCell()
		}
	}
}

func RandomCell() Cell {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := r.Int()
	if randomInt%4 == 0 {
		return true
	}

	return false
}

func CheckIfSequence(boards []Board, b Board) bool {
	for i := len(boards) - 1; i >= 0; i-- {
		if b.Equals(boards[i]) {
			fmt.Printf("Sequence detected %v iterations back.\n", len(boards)-i)
			return true
		}
	}
	return false
}

func main() {

	app := cli.NewApp()
	app.Name = "Game of Life"
	app.Usage = "play the game of life"

	app.Action = func(c *cli.Context) {
		if len(c.Args()) <= 0 {
			log.Fatalf("Usage: %v <length> <width> <num_rounds>", app.Name)
		}

		length, err := strconv.Atoi(c.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		width, err := strconv.Atoi(c.Args()[1])
		if err != nil {
			log.Fatal(err)
		}

		rounds, err := strconv.Atoi(c.Args()[2])
		if err != nil {
			log.Fatal(err)
		}

		game := NewGame(length, width)

		game.SeedBoard()
		game.PrintBoard()

		previousBoards := make([]Board, 0, 0)
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
