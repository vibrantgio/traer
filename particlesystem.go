// SPDX-License-Identifier: Unlicense OR MIT

package traer

import "math"

const (
	DefaultMass = 1

	// DefaultGravity is no gravity.
	DefaultGravity = 0

	DefaultDrag = 0.001
)

// IntegrationStep is the interface common to all integration algorithms.
// The return value for the function is a measure of how much speed is still
// in particles the system.
type IntegrationStep func(t float64) float64

// ParticleSystem is in charge of everything. It makes particles and forces
// for you and you tell it to advance the simulation using Tick().
type ParticleSystem struct {
	Particles   []*Particle
	Springs     []*Spring
	Attractions []*Attraction

	// Gravity contains the strength of gravity, down (in the positive y
	// direction) or in whatever 3D direction you feel like. You probably want
	// the magnitude of this to be in the range of 0-5.
	Gravity Vec3
	// Drag is the drag force that acts on all objects equally, and proportional
	// to velocity.
	Drag float64

	// Step holds a Velocity Verlet integration step function by default. This
	// implementation is practically identical to the Modified Euler integrator
	// from Traer Physics 3.0 but has bug-fixes to the algorithm.
	Step IntegrationStep
}

// NewParticleSystem constructs a new particle system with some downward
// (positive y) or 3D gravity and some drag. You can make as many of these as
// you'd like as long as forces from one system don't refer to particles from
// another. I don't know what would happen if you connected particles from one
// system to another.
func NewParticleSystem(g, drag float64) *ParticleSystem {
	ps := &ParticleSystem{Gravity: Vec3{0, g, 0}, Drag: drag}
	ps.Step = NewVelocityVerletIntegrator(ps)
	return ps
}

// Construct a new particle system with DefaultGravity (0) and
// DefaultDrag (0.001).
func NewDefaultParticleSystem() *ParticleSystem {
	return NewParticleSystem(DefaultGravity, DefaultDrag)
}

// Clear deletes all the particles and all the forces in the system (except
// the omnipresent gravity and drag even if they are 0).
func (ps *ParticleSystem) Clear() {
	ps.Particles = nil
	ps.Attractions = nil
	ps.Springs = nil
}

// NewParticle creates a new particle in the system with some mass and at
// some x, y, z position.
func (ps *ParticleSystem) NewParticle(mass, x, y, z float64) *Particle {
	particle := &Particle{Mass: mass, Position: Vec3{x, y, z}}
	ps.Particles = append(ps.Particles, particle)
	return particle
}

// NewDefaultParticle creates a new particle in the system with mass 1.0 at
// x, y, z position (0, 0, 0).
func (ps *ParticleSystem) NewDefaultParticle() *Particle {
	return ps.NewParticle(DefaultMass, 0, 0, 0)
}

// NewAttraction makes an attraction (or repulsion) force between two
// particles. If the strength is negative they repel each other, if the
// strength is positive they attract. There is also a minimum distance that
// limits how strong this force can get close up.
func (ps *ParticleSystem) NewAttraction(a, b *Particle, strength, minimumDistance float64) *Attraction {
	attraction := &Attraction{
		A:                      a,
		B:                      b,
		Strength:               strength,
		MinimumDistanceSquared: minimumDistance * minimumDistance,
		On:                     true,
	}
	ps.Attractions = append(ps.Attractions, attraction)
	return attraction
}

// NewSpring makes a spring in the system between 2 particles you have
// previously created.
//  strength -  A strong spring acts like a stick. A weak one takes a
//    long time to return to its rest length.
//  damping - A spring with high damping doesn't overshoot and settles
//     down quickly, while a low damping spring oscillates.
//  restLength - A spring wants to be at this length and acts on the
//      particles to push or pull them exactly this far apart at all times.
func (ps *ParticleSystem) NewSpring(a, b *Particle, strength, damping, restLength float64) *Spring {
	spring := &Spring{
		A:          a,
		B:          b,
		Strength:   strength,
		Damping:    damping,
		RestLength: restLength,
		On:         true,
	}
	ps.Springs = append(ps.Springs, spring)
	return spring
}

// Tick advances the simulation by a 1/t seconds (t is the argument to Tick).
// By default use a t of 1.0 indicating a simulation duration of a second for
// that Tick call. Increase t to a higher value in order to make the
// simulation run SLOWER, as a higher t will lead to a lower 1/t value forcing
// the simulation to run smaller time increments for every call to Tick.
//
// Note that target framerate in TraerAS3 was 31fps and it used Tick(1). We
// usually get 60fps so we can double the step size and by doing so splitting
// the step time in half.
func (ps *ParticleSystem) Tick(t float64) float64 {
	return ps.Step(t)
}

// ApplyForces wil apply gravity, drag, spring and attraction forces to
// particles
func (ps *ParticleSystem) ApplyForces() {
	for _, particle := range ps.Particles {
		if ps.Gravity.Length() > 0.0 {
			// Original Traer Physics version 3.0 does not take particle mass
			// into account for gravity. Only matters for particles with mass
			// different from the default 1.0 value though.
			particle.Force = ps.Gravity.Scale(particle.Mass)
		} else {
			particle.Force = Vec3{}
		}
	}

	for _, particle := range ps.Particles {
		particle.Force.SubtractAssign(particle.Velocity.Scale(ps.Drag))
	}

	for _, s := range ps.Springs {
		if !s.On || (s.A.Fixed && s.B.Fixed) {
			continue
		}

		a2b := s.A.Position.Subtract(s.B.Position)
		a2bDistance := a2b.Length()
		if a2bDistance > 0.0 {
			a2b.ScaleAssign(1 / a2bDistance) // normalize a2b
		}

		// spring force is proportional to how much it stretched
		springForce := -(a2bDistance - s.RestLength) * s.Strength

		// want velocity along line b/w a & b, damping force is proportional to this
		va2b := s.A.Velocity.Subtract(s.B.Velocity)
		dampingForce := -s.Damping * a2b.Dot(va2b)

		// forceB is same as forceA in opposite direction
		a2b.ScaleAssign(springForce + dampingForce)

		if !s.A.Fixed {
			s.A.Force.AddAssign(a2b)
		}
		if !s.B.Fixed {
			s.B.Force.SubtractAssign(a2b)
		}
	}

	for _, a := range ps.Attractions {
		if !a.On || (a.A.Fixed && a.B.Fixed) {
			continue
		}
		a2b := a.A.Position.Subtract(a.B.Position)
		a2bDistanceSquared := math.Max(a2b.LengthSquared(), a.MinimumDistanceSquared)
		force := a.Strength * a.A.Mass * a.B.Mass / a2bDistanceSquared
		length := math.Sqrt(a2bDistanceSquared)

		a2b.ScaleAssign(force / length)
		if !a.A.Fixed {
			a.A.Force.SubtractAssign(a2b)
		}
		if !a.B.Fixed {
			a.B.Force.AddAssign(a2b)
		}
	}
}
