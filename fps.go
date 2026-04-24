// SPDX-License-Identifier: Unlicense OR MIT

package traer

import (
	"strconv"
	"time"
)

// FPS is a simple frame-per-second meter. It recomputes Value roughly once per
// second from the number of Tick calls observed in that window.
type FPS struct {
	start time.Time
	count float64

	Value float64
}

// Tick records one frame. On the first call it initializes the measurement
// window; Value stays 0 until at least one second has elapsed.
func (f *FPS) Tick() *FPS {
	if f.start.IsZero() {
		f.start = time.Now()
		return f
	}
	elapsed := time.Since(f.start)
	if elapsed > time.Second {
		f.Value = f.count / elapsed.Seconds()
		f.start = time.Now()
		f.count = 0
	} else {
		f.count++
	}
	return f
}

func (f FPS) String() string {
	return strconv.FormatFloat(f.Value, 'f', 0, 64)
}
