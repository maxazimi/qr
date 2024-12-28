package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

var (
	th = material.NewTheme()
)

func init() {
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette.ContrastFg = WhiteColor
	th.Palette.ContrastBg = BlueGreyColor
}

func CurrentTheme() *material.Theme {
	return th
}

func SwitchColors() {
	if Themes[current] != nil {
		th.Palette.Bg = th.BackgroundColor
		th.Palette.Fg = th.ForegroundColor
	}
}

func Fill(gtx C, col color.NRGBA) D {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	track := image.Rectangle{
		Max: d,
	}
	defer clip.Rect(track).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, col)

	return D{Size: d}
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
