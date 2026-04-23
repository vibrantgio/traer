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
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/reactivego/gio"
	"github.com/reactivego/gio/style"
	"github.com/vibrantgio/traer"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

func main() {
	go Gravity()
	app.Main()
}

func Gravity() {
	const WidthDp, HeightDp = 800, 600
	const BallVelocity, BallRadius, NumBalls = 50, 8, 60
	const AutoScale = true

	window := app.NewWindow(
		app.Title("Traer Physics: Gravity Well"),
		app.Size(WidthDp, HeightDp))

	Grey50 := color.NRGBAModel.Convert(colornames.Grey50).(color.NRGBA)
	Grey200 := color.NRGBAModel.Convert(colornames.Grey200).(color.NRGBA)
	Grey900 := color.NRGBAModel.Convert(colornames.Grey900).(color.NRGBA)
	LightBlue500 := color.NRGBAModel.Convert(colornames.LightBlue500).(color.NRGBA)

	oops := new(op.Ops)
	field := MakeField(WidthDp, HeightDp, BallVelocity, BallRadius, NumBalls)
	fps := traer.FPS{}
	shaper := text.NewShaper(style.FontFaces())
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			gtx := layout.NewContext(oops, frame)

			// Target framerate in TraerAS3 was 31fps it used Tick(1).
			// We usually get 60fps so we can double the step size and by doing so splitting
			// the step time in half.
			activity := field.Tick(math.Max(1, fps.Value/30))

			// Fill backdrop
			paint.Fill(gtx.Ops, Grey50)

			metric := frame.Metric
			if !AutoScale {
				metric = unit.Metric{PxPerDp: 1.0, PxPerSp: 1.0}
			}

			dx, dy := float64(frame.Size.X), float64(frame.Size.Y)
			field.Constrain(dx, dy)
			field.Contour(dx, dy, metric)

			// Render contours
			shape := clip.Outline{Path: field.Render(gtx.Ops)}.Op()
			paint.FillShape(gtx.Ops, LightBlue500, shape)

			// Render attractor
			radius := float32(metric.Dp(20))
			color := Grey900
			if field.AttractorStrength < 0 {
				radius = float32(metric.Dp(50))
				color = Grey200
			}
			fap := field.Attractor.Position
			ap := f32.Pt(float32(fap.X), float32(fap.Y))
			shape = clip.Outline{Path: gio.CirclePath(gtx.Ops, ap, radius)}.Op()
			paint.FillShape(gtx.Ops, color, shape)

			pointer.InputOp{Tag: field, Types: pointer.Press | pointer.Release | pointer.Drag}.Add(gtx.Ops)
			for _, e := range frame.Queue.Events(field) {
				if point, ok := e.(pointer.Event); ok {
					field.Pointer(point)
				}
			}

			layout.UniformInset(12).Layout(gtx, gio.Text(shaper, style.H3, 0.0, 0.0, Grey900, "Gravity Well"))
			fps.Tick()
			if activity > 0.01 {
				layout.UniformInset(12).Layout(gtx, gio.Text(shaper, style.H4, 1.0, 1.0, Grey900, fmt.Sprint(fps, "fps")))
				op.InvalidateOp{}.Add(gtx.Ops)
			}

			frame.Frame(gtx.Ops)
		}
	}
	os.Exit(0)
}
