// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image/color"
	"sync"

	"golang.org/x/image/math/fixed"

	"eliasnaur.com/font/roboto/robotoblack"
	"eliasnaur.com/font/roboto/robotoblackitalic"
	"eliasnaur.com/font/roboto/robotobold"
	"eliasnaur.com/font/roboto/robotobolditalic"
	"eliasnaur.com/font/roboto/robotoitalic"
	"eliasnaur.com/font/roboto/robotolight"
	"eliasnaur.com/font/roboto/robotolightitalic"
	"eliasnaur.com/font/roboto/robotomedium"
	"eliasnaur.com/font/roboto/robotomediumitalic"
	"eliasnaur.com/font/roboto/robotoregular"
	"eliasnaur.com/font/roboto/robotothin"
	"eliasnaur.com/font/roboto/robotothinitalic"

	"gioui.org/f32"
	"gioui.org/font/opentype"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
)

var (
	once   sync.Once
	roboto []text.FontFace
)

const (
	Thin   = 100 - 400
	Light  = 200 - 400
	Normal = text.Normal
	Medium = text.Medium
	Bold   = text.Bold
	Black  = 800 - 400
)

func RobotoFontFaces() []text.FontFace {
	register := func(fnt text.Font, ttf []byte) {
		face, err := opentype.Parse(ttf)
		if err != nil {
			panic(fmt.Sprintf("failed to parse font: %v", err))
		}
		fnt.Typeface = "Roboto"
		roboto = append(roboto, text.FontFace{Font: fnt, Face: face})
	}
	once.Do(func() {
		// Normal (400)
		register(text.Font{Weight: Normal}, robotoregular.TTF)
		register(text.Font{Weight: Normal, Style: text.Italic}, robotoitalic.TTF)

		// Thin (100)
		register(text.Font{Weight: Thin}, robotothin.TTF)
		register(text.Font{Weight: Thin, Style: text.Italic}, robotothinitalic.TTF)

		// Light (200)
		register(text.Font{Weight: Light}, robotolight.TTF)
		register(text.Font{Weight: Light, Style: text.Italic}, robotolightitalic.TTF)

		// Medium (500)
		register(text.Font{Weight: Medium}, robotomedium.TTF)
		register(text.Font{Weight: Medium, Style: text.Italic}, robotomediumitalic.TTF)

		// Bold (600)
		register(text.Font{Weight: Bold}, robotobold.TTF)
		register(text.Font{Weight: Bold, Style: text.Italic}, robotobolditalic.TTF)

		// Black (800)
		register(text.Font{Weight: Black}, robotoblack.TTF)
		register(text.Font{Weight: Black, Style: text.Italic}, robotoblackitalic.TTF)
	})
	return roboto
}

var (
	RobotoThin   = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: Thin}
	RobotoLight  = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: Light}
	RobotoNormal = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: Normal}
	RobotoMedium = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: Medium}
	RobotoBold   = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: Bold}
	RobotoBlack  = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: Black}
)

type TextStyle struct {
	Font text.Font
	Size int
}

var (
	H1        = TextStyle{RobotoThin, 96}   // w300
	H2        = TextStyle{RobotoLight, 60}  // w300
	H3        = TextStyle{RobotoNormal, 48} // w400
	H4        = TextStyle{RobotoNormal, 34} // w400
	H5        = TextStyle{RobotoNormal, 24} // w400
	H6        = TextStyle{RobotoMedium, 20} // w500
	Subtitle1 = TextStyle{RobotoNormal, 16} // w400
	Subtitle2 = TextStyle{RobotoMedium, 14} // w500
	BodyText1 = TextStyle{RobotoNormal, 16} // w400
	BodyText2 = TextStyle{RobotoNormal, 14} // w400
	Button    = TextStyle{RobotoMedium, 14} // w500
	Caption   = TextStyle{RobotoNormal, 12} // w400
	Overline  = TextStyle{RobotoNormal, 10} // w400
)

var shaper = text.NewCache(RobotoFontFaces())

func TextSize(txt string, width float32, style TextStyle) (dx, dy float32) {
	lines := shaper.LayoutString(style.Font, fixed.I(style.Size), int(width), txt)
	for _, line := range lines {
		dy += float32(line.Ascent.Ceil() + line.Descent.Ceil())
		lineWidth := float32(line.Width.Ceil())
		if dx < lineWidth {
			dx = lineWidth
		}
	}
	return
}

func PrintText(txt string, r f32.Rectangle, ax, ay float32, style TextStyle, col color.NRGBA, ops *op.Ops) (dx, dy float32) {
	lines := shaper.LayoutString(style.Font, fixed.I(style.Size), int(r.Dx()), txt)
	for _, line := range lines {
		dy += float32(line.Ascent.Ceil() + line.Descent.Ceil())
		lineWidth := float32(line.Width.Ceil())
		if dx < lineWidth {
			dx = lineWidth
		}
	}
	offset := f32.Pt(r.Min.X+ax*(r.Dx()-dx), r.Min.Y+ay*(r.Dy()-dy))
	for _, line := range lines {
		state := op.Save(ops)
		offset.Y += float32(line.Ascent.Ceil())
		op.Offset(offset).Add(ops)
		offset.Y += float32(line.Descent.Ceil())
		shaper.Shape(style.Font, fixed.I(style.Size), line.Layout).Add(ops)
		paint.ColorOp{Color: col}.Add(ops)
		paint.PaintOp{}.Add(ops)
		state.Load()
	}
	return
}
