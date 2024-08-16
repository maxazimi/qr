package anim

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"math"
)

type TransFunc func(gtx layout.Context, value float32) op.TransformOp

func TransScale(gtx layout.Context, value float32) op.TransformOp {
	pt := gtx.Constraints.Min.Div(2)
	origin := f32.Pt(float32(pt.X), float32(pt.Y))
	trans := f32.Affine2D{}.Scale(origin, f32.Point{X: value, Y: value})
	return op.Affine(trans)
}

func TransY(gtx layout.Context, value float32) op.TransformOp {
	pt := f32.Pt(0, float32(gtx.Constraints.Max.Y)*value)
	trans := f32.Affine2D{}.Offset(pt)
	return op.Affine(trans)
}

func TransX(gtx layout.Context, value float32) op.TransformOp {
	pt := f32.Pt(float32(gtx.Constraints.Max.X)*value, 0)
	trans := f32.Affine2D{}.Offset(pt)
	return op.Affine(trans)
}

func TransXY(gtx layout.Context, value float32) op.TransformOp {
	pt := f32.Pt(-float32(gtx.Constraints.Max.X)*value, float32(gtx.Constraints.Max.Y)*value)
	trans := f32.Affine2D{}.Offset(pt)
	return op.Affine(trans)
}

func TransRotate(gtx layout.Context, value float32) op.TransformOp {
	pt := gtx.Constraints.Min.Div(2)
	origin := f32.Pt(float32(pt.X), float32(pt.Y))
	trans := f32.Affine2D{}.Rotate(origin, value*2*math.Pi)
	return op.Affine(trans)
}
