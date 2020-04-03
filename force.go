package traer

type Force interface {
	TurnOn()
	TurnOff()
	IsOn() bool
	IsOff() bool
	Apply()
}
