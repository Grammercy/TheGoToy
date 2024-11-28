package main

import (
  "fmt"
)

type Board struct {
  arr [][]Particle
}

type Particle struct {
  isEmpty bool
  name string
  xVelocity float64
  yVelocity float64
}

func main() {
  var board Board
  board.arr = make([][]Particle, 20)
  for i := range board.arr {
    board.arr[i] = make([]Particle, 20)
  }
  fmt.Println(board)
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
