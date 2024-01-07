// from https://github.com/g45t345rt/g45w/blob/master/components/notification.go

package components

import (
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"time"
)

type NotificationColors struct {
	TextColor       color.NRGBA
	BackgroundColor color.NRGBA
}

type NotificationStyle struct {
	NotificationColors
	Icon       *widget.Icon
	Direction  layout.Direction
	OuterInset layout.Inset
	InnerInset layout.Inset
	Radius     unit.Dp
	Animation  ModalAnimation
}

type NotificationModal struct {
	NotificationStyle
	Modal      *Modal
	title      string
	text       string
	textEditor *widget.Editor
	timer      *time.Timer
}

func NewNotificationModal(style NotificationStyle) *NotificationModal {
	modal := NewModal(ModalStyle{
		CloseOnInsideClick: true,
		Direction:          style.Direction,
		Inset:              style.OuterInset,
		Radius:             style.Radius,
		Animation:          style.Animation,
	})

	editor := new(widget.Editor)
	editor.ReadOnly = true

	return &NotificationModal{
		NotificationStyle: style,
		Modal:             modal,
		textEditor:        editor,
	}
}

func (n *NotificationModal) SetText(title, text string) *NotificationModal {
	n.title = title
	n.text = text
	return n
}

func (n *NotificationModal) SetVisible(w *app.Window, visible bool, closeAfter time.Duration) {
	if visible {
		if n.timer != nil {
			n.timer.Stop()
		}
		if closeAfter > 0 {
			n.timer = time.AfterFunc(closeAfter, func() {
				n.Modal.SetVisible(false)
				w.Invalidate()
			})
		}
	}
	n.Modal.SetVisible(visible)
	w.Invalidate()
}

func (n *NotificationModal) Closed() bool {
	return n.Closed()
}

func (n *NotificationModal) Layout(gtx C, th *material.Theme) D {
	n.Modal.BackgroundColor = n.BackgroundColor
	return n.Modal.Layout(gtx, func(gtx C) D {
		return n.InnerInset.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if n.Icon != nil {
						return n.Icon.Layout(gtx, n.TextColor)
					}
					return D{}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							label := material.Label(th, unit.Sp(18), n.title)
							label.Font.Weight = font.Bold
							label.Color = n.TextColor
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							editor := material.Editor(th, n.textEditor, "")
							editor.Color = n.TextColor
							if n.textEditor.Text() != n.text {
								n.textEditor.SetText(n.text)
							}
							gtx.Constraints.Max.Y = gtx.Dp(150)
							return editor.Layout(gtx)
						}),
					)
				}),
			)
		})
	})
}
