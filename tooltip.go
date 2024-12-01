package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/ui/theme"
)

type Tooltip struct {
	shadow  *Shadow
	padding layout.Inset
}

func NewTooltip(padding unit.Dp) *Tooltip {
	return &Tooltip{
		shadow:  NewShadow(),
		padding: layout.UniformInset(padding),
	}
}

func (t *Tooltip) layout(gtx C, pos layout.Inset, w layout.Widget) D {
	th := theme.Current()
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
					Padding:    t.padding,
					Background: th.SurfaceColor,
					Border:     border,
					Shadow:     t.shadow,
				}.Layout(gtx, layout.Rigid(w))
			}),
		)
	})
}

func (t *Tooltip) Layout(gtx C, pos layout.Inset, w layout.Widget) D {
	m := op.Record(gtx.Ops)
	t.layout(gtx, pos, w)
	call := m.Stop()
	ops := gtx.Ops
	op.Defer(ops, call)
	return D{Size: gtx.Constraints.Min}
}
