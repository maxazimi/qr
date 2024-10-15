package components

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"image"
	"image/color"
)

type CornerRadius struct {
	TopLeft     int
	TopRight    int
	BottomRight int
	BottomLeft  int
}

func Radius(radius int) CornerRadius {
	return CornerRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

func TopRadius(radius int) CornerRadius {
	return CornerRadius{
		TopLeft:  radius,
		TopRight: radius,
	}
}

func BottomRadius(radius int) CornerRadius {
	return CornerRadius{
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

const (
	defaultRadius = 14
)

type Card struct {
	CardColor      color.NRGBA
	CardHoverColor color.NRGBA
	layout.Inset
	Radius CornerRadius
}

func NewCard() Card {
	return Card{
		Radius: Radius(defaultRadius),
	}
}

func (c Card) Layout(gtx C, w layout.Widget) D {
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := gtx.Dp(unit.Dp(c.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(c.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(c.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(c.Radius.BottomLeft))
				defer clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()
				return fill(gtx, c.CardColor)
			}),
			layout.Stacked(w),
		)
	})

	return dims
}

func (c Card) HoverableLayout(gtx C, btn *Clickable, w layout.Widget) D {
	th := theme.Current()
	background := c.CardColor
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := gtx.Dp(unit.Dp(c.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(c.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(c.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(c.Radius.BottomLeft))
				defer clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()

				if btn.Hoverable && btn.button.Hovered() {
					background = th.ClickableHoverColor
				}

				return fill(gtx, background)
			}),
			layout.Stacked(w),
		)
	})

	return dims
}

func (c Card) GradientLayout(gtx C, w layout.Widget) D {
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := gtx.Dp(unit.Dp(c.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(c.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(c.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(c.Radius.BottomLeft))

				dr := image.Rectangle{Max: gtx.Constraints.Min}

				paint.LinearGradientOp{
					Stop1:  layout.FPt(dr.Min),
					Stop2:  layout.FPt(dr.Max),
					Color1: color.NRGBA{R: 0x10, G: 0xff, B: 0x10, A: 0xFF},
					Color2: color.NRGBA{R: 0x10, G: 0x10, B: 0xff, A: 0xFF},
				}.Add(gtx.Ops)
				defer clip.RRect{
					Rect: dr,
					NW:   tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{
					Size: gtx.Constraints.Max,
				}
			}),
			layout.Stacked(w),
		)
	})

	return dims
}
