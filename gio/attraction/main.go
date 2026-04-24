package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"golang.org/x/exp/shiny/materialdesign/colornames"

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

	"github.com/vibrantgio/circle"
	"github.com/vibrantgio/style"
	"github.com/vibrantgio/textdraw"

	"github.com/vibrantgio/traer"
)

const WindowWidthDp = 1280 / 2
const WindowPaddingDp = 12

func main() {
	go Attraction()
	app.Main()
}

func Attraction() {
	window := app.NewWindow(
		app.Title("Traer Physics: Weak Attraction"),
		app.Size(unit.Dp(WindowWidthDp), unit.Dp((WindowWidthDp*9)/16)),
	)

	// Grey100 := color.NRGBAModel.Convert(colornames.Grey100).(color.NRGBA)
	Grey900 := color.NRGBAModel.Convert(colornames.Grey900).(color.NRGBA)
	Red400 := color.NRGBAModel.Convert(colornames.Red400).(color.NRGBA)

	partsys := traer.NewParticleSystem(0.0, 0.3)
	partsys.Clear()

	var anchor, particle, attractor *traer.Particle

	tag := new(int)

	var fps traer.FPS
	oops := new(op.Ops)
	shaper := text.NewShaper(style.FontFaces())
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			gtx := layout.NewContext(oops, frame)

			if anchor == nil {
				cx, cy := float64(frame.Size.X/2), float64(frame.Size.Y/2)
				anchor = partsys.NewParticle(1.0, cx, cy, 0)
				anchor.Fixed = true
				particle = partsys.NewParticle(1.0, cx, cy, 0)
				attractor = partsys.NewParticle(1.0, cx, cy, 0)
				attractor.Fixed = true
				partsys.NewSpring(anchor, particle, 0.1, 0.01, 0)
				partsys.NewAttraction(attractor, particle, 9000, 30)
			}
			activity := partsys.Tick(math.Max(1.2, fps.Value/30))

			anchor.Position = traer.Vec3{X: float64(frame.Size.X / 2), Y: float64(frame.Size.Y / 2)}

			pointer.InputOp{Tag: tag, Types: pointer.Move}.Add(gtx.Ops)
			for _, event := range frame.Queue.Events(tag) {
				if move, ok := event.(pointer.Event); ok {
					position := traer.Vec3{X: float64(move.Position.X), Y: float64(move.Position.Y)}
					attractor.Position = position
				}
			}

			pos := f32.Pt(float32(particle.Position.X), float32(particle.Position.Y))
			cstack := clip.Outline{Path: circle.CirclePath(gtx.Ops, pos, 30)}.Op().Push(gtx.Ops)
			paint.ColorOp{Color: Red400}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cstack.Pop()

			text := textdraw.Text(shaper, style.H3, 0.0, 0.0, Grey900, "Weak Attraction")
			layout.UniformInset(12).Layout(gtx, text)
			fps.Tick()
			if activity > 0.1 {
				text := textdraw.Text(shaper, style.H4, 1.0, 1.0, Grey900, fmt.Sprint(fps, "fps"))
				layout.UniformInset(12).Layout(gtx, text)
				op.InvalidateOp{}.Add(gtx.Ops)
			}

			frame.Frame(gtx.Ops)
		}
	}
	os.Exit(0)
}
