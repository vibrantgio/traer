package traer

func MakeVelocityVerletIntegrator(ps *ParticleSystem) Integrator {
	return &velocityVerletIntegrator{ps: ps}
}

type velocityVerletIntegrator struct {
	ps *ParticleSystem
}

func (i *velocityVerletIntegrator) Step(t float64) {
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
