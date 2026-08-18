# AGENTS.md — traer

A particle-system physics engine, ported from the Processing library
TRAER.PHYSICS 3.0 and deliberately rendering-agnostic: `ParticleSystem`
owns gravity, drag and the particles, `Particle` is a point mass in 3-D
that can be `Fixed`, `Spring` and `Attraction` act on pairs, and `Tick`
advances the simulation through a Verlet integrator —
`NewDefaultVerletIntegrator` or `NewVelocityVerletIntegrator`. Drawing the
result is the caller's problem.

**Layer.** Outside ADR-001's tier table: a support library, which the rule
binds in one direction only — every tier may import it, and it may import
nothing in the table itself. It is a particle-system library: a system
integrates particles under springs, attractions and drag, and everything it
exports is in service of that one loop. Its root module imports nothing
else in the organization. Its nested `traer/gio` module adds `circle`,
`font`, `style` and `textdraw` — those edges are the nested module's and
not the root's. That direction is measured rather than typed —
`scripts/check-layers.sh --edges` reports the graph and
`scripts/sync-agents.sh` renders these sentences from it — so correcting
them here changes nothing. The other direction is measured too and
deliberately not written down: the gate checks the graph both ways, but a
public API's consumers are unknowable, so this file says what its module
needs and never who needs it.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Modules.** `github.com/vibrantgio/traer` at the repository root, and one
nested module: `gio/` (`github.com/vibrantgio/traer/gio`). Nested-module
tags carry the directory as a prefix — `gio/v0.0.8`, not `v0.0.8`.

**Build and test.** From the repository root, and again inside each nested
module directory — `./...` does not cross a module boundary:

    go build ./... && go test ./...
