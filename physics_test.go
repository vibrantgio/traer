// SPDX-License-Identifier: Unlicense OR MIT

package traer

import (
	"math"
	"testing"
)

// TestFreeFall drops a unit-mass particle under gravity=10 (no drag) and
// checks the velocity-verlet integrator reproduces the closed-form result
// it should for constant acceleration: after N steps of dt the particle
// reaches position 0.5*a*(N*dt)^2 + a*dt*dt/2 (the +a*dt*dt/2 being the
// well-known half-step offset from velocity verlet with v0=0).
func TestFreeFall(t *testing.T) {
	ps := NewParticleSystem(10, 0)
	p := ps.NewParticle(1, 0, 0, 0)

	const steps = 10
	const tInv = 10.0 // dt = 1/tInv = 0.1
	for range steps {
		ps.Tick(tInv)
	}

	dt := 1.0 / tInv
	want := 0.5 * 10 * math.Pow(float64(steps)*dt, 2)
	if math.Abs(p.Position.Y-want) > 1e-9 {
		t.Errorf("y = %v after %d steps, want %v", p.Position.Y, steps, want)
	}
	if p.Position.X != 0 || p.Position.Z != 0 {
		t.Errorf("particle drifted off-axis: %v", p.Position)
	}
}

func TestFixedParticleDoesNotMove(t *testing.T) {
	ps := NewParticleSystem(10, 0.5)
	p := ps.NewParticle(1, 3, 4, 5)
	p.Fixed = true

	for range 100 {
		ps.Tick(10)
	}
	if p.Position != (Vec3{3, 4, 5}) {
		t.Errorf("fixed particle moved to %v", p.Position)
	}
}

// TestSpringSettles connects a free particle to a fixed one with a strong,
// well-damped spring; the free particle should converge toward the
// rest-length position.
func TestSpringSettles(t *testing.T) {
	ps := NewParticleSystem(0, 0.1)
	anchor := ps.NewParticle(1, 0, 0, 0)
	anchor.Fixed = true
	free := ps.NewParticle(1, 10, 0, 0)
	ps.NewSpring(anchor, free, 0.5, 0.9, 3)

	for range 2000 {
		ps.Tick(5)
	}
	d := free.Position.Subtract(anchor.Position).Length()
	if math.Abs(d-3) > 0.05 {
		t.Errorf("spring settled at distance %v, want ~3", d)
	}
	if free.Velocity.Length() > 0.01 {
		t.Errorf("free particle not at rest, |v| = %v", free.Velocity.Length())
	}
}

// TestAttractionPullsTogether: two free particles, mutually attracting,
// should end up closer than they started. Uses a short horizon and a
// minimum-distance clamp well above the converged radius so the force
// stays bounded and the pair can't slingshot past each other.
func TestAttractionPullsTogether(t *testing.T) {
	ps := NewParticleSystem(0, 0)
	a := ps.NewParticle(1, 0, 0, 0)
	b := ps.NewParticle(1, 10, 0, 0)
	ps.NewAttraction(a, b, 10, 5)

	start := b.Position.Subtract(a.Position).Length()
	for range 10 {
		ps.Tick(10)
	}
	end := b.Position.Subtract(a.Position).Length()
	if end >= start {
		t.Errorf("attraction did not pull particles together: start=%v end=%v", start, end)
	}
}

func TestClearRemovesEverything(t *testing.T) {
	ps := NewDefaultParticleSystem()
	a := ps.NewDefaultParticle()
	b := ps.NewDefaultParticle()
	ps.NewSpring(a, b, 1, 0.1, 1)
	ps.NewAttraction(a, b, 1, 1)

	ps.Clear()
	if len(ps.Particles) != 0 || len(ps.Springs) != 0 || len(ps.Attractions) != 0 {
		t.Errorf("Clear left state: particles=%d springs=%d attractions=%d",
			len(ps.Particles), len(ps.Springs), len(ps.Attractions))
	}
}

func BenchmarkTick(b *testing.B) {
	ps := NewParticleSystem(1, 0.01)
	const n = 50
	particles := make([]*Particle, n)
	for i := range n {
		particles[i] = ps.NewParticle(1, float64(i), 0, 0)
	}
	for i := range n - 1 {
		ps.NewSpring(particles[i], particles[i+1], 0.1, 0.1, 1)
	}
	b.ResetTimer()
	for range b.N {
		ps.Tick(10)
	}
}
