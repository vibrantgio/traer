// SPDX-License-Identifier: Unlicense OR MIT

package traer

// Attraction or repulsion (negative attraction) acts on 2 particles and either
// constantly pulls them together or constantly pushes them apart by applying a
// force to each particle:
//
//	G*m1*m2/d^2
//
// Because in the formula the d (distance between A,B) is squared, the force is
// much stronger close up than far away.
type Attraction struct {
	A *Particle
	B *Particle

	// Strength, the G in the formula G*m1*m2/d^2
	Strength float64
	// Minimum Distance, the force does not get stronger closer than this.
	// The value is stored squared as that is needed for applying attractions.
	MinimumDistanceSquared float64

	On bool
}

// SetMinimumDistance will square the minimumDistance argument and set the
// MinimumDistanceSquared field of the attraction.
func (a *Attraction) SetMinimumDistance(minimumDistance float64) {
	a.MinimumDistanceSquared = minimumDistance * minimumDistance
}
