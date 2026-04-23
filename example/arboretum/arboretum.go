package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/reactivego/gio"

	"github.com/vibrantgio/traer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Arboretum struct {
	*traer.ParticleSystem
}

func NewArboretum() *Arboretum {
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

	p.Position = traer.Vec3{X: q.Position.X + 2.0*rand.Float64() - 1.0, Y: q.Position.Y + 2.0*rand.Float64() - 1.0, Z: 0}
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

func (ps *Arboretum) DrawNetwork(rect image.Rectangle, metric unit.Metric, ops *op.Ops) {
	dp := func(v float32) float32 { return float32(metric.Dp(unit.Dp(v))) }

	DeepPurple500 := color.NRGBAModel.Convert(colornames.DeepPurple500).(color.NRGBA)
	DeepPurple800 := color.NRGBAModel.Convert(colornames.DeepPurple800).(color.NRGBA)
	DeepOrange500 := color.NRGBAModel.Convert(colornames.DeepOrange500).(color.NRGBA)

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
	contentCentroid := f32.Point{X: float32(minX + maxX), Y: float32(minY + maxY)}.Mul(0.5)

	rect = rect.Inset(metric.Dp(20))
	centre := rect.Min.Add(rect.Size().Div(2))
	screenCentre := f32.Point{X: float32(centre.X), Y: float32(centre.Y)}

	scale := float32(math.Min(float64(rect.Dx()), float64(rect.Dy())) / math.Max(maxX-minX, maxY-minY))
	var pen f32.Point
	to := func(p f32.Point) f32.Point {
		absolutepoint := p.Sub(contentCentroid).Mul(scale).Add(screenCentre)
		relativePoint := absolutepoint.Sub(pen)
		pen = absolutepoint
		return relativePoint
	}

	// render edges
	path := &clip.Path{}
	path.Begin(ops)
	for _, spring := range ps.Springs {
		_ = spring
		a := f32.Point{X: float32(spring.A.Position.X), Y: float32(spring.A.Position.Y)}
		b := f32.Point{X: float32(spring.B.Position.X), Y: float32(spring.B.Position.Y)}
		d := b.Sub(a)
		d = d.Mul(float32(1.0 / math.Hypot(float64(d.X), float64(d.Y))))
		nccw := f32.Point{X: -dp(d.Y), Y: dp(d.X)}
		ncw := f32.Point{X: dp(d.Y), Y: -dp(d.X)}
		path.Move(to(a.Add(nccw)))
		path.Line(to(b.Add(nccw)))
		path.Line(to(b.Add(ncw)))
		path.Line(to(a.Add(ncw)))
		path.Line(to(a.Add(nccw)))
		path.Close()
	}
	cstack := clip.Outline{Path: path.End()}.Op().Push(ops)
	paint.ColorOp{Color: DeepPurple500}.Add(ops)
	paint.PaintOp{}.Add(ops)
	cstack.Pop()

	// render nodes
	pen = f32.Point{X: 0, Y: 0}
	path.Begin(ops)
	for _, particle := range ps.Particles[1:] {
		p := f32.Point{X: float32(particle.Position.X), Y: float32(particle.Position.Y)}
		var nodesize = dp(5)
		path.Move(to(p.Add(f32.Point{X: -nodesize, Y: -nodesize})))
		path.Line(to(p.Add(f32.Point{X: nodesize, Y: -nodesize})))
		path.Line(to(p.Add(f32.Point{X: nodesize, Y: nodesize})))
		path.Line(to(p.Add(f32.Point{X: -nodesize, Y: nodesize})))
		path.Line(to(p.Add(f32.Point{X: -nodesize, Y: -nodesize})))
		path.Close()
	}
	cstack = clip.Outline{Path: path.End()}.Op().Push(ops)
	paint.ColorOp{Color: DeepOrange500}.Add(ops)
	paint.PaintOp{}.Add(ops)
	cstack.Pop()

	// render root node
	pen = f32.Point{X: 0, Y: 0}
	particle := ps.Particles[0]
	p := f32.Point{X: float32(particle.Position.X), Y: float32(particle.Position.Y)}
	var nodesize = dp(5)
	cstack = clip.Outline{Path: gio.CirclePath(ops, to(p), 3*nodesize*scale)}.Op().Push(ops)
	paint.ColorOp{Color: DeepPurple800}.Add(ops)
	paint.PaintOp{}.Add(ops)
	cstack.Pop()
}
