package main

import (
  "fmt"
)

type screen struct {
}

func (screen *screen) passGravity

func main() {
  powder := Powder{Particle{"Sand", 0.0, 0.0}, Velocity{0.0, 0.0}}
  for i := 0; i < 10; i++ {
    powder.passGravity()
    fmt.Println(powder)
  }
}
