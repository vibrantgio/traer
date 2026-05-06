package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
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
	WindowWidthDp  = 1800 // 640
	WindowHeightDp = 800  // 480

	SystemAcceleration    = 0.0
	SystemDrag            = 0.1
	SystemMinActivity     = 0.01
	FixedSpringConstant   = 0.4
	PointerSpringConstant = 0.4
	SpringDamping         = 1.1
	MinScrollLength       = 0.1
)

func main() {
	go Scrolling()
	app.Main()
}

func Scrolling() {
	window := new(app.Window)
	window.Option(
		app.Title("Traer Physics: Kinetic Scrolling"),
		app.Size(unit.Dp(WindowWidthDp), unit.Dp(WindowHeightDp)),
	)
	oops := new(op.Ops)
	unsplash, err := img("bridge-unsplash.jpg", "https://images.unsplash.com/photo-1423347834838-5162bb452ca7?ixlib=rb-1.2.1&q=80&fm=png&crop=entropy&cs=tinysrgb&w=1080&fit=max&ixid=eyJhcHBfaWQiOjF9")
	if err != nil {
		log.Fatal(err)
	}

	// Instead of the full image use a rectangular sub image area inside the full image bounds.
	unsplash = SubImage(unsplash, 100, 0, 740, 720)

	scroller := NewScroller()
	scroller.Content = image.Rectangle{Max: unsplash.Bounds().Size()}
	shaper := text.NewShaper(text.WithCollection(style.FontFaces()))
	fps := traer.FPS{}
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			e.Insets.Top += 12
			e.Insets.Bottom += 12
			e.Insets.Left += 12
			e.Insets.Right += 12
			gtx := app.NewContext(oops, e)

			event.Op(gtx.Ops, scroller)
			for {
				ev, ok := gtx.Source.Event(pointer.Filter{Target: scroller, Kinds: pointer.Press | pointer.Release | pointer.Drag | pointer.Scroll, ScrollX: pointer.ScrollRange{Min: -1e6, Max: 1e6}, ScrollY: pointer.ScrollRange{Min: -1e6, Max: 1e6}})
				if !ok {
					break
				}
				if point, ok := ev.(pointer.Event); ok {
					scroller.Pointer(point)
				}
			}

			scroller.View = image.Rectangle{Max: gtx.Constraints.Constrain(e.Size)}

			// t values of 1.2 and higher provide a stable physics simulation
			activity := scroller.Tick(math.Max(1.2, fps.Value/30))

			scrollOffset := image.Pt(scroller.Content.Min.X, scroller.Content.Min.Y)

			// Draw image at offset
			imageOp := paint.NewImageOp(unsplash)
			// imageOp.Filter = paint.FilterNearest

			stack := op.Offset(scrollOffset).Push(gtx.Ops)
			imageOp.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Pop()

			scroller.Draw(scroller.View, gtx.Ops)

			txt := textdraw.Text(shaper, style.H3, 0.0, 0.0, colornames.Grey900, "Kinetic Scrolling")
			layout.UniformInset(12).Layout(gtx, txt)
			fps.Tick()
			if activity > SystemMinActivity {
				txt := textdraw.Text(shaper, style.H4, 1.0, 1.0, colornames.Grey900, fmt.Sprint(fps, "fps"))
				layout.UniformInset(12).Layout(gtx, txt)
				gtx.Execute(op.InvalidateCmd{})
			}

			e.Frame(gtx.Ops)
		}
	}
}
