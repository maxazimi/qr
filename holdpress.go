package components

import (
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"image"
	"time"
)

type HoldPress struct {
	TriggerDuration time.Duration
	Triggered       bool
	hold            gesture.Click
	pressTime       *time.Time
}

func NewHoldPress(duration time.Duration) *HoldPress {
	return &HoldPress{TriggerDuration: duration}
}

func (h *HoldPress) Layout(gtx C, w layout.Widget) D {
	h.Triggered = false

	if h.pressTime != nil {
		if h.pressTime.Add(h.TriggerDuration).Before(gtx.Now) {
			h.pressTime = nil
			h.Triggered = true
		}
		gtx.Execute(op.InvalidateCmd{})
	}

	//for _, e := range h.hold.Events(gtx.Queue) {
	//	switch e.Type {
	//	case gesture.TypeClick:
	//		h.pressTime = nil
	//	case gesture.TypePress:
	//		h.pressTime = &gtx.Now
	//	case gesture.TypeCancel:
	//		h.pressTime = nil
	//	}
	//}

	for {
		event, ok := gtx.Event(
			pointer.Filter{
				Target: &h.hold,
				Kinds:  pointer.Press | pointer.Release | pointer.Cancel,
			},
		)
		if !ok {
			break
		}

		e, ok := event.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Kind {
		case pointer.Press:
			h.pressTime = &gtx.Now
		case pointer.Release:
			fallthrough
		case pointer.Cancel:
			h.pressTime = nil
		}
	}

	r := op.Record(gtx.Ops)
	dims := w(gtx)
	c := r.Stop()

	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	h.hold.Add(gtx.Ops)
	c.Add(gtx.Ops)

	return dims
}
