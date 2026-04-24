// SPDX-License-Identifier: Unlicense OR MIT

package traer

import (
	"math"
	"testing"
)

const eps = 1e-9

func vec3Equal(a, b Vec3) bool {
	return math.Abs(a.X-b.X) < eps &&
		math.Abs(a.Y-b.Y) < eps &&
		math.Abs(a.Z-b.Z) < eps
}

func TestVec3Dot(t *testing.T) {
	tests := []struct {
		a, b Vec3
		want float64
	}{
		{Vec3{1, 2, 3}, Vec3{4, -5, 6}, 4 - 10 + 18},
		{Vec3{}, Vec3{1, 2, 3}, 0},
		{Vec3{1, 0, 0}, Vec3{0, 1, 0}, 0},
	}
	for _, tc := range tests {
		if got := tc.a.Dot(tc.b); got != tc.want {
			t.Errorf("%v.Dot(%v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestVec3Length(t *testing.T) {
	v := Vec3{3, 4, 0}
	if got := v.LengthSquared(); got != 25 {
		t.Errorf("LengthSquared = %v, want 25", got)
	}
	if got := v.Length(); math.Abs(got-5) > eps {
		t.Errorf("Length = %v, want 5", got)
	}
}

func TestVec3Arithmetic(t *testing.T) {
	a := Vec3{1, 2, 3}
	b := Vec3{4, 5, 6}
	if got := a.Add(b); !vec3Equal(got, Vec3{5, 7, 9}) {
		t.Errorf("Add: got %v", got)
	}
	if got := a.Subtract(b); !vec3Equal(got, Vec3{-3, -3, -3}) {
		t.Errorf("Subtract: got %v", got)
	}
	if got := a.Scale(2); !vec3Equal(got, Vec3{2, 4, 6}) {
		t.Errorf("Scale: got %v", got)
	}
}

func TestVec3AssignOps(t *testing.T) {
	v := Vec3{1, 2, 3}
	v.AddAssign(Vec3{10, 20, 30})
	if !vec3Equal(v, Vec3{11, 22, 33}) {
		t.Errorf("AddAssign: got %v", v)
	}
	v.SubtractAssign(Vec3{1, 2, 3})
	if !vec3Equal(v, Vec3{10, 20, 30}) {
		t.Errorf("SubtractAssign: got %v", v)
	}
	v.ScaleAssign(0.5)
	if !vec3Equal(v, Vec3{5, 10, 15}) {
		t.Errorf("ScaleAssign: got %v", v)
	}
}
