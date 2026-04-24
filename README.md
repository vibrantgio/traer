# traer

A simple particle-system physics engine for Go, ported from the Processing
library [TRAER.PHYSICS 3.0](http://murderandcreate.com/physics/).

The package is deliberately small and rendering-agnostic: you create particles,
connect them with springs and attractions, and call `Tick` to advance the
simulation. Drawing the result is up to you.

## Install

```
go get github.com/vibrantgio/traer
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/vibrantgio/traer"
)

func main() {
    // Particle system with gravity and a bit of drag.
    ps := traer.NewParticleSystem(9.8, 0.01)

    anchor := ps.NewParticle(1, 0, 0, 0)
    anchor.Fixed = true

    bob := ps.NewParticle(1, 10, 0, 0)
    ps.NewSpring(anchor, bob, 0.5, 0.1, 5) // strength, damping, rest length

    // Advance the simulation; pass the inverse of the desired dt.
    for range 120 {
        ps.Tick(10) // dt = 0.1s
    }
    fmt.Printf("bob at %+v\n", bob.Position)
}
```

## Concepts

- **ParticleSystem** — owns the particles and forces, applies gravity and
  drag, and advances the simulation each `Tick`.
- **Particle** — a point mass in 3D. Set `Fixed = true` to pin it in place.
- **Spring** — pulls two particles toward a rest length; tune `Strength` and
  `Damping` to control stiffness and overshoot.
- **Attraction** — pulls (or with negative strength, repels) two particles
  using an inverse-square force with a minimum distance clamp.

`Tick` takes the *inverse* of `dt` (so `Tick(10)` simulates a 0.1 s step).
Traer AS3 targeted 31 fps with `Tick(1)`; at 60 fps, `Tick(2)` gives a
comparable simulation rate.

## Examples

Runnable Gio examples live under [`gio/`](./gio):

- `gio/attraction` — pointer-driven attraction force
- `gio/gravity` — many balls pulled toward a cursor attractor
- `gio/arboretum` — procedurally grown graph with springs
- `gio/scrolling` — kinetic scrolling built out of springs

## License

Unlicense OR MIT
