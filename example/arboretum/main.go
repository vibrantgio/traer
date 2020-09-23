package main

import (
	"fmt"
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
)

func main() {
	rand.Seed(time.Now().UnixNano())
	go Arboretum(NumNodes)
	app.Main()
}

func Arboretum(NumNodes int) {
	window := app.NewWindow(
		app.Title("Traer Physics: Random Arboretum"),
		app.Size(unit.Dp(WindowWidthDp), unit.Dp(WindowHeightDp)),
	)
	arboretum := MakeArboretum()
	var fps FPS
	ops := new(op.Ops)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()
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
			activity := arboretum.Tick(math.Max(1, fps.Value/30))
			rect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))
			paint.ColorOp{colornames.Grey100}.Add(ops)
			paint.PaintOp{Rect: rect}.Add(ops)
			arboretum.DrawNetwork(rect).Add(ops)
			inset := f32.Pt(12, 12)
			rect = f32.Rectangle{Min: rect.Min.Add(inset), Max: rect.Max.Sub(inset)}
			PrintText("Random Arboretum", rect.Min, 0.0, 0.0, rect.Dx(), H2, ops)
			if activity > 2 {
				fps.Tick()
				lb := f32.Pt(rect.Min.X, rect.Max.Y)
				PrintText(fmt.Sprint(fps, "fps"), lb, 0.0, 1.0, rect.Dx(), H4, ops)
				op.InvalidateOp{}.Add(ops)
			}
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
