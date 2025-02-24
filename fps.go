package traer

import (
	"strconv"
	"time"
)

type FPS struct {
	start time.Time
	count float64

	Value float64
}

func (f *FPS) Tick() *FPS {
	if time.Since(f.start) > time.Second {
		f.Value = f.count / time.Since(f.start).Seconds()
		if f.Value == 0 {
			f.Value = 60
		}
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
