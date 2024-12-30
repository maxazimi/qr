package components

import (
	"gioui.org/f32"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

func paintColor(gtx C, max image.Point, r int, color color.NRGBA) {
	bounds := image.Rectangle{Max: max}
	defer clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func paintGradient(gtx C, max image.Point, r int, color1, color2 color.NRGBA) {
	bounds := image.Rectangle{Max: max}
	defer clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Push(gtx.Ops).Pop()

	paint.LinearGradientOp{
		Color1: color1,
		Color2: color2,
		Stop2:  f32.Pt(float32(max.X), float32(max.Y)),
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
