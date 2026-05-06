package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
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
	window := new(app.Window)
	window.Option(
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
	shaper := text.NewShaper(text.WithCollection(style.FontFaces()))
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(oops, e)

			if anchor == nil {
				cx, cy := float64(e.Size.X/2), float64(e.Size.Y/2)
				anchor = partsys.NewParticle(1.0, cx, cy, 0)
				anchor.Fixed = true
				particle = partsys.NewParticle(1.0, cx, cy, 0)
				attractor = partsys.NewParticle(1.0, cx, cy, 0)
				attractor.Fixed = true
				partsys.NewSpring(anchor, particle, 0.1, 0.01, 0)
				partsys.NewAttraction(attractor, particle, 9000, 30)
			}
			activity := partsys.Tick(math.Max(1.2, fps.Value/30))

			anchor.Position = traer.Vec3{X: float64(e.Size.X / 2), Y: float64(e.Size.Y / 2)}

			event.Op(gtx.Ops, tag)
			for {
				ev, ok := gtx.Source.Event(pointer.Filter{Target: tag, Kinds: pointer.Move})
				if !ok {
					break
				}
				if move, ok := ev.(pointer.Event); ok {
					position := traer.Vec3{X: float64(move.Position.X), Y: float64(move.Position.Y)}
					attractor.Position = position
				}
			}

			pos := f32.Pt(float32(particle.Position.X), float32(particle.Position.Y))
			cstack := clip.Outline{Path: circle.CirclePath(gtx.Ops, pos, 30)}.Op().Push(gtx.Ops)
			paint.ColorOp{Color: Red400}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cstack.Pop()

			txt := textdraw.Text(shaper, style.H3, 0.0, 0.0, Grey900, "Weak Attraction")
			layout.UniformInset(12).Layout(gtx, txt)
			fps.Tick()
			if activity > 0.1 {
				txt := textdraw.Text(shaper, style.H4, 1.0, 1.0, Grey900, fmt.Sprint(fps, "fps"))
				layout.UniformInset(12).Layout(gtx, txt)
				gtx.Execute(op.InvalidateCmd{})
			}

			e.Frame(gtx.Ops)
		}
	}
}
