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

	halftt := 0.5 * t * t
	oneovert := 1 / t

	for _, p := range i.ps.Particles {
		if !p.Fixed {
			a := p.Force.Scale(1.0 / p.Mass)
			p.Position.AddAssign(p.Velocity.Scale(oneovert))
			p.Position.AddAssign(a.Scale(halftt))
			p.Velocity.AddAssign(a.Scale(oneovert))
		}
	}
}
