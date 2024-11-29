package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type Board struct {
	arr [][]Particle
}

type Particle struct {
	isFull   bool
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
  running := true
  renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
  board.arr[19][0].isFull = true
  board.arr[19][0].yVelocity = -1.5
  for running {
    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

    board.updatePositions()
    board.render(renderer)
    renderer.Present()
    sdl.Delay(16)
  }
  fmt.Println("done")
}

func (board Board) render (renderer *sdl.Renderer) {
  for i := range board.arr {
    for j := range board.arr[i] {
      x := int32(j * gridSize)
      y := int32((len(board.arr)-i) * gridSize)
      renderer.SetDrawColor(176, 176, 176, 255)
      renderer.FillRect(&sdl.Rect{X: x, Y: y, W: int32(gridSize), H: int32(gridSize)})

      if board.arr[i][j].isFull {
        renderer.SetDrawColor(255, 255, 255, 255)
      } else {
        renderer.SetDrawColor(0, 0, 0, 255)
      }
      renderer.FillRect(&sdl.Rect{X: x + borderWidth, Y: y + borderWidth, W: int32(gridSize) - borderWidth, H: int32(gridSize) - borderWidth})
    }
  }
}

func (board *Board) updatePositions () {
  for i := range board.arr {
    for j := range board.arr[i] {
      if !board.arr[i][j].isFull {
        continue
      }
      adjustPostition(board, i, j)
    }
  }
}

func adjustPostition(board *Board, i, j int) {
  if int(board.arr[i][j].xVelocity) == 0 && int(board.arr[i][j].yVelocity) == 0 {
    return
  }
  if j + int(board.arr[i][j].xVelocity) >= 0 && j + int(board.arr[i][j].xVelocity) < len(board.arr[i]) && i + int(board.arr[i][j].yVelocity) >= 0 && i + int(board.arr[i][j].yVelocity) < len(board.arr) {
    if !board.arr[i + int(board.arr[i][j].yVelocity)][j + int(board.arr[i][j].xVelocity)].isFull {
      board.arr[i + int(board.arr[i][j].yVelocity)][j + int(board.arr[i][j].xVelocity)] = board.arr[i][j].clone()
      board.arr[i][j] = Particle{}
      return
    }
    board.arr[i][j].xVelocity *= 0.9
    board.arr[i][j].yVelocity *= 0.9
    adjustPostition(board, i, j)
  }
}

func (particle *Particle) passGravity () {
  particle.yVelocity -= 0.5
}

func (particle Particle) clone () Particle{
  return particle
}

func determineSquareSize(totalX, totalY, gridX, gridY int) int {
  gridX += 2 * borderWidth
  gridY += 2 * borderWidth
  gridSize := totalX/gridX
  if totalY/gridY < gridSize {
    gridSize = totalY/gridY
  }
  return gridSize
}

func (board *Board) setupBoard() {
	board.arr = make([][]Particle, 20)
	for i := range board.arr {
		board.arr[i] = make([]Particle, 20)
	}
}

func (board Board) String() string {
	str := ""
	for i := range board.arr {
		for j := range board.arr {
			if !board.arr[i][j].isFull {
				str += " "
			} else {
				str += "󰝤 "
			}
		}
		str += "\n"
	}
	return str
}
