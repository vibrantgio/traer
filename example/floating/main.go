package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const (
	WindowWidthDp  = 640
	WindowHeightDp = 480

	BallVelocity = 50
	BallRadius   = 8
	NumBalls     = 60
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lmicroseconds)
	go Floating()
	app.Main()
}

func Floating() {
	window := app.NewWindow(
		app.Title("Traer Physics: Free Floating"),
		app.Size(unit.Dp(WindowWidthDp), unit.Dp(WindowHeightDp)),
	)
	ops := new(op.Ops)
	floaters := MakeFloaters(WindowWidthDp, WindowHeightDp, BallVelocity, BallRadius, NumBalls)
	fps := FPS{}
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()

			// Target framerate in TraerAS3 was 31fps it used Tick(1).
			// We usually get 60fps so we can double the step size and by doing so splitting
			// the step time in half.
			activity := floaters.Tick(math.Max(1, fps.Value/30))

			cw := float64(frame.Size.X)
			ch := float64(frame.Size.Y)
			floaters.Position(cw, ch)
			floaters.Contour(cw, ch)

			// Fill backdrop
			backdrop := f32.Rect(0, 0, float32(cw), float32(ch))
			paint.ColorOp{colornames.Grey50}.Add(ops)
			paint.PaintOp{Rect: backdrop}.Add(ops)

			// Render contours
			stack := op.Push(ops)
			floaters.Render().Add(ops)
			paint.ColorOp{colornames.LightBlue500}.Add(ops)
			paint.PaintOp{Rect: backdrop}.Add(ops)
			stack.Pop()

			// Render attractor
			stack = op.Push(ops)
			radius := float32(20)
			color := colornames.Grey900
			if floaters.AttractorStrength < 0 {
				radius = 50
				color = colornames.Grey200
			}
			fap := floaters.Attractor.Position
			ap := f32.Pt(float32(fap.X), float32(fap.Y))
			area := CircleClip(ap, radius, ops)
			paint.ColorOp{Color: color}.Add(ops)
			paint.PaintOp{Rect: area}.Add(ops)
			stack.Pop()

			stack = op.Push(ops)
			clip.Rect(image.Rect(0, 0, frame.Size.X, frame.Size.Y)).Add(ops)
			pointer.InputOp{Tag: floaters, Types: pointer.Press | pointer.Release | pointer.Drag}.Add(ops)
			for _, e := range frame.Queue.Events(floaters) {
				if point, ok := e.(pointer.Event); ok {
					floaters.Pointer(point)
				}
			}
			stack.Pop()

			PrintText("Free Floating", backdrop.Min, 0.0, 0.0, backdrop.Dx(), H2, ops)
			fps.Tick()
			if activity > 0.01 {
				PrintText(fmt.Sprint(fps, "fps"), f32.Pt(backdrop.Min.X, backdrop.Max.Y), 0.0, 1.0, backdrop.Dx(), H4, ops)
				op.InvalidateOp{}.Add(ops)
			}
			frame.Frame(ops)
		} else if _, ok := event.(pointer.Event); !ok {
			log.Printf("%T%+v\n", event, event)
		}
	}
	os.Exit(0)
}
