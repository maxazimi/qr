// Inspired from https://github.com/g45t345rt/g45w/blob/master/components/notification.go

package components

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/ui/instance"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"image/color"
	"time"
)

type NotificationStyle struct {
	TitleColor color.NRGBA
	Icon       *widget.Icon
	Direction  layout.Direction
	OuterInset layout.Inset
	InnerInset layout.Inset
	Radius     unit.Dp
	Animation  ModalAnimation
}

type NotificationModal struct {
	*Modal
	NotificationStyle
	title      string
	text       string
	textEditor *widget.Editor
	timer      *time.Timer
}

func NewNotificationModal(s NotificationStyle) *NotificationModal {
	modal := NewModal(s.Direction, s.OuterInset, s.InnerInset, s.Radius, s.Animation)
	editor := new(widget.Editor)
	editor.ReadOnly = true

	return &NotificationModal{
		NotificationStyle: s,
		Modal:             modal,
		textEditor:        editor,
	}
}

func (n *NotificationModal) SetText(title, text string) *NotificationModal {
	n.title = title
	n.text = text
	return n
}

func (n *NotificationModal) SetVisible(visible bool, closeAfter time.Duration) {
	if visible {
		if n.timer != nil {
			n.timer.Stop()
		}
		if closeAfter > 0 {
			n.timer = time.AfterFunc(closeAfter, func() {
				n.Modal.SetVisible(false)
				instance.Window().Invalidate()
			})
		}
	}
	n.Modal.SetVisible(visible)
	instance.Window().Invalidate()
}

func (n *NotificationModal) Closed() bool {
	return n.Modal.Closed()
}

func (n *NotificationModal) Layout(gtx C) D {
	return n.Modal.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if n.Icon != nil {
					return n.Icon.Layout(gtx, n.TitleColor)
				}
				return D{}
			}),
			layout.Rigid(layout.Spacer{Width: 10}.Layout),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						label := material.Label(theme.Current().Theme, unit.Sp(18), n.title)
						label.Font.Weight = font.Bold
						label.Color = n.TitleColor
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						editor := material.Editor(theme.Current().Theme, n.textEditor, "")
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
}
