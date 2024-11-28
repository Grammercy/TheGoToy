package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type Board struct {
	arr [][]Particle
}

type Particle struct {
	isEmpty   bool
	name      string
	xVelocity float64
	yVelocity float64
}

const (
  borderWidth  = 2
	windowWidth  = 2460
	windowHeight = 1500
	gridX        = 20
	gridY        = 20
)

var gridSize int = determineSquareSize(windowWidth, windowHeight, gridX, gridY)

func main() {
	var board Board
  board.setupBoard()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, renderer, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()
}

func determineSquareSize(totalX, totalY, gridX, gridY int) int {
  gridX += borderWidth
  gridY += borderWidth
  gridSize := totalX/gridX
  if totalY/gridY < gridSize {
    gridSize = totalY/gridY
  }
  return gridSize
}

func (board Board) setupBoard() {
	board.arr = make([][]Particle, 20)
	for i := range board.arr {
		board.arr[i] = make([]Particle, 20)
	}
}

func (board Board) String() string {
	str := ""
	for i := range board.arr {
		for j := range board.arr {
			if !board.arr[i][j].isEmpty {
				str += " "
			} else {
				str += "󰝤 "
			}
		}
		str += "\n"
	}
	return str
}
