package main

import (
	"math"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/fogleman/contourmap"
	"github.com/reactivego/traer"
)

// AS3 defaults
const AttractorInitialStrength = 10000.0
const AttractorInitialMinDistance = 30.0
const AttractorStrengthFactor = 60.0

// iPhone
// const AttractorInitialStrength = 50.0
// const AttractorInitialMinDistance = 15.0
// const AttractorStrengthFactor = 30.0

const ContouringScale = 1.0 / 7.0

type Ball struct {
	P *traer.Particle
	A *traer.Attraction
	R float64
}

type Floaters struct {
	*traer.ParticleSystem
	Attractor            *traer.Particle
	AttractorMinDistance float64
	AttractorStrength    float64

	Balls    []Ball
	Grid     []float64
	Contours []contourmap.Contour
}

func MakeFloaters(InitialWidth, InitialHeight, InitialVelocity, Radius float64, Count int) *Floaters {
	s := &Floaters{}
	s.ParticleSystem = traer.MakeParticleSystem(traer.DEFAULT_GRAVITY, 0.02)
	s.Attractor = s.MakeParticle(0.8, InitialWidth/2, InitialHeight/2, 0)
	s.Attractor.Fixed = true
	s.AttractorMinDistance = AttractorInitialMinDistance
	s.AttractorStrength = AttractorInitialStrength
	s.Balls = make([]Ball, Count)
	for i := 0; i < len(s.Balls); i++ {
		randX := rand.Float64() * InitialWidth
		randY := rand.Float64() * InitialHeight
		randVX := (rand.Float64() - 0.5) * InitialVelocity
		randVY := (rand.Float64() - 0.5) * InitialVelocity
		randR := Radius
		p := s.MakeParticle(0.8, randX, randY, 0)
		p.Velocity = traer.Vec3{X: randVX, Y: randVY}
		a := s.MakeAttraction(p, s.Attractor, s.AttractorStrength, s.AttractorMinDistance)
		s.Balls[i] = Ball{P: p, A: a, R: randR}
	}
	return s
}

func (s *Floaters) Pointer(event pointer.Event) {
	switch event.Type {
	case pointer.Press, pointer.Drag:
		px, py := float64(event.Position.X), float64(event.Position.Y)
		d := math.Hypot(px-s.Attractor.Position.X, py-s.Attractor.Position.Y)
		s.Attractor.Position = traer.Vec3{X: px, Y: py}
		if d > 0.001 {
			s.AttractorMinDistance = d * 1.2
			s.AttractorStrength = -(10 + (AttractorStrengthFactor * d * d))
		}
	case pointer.Release:
		s.AttractorMinDistance = AttractorInitialMinDistance
		s.AttractorStrength = AttractorInitialStrength
	}
}

func (s *Floaters) Position(Width, Height float64) {
	for _, ball := range s.Balls {
		p := ball.P.Position
		wrapX := math.Mod(Width+p.X, Width)
		wrapY := math.Mod(Height+p.Y, Height)
		ball.P.Position = traer.Vec3{X: wrapX, Y: wrapY}
		ball.A.SetMinimumDistance(s.AttractorMinDistance)
		ball.A.Strength = s.AttractorStrength
	}
}

func (s *Floaters) Contour(Width, Height float64) {
	w := int(Width * ContouringScale)
	h := int(Height * ContouringScale)
	if cap(s.Grid) < w*h {
		s.Grid = make([]float64, w*h)
	}
	if len(s.Grid) < w*h {
		s.Grid = s.Grid[:w*h]
	}
	for y := 0; y < h; y++ {
		rowstart := y * w
		for x := 0; x < w; x++ {
			sum := 0.0
			for _, b := range s.Balls {
				dx := float64(x) - b.P.Position.X*ContouringScale
				dy := float64(y) - b.P.Position.Y*ContouringScale
				sum += (b.R * b.R * ContouringScale) / (dx*dx + dy*dy)
			}
			s.Grid[rowstart+x] = sum
		}
	}
	pointCount := 0
	m := contourmap.FromFloat64s(w, h, s.Grid).Closed()
	s.Contours = m.Contours(1)
	for _, contour := range s.Contours {
		for i, p := range contour {
			contour[i] = contourmap.Point{p.X / ContouringScale, p.Y / ContouringScale}
			pointCount++
		}
	}
	// log.Printf("Contour resolution:w%dxh%d, points:%d\n", w, h, pointCount)
}

func (s *Floaters) Render() op.CallOp {
	curveCount := 0
	var pen contourmap.Point
	move := func(p contourmap.Point) f32.Point {
		px, py := float32(p.X-pen.X), float32(p.Y-pen.Y)
		pen = p
		return f32.Point{px, py}
	}
	quad := func(c, p contourmap.Point) (f32.Point, f32.Point) {
		cx, cy := float32(c.X-pen.X), float32(c.Y-pen.Y)
		px, py := float32(p.X-pen.X), float32(p.Y-pen.Y)
		pen = p
		curveCount++
		return f32.Point{cx, cy}, f32.Point{px, py}
	}
	ops := &op.Ops{}
	macro := op.Record(ops)
	path := &clip.Path{}
	path.Begin(ops)
	for _, c := range s.Contours {
		lenc := (len(c) >> 1) << 1
		if lenc < 4 {
			continue
		}
		path.Move(move(c[0]))
		for i := 1; i < lenc-1; i += 2 {
			path.Quad(quad(c[i], c[i+1]))
		}
		path.Quad(quad(c[lenc-1], c[0]))
	}
	path.End().Add(ops)

	return macro.Stop()
}
