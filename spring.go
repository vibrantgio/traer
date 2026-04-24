// SPDX-License-Identifier: Unlicense OR MIT

package traer

// Spring connects 2 particles A,B and tries to keep them a certain distance apart.
//
// A spring has 3 properties Strength, Damping and RestLength that determine its
// behavior.
type Spring struct {
	A *Particle
	B *Particle

	// Strength, when high makes a spring strong and act like a stick. When
	// set low, will make a spring weak, causing it to take a long time to
	// return to its rest length (see below).
	Strength float64
	// Damping, when set high will prevent the spring from overshooting and
	// cause it to settle down quickly. When set low, a spring oscillates.
	Damping float64
	// Rest Length is the length the spring wants to be at. It acts on the
	// particles to push or pull them exactly this far apart at all times.
	RestLength float64

	On bool
}
