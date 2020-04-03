package traer

type Particle struct {
	Position Vec3
	Velocity Vec3
	Force    Vec3
	Mass     float64
	Fixed    bool
}

func (p *Particle) Reset() {
	p.Mass = 1.0
	p.Fixed = false
	p.Position = Vec3{}
	p.Velocity = Vec3{}
	p.Force = Vec3{}
}
