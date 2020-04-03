package traer

const DEFAULT_MASS = 1
const DEFAULT_GRAVITY = 0
const DEFAULT_DRAG = 0.001

type ParticleSystem struct {
	Particles   []*Particle
	Springs     []*Spring
	Attractions []*Attraction

	// Gravity contains the strength of gravity, down (in the positive y
	// direction) or in whatever 3D direction you feel like. You probably
	// want the magnitude of this to be in the range of 0-5.
	Gravity     Vec3
	Drag        float64
	Integrator  Integrator
}

func MakeParticleSystem(g, drag float64) *ParticleSystem {
	ps := &ParticleSystem{Gravity: Vec3{0, g, 0}, Drag: drag}
	ps.Integrator = MakeVelocityVerletIntegrator(ps)
	return ps
}

func MakeDefaultParticleSystem() *ParticleSystem {
	return MakeParticleSystem(DEFAULT_GRAVITY, DEFAULT_DRAG)
}

func (ps *ParticleSystem) Clear() {
	ps.Particles = nil
	ps.Attractions = nil
	ps.Springs = nil
}

func (ps *ParticleSystem) MakeParticle(mass, x, y, z float64) *Particle {
	particle := &Particle{Mass: mass, Position: Vec3{x, y, z}}
	ps.Particles = append(ps.Particles, particle)
	return particle
}

func (ps *ParticleSystem) MakeDefaultParticle() *Particle {
	return ps.MakeParticle(DEFAULT_MASS, 0, 0, 0)
}

func (ps *ParticleSystem) MakeAttraction(a, b *Particle, k, distanceMin float64) *Attraction {
	attraction := MakeAttraction(a, b, k, distanceMin)
	ps.Attractions = append(ps.Attractions, attraction)
	return attraction
}

func (ps *ParticleSystem) MakeSpring(pa, pb *Particle, ks, d, r float64) *Spring {
	spring := MakeSpring(pa, pb, ks, d, r)
	ps.Springs = append(ps.Springs, spring)
	return spring
}

// Tick advances the simulation by a 1/t seconds (t is the argument to Tick).
// By default use a t of 1.0 indicating a simulation duration of a second for
// that Tick call. Increase t to a higher value in order to make the simulation
// run SLOWER, as a higher t will lead to a lower 1/t value forcing the
// simulation to run smaller time increments for every call to Tick.
func (ps *ParticleSystem) Tick(t float64) {
	ps.Integrator.Step(t)
}

func (ps *ParticleSystem) ClearForces() {
	for _, p := range ps.Particles {
		p.Force = Vec3{}
	}
}

func (ps *ParticleSystem) ApplyForces() {
	if ps.Gravity.Length() > 0.0 {
		for _, particle := range ps.Particles {
			particle.Force.AddAssign(ps.Gravity.Scale(particle.Mass))
		}
	}

	for _, particle := range ps.Particles {
		particle.Force.SubtractAssign(particle.Velocity.Scale(ps.Drag))
	}

	for _, spring := range ps.Springs {
		spring.Apply()
	}

	for _, attraction := range ps.Attractions {
		attraction.Apply()
	}
}
