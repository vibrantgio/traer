package traer

func MakeModifiedEulerIntegrator(ps *ParticleSystem) Integrator {
	return &modifiedEulerIntegrator{ps: ps}
}

type modifiedEulerIntegrator struct {
	ps *ParticleSystem
}

func (i *modifiedEulerIntegrator) Step(t float64) {
	i.ps.ClearForces()
	i.ps.ApplyForces()

	dt := 1 / t
	halfdtdt := 0.5 * dt * dt

	for _, p := range i.ps.Particles {
		if !p.Fixed {
			a := p.Force.Scale(1.0 / p.Mass)
			p.Position.AddAssign(p.Velocity.Scale(dt))
			p.Position.AddAssign(a.Scale(halfdtdt))
			p.Velocity.AddAssign(a.Scale(dt))
		}
	}
}
