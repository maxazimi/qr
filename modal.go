// https://github.com/g45t345rt/g45w/blob/master/components/modal.go

package components

import (
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/ui/anim"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	defaultAnimDuration = .25
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
			gween.New(0, 1, defaultAnimDuration, ease.OutBounce),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(1, 0, defaultAnimDuration, ease.OutBounce),
		)),
		transIn:  anim.TransScale,
		transOut: anim.TransScale,
	}
}

func NewModalAnimationUp() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(1, 0, defaultAnimDuration, ease.OutCubic),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(0, 1, defaultAnimDuration, ease.InCubic),
		)),
		transIn:  anim.TransY,
		transOut: anim.TransY,
	}
}

func NewModalAnimationDown() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(-1, 0, defaultAnimDuration, ease.OutCubic),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(0, -1, defaultAnimDuration, ease.InCubic),
		)),
		transIn:  anim.TransY,
		transOut: anim.TransY,
	}
}

func NewModalAnimationRight() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(-1, 0, defaultAnimDuration, ease.OutCubic),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(0, -1, defaultAnimDuration, ease.InCubic),
		)),
		transIn:  anim.TransX,
		transOut: anim.TransX,
	}
}

func NewModalAnimationLeftDown() ModalAnimation {
	return ModalAnimation{
		animIn: anim.New(false, gween.NewSequence(
			gween.New(-1, 0, defaultAnimDuration, ease.OutCubic),
		)),
		animOut: anim.New(false, gween.NewSequence(
			gween.New(0, -1, defaultAnimDuration, ease.InCubic),
		)),
		transIn:  anim.TransXY,
		transOut: anim.TransXY,
	}
}

type ModalStyle struct {
	theme.ModalColors
	Direction  layout.Direction
	OuterInset layout.Inset
	InnerInset layout.Inset
	Radius     unit.Dp
	Animation  ModalAnimation
}

type Modal struct {
	ModalStyle
	Visible bool
	//CloseKeySet  string
	clickableOut *widget.Clickable
	clickableIn  *widget.Clickable
	closed       bool
}

func NewModal(direction layout.Direction, outerInset, innerInset layout.Inset, radius unit.Dp,
	animation ModalAnimation) *Modal {
	modal := &Modal{
		//CloseKeySet:  key.NameEscape + "|" + key.NameBack,
		Visible:      false,
		clickableOut: new(widget.Clickable),
		clickableIn:  new(widget.Clickable),
	}

	modal.Direction = direction
	modal.OuterInset = outerInset
	modal.InnerInset = innerInset
	modal.Radius = radius
	modal.ModalColors = theme.Current().ModalColors
	modal.Animation = animation

	return modal
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

func (m *Modal) IsVisible() bool {
	return m.Visible
}

func (m *Modal) Closed() bool {
	return m.closed
}

func (m *Modal) handleKeyClose(gtx C) {
	event.Op(gtx.Ops, m)
	for {
		ev, ok := gtx.Event(
			key.FocusFilter{
				Target: m,
			},
			key.Filter{
				Focus: m,
				Name:  key.NameEscape,
			},
			key.Filter{
				Focus: m,
				Name:  key.NameBack,
			},
		)
		if !ok {
			break
		}
		_, ok = ev.(key.Event)
		if !ok {
			continue
		}
		m.SetVisible(false)
	}

	//if m.CloseKeySet != "" {
	//	key.InputOp{
	//		Tag:  m,
	//		Keys: m.CloseKeySet,
	//	}.Add(gtx.ops)
	//
	//	for _, e := range gtx.Events(m) {
	//		switch e := e.(type) {
	//		case key.Event:
	//			if e.State == key.Press {
	//				m.SetVisible(false)
	//			}
	//		}
	//	}
	//}
}

func (m *Modal) Layout(gtx C, w layout.Widget) D {
	m.closed = false
	if !m.Visible {
		return D{Size: gtx.Constraints.Max}
	}

	m.ModalColors = theme.Current().ModalColors
	m.handleKeyClose(gtx)

	animIn := m.Animation.animIn
	animOut := m.Animation.animOut
	transIn := m.Animation.transIn
	transOut := m.Animation.transOut

	if m.clickableOut.Clicked(gtx) && !m.clickableIn.Clicked(gtx) {
		animOut.Start()
	}

	if m.BackdropColor != nil {
		bgColor := *m.BackdropColor
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}

	if animIn != nil && animIn.IsActive() {
		value, finished := animIn.Update(gtx)
		if !finished {
			defer transIn(gtx, value).Push(gtx.Ops).Pop()
		}
	}

	if animOut != nil && animOut.IsActive() {
		value, finished := animOut.Update(gtx)
		if finished {
			m.Visible = false
			m.closed = true
			gtx.Execute(op.InvalidateCmd{})
			return D{Size: gtx.Constraints.Max}
		}
		defer transOut(gtx, value).Push(gtx.Ops).Pop()
	}

	r := op.Record(gtx.Ops)
	dims := m.OuterInset.Layout(gtx, func(gtx C) D {
		return m.Direction.Layout(gtx, func(gtx C) D {
			r := op.Record(gtx.Ops)
			dims := m.clickableIn.Layout(gtx, func(gtx C) D {
				return m.InnerInset.Layout(gtx, w)
			})
			c := r.Stop()

			paintGradient(gtx, dims.Size, int(m.Radius), m.BackgroundColor, m.BackgroundColor2)

			c.Add(gtx.Ops)
			return dims
		})
	})
	c := r.Stop()

	return m.clickableOut.Layout(gtx, func(gtx C) D {
		c.Add(gtx.Ops)
		return dims
	})
}

func (m *Modal) Appear() {
	m.SetVisible(true)
}

func (m *Modal) Disappear() {
	m.SetVisible(false)
}
