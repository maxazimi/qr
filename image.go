package components

import (
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/image/draw"
	"image"
	"sync"
)

type Image struct {
	image.Image
	Size unit.Dp

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
