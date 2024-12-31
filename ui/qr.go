package ui

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/maxazimi/camera"
	"github.com/maxazimi/qr/ui/anim"
	"github.com/maxazimi/qr/ui/components"
	"github.com/maxazimi/qr/ui/instance"
	"github.com/maxazimi/qr/ui/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"image"
	"math"
	"time"
)

type QRScanner struct {
	cameraImage       *components.Image
	animation         *anim.Animation
	width, height     int
	cameraOrientation int
	value             string
	err               error
	scanning          bool
	opened            bool
	resultChan        chan string
}

var (
	qrScannerInstance *QRScanner
	loading           = true
)

func NewQRScanner() *QRScanner {
	if qrScannerInstance != nil {
		return qrScannerInstance
	}

	cameraImage := &components.Image{
		Fit:    components.Contain,
		Radius: 10,
	}

	animation := anim.New(false, gween.NewSequence(
		gween.New(0, 1, 2.0, ease.OutCubic),
		gween.New(1, 0, 2.0, ease.OutCubic),
	))

	qrScannerInstance = &QRScanner{
		cameraImage: cameraImage,
		animation:   animation,
	}
	return qrScannerInstance
}

func (q *QRScanner) scan() {
	if q.scanning {
		return
	}

	if err := camera.Open(0, q.width, q.height); err != nil {
		q.err = err
		return
	}

	if err := camera.StartPreview(); err != nil {
		q.err = err
		return
	}

	q.err = nil
	q.scanning = true
	frameBufferChan := camera.GetCameraFrameChan()
	qrReader := qrcode.NewQRCodeReader()

	go func() {
		for q.scanning {
			var frame image.Image
			select {
			case frame = <-frameBufferChan:
				if frame == nil || frame.(*image.RGBA) == nil || frame.Bounds().Empty() {
					continue
				}
			case <-time.After(500 * time.Millisecond): // No new frame available
				continue
			}
			q.cameraImage.Src = paint.NewImageOp(frame)

			bmp, _ := gozxing.NewBinaryBitmapFromImage(frame)
			result, err := qrReader.Decode(bmp, nil)

			if err == nil {
				q.value = result.String()
				q.Close()
			}

			instance.Window().Invalidate()
		}
		q.scanning = false
	}()
}

func (q *QRScanner) Open(width, height int) <-chan string {
	if q.opened {
		return nil
	}
	time.AfterFunc(3*time.Second, func() {
		loading = false
	})

	q.width = width
	q.height = height
	q.scan()
	q.animation.Start()
	q.opened = true
	q.resultChan = make(chan string)
	return q.resultChan
}

func (q *QRScanner) Opened() bool {
	return q.opened
}

func (q *QRScanner) Close() {
	if !q.opened {
		return
	}

	q.opened = false
	q.scanning = false
	camera.Close()

	q.resultChan <- q.value
	q.value = ""
	close(q.resultChan)
}

func (q *QRScanner) Layout(gtx C) D {
	th := theme.Current()

	var children []layout.FlexChild
	if q.scanning {
		children = append(children, layout.Rigid(q.layoutPreview))
	}

	if q.err != nil {
		children = append(children,
			layout.Rigid(func(gtx C) D {
				lbl := material.Label(th.Theme, unit.Sp(14), q.err.Error())
				return lbl.Layout(gtx)
			}),
		)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (q *QRScanner) layoutPreview(gtx C) D {
	var imageDims D
	videoWidget := func(gtx C) D {
		return q.cameraImage.LayoutTransform(gtx, func(dims D, trans f32.Affine2D) f32.Affine2D {
			imageDims = dims
			pt := dims.Size.Div(2)
			origin := f32.Pt(float32(pt.X), float32(pt.Y))
			rotate := float32(q.cameraOrientation) * (math.Pi / 180)
			return trans.Rotate(origin, rotate)
		})
	}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(videoWidget),
		layout.Stacked(func(gtx C) D {
			const offset = 30
			bounds := image.Rect(offset*2, offset, imageDims.Size.X-offset*2, imageDims.Size.Y-offset)
			makeSquare(&bounds)
			lineHeight := 0

			if q.animation.IsActive() {
				value, finished := q.animation.Update(gtx)
				lineHeight = int(float32(bounds.Max.Y-offset)*value) + offset
				if finished {
					q.animation.Start()
				}
			} else {
				lineHeight = bounds.Max.Y / 2
			}

			line := clip.Rect{
				Min: image.Point{X: bounds.Min.X + 5, Y: lineHeight - 1},
				Max: image.Point{X: bounds.Max.X - 5, Y: lineHeight + 1},
			}.Op()
			paint.FillShape(gtx.Ops, theme.WhiteColor, line)

			drawRectangle(gtx, bounds)
			return D{Size: imageDims.Size}
		}),
	)
}

var (
	ops         op.Ops
	cachedShape op.CallOp
	cachedSize  image.Point
)

func drawRectangle(gtx C, bounds image.Rectangle) {
	var path clip.Path

	if !gtx.Constraints.Max.Eq(cachedSize) || loading {
		cachedSize = gtx.Constraints.Max

		macro := op.Record(&ops)
		path.Begin(&ops)

		// Calculate the partials
		const factor = 6
		width := float32(bounds.Dx())
		height := float32(bounds.Dy())
		partialWidth := width / factor
		partialHeight := height / factor

		// Top-left corner
		path.MoveTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y)))
		path.LineTo(f32.Pt(float32(bounds.Min.X)+partialWidth, float32(bounds.Min.Y)))

		// Top-right corner
		path.MoveTo(f32.Pt(float32(bounds.Max.X)-partialWidth, float32(bounds.Min.Y)))
		path.LineTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Min.Y)))

		// Right-top corner
		path.MoveTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Min.Y)))
		path.LineTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Min.Y)+partialHeight))

		// Right-bottom corner
		path.MoveTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y)-partialHeight))
		path.LineTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y)))

		// Bottom-right corner
		path.MoveTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y)))
		path.LineTo(f32.Pt(float32(bounds.Max.X)-partialWidth, float32(bounds.Max.Y)))

		// Bottom-left corner
		path.MoveTo(f32.Pt(float32(bounds.Min.X)+partialWidth, float32(bounds.Max.Y)))
		path.LineTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Max.Y)))

		// Left-bottom corner
		path.MoveTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Max.Y)))
		path.LineTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Max.Y)-partialHeight))

		// Left-top corner
		path.MoveTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y)+partialHeight))
		path.LineTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y)))

		path.Close()
		paint.FillShape(&ops, theme.WhiteColor, clip.Stroke{Path: path.End(), Width: 2}.Op())

		cachedShape = macro.Stop()
	}

	cachedShape.Add(gtx.Ops)
}

func makeSquare(b *image.Rectangle) {
	deltaX := sub(b.Max.X, b.Min.X)
	deltaY := sub(b.Max.Y, b.Min.Y)

	if deltaX > deltaY {
		diff := (deltaX - deltaY) / 2
		b.Min.X += diff
		b.Max.X -= diff
	} else {
		diff := (deltaY - deltaX) / 2
		b.Min.Y += diff
		b.Max.Y -= diff
	}
}

func sub(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
