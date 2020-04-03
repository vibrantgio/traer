package traer

type Particle struct {
	Position Vec3
	Velocity Vec3
	Force    Vec3
	Mass     float64

	Fixed    bool
}
