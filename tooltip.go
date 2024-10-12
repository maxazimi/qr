package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"image"
)

type Tooltip struct {
	*theme.Theme
	hoverable *Hoverable
	shadow    *Shadow
}

func NewTooltip() *Tooltip {
	return &Tooltip{
		Theme:     theme.Current(),
		hoverable: NewHoverable(),
		shadow:    NewShadow(),
	}
}

func (t *Tooltip) layout(gtx C, pos layout.Inset, w layout.Widget) D {

	border := widget.Border{
		Width:        1,
		CornerRadius: 8,
	}

	return pos.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return LinearLayout{
					Width:      WrapContent,
					Height:     WrapContent,
					Padding:    layout.UniformInset(12),
					Background: t.Theme.SurfaceColor,
					Border:     border,
					Shadow:     t.shadow,
				}.Layout(gtx, layout.Rigid(w))
			}),
		)
	})
}

func (t *Tooltip) Layout(gtx C, rect image.Rectangle, pos layout.Inset, w layout.Widget) D {
	t.Theme = theme.Current()
	if t.hoverable.Hovered() {
		m := op.Record(gtx.Ops)
		t.layout(gtx, pos, w)
		call := m.Stop()
		ops := gtx.Ops
		op.Defer(ops, call)
	}
	t.hoverable.Layout(gtx, rect)
	return D{Size: rect.Min}
}
