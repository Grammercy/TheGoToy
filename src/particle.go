package main

type Particle interface {
  passGravity()
}

type Velocity struct {
  Xspeed float64
  Yspeed float64
}

type Powder struct {
  Name string
  X float64
  Y float64
  Velocity
}

func (p *Powder) passGravity () {
  p.Yspeed -= 0.1
}
