// Inspired from https://github.com/g45t345rt/g45w/blob/master/components/button.go

package components

import (
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/ui/anim"
	"github.com/maxazimi/v2ray-gio/ui/instance"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"github.com/maxazimi/v2ray-gio/ui/values"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"image"
	"image/color"
	"time"
)

const (
	ratio          float32 = 0.95
	tooltipPadding         = 12
)

type ButtonAnimation struct {
	clickable    *widget.Clickable
	animIn       *anim.Animation
	transFunc    anim.TransFunc
	animOut      *anim.Animation
	animClick    *anim.Animation
	animLoading  *anim.Animation
	transLoading anim.TransFunc
}

func NewButtonAnimationDefault() *ButtonAnimation {
	return NewButtonAnimationScale(ratio)
}

func NewButtonAnimationScale(v float32) *ButtonAnimation {
	animIn := anim.New(false, gween.NewSequence(gween.New(1, v, .1, ease.Linear)))
	animOut := anim.New(false, gween.NewSequence(gween.New(v, 1, .1, ease.Linear)))

	animClick := anim.New(false, gween.NewSequence(
		gween.New(1, v*ratio, .1, ease.Linear),
		gween.New(v*ratio, 1, .4, ease.OutBounce),
	))

	animLoading := anim.New(false, gween.NewSequence(gween.New(0, 1, 1, ease.Linear)))
	animLoading.Sequence.SetLoop(-1)

	return &ButtonAnimation{
		clickable:    new(widget.Clickable),
		animIn:       animIn,
		transFunc:    anim.TransScale,
		animOut:      animOut,
		animClick:    animClick,
		animLoading:  animLoading,
		transLoading: anim.TransRotate,
	}
}

func (b *ButtonAnimation) Hovered() bool {
	return b.clickable.Hovered()
}

func (b *ButtonAnimation) Clicked(gtx C) bool {
	return b.clickable.Clicked(gtx)
}

func (b *ButtonAnimation) Layout(gtx C, w layout.Widget) D {
	return b.clickable.Layout(gtx, func(gtx C) D {
		if b.animIn != nil {
			value, finished := b.animIn.Update(gtx)
			if !finished {
				defer b.transFunc(gtx, value).Push(gtx.Ops).Pop()
			}
		}

		if b.animOut != nil {
			value, finished := b.animOut.Update(gtx)
			if !finished {
				defer b.transFunc(gtx, value).Push(gtx.Ops).Pop()
			}
		}

		if b.animClick != nil {
			value, finished := b.animClick.Update(gtx)
			if !finished {
				defer b.transFunc(gtx, value).Push(gtx.Ops).Pop()
			}
		}

		return w(gtx)
	})
}

type ButtonStyle struct {
	Tag         interface{}
	Text        string
	Description string
	tipText     string
	Colors      theme.ButtonColors
	Radius      unit.Dp
	TextSize    unit.Sp
	Inset       layout.Inset
	Font        font.Font
	IconGap     unit.Dp
	Border      widget.Border
	Icon        *Icon
	LoadingIcon *Icon
	ImgKey      string
	img         *widget.Image
	Animation   *ButtonAnimation
}

type Button struct {
	ButtonStyle
	Clickable  *widget.Clickable
	Label      *widget.Label
	tooltip    *Tooltip
	Focused    bool
	Disabled   bool
	Loading    bool
	Flex       bool
	hoverState bool
	colorFixed bool
}

func NewButton(style ButtonStyle) *Button {
	if style.Animation == nil {
		style.Animation = NewButtonAnimationDefault()
	}

	b := &Button{
		ButtonStyle: style,
		Clickable:   new(widget.Clickable),
		Label:       new(widget.Label),
		tooltip:     NewTooltip(tooltipPadding),
		Focused:     false,
		hoverState:  true,
	}

	if b.Colors != (theme.ButtonColors{}) {
		b.colorFixed = true
	}
	return b
}

func (b *Button) SetLoading(loading bool) {
	b.Loading = loading
	b.Disabled = loading

	animLoading := b.Animation.animLoading
	if loading {
		animLoading.Reset().Start()
	} else {
		animLoading.Pause()
	}
}

func (b *Button) SetTipText(text string) {
	b.tipText = text
	time.AfterFunc(time.Second, func() {
		b.tipText = ""
		instance.Window().Invalidate()
	})
}

func (b *Button) Clicked(gtx C) bool {
	if b.Disabled {
		return false
	}

	if b.Clickable.Clicked(gtx) {
		if b.Animation.animClick != nil {
			b.Animation.animClick.Reset().Start()
		}
		return true
	}
	return false
}

func (b *Button) handleEvents(gtx C) {
	semantic.Button.Add(gtx.Ops)

	if b.Disabled {
		b.Animation.animOut.Reset()
		b.Animation.animIn.Reset()
		b.Animation.animClick.Reset()
		return
	}

	if b.Animation.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
		if b.Colors.HoverBackgroundColor != nil {
			b.Colors.BackgroundColor = *b.Colors.HoverBackgroundColor
		}
		if b.Colors.HoverTextColor != nil {
			b.Colors.TextColor = *b.Colors.HoverTextColor
		}

		if !b.hoverState {
			b.hoverState = true
			b.tipText = b.Description

			if b.Animation.animIn != nil {
				b.Animation.animIn.Start()
				gtx.Execute(op.InvalidateCmd{})
			}
			if b.Animation.animOut != nil {
				b.Animation.animOut.Reset()
			}
		}
	} else if b.hoverState {
		b.hoverState = false
		b.tipText = ""

		if b.Animation.animOut != nil {
			b.Animation.animOut.Start()
			gtx.Execute(op.InvalidateCmd{})
		}
		if b.Animation.animIn != nil {
			b.Animation.animIn.Reset()
		}
	}
}

func (b *Button) Layout(gtx C) D {
	th := theme.Current()

	if b.Loading {
		b.img = theme.GetImage(theme.ICLoading)
	} else if b.ImgKey != "" {
		b.img = theme.GetImage(b.ImgKey)
	}

	if b.img == nil && !b.colorFixed {
		b.Colors = th.ButtonColors
	}

	return b.Clickable.Layout(gtx, func(gtx C) D {
		return b.Animation.Layout(gtx, func(gtx C) D {
			b.handleEvents(gtx)

			c := op.Record(gtx.Ops)
			b.Border.Color = b.Colors.BorderColor
			dims := b.Border.Layout(gtx, func(gtx C) D {
				return b.Inset.Layout(gtx, func(gtx C) D {
					if b.Text != "" {
						return layout.Center.Layout(gtx, b.layoutText)
					} else {
						return layout.Center.Layout(gtx, b.layoutIcon)
					}
				})
			})
			m := c.Stop()

			if b.Flex {
				dims = D{Size: gtx.Constraints.Max}
			}

			bounds := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, b.Colors.BackgroundColor,
				clip.RRect{
					Rect: bounds,
					SE:   gtx.Dp(b.Radius),
					SW:   gtx.Dp(b.Radius),
					NE:   gtx.Dp(b.Radius),
					NW:   gtx.Dp(b.Radius),
				}.Op(gtx.Ops),
			)
			m.Add(gtx.Ops)

			// Tooltip
			if b.tipText != "" {
				ops := gtx.Ops
				c := op.Record(ops)
				labelDims := material.Label(th.Theme, values.TextSize10, b.tipText).Layout(gtx)
				m := c.Stop()

				left := unit.Dp(subtract(dims.Size.X, labelDims.Size.X+tooltipPadding*2)) / 2
				b.tooltip.Layout(gtx, layout.Inset{Top: -45, Left: -left}, func(gtx C) D {
					m.Add(ops)
					return labelDims
				})
			}

			return dims
		})
	})
}

func (b *Button) layoutIcon(gtx C) D {
	if b.Icon == nil && b.img == nil {
		return D{}
	}

	w := func(gtx C, c color.NRGBA) D {
		if b.Icon != nil {
			icon := b.Icon
			if b.LoadingIcon != nil && b.Loading {
				icon = b.LoadingIcon
			}
			if b.Flex {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return b.Icon.Layout(gtx, c)
				})
			}
			return icon.Layout(gtx, c)
		}

		if b.Focused {
			return widget.Image{Src: b.img.Src, Scale: b.img.Scale * 1.6}.Layout(gtx)
		}
		return b.img.Layout(gtx)
	}

	r := op.Record(gtx.Ops)
	dims := w(gtx, b.Colors.TextColor)
	c := r.Stop()

	gtx.Constraints.Min = dims.Size

	if b.Animation.animLoading != nil {
		value, finished := b.Animation.animLoading.Update(gtx)
		if !finished {
			defer b.Animation.transLoading(gtx, value).Push(gtx.Ops).Pop()
		}
	}

	c.Add(gtx.Ops)
	return dims
}

func (b *Button) layoutText(gtx C) D {
	var children []layout.FlexChild
	if b.Icon != nil {
		children = append(children,
			layout.Rigid(b.layoutIcon),
			layout.Rigid(layout.Spacer{Width: b.IconGap}.Layout),
		)
	}

	children = append(children,
		layout.Rigid(func(gtx C) D {
			paint.ColorOp{Color: b.Colors.TextColor}.Add(gtx.Ops)
			return b.Label.Layout(gtx, theme.Current().Theme.Shaper, b.Font,
				b.TextSize, b.Text, op.CallOp{})
		}),
	)

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx, children...)
}

func subtract(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
