package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
)

type Shadow struct {
	shadowElevation float32
	shadowRadius    float32
	bgColor         color.NRGBA
}

func NewShadow() *Shadow {
	return &Shadow{
		shadowRadius:    8, // radius of the shadow
		shadowElevation: 7, // height/spread of the shadow
	}
}

func (s *Shadow) Layout(gtx C, w layout.Widget) D {
	r := op.Record(gtx.Ops)
	dims := w(gtx)
	c := r.Stop()

	offset := pxf(gtx.Metric, unit.Dp(s.shadowElevation))
	rr := gtx.Dp(unit.Dp(s.shadowRadius))
	rect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)

	// add background color
	gradientBox(gtx.Ops, rect, rr, 0, s.bgColor)

	// shadow layers arranged from the biggest to the smallest
	gradientBox(gtx.Ops, rect, rr, int(float32(offset)/float32(0.8)), color.NRGBA{A: 0x3})
	gradientBox(gtx.Ops, rect, rr, offset, color.NRGBA{A: 0x4})
	gradientBox(gtx.Ops, rect, rr, int(float32(offset)/float32(1.5)), color.NRGBA{A: 0x7})
	gradientBox(gtx.Ops, rect, rr, int(float32(offset)/float32(2.5)), color.NRGBA{A: 0x10})

	c.Add(gtx.Ops)
	return dims
}

func (s *Shadow) SetShadowRadius(radius float32) *Shadow {
	s.shadowRadius = radius
	return s
}

func (s *Shadow) SetShadowElevation(elev float32) *Shadow {
	s.shadowElevation = elev
	return s
}

func (s *Shadow) SetBackgroundColor(col color.NRGBA) *Shadow {
	s.bgColor = col
	return s
}

func gradientBox(ops *op.Ops, r image.Rectangle, rr, spread int, col color.NRGBA) {
	paint.FillShape(ops, col, clip.RRect{
		Rect: outset(r, spread),
		SE:   rr + spread, SW: rr + spread, NW: rr + spread, NE: rr + spread,
	}.Op(ops))
}

func outset(r image.Rectangle, rr int) image.Rectangle {
	r.Min.X -= rr
	r.Min.Y -= rr
	r.Max.X += rr
	r.Max.Y += rr
	return r
}

func pxf(c unit.Metric, v unit.Dp) int {
	s := c.PxPerDp
	if s == 0 {
		s = 1
	}
	return int(s) * int(v)
}
