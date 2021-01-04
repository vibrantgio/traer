package main

import (
	"image/color"
	"math"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"github.com/reactivego/traer"
)

var (
	DeepPurple500 = color.NRGBAModel.Convert(colornames.DeepPurple500).(color.NRGBA)
	DeepPurple800 = color.NRGBAModel.Convert(colornames.DeepPurple800).(color.NRGBA)
	DeepOrange500 = color.NRGBAModel.Convert(colornames.DeepOrange500).(color.NRGBA)
)

type Arboretum struct {
	*traer.ParticleSystem
}

func MakeArboretum() *Arboretum {
	ps := traer.MakeParticleSystem(0.0, 0.3)
	ps.Clear()
	ps.MakeDefaultParticle().Fixed = true
	return &Arboretum{ps}
}

func (ps *Arboretum) AddNode() {
	p := ps.MakeDefaultParticle()
	maxParticle := len(ps.Particles) - 1
	if maxParticle == 0 {
		return
	}

	q := ps.Particles[rand.Intn(maxParticle)]
	for p == q {
		q = ps.Particles[rand.Intn(maxParticle)]
	}

	for _, r := range ps.Particles {
		if p != r {
			ps.MakeAttraction(p, r, -1000, 20)
		}
	}

	ps.MakeSpring(p, q, 0.2, 0.2, 20)

	p.Position = traer.Vec3{q.Position.X + 2.0*rand.Float64() - 1.0, q.Position.Y + 2.0*rand.Float64() - 1.0, 0}
}

func (ps *Arboretum) aabb() (float64, float64, float64, float64) {
	maxX := -math.MaxFloat64
	minX := math.MaxFloat64
	maxY := -math.MaxFloat64
	minY := math.MaxFloat64
	for _, p := range ps.Particles {
		maxX = math.Max(maxX, p.Position.X)
		minX = math.Min(minX, p.Position.X)
		maxY = math.Max(maxY, p.Position.Y)
		minY = math.Min(minY, p.Position.Y)
	}
	return minX, minY, maxX, maxY
}

func (ps *Arboretum) DrawNetwork(rect f32.Rectangle) op.CallOp {
	minX, minY, maxX, maxY := ps.aabb()
	if MinWidthDp > (maxX - minX) {
		outsetX := (MinWidthDp - maxX + minX) / 2
		minX -= outsetX
		maxX += outsetX
	}
	if MinHeightDp > (maxY - minY) {
		outsetY := (MinHeightDp - maxY + minY) / 2
		minY -= outsetY
		maxY += outsetY
	}
	contentCentroid := f32.Point{float32(minX + maxX), float32(minY + maxY)}.Mul(0.5)
	inset := f32.Point{20, 20}
	insets := f32.Rectangle{rect.Min.Add(inset), rect.Max.Sub(inset)}

	screenCentre := insets.Min.Add(insets.Size().Mul(0.5))
	scale := float32(math.Min(float64(insets.Dx()), float64(insets.Dy())) / math.Max(maxX-minX, maxY-minY))
	var pen f32.Point
	to := func(p f32.Point) f32.Point {
		absolutepoint := p.Sub(contentCentroid).Mul(scale).Add(screenCentre)
		relativePoint := absolutepoint.Sub(pen)
		pen = absolutepoint
		return relativePoint
	}

	ops := &op.Ops{}
	macro := op.Record(ops)
	// render nodes
	stack := op.Push(ops)
	path := &clip.Path{}
	path.Begin(ops)
	for _, spring := range ps.Springs {
		_ = spring
		a := f32.Point{float32(spring.A.Position.X), float32(spring.A.Position.Y)}
		b := f32.Point{float32(spring.B.Position.X), float32(spring.B.Position.Y)}
		d := b.Sub(a)
		d = d.Mul(float32(1.0 / math.Hypot(float64(d.X), float64(d.Y))))
		nccw := f32.Point{-d.Y, d.X}
		ncw := f32.Point{d.Y, -d.X}
		path.Move(to(a.Add(nccw)))
		path.Line(to(b.Add(nccw)))
		path.Line(to(b.Add(ncw)))
		path.Line(to(a.Add(ncw)))
		path.Line(to(a.Add(nccw)))
	}
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: DeepPurple500}.Add(ops)
	paint.PaintOp{}.Add(ops)
	stack.Pop()

	// render edges
	stack = op.Push(ops)
	path.Begin(ops)
	for _, particle := range ps.Particles[1:] {
		p := f32.Point{float32(particle.Position.X), float32(particle.Position.Y)}
		const nodesize = 5
		path.Move(to(p.Add(f32.Point{-nodesize, -nodesize})))
		path.Line(to(p.Add(f32.Point{nodesize, -nodesize})))
		path.Line(to(p.Add(f32.Point{nodesize, nodesize})))
		path.Line(to(p.Add(f32.Point{-nodesize, nodesize})))
		path.Line(to(p.Add(f32.Point{-nodesize, -nodesize})))
	}
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: DeepOrange500}.Add(ops)
	paint.PaintOp{}.Add(ops)
	stack.Pop()

	// render root node
	stack = op.Push(ops)
	pen = f32.Point{0, 0}
	particle := ps.Particles[0]
	p := f32.Point{float32(particle.Position.X), float32(particle.Position.Y)}
	const nodesize = 5
	clip.Outline{Path: Circle(to(p), 3*nodesize*scale, ops)}.Op().Add(ops)
	paint.ColorOp{Color: DeepPurple800}.Add(ops)
	paint.PaintOp{}.Add(ops)
	stack.Pop()

	return macro.Stop()
}
