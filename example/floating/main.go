package main

import (
	"fmt"
	"image/color"
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

var (
	Grey50       = color.NRGBAModel.Convert(colornames.Grey50).(color.NRGBA)
	Grey200      = color.NRGBAModel.Convert(colornames.Grey200).(color.NRGBA)
	Grey900      = color.NRGBAModel.Convert(colornames.Grey900).(color.NRGBA)
	LightBlue500 = color.NRGBAModel.Convert(colornames.LightBlue500).(color.NRGBA)
)

func main() {
	go Floating()
	app.Main()
}

func Floating() {
	window := app.NewWindow(
		app.Title("Traer Physics: Free Floating"),
		app.Size(unit.Dp(WindowWidthDp), unit.Dp(WindowHeightDp)),
	)
	floaters := MakeFloaters(WindowWidthDp, WindowHeightDp, BallVelocity, BallRadius, NumBalls)
	fps := FPS{}
	ops := new(op.Ops)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()

			// Target framerate in TraerAS3 was 31fps it used Tick(1).
			// We usually get 60fps so we can double the step size and by doing so splitting
			// the step time in half.
			activity := floaters.Tick(math.Max(1, fps.Value/30))

			// Fill backdrop
			paint.ColorOp{Grey50}.Add(ops)
			paint.PaintOp{}.Add(ops)

			dx, dy := float64(frame.Size.X), float64(frame.Size.Y)
			floaters.Position(dx, dy)
			floaters.Contour(dx, dy)

			// Render contours
			stack := op.Push(ops)
			floaters.Render().Add(ops)
			paint.ColorOp{LightBlue500}.Add(ops)
			paint.PaintOp{}.Add(ops)
			stack.Pop()

			// Render attractor
			stack = op.Push(ops)
			radius := float32(20)
			color := Grey900
			if floaters.AttractorStrength < 0 {
				radius = 50
				color = Grey200
			}
			fap := floaters.Attractor.Position
			ap := f32.Pt(float32(fap.X), float32(fap.Y))
			clip.Outline{Path: Circle(ap, radius, ops)}.Op().Add(ops)
			paint.ColorOp{Color: color}.Add(ops)
			paint.PaintOp{}.Add(ops)
			stack.Pop()

			stack = op.Push(ops)
			pointer.InputOp{Tag: floaters, Types: pointer.Press | pointer.Release | pointer.Drag}.Add(ops)
			for _, e := range frame.Queue.Events(floaters) {
				if point, ok := e.(pointer.Event); ok {
					floaters.Pointer(point)
				}
			}
			stack.Pop()

			rect := f32.Rect(12, 12, float32(dx-12), float32(dy-12))
			PrintText("Free Floating", rect, 0.0, 0.0, H2, Grey900, ops)
			fps.Tick()
			if activity > 0.01 {
				PrintText(fmt.Sprint(fps, "fps"), rect, 1.0, 1.0, H4, Grey900, ops)
				op.InvalidateOp{}.Add(ops)
			}
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
