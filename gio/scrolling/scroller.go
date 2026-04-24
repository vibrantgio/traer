package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"github.com/vibrantgio/circle"

	"github.com/vibrantgio/traer"
)

type KineticScroller struct {
	View    image.Rectangle
	Content image.Rectangle

	Physics *traer.ParticleSystem

	// Horizontal

	FixedParticleH   *traer.Particle
	SpringH1         *traer.Spring
	ContentParticleH *traer.Particle
	SpringH2         *traer.Spring
	PointerParticleH *traer.Particle

	// Vertical

	FixedParticleV   *traer.Particle
	SpringV1         *traer.Spring
	ContentParticleV *traer.Particle
	SpringV2         *traer.Spring
	PointerParticleV *traer.Particle

	// Mouse/Touch drag tracking

	IsManipulating       bool
	PointerStartLocation traer.Vec3
	PointerLocation      traer.Vec3

	// Scroll tracking

	Scroll traer.Vec3
}

func NewScroller() *KineticScroller {
	s := &KineticScroller{}
	s.Physics = traer.NewParticleSystem(SystemAcceleration, SystemDrag)

	// Horizontal
	s.FixedParticleH = s.Physics.NewDefaultParticle()
	s.FixedParticleH.Fixed = true
	s.ContentParticleH = s.Physics.NewDefaultParticle()
	s.SpringH1 = s.Physics.NewSpring(s.FixedParticleH, s.ContentParticleH, FixedSpringConstant, SpringDamping, 0)
	s.PointerParticleH = s.Physics.NewDefaultParticle()
	s.PointerParticleH.Fixed = true
	s.SpringH2 = s.Physics.NewSpring(s.ContentParticleH, s.PointerParticleH, PointerSpringConstant, SpringDamping, 0)

	// Vertical
	s.FixedParticleV = s.Physics.NewDefaultParticle()
	s.FixedParticleV.Fixed = true
	s.ContentParticleV = s.Physics.NewDefaultParticle()
	s.SpringV1 = s.Physics.NewSpring(s.FixedParticleV, s.ContentParticleV, FixedSpringConstant, SpringDamping, 0)
	s.PointerParticleV = s.Physics.NewDefaultParticle()
	s.PointerParticleV.Fixed = true
	s.SpringV2 = s.Physics.NewSpring(s.ContentParticleV, s.PointerParticleV, PointerSpringConstant, SpringDamping, 0)

	return s
}

func (s *KineticScroller) Pointer(event pointer.Event) {
	switch event.Type {
	case pointer.Press:
		s.IsManipulating = true
		px, py := float64(event.Position.X), float64(event.Position.Y)
		s.PointerStartLocation = traer.Vec3{X: px, Y: py}
		s.PointerLocation = traer.Vec3{X: px, Y: py}

	case pointer.Drag:
		px, py := float64(event.Position.X), float64(event.Position.Y)
		s.PointerLocation = traer.Vec3{X: px, Y: py}

	case pointer.Release:
		s.IsManipulating = false

	case pointer.Scroll:
		s.Scroll = traer.Vec3{X: float64(event.Scroll.X), Y: float64(event.Scroll.Y)}

	default:
		fmt.Println(event)
	}
}

func (s *KineticScroller) Tick(t float64) float64 {
	s.SpringH1.TurnOff()
	s.SpringH2.TurnOff()
	s.SpringV1.TurnOff()
	s.SpringV2.TurnOff()

	// Handle scroll events with no bouncing. Scroll events are kinetic in their
	// own right and a fling keeps on generating scroll events after the gesture
	// has stopped.
	if s.Scroll.Length() > MinScrollLength {
		contentMinX := s.ContentParticleH.Position.X - s.Scroll.X
		contentMinY := s.ContentParticleV.Position.Y - s.Scroll.Y

		ddx := s.View.Dx() - s.Content.Dx()
		if ddx == 0 || contentMinX/float64(ddx) <= 0 {
			contentMinX = 0
		} else if contentMinX/float64(ddx) >= 1 {
			contentMinX = float64(ddx)
		}
		ddy := s.View.Dy() - s.Content.Dy()
		if ddy == 0 || contentMinY/float64(ddy) <= 0 {
			contentMinY = 0
		} else if contentMinY/float64(ddy) >= 1 {
			contentMinY = float64(ddy)
		}
		mx := int(math.Round(contentMinX))
		my := int(math.Round(contentMinY))
		s.Content.Min.X, s.Content.Max.X = mx, mx+s.Content.Dx()
		s.Content.Min.Y, s.Content.Max.Y = my, my+s.Content.Dy()

		s.ContentParticleH.Position.X = contentMinX
		s.ContentParticleV.Position.Y = contentMinY
	}

	contentMinX := s.ContentParticleH.Position.X
	contentMinY := s.ContentParticleV.Position.Y

	sx := 0.0
	ddx := s.View.Dx() - s.Content.Dx()
	if ddx != 0 {
		sx = contentMinX / float64(ddx)
	}
	if sx <= 0 || sx >= 1 {
		if sx <= 0 {
			contentMinX = 0 // scrolled too far to the left, so close gap
		} else {
			contentMinX = float64(ddx) // scrolled too far to the right, so close gap
		}
		s.SpringH1.TurnOn()
	}

	sy := 0.0
	ddy := s.View.Dy() - s.Content.Dy()
	if ddy != 0 {
		sy = contentMinY / float64(ddy)
	}
	if sy <= 0 || sy >= 1 {
		if sy <= 0 {
			contentMinY = 0 // scrolled too far up, close gap
		} else {
			contentMinY = float64(ddy) // scrolled too far down, close gap
		}
		s.SpringV1.TurnOn()
	}

	s.FixedParticleH.Position.X = contentMinX
	s.FixedParticleV.Position.Y = contentMinY

	if s.IsManipulating {
		s.SpringH2.TurnOn()
		s.SpringV2.TurnOn()
		d := s.PointerLocation.Subtract(s.PointerStartLocation)
		s.PointerStartLocation = s.PointerLocation
		s.PointerParticleH.Position.X += d.X
		s.PointerParticleV.Position.Y += d.Y
	} else {
		s.PointerParticleH.Position = s.ContentParticleH.Position
		s.PointerParticleV.Position = s.ContentParticleV.Position
	}

	activity := s.Physics.Tick(t)

	mx := int(math.Round(s.ContentParticleH.Position.X))
	my := int(math.Round(s.ContentParticleV.Position.Y))
	s.Content.Min.X, s.Content.Max.X = mx, mx+s.Content.Dx()
	s.Content.Min.Y, s.Content.Max.Y = my, my+s.Content.Dy()

	return activity
}

func (s *KineticScroller) Draw(rect image.Rectangle, ops *op.Ops) {
	draw := func(particle *traer.Particle, color color.Color, ax, ay float32, ops *op.Ops) {
		pt := f32.Pt(float32(particle.Position.X), float32(particle.Position.Y))
		pt.X += ax * float32(rect.Min.X+rect.Max.X)
		pt.Y += ay * float32(rect.Min.Y+rect.Max.Y)
		circle.FillCircle(ops, pt, 10, color)
	}
	draw(s.PointerParticleH, colornames.Red500, 0.0, 0.5, ops)
	draw(s.PointerParticleV, colornames.Red500, 0.5, 0.0, ops)
	draw(s.FixedParticleH, colornames.Green500, 0.0, 0.5, ops)
	draw(s.FixedParticleV, colornames.Green500, 0.5, 0.0, ops)
	draw(s.ContentParticleH, colornames.Blue500, 0.0, 0.5, ops)
	draw(s.ContentParticleV, colornames.Blue500, 0.5, 0.0, ops)
}
