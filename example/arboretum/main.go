package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/vibrantgio/style"
	"github.com/vibrantgio/textdraw"
	"github.com/vibrantgio/traer"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const (
	WidthDp     = 800
	HeightDp    = 600
	MinWidthDp  = WidthDp
	MinHeightDp = HeightDp

	NumNodes = 200

	AutoScale = false
)

func main() {
	go RandomArboretum()
	app.Main()
}

func RandomArboretum() {
	window := app.NewWindow(
		app.Title("Traer Physics: Random Arboretum"),
		app.Size(WidthDp, HeightDp))

	Grey100 := color.NRGBAModel.Convert(colornames.Grey100).(color.NRGBA)
	Grey900 := color.NRGBAModel.Convert(colornames.Grey900).(color.NRGBA)

	arboretum := NewArboretum()
	fps := traer.FPS{}
	oops := new(op.Ops)
	shaper := text.NewShaper(style.FontFaces())
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			gtx := layout.NewContext(oops, frame)

			// backdrop
			pointer.InputOp{Tag: arboretum, Types: pointer.Press}.Add(gtx.Ops)
			for _, event := range frame.Queue.Events(arboretum) {
				if point, ok := event.(pointer.Event); ok {
					if point.Type == pointer.Press {
						arboretum = NewArboretum()
					}
				}
			}
			paint.Fill(gtx.Ops, Grey100)

			// add child nodes to an arboretum until the total number of nodes is NumNodes.
			for i := len(arboretum.Particles); i < NumNodes; i++ {
				arboretum.AddNode()
			}

			// Target framerate in TraerAS3 was 31fps it used Tick(1).
			// We usually get 60fps so we can double the step size and by doing so splitting
			// the step time in half.
			activity := arboretum.Tick(math.Max(1, fps.Value/30))

			rect := image.Rectangle{Max: gtx.Constraints.Constrain(frame.Size)}
			metric := gtx.Metric
			if !AutoScale {
				metric = unit.Metric{PxPerDp: 1.0, PxPerSp: 1.0}
			}
			arboretum.DrawNetwork(rect, metric, gtx.Ops)

			text := textdraw.Text(shaper, style.H3, 0.0, 0.0, Grey900, "Random Arboretum")
			layout.UniformInset(12).Layout(gtx, text)
			fps.Tick()
			if activity > 2 {
				text := textdraw.Text(shaper, style.H4, 1.0, 1.0, Grey900, fmt.Sprint(fps, "fps"))
				layout.UniformInset(12).Layout(gtx, text)
				op.InvalidateOp{}.Add(gtx.Ops)
			}

			frame.Frame(gtx.Ops)
		}
	}
	os.Exit(0)
}
