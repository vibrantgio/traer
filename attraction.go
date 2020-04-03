package traer

import (
	"math"
)

type Attraction struct {
	A                  *Particle
	B                  *Particle
	K                  float64
	On                 bool
	DistanceMin        float64
	DistanceMinSquared float64
}

func MakeAttraction(a, b *Particle, k, distanceMin float64) *Attraction {
	return &Attraction{a, b, k, true, distanceMin, distanceMin * distanceMin}
}

func (a *Attraction) TurnOn() {
	a.On = true
}

func (a *Attraction) TurnOff() {
	a.On = false
}

func (a *Attraction) IsOn() bool {
	return a.On
}

func (a *Attraction) IsOff() bool {
	return !a.On
}

func (a *Attraction) Apply() {
	if !a.On || (a.A.Fixed && a.B.Fixed) {
		return
	}
	a2b := a.A.Position.Subtract(a.B.Position)
	a2bDistanceSquared := math.Max(a2b.LengthSquared(), a.DistanceMinSquared)
	force := a.K * a.A.Mass * a.B.Mass / a2bDistanceSquared
	a2b.ScaleAssign(force / math.Sqrt(a2bDistanceSquared))
	if !a.A.Fixed {
		a.A.Force.SubtractAssign(a2b)
	}
	if !a.B.Fixed {
		a.B.Force.AddAssign(a2b)
	}
}
