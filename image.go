package components

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/image/draw"
	"image"
	"sync"
)

type Transform func(dims D, trans f32.Affine2D) f32.Affine2D

type Image struct {
	image.Image
	Size     unit.Dp
	Src      paint.ImageOp
	Fit      Fit
	Position layout.Direction
	Scale    float32
	Radius   unit.Dp

	aspectRatio int

	// Keep a cache for scaled images to reduce resource use.
	layoutSizeMtx sync.Mutex
	layoutSizeDp  unit.Dp
	layoutSizeImg *image.RGBA

	layoutSize2Mtx                 sync.Mutex
	layoutSize2DpX, layoutSize2DpY unit.Dp
	layoutSize2Img                 *image.RGBA
}

func NewImage(src image.Image) *Image {
	imageBounds := src.Bounds()
	return &Image{
		Image:       src,
		aspectRatio: imageBounds.Dx() / imageBounds.Dy(),
	}
}

func (img *Image) LayoutSize(gtx C, size unit.Dp) D {
	var dst *image.RGBA
	img.layoutSizeMtx.Lock()
	if img.layoutSizeDp == size {
		dst = img.layoutSizeImg
	} else {
		dst = image.NewRGBA(image.Rectangle{Max: image.Point{X: int(size * 2), Y: int(size * 2)}})
		img.layoutSizeImg = dst
		img.layoutSizeDp = size
		draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)
	}
	img.layoutSizeMtx.Unlock()

	i := widget.Image{Src: paint.NewImageOp(dst)}
	i.Scale = .5 // reduced the original scale of 1 by half to fix blurry images
	return i.Layout(gtx)
}

func (img *Image) Layout(gtx C) D {
	return img.LayoutSize(gtx, img.Size)
}

func (img *Image) LayoutTransform(gtx C, transform Transform) D {
	scale := img.Scale
	if scale == 0 {
		scale = 1
	}

	size := img.Src.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Dp(unit.Dp(wf*scale)), gtx.Dp(unit.Dp(hf*scale))

	dims, trans := img.Fit.scale(gtx.Constraints, img.Position, layout.Dimensions{Size: image.Pt(w, h)})

	defer clip.RRect{
		Rect: image.Rectangle{Max: dims.Size},
		NW:   gtx.Dp(img.Radius), NE: gtx.Dp(img.Radius),
		SE: gtx.Dp(img.Radius), SW: gtx.Dp(img.Radius),
	}.Push(gtx.Ops).Pop()

	if transform != nil {
		trans = transform(dims, trans)
	}

	pixelScale := scale * gtx.Metric.PxPerDp
	trans = trans.Mul(f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(pixelScale, pixelScale)))
	defer op.Affine(trans).Push(gtx.Ops).Pop()

	img.Src.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return dims
}
