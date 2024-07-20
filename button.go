// from https://github.com/g45t345rt/g45w/blob/master/components/button.go

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
	"github.com/maxazimi/v2ray-gio/ui/anim"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"image"
)

type ButtonAnimation struct {
	animIn       *anim.Animation
	transIn      anim.TransFunc
	animOut      *anim.Animation
	transOut     anim.TransFunc
	animClick    *anim.Animation
	transClick   anim.TransFunc
	animLoading  *anim.Animation
	transLoading anim.TransFunc
}

type ButtonStyle struct {
	Text        string
	Colors      theme.ButtonColors
	Radius      unit.Dp
	TextSize    unit.Sp
	Inset       layout.Inset
	Font        font.Font
	Icon        *widget.Icon
	IconGap     unit.Dp
	Animation   ButtonAnimation
	Border      widget.Border
	LoadingIcon *widget.Icon
}

type Button struct {
	ButtonStyle
	Clickable        *widget.Clickable
	Label            *widget.Label
	Focused          bool
	Disabled         bool
	Loading          bool
	Flex             bool
	animClickable    *widget.Clickable
	hoverSwitchState bool
}

func NewButtonAnimationDefault() ButtonAnimation {
	return NewButtonAnimationScale(.98)
}

func NewButtonAnimationScale(v float32) ButtonAnimation {
	animIn := anim.New(false,
		gween.NewSequence(
			gween.New(1, v, .1, ease.Linear),
		),
	)

	animOut := anim.New(false,
		gween.NewSequence(
			gween.New(v, 1, .1, ease.Linear),
		),
	)

	animClick := anim.New(false,
		gween.NewSequence(
			gween.New(1, v, .1, ease.Linear),
			gween.New(v, 1, .4, ease.OutBounce),
		),
	)

	animLoading := anim.New(false,
		gween.NewSequence(
			gween.New(0, 1, 1, ease.Linear),
		),
	)
	animLoading.Sequence.SetLoop(-1)

	return ButtonAnimation{
		animIn:       animIn,
		transIn:      anim.TransScale,
		animOut:      animOut,
		transOut:     anim.TransScale,
		animClick:    animClick,
		transClick:   anim.TransScale,
		animLoading:  animLoading,
		transLoading: anim.TransRotate,
	}
}

func NewButton(style ButtonStyle) *Button {
	return &Button{
		ButtonStyle:      style,
		Clickable:        new(widget.Clickable),
		Label:            new(widget.Label),
		animClickable:    new(widget.Clickable),
		Focused:          false,
		hoverSwitchState: false,
	}
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

func (b *Button) Clicked() bool {
	if b.Disabled {
		return false
	}

	return b.Clickable.Clicked()
}

func (b *Button) Layout(gtx C) D {
	return b.Clickable.Layout(gtx, func(gtx C) D {
		return b.animClickable.Layout(gtx, func(gtx C) D {
			semantic.Button.Add(gtx.Ops)

			if b.Animation.animIn != nil {
				state := b.Animation.animIn.Update(gtx)
				if state.Active {
					defer b.Animation.transIn(gtx, state.Value).Push(gtx.Ops).Pop()
				}
			}

			if b.Animation.animOut != nil {
				state := b.Animation.animOut.Update(gtx)
				if state.Active {
					defer b.Animation.transOut(gtx, state.Value).Push(gtx.Ops).Pop()
				}
			}

			if b.Animation.animClick != nil {
				state := b.Animation.animClick.Update(gtx)
				if state.Active {
					defer b.Animation.transClick(gtx, state.Value).Push(gtx.Ops).Pop()
				}
			}

			backgroundColor := b.Colors.BackgroundColor
			textColor := b.Colors.TextColor

			if !b.Disabled {
				if b.animClickable.Hovered() {
					pointer.CursorPointer.Add(gtx.Ops)
					if b.Colors.HoverBackgroundColor != nil {
						backgroundColor = *b.Colors.HoverBackgroundColor
					}

					if b.Colors.HoverTextColor != nil {
						textColor = *b.Colors.HoverTextColor
					}
				}

				if b.animClickable.Hovered() && !b.hoverSwitchState {
					b.hoverSwitchState = true

					if b.Animation.animIn != nil {
						b.Animation.animIn.Start()
					}

					if b.Animation.animOut != nil {
						b.Animation.animOut.Reset()
					}
				}

				if !b.animClickable.Hovered() && b.hoverSwitchState {
					b.hoverSwitchState = false

					if b.Animation.animOut != nil {
						b.Animation.animOut.Start()
					}

					if b.Animation.animIn != nil {
						b.Animation.animIn.Reset()
					}
				}

				if b.animClickable.Clicked() {
					if b.Animation.animClick != nil {
						b.Animation.animClick.Reset().Start()
					}
				}
			} else {
				b.Animation.animOut.Reset()
				b.Animation.animIn.Reset()
				b.Animation.animClick.Reset()
			}

			c := op.Record(gtx.Ops)
			b.Border.Color = b.Colors.BorderColor
			dims := b.Border.Layout(gtx, func(gtx C) D {
				return b.Inset.Layout(gtx, func(gtx C) D {
					var iconWidget layout.Widget
					if b.Icon != nil {
						iconWidget = func(gtx C) D {
							icon := b.Icon

							if b.LoadingIcon != nil && b.Loading {
								icon = b.LoadingIcon
							}

							var dims D
							r := op.Record(gtx.Ops)
							if b.Flex {
								dims = layout.Center.Layout(gtx, func(gtx C) D {
									return b.Icon.Layout(gtx, textColor)
								})
							} else {
								dims = icon.Layout(gtx, textColor)
							}
							c := r.Stop()

							gtx.Constraints.Min = dims.Size

							if b.Animation.animLoading != nil {
								state := b.Animation.animLoading.Update(gtx)
								if state.Active {
									defer b.Animation.transLoading(gtx, state.Value).Push(gtx.Ops).Pop()
								}
							}

							c.Add(gtx.Ops)
							return dims
						}
					}

					if b.Text != "" {
						var children []layout.FlexChild

						if iconWidget != nil {
							children = append(children,
								layout.Rigid(iconWidget),
								layout.Rigid(layout.Spacer{Width: b.IconGap}.Layout),
							)
						}

						children = append(children,
							layout.Rigid(func(gtx C) D {
								paint.ColorOp{Color: textColor}.Add(gtx.Ops)
								return b.Label.Layout(gtx, theme.Current().Theme.Shaper, b.Font,
									b.TextSize, b.Text, op.CallOp{})
							}),
						)

						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
						}.Layout(gtx, children...)
					} else {
						return iconWidget(gtx)
					}
				})
			})
			m := c.Stop()

			if b.Flex {
				dims = D{Size: gtx.Constraints.Max}
			}

			bounds := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor,
				clip.RRect{
					Rect: bounds,
					SE:   gtx.Dp(b.Radius),
					SW:   gtx.Dp(b.Radius),
					NE:   gtx.Dp(b.Radius),
					NW:   gtx.Dp(b.Radius),
				}.Op(gtx.Ops),
			)

			m.Add(gtx.Ops)
			return dims
		})
	})
}
