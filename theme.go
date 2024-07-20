package theme

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/assets"
	"image"
	"image/color"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	WhiteColor     = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	BlackColor     = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	GreyColor      = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	BlueColor      = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	BlueColor1     = rgb(0x3377FF)
	BlueGreyColor  = rgb(0x496495)
	BlueGreyColor1 = rgb(0x3a517b)
)

var (
	th           = material.NewTheme()
	current      = Light
	Themes       = []*Theme{Light, Dark}
	BackdropInst = &Backdrop{}
)

func init() {
	th.Shaper = text.NewShaper(text.WithCollection(assets.FontCollection()))
	th.Palette.ContrastFg = WhiteColor
	th.Palette.ContrastBg = BlueGreyColor
}

func Current() *Theme {
	return current
}

func SetCurrent(key string) {
	th := Get(key)
	if th != nil {
		current = th
		current.Theme.Palette.Bg = current.BackgroundColor
		current.Theme.Palette.Fg = current.ForegroundColor
	}
}

func Get(key string) *Theme {
	for _, theme := range Themes {
		if theme.Key == key {
			return theme
		}
	}
	return nil
}

func IsDarkModeOn() bool {
	if current == Dark {
		return true
	}
	return false
}

type Backdrop struct {
	widget.Clickable
}

func (b *Backdrop) Layout(gtx C) D {
	return b.Clickable.Layout(gtx, func(gtx C) D {
		return Fill(gtx, current.BackgroundColor)
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
	temp := b.TextColor
	b.TextColor = b.BackgroundColor
	b.BackgroundColor = temp
	return b
}

type ModalColors struct {
	BackgroundColor color.NRGBA
	BackdropColor   *color.NRGBA
}

type NotificationColors struct {
	TextColor       color.NRGBA
	BackgroundColor color.NRGBA
}

type Theme struct {
	*material.Theme

	Key  string
	Name string

	// generic colors
	BackgroundColor color.NRGBA
	ForegroundColor color.NRGBA
	PrimaryColor    color.NRGBA
	SurfaceColor    color.NRGBA
	DeepBlueColor   color.NRGBA

	// text colors
	TextColor     color.NRGBA
	TextMuteColor color.NRGBA

	IndicatorColor       color.NRGBA
	DividerColor         color.NRGBA
	BgColor              color.NRGBA
	BgGradientStartColor color.NRGBA
	BgGradientEndColor   color.NRGBA
	HideBalanceBgColor   color.NRGBA

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
	SwitchActiveColor       color.NRGBA
	SwitchInactiveColor     color.NRGBA
	SwitchThumbColor        color.NRGBA
	SwitchActiveTextColor   color.NRGBA
	SwitchInactiveTextColor color.NRGBA

	// Modal
	ModalColors       ModalColors
	ModalButtonColors ButtonColors

	// Notifications
	NotificationSuccessColors NotificationColors
	NotificationErrorColors   NotificationColors
	NotificationInfoColors    NotificationColors

	// List
	ListTextColor        color.NRGBA
	ListBgColor          color.NRGBA
	ListItemHoverBgColor color.NRGBA
	ListScrollBarBgColor color.NRGBA
	ListItemTagBgColor   color.NRGBA
	ListItemTagTextColor color.NRGBA

	// Images
	ArrowDownArcImage paint.ImageOp
	ArrowUpArcImage   paint.ImageOp
	ManageFilesImage  paint.ImageOp
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
