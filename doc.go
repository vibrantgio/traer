// SPDX-License-Identifier: Unlicense OR MIT

// Package traer provides a simple particle system physics engine for Go.
// Designed to be application / domain agnostic. All this is supposed to do is
// let you make particles, apply forces and calculate the positions of
// particles over time in real-time. Anything else you need to handle
// yourself.
//
// There are four parts
//	ParticleSystem - takes care of gravity, drag, making particles, applying forces and advancing the simulation
//	Particles - they move around in 3D space based on forces you've applied to them
//	Springs - they act on two particles
//	Attractions - which also act on two particles
//
// Acknowledgement
//
// This package is a port of the processing library TRAER.PHYSICIS 3.0
//
// For the orginal library see http://murderandcreate.com/physics/
package traer
