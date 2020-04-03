package traer

type Spring struct {
	A              *Particle
	B              *Particle
	SpringConstant float64
	Damping        float64
	RestLength     float64
	On             bool
}

func MakeSpring(a, b *Particle, ks, d, r float64) *Spring {
	return &Spring{a, b, ks, d, r, true}
}

func (s *Spring) TurnOn() {
	s.On = true
}

func (s *Spring) TurnOff() {
	s.On = false
}

func (s *Spring) IsOn() bool {
	return s.On
}

func (s *Spring) IsOff() bool {
	return !s.On
}

func (s *Spring) Apply() {
	if !s.On || (s.A.Fixed && s.B.Fixed) {
		return
	}

	a2b := s.A.Position.Subtract(s.B.Position)
	a2bDistance := a2b.Length()
	if a2bDistance > 0.0 {
		a2b.ScaleAssign(1 / a2bDistance) // normalize a2b
	}

	// spring force is proportional to how much it stretched
	springForce := -(a2bDistance - s.RestLength) * s.SpringConstant

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
