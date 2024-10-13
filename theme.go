package theme

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/assets"
	"github.com/maxazimi/v2ray-gio/config"
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
	th.Shaper = text.NewShaper(text.WithCollection(assets.FontCollection()))
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

	config.Get().ThemeName = name
	_ = config.Get().Save()
}

func IsDarkModeOn() bool {
	return current == DARK
}

type Backdrop struct {
	widget.Clickable
}

func (b *Backdrop) Layout(gtx C) D {
	return b.Clickable.Layout(gtx, func(gtx C) D {
		return Fill(gtx, Themes[current].BackgroundColor)
	})
}

type InputColors struct {
	BorderColor     color.NRGBA
	BackgroundColor color.NRGBA
	TextColor       color.NRGBA
	HintColor       color.NRGBA
}

type ButtonColors struct {
	TextColor            color.NRGBA
	BackgroundColor      color.NRGBA
	HoverBackgroundColor *color.NRGBA
	HoverTextColor       *color.NRGBA
	BorderColor          color.NRGBA
}

func (b ButtonColors) Reverse() ButtonColors {
	b.TextColor, b.BackgroundColor = b.BackgroundColor, b.TextColor
	b.BorderColor = b.TextColor
	return b
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

type SwitchColors struct {
	ActiveColor       color.NRGBA
	InactiveColor     color.NRGBA
	ThumbColor        color.NRGBA
	ActiveTextColor   color.NRGBA
	InactiveTextColor color.NRGBA
}

type Colors struct {
}

type Theme struct {
	*material.Theme
	Name string

	// generic colors
	BackgroundColor color.NRGBA
	ForegroundColor color.NRGBA
	//PrimaryColor    color.NRGBA
	//SurfaceColor    color.NRGBA
	//DeepBlueColor   color.NRGBA
	GreenColor color.NRGBA
	RedColor   color.NRGBA

	// specific colors
	PrimaryColor     color.NRGBA
	PageNavTextColor color.NRGBA
	TextColor        color.NRGBA
	GrayText1Color   color.NRGBA
	GrayText2Color   color.NRGBA
	GrayText3Color   color.NRGBA
	GrayText4Color   color.NRGBA
	DeepBlueColor    color.NRGBA
	Gray1Color       color.NRGBA
	Gray2Color       color.NRGBA
	Gray3Color       color.NRGBA
	Gray4Color       color.NRGBA
	Gray5Color       color.NRGBA
	SurfaceColor     color.NRGBA
	LightGrayColor   color.NRGBA

	// text colors
	//TextColor     color.NRGBA
	TextMuteColor color.NRGBA

	//IndicatorColor       color.NRGBA
	//DividerColor         color.NRGBA
	//BgColor              color.NRGBA
	//BgGradientStartColor color.NRGBA
	//BgGradientEndColor   color.NRGBA
	//HideBalanceBgColor   color.NRGBA

	// Button
	ButtonColors ButtonColors

	// Input
	InputColors InputColors

	// Clickable
	ClickableColor      color.NRGBA
	ClickableHoverColor color.NRGBA

	// Card
	CardColor      color.NRGBA
	CardHoverColor color.NRGBA

	// Switch
	SwitchColors SwitchColors

	// Modal
	ModalColors       ModalColors
	ModalButtonColors ButtonColors

	// AppBar
	AppBarColors ButtonColors

	// Bottom Bar
	BottomBarBgColor          color.NRGBA
	BottomButtonColors        ButtonColors
	BottomButtonSelectedColor color.NRGBA
	BottomShadowColor         color.NRGBA

	// List
	ListTextColor        color.NRGBA
	ListBgColor          color.NRGBA
	ListItemHoverBgColor color.NRGBA
	ListScrollBarBgColor color.NRGBA
	ListItemTagBgColor   color.NRGBA
	ListItemTagTextColor color.NRGBA

	// Images
	Images []*widget.Image
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
