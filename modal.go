// from https://github.com/g45t345rt/g45w/blob/master/components/modal.go

package components

import (
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/ui/anim"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"image"
	"image/color"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type ModalAnimation struct {
	animIn   *anim.Animation
	transIn  anim.TransFunc
	animOut  *anim.Animation
	transOut anim.TransFunc
}

func NewModalAnimationScaleBounce() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(0, 1, .25, ease.OutBounce),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(1, 0, .25, ease.OutBounce),
		)),
		transIn:  anim.TransScale,
		transOut: anim.TransScale,
	}
}

func NewModalAnimationUp() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(1, 0, .25, ease.OutCubic),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(0, 1, .25, ease.InCubic),
		)),
		transIn:  anim.TransY,
		transOut: anim.TransY,
	}
}

func NewModalAnimationDown() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(-1, 0, .25, ease.OutCubic),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(0, -1, .25, ease.InCubic),
		)),
		transIn:  anim.TransY,
		transOut: anim.TransY,
	}
}

type ModalColors struct {
	BackgroundColor color.NRGBA
	BackdropColor   *color.NRGBA
}

type ModalStyle struct {
	ModalColors
	CloseOnInsideClick bool
	Direction          layout.Direction
	Inset              layout.Inset
	Radius             unit.Dp
	Animation          ModalAnimation
}

type Modal struct {
	ModalStyle
	Visible      bool
	CloseKeySet  key.Set
	clickableOut *widget.Clickable
	clickableIn  *widget.Clickable
	closed       bool
}

func NewModal(style ModalStyle) *Modal {
	return &Modal{
		CloseKeySet:  key.NameEscape + "|" + key.NameBack,
		ModalStyle:   style,
		Visible:      false,
		clickableOut: new(widget.Clickable),
		clickableIn:  new(widget.Clickable),
	}
}

func (m *Modal) SetVisible(visible bool) {
	if visible {
		m.Visible = true
		m.Animation.animIn.Start()
		m.Animation.animOut.Reset()
	} else {
		m.Animation.animOut.Start()
		m.Animation.animIn.Reset()
	}
}

func (m *Modal) Closed() bool {
	return m.closed
}

func (m *Modal) handleKeyClose(gtx C) {
	if m.CloseKeySet != "" {
		key.InputOp{
			Tag:  m,
			Keys: m.CloseKeySet,
		}.Add(gtx.Ops)

		for _, e := range gtx.Events(m) {
			switch e := e.(type) {
			case key.Event:
				if e.State == key.Press {
					m.SetVisible(false)
				}
			}
		}
	}
}

func (m *Modal) Layout(gtx C, w layout.Widget) D {
	m.closed = false
	if !m.Visible {
		return D{Size: gtx.Constraints.Max}
	}

	m.handleKeyClose(gtx)

	animIn := m.Animation.animIn
	animOut := m.Animation.animOut
	transIn := m.Animation.transIn
	transOut := m.Animation.transOut

	clickedOut := m.clickableOut.Clicked()
	clickedIn := m.clickableIn.Clicked()

	if !m.CloseOnInsideClick {
		if clickedOut && !clickedIn {
			animOut.Start()
		}
	} else {
		if clickedIn {
			animOut.Start()
		}

		if m.clickableIn.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}
	}

	if m.BackdropColor != nil {
		bgColor := *m.BackdropColor
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}

	if animIn != nil {
		state := animIn.Update(gtx)
		if state.Active {
			defer transIn(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	if animOut != nil {
		state := animOut.Update(gtx)
		if state.Active {
			defer transOut(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			m.Visible = false
			m.closed = true
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	r := op.Record(gtx.Ops)
	dims := m.Inset.Layout(gtx, func(gtx C) D {
		return m.Direction.Layout(gtx, func(gtx C) D {
			r := op.Record(gtx.Ops)
			dims := m.clickableIn.Layout(gtx, w)
			c := r.Stop()

			bgColor := m.BackgroundColor
			paint.FillShape(gtx.Ops, bgColor,
				clip.RRect{
					Rect: image.Rectangle{Max: dims.Size},
					SE:   gtx.Dp(m.Radius),
					SW:   gtx.Dp(m.Radius),
					NW:   gtx.Dp(m.Radius),
					NE:   gtx.Dp(m.Radius),
				}.Op(gtx.Ops),
			)

			c.Add(gtx.Ops)
			return dims
		})
	})
	c := r.Stop()

	if !m.CloseOnInsideClick {
		return m.clickableOut.Layout(gtx, func(gtx C) D {
			c.Add(gtx.Ops)
			return dims
		})
	}

	c.Add(gtx.Ops)
	return dims
}
