package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
	"sync"
	"time"
)

type Board struct {
	arr [][]Particle
}

type Particle struct {
	isFull    bool
	partType  uint64
	xVelocity float64
	yVelocity float64
}

const (
	boardWidth      = 20//604
	boardHeight     = 30//376
	borderWidth     = 0
	windowWidth     = 2550
	windowHeight    = 1550
	gravityConstant = 0.2
)

const (
  POWDER = 100
)

var gridSize int = determineSquareSize(windowWidth, windowHeight, boardWidth, boardHeight)

func main() {
	start := time.Now()
	var board Board
	board.setupBoard()
	running := true
	paused := false
	averageFrameTime := time.Since(start)
	// numFrames := 0
	var boardSurface *sdl.Surface

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, renderer, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	boardSurface, err = sdl.CreateRGBSurface(0, windowWidth, windowHeight, 32, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("Time to setup was " + time.Since(start).String())

	for running {
		startOfFrameTime := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_SPACE:
						paused = !paused
					}
				}
			}
		}
		x, y, mouseHeld := sdl.GetMouseState()
		if mouseHeld != 0 {
			col := int(x) / gridSize
			row := len(board.arr) - (int(y) / gridSize) - 1
			if row >= 0 && col >= 0 && col < len(board.arr[0]) && row < len(board.arr) {
				if !board.arr[row][col].isFull {
          board.arr[row][col] = Particle{isFull: true, partType: POWDER}
				}
			}
		}
		// timeToHandleEvents := time.Since(startOfFrameTime)
		// _ = timeToHandleEvents
		// startOfCalcTime := time.Now()
		if !paused {
			board.passGravity()
      board.movePowders()
			board.updatePositions()
		}
		// startOfRenderTime := time.Now()
		board.render(renderer, boardSurface)
		renderer.Present()
		frameTime := time.Since(startOfFrameTime)
		averageFrameTime = averageFrameTime*time.Duration(199) + frameTime
		averageFrameTime /= time.Duration(200)
		runtime.GC()
    time.Sleep(16 * time.Millisecond - time.Since(startOfFrameTime))
	}
	fmt.Println("Average frame time was " + averageFrameTime.String())
	fmt.Println("done")
}

func (board Board) render(renderer *sdl.Renderer, boardSurface *sdl.Surface) {
	var wg sync.WaitGroup
	length := len(board.arr)
	for i := range board.arr {
		wg.Add(1)
		go func() {
			var color sdl.Color
			defer wg.Done()
			for j := range board.arr[i] {
				x := int32(j * gridSize)
				y := int32((length - i) * gridSize - gridSize) 
				if board.arr[i][j].isFull {
					color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
				} else {
					color = sdl.Color{R: 255, G: 20, B: 20, A: 20}
				}

				boardSurface.FillRect(&sdl.Rect{X: x + borderWidth, Y: y + borderWidth, W: int32(gridSize) - borderWidth, H: int32(gridSize) - borderWidth}, color.Uint32())
			}
		}()
	}
	wg.Wait()
	boardTexture, err := renderer.CreateTextureFromSurface(boardSurface)
	defer boardTexture.Destroy()
	if err != nil {
		panic(err)
	}
	renderer.Copy(boardTexture, nil, nil)
}

func (board *Board) movePowders() {
  for i := range board.arr {
    if i == 0 {
      continue
    }
    for j := range board.arr[i] {
      if board.arr[i][j].partType == POWDER && board.arr[i][j].isFull {
        if board.arr[i-1][j].isFull && board.arr[i-1][j].partType == POWDER {
          if j-1 >= 0 && !board.arr[i-1][j-1].isFull {
            board.arr[i-1][j-1], board.arr[i][j] = board.arr[i][j], board.arr[i-1][j-1]
          } else if j + 1 < len(board.arr[i]) && !board.arr[i-1][j+1].isFull {
            board.arr[i-1][j+1], board.arr[i][j] = board.arr[i][j], board.arr[i-1][j+1]
          }
        }
      }
    }
  }
}

func (board *Board) updatePositions() {
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
	if j+int(board.arr[i][j].xVelocity) >= 0 && j+int(board.arr[i][j].xVelocity) < len(board.arr[i]) && i+int(board.arr[i][j].yVelocity) >= 0 && i+int(board.arr[i][j].yVelocity) < len(board.arr) {
		if !board.arr[i+int(board.arr[i][j].yVelocity)][j+int(board.arr[i][j].xVelocity)].isFull {
			board.arr[i+int(board.arr[i][j].yVelocity)][j+int(board.arr[i][j].xVelocity)] = board.arr[i][j].clone()
			board.arr[i][j] = Particle{}
			return
		}
	}
	board.arr[i][j].xVelocity *= 0.9
	board.arr[i][j].yVelocity *= 0.9
	adjustPostition(board, i, j)
}

func (board *Board) passGravity() {
	for i := range board.arr {
		for j := range board.arr[i] {
			if board.arr[i][j].isFull {
				board.arr[i][j].passGravity()
			}
		}
	}
}

func (particle *Particle) passGravity() {
	particle.yVelocity -= gravityConstant
}

func (particle Particle) clone() Particle {
	return particle
}

func determineSquareSize(totalX, totalY, gridX, gridY int) int {
	gridX += 2 * borderWidth
	gridY += 2 * borderWidth
	gridSize := totalX / gridX
	if totalY/gridY < gridSize {
		gridSize = totalY / gridY
	}
	return gridSize
}

func (board *Board) setupBoard() {
	board.arr = make([][]Particle, boardHeight)
	for i := range board.arr {
		board.arr[i] = make([]Particle, boardWidth)
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
