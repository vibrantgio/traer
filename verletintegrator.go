// SPDX-License-Identifier: Unlicense OR MIT

package traer

import "math"

// NewDefaultVerletIntegrator creates an integrator that performs the
// following calculation for every particle p that is not fixed.
//
//	a := p.Force.Scale(1.0 / p.Mass)
//	position := p.Position.Add(p.Velocity.Scale(1.0 / invDt)).Add(a.Scale(1.0 / (invDt * invDt)))
//	p.Velocity = position.Subtract(p.Position).Scale(invDt)
//	p.Position = position
func NewDefaultVerletIntegrator(ps *ParticleSystem) IntegrationStep {
	step := func(invDt float64) float64 {
		ps.ApplyForces()

		dt := 1.0 / invDt
		dtdt := dt * dt

		activity := 0.0
		for _, p := range ps.Particles {
			if !p.Fixed {
				a := p.Force.Scale(1.0 / p.Mass)
				position := p.Position.Add(p.Velocity.Scale(dt)).Add(a.Scale(dtdt))
				p.Velocity = position.Subtract(p.Position).Scale(invDt)
				p.Position = position
				activity += p.Velocity.LengthSquared()
			}
		}
		return math.Sqrt(activity)
	}
	return step
}

// NewVelocityVerletIntegrator creates an integrator that performs the
// following calculation for every particle p that is not fixed.
//
//	a := p.Force.Scale(1.0 / p.Mass)
//	p.Position.AddAssign(p.Velocity.Scale(1.0 / invDt))
//	p.Position.AddAssign(a.Scale(1.0 / (2.0 * invDt * invDt)))
//	p.Velocity.AddAssign(a.Scale(1.0 / invDt))
func NewVelocityVerletIntegrator(ps *ParticleSystem) IntegrationStep {
	step := func(invDt float64) float64 {
		ps.ApplyForces()

		dt := 1.0 / invDt
		halfdtdt := 0.5 * dt * dt

		activity := 0.0
		for _, p := range ps.Particles {
			if !p.Fixed {
				a := p.Force.Scale(1.0 / p.Mass)
				p.Position.AddAssign(p.Velocity.Scale(dt))
				p.Position.AddAssign(a.Scale(halfdtdt))
				p.Velocity.AddAssign(a.Scale(dt))
				activity += p.Velocity.LengthSquared()
			}
		}
		return math.Sqrt(activity)
	}
	return step
}
