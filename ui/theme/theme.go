package theme

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"strings"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	LIGHT = iota
	DARK
)

var (
	th           = material.NewTheme()
	Themes       = []*Theme{Light, Dark}
	current      = LIGHT
	BackdropInst = &Backdrop{}
)

func init() {
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette.ContrastFg = WhiteColor
	th.Palette.ContrastBg = BlueGreyColor
}

func Current() *Theme {
	return Themes[current]
}

func SetCurrent(key string) {
	name := strings.ToLower(key)
	switch name {
	case "light":
		current = LIGHT
	case "dark":
		current = DARK
	default:
		return
	}

	if Themes[current] != nil {
		Themes[current].Theme.Palette.Bg = Themes[current].BackgroundColor
		Themes[current].Theme.Palette.Fg = Themes[current].ForegroundColor
	}
}

type Backdrop struct {
	widget.Clickable
}

func (b *Backdrop) Layout(gtx C) D {
	return b.Clickable.Layout(gtx, func(gtx C) D {
		return Fill(gtx, Themes[current].BackgroundColor)
	})
}

type ModalColors struct {
	BackgroundColor  color.NRGBA
	BackgroundColor2 color.NRGBA
	TextColor        color.NRGBA
	BackdropColor    *color.NRGBA
}

type NotificationColors struct {
	TitleColor      color.NRGBA
	TextColor       color.NRGBA
	BackgroundColor color.NRGBA
}

type Theme struct {
	*material.Theme
	Name string

	// generic colors
	BackgroundColor color.NRGBA
	ForegroundColor color.NRGBA
	GreenColor      color.NRGBA
	RedColor        color.NRGBA

	// specific colors
	PrimaryColor   color.NRGBA
	DeepBlueColor  color.NRGBA
	Gray1Color     color.NRGBA
	Gray2Color     color.NRGBA
	Gray3Color     color.NRGBA
	Gray4Color     color.NRGBA
	Gray5Color     color.NRGBA
	SurfaceColor   color.NRGBA
	LightGrayColor color.NRGBA

	// text colors
	TextColor        color.NRGBA
	TextMuteColor    color.NRGBA
	PageNavTextColor color.NRGBA
	GrayText1Color   color.NRGBA
	GrayText2Color   color.NRGBA
	GrayText3Color   color.NRGBA
	GrayText4Color   color.NRGBA

	// Modal
	ModalColors ModalColors
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
