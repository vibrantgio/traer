package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const (
	BaseSizeDp     = 640
	WindowWidthDp  = BaseSizeDp
	WindowHeightDp = BaseSizeDp
	MinWidthDp     = WindowWidthDp
	MinHeightDp    = WindowHeightDp

	NumNodes = 200

	AutoScale = true
)

var (
	Grey100 = color.NRGBAModel.Convert(colornames.Grey100).(color.NRGBA)
	Grey900 = color.NRGBAModel.Convert(colornames.Grey900).(color.NRGBA)
)

func main() {
	rand.Seed(time.Now().UnixNano())
	go RandomArboretum(NumNodes)
	app.Main()
}

func RandomArboretum(NumNodes int) {
	window := app.NewWindow(
		app.Title("Traer Physics: Random Arboretum"),
		app.Size(unit.Dp(WindowWidthDp), unit.Dp(WindowHeightDp)),
	)
	arboretum := MakeArboretum()
	fps := FPS{}
	ops := new(op.Ops)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()

			metric := frame.Metric
			if !AutoScale {
				metric = unit.Metric{PxPerDp: 1.0, PxPerSp: 1.0}
			}

			pointer.InputOp{Tag: arboretum, Types: pointer.Press}.Add(ops)
			for _, event := range frame.Queue.Events(arboretum) {
				if point, ok := event.(pointer.Event); ok {
					if point.Type == pointer.Press {
						arboretum = MakeArboretum()
					}
				}
			}
			if len(arboretum.Particles) == 1 {
				for i := 0; i < NumNodes; i++ {
					arboretum.AddNode()
				}
			}

			// Target framerate in TraerAS3 was 31fps it used Tick(1).
			// We usually get 60fps so we can double the step size and by doing so splitting
			// the step time in half.
			activity := arboretum.Tick(math.Max(1, fps.Value/30))

			// Fill backdrop
			paint.ColorOp{Color: Grey100}.Add(ops)
			paint.PaintOp{}.Add(ops)

			rect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))
			arboretum.DrawNetwork(rect, metric, ops)

			inset := float32(metric.Px(unit.Dp(12)))

			rect = f32.Rectangle{Min: rect.Min.Add(f32.Pt(inset, inset)), Max: rect.Max.Sub(f32.Pt(inset, inset))}
			PrintText("Random Arboretum", rect, 0.0, 0.0, H2, Grey900, metric, ops)
			fps.Tick()
			if activity > 2 {
				PrintText(fmt.Sprint(fps, "fps"), rect, 1.0, 1.0, H4, Grey900, metric, ops)
				op.InvalidateOp{}.Add(ops)
			}
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
