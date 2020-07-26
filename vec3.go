// SPDX-License-Identifier: Unlicense OR MIT

package traer

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Dot(a Vec3) float64 {
	return v.X*a.X + v.Y*a.Y + v.Z*a.Z
}

func (v Vec3) LengthSquared() float64 {
	return v.Dot(v)
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec3) Scale(scalar float64) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec3) Add(a Vec3) Vec3 {
	return Vec3{v.X + a.X, v.Y + a.Y, v.Z + a.Z}
}

func (v Vec3) Subtract(a Vec3) Vec3 {
	return Vec3{v.X - a.X, v.Y - a.Y, v.Z - a.Z}
}

func (v *Vec3) ScaleAssign(scalar float64) {
	v.X, v.Y, v.Z = v.X*scalar, v.Y*scalar, v.Z*scalar
}

func (v *Vec3) AddAssign(a Vec3) {
	v.X, v.Y, v.Z = v.X+a.X, v.Y+a.Y, v.Z+a.Z
}

func (v *Vec3) SubtractAssign(a Vec3) {
	v.X, v.Y, v.Z = v.X-a.X, v.Y-a.Y, v.Z-a.Z
}
