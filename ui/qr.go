package ui

import (
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/maxazimi/qr/camera"
	"github.com/maxazimi/qr/ui/anim"
	"github.com/maxazimi/qr/ui/components"
	"github.com/maxazimi/qr/ui/instance"
	"github.com/maxazimi/qr/ui/lang"
	"github.com/maxazimi/qr/ui/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"image"
	"math"
	"time"
)

type QRScanner struct {
	cameraImage  *components.Image
	buttonCancel *components.Button
	animation    *anim.Animation

	scanning          bool
	value             string
	err               error
	cameraOrientation int
	opened            bool

	resultChan chan string
}

var (
	qrScannerInstance *QRScanner
)

func NewQRScanner() *QRScanner {
	if qrScannerInstance != nil {
		return qrScannerInstance
	}

	buttonCancel := components.NewButton(components.ButtonStyle{
		Radius:   5,
		TextSize: 14,
		Inset:    layout.UniformInset(10),
	})
	buttonCancel.Font.Weight = font.Bold

	cameraImage := &components.Image{
		Fit:    components.Contain,
		Radius: 10,
	}

	animation := anim.New(false, gween.NewSequence(
		gween.New(0, 1, 2.0, ease.OutCubic),
		gween.New(1, 0, 2.0, ease.OutCubic),
	))

	qrScannerInstance = &QRScanner{
		cameraImage:  cameraImage,
		buttonCancel: buttonCancel,
		animation:    animation,
	}
	return qrScannerInstance
}

func (q *QRScanner) scan() {
	if q.scanning {
		return
	}

	if err := camera.Open(0, 640, 480); err != nil {
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

func (q *QRScanner) Open() <-chan string {
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

	q.scanning = false
	camera.Close()

	select {
	case q.resultChan <- q.value:
	}

	q.value = ""
	close(q.resultChan)
	q.opened = false
}

func (q *QRScanner) Layout(gtx C) D {
	if q.buttonCancel.Clicked(gtx) {
		q.Close()
	}

	th := theme.Current().Theme
	return layout.UniformInset(15).Layout(gtx, func(gtx C) D {
		var children []layout.FlexChild
		if q.scanning {
			children = append(children, layout.Rigid(q.layoutPreview))
		}

		if q.err != nil {
			children = append(children,
				layout.Rigid(func(gtx C) D {
					lbl := material.Label(th, unit.Sp(14), q.err.Error())
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							q.buttonCancel.Text = lang.Str("cancel")
							q.buttonCancel.Colors = theme.Current().ButtonColors
							return q.buttonCancel.Layout(gtx)
						}),
					)
				}),
			)
		}

		if q.scanning {
			children = append(children,
				layout.Rigid(layout.Spacer{Height: 15}.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx, layout.Flexed(1, func(gtx C) D {
						q.buttonCancel.Text = lang.Str("cancel")
						q.buttonCancel.Colors = theme.Current().ButtonColors
						return q.buttonCancel.Layout(gtx)
					}))
				}),
			)
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
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
			const offset = 20
			bounds := image.Rect(offset*2, offset, imageDims.Size.X-offset*2, imageDims.Size.Y-offset)
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
				Min: image.Point{X: bounds.Min.X, Y: lineHeight - 1},
				Max: image.Point{X: bounds.Max.X, Y: lineHeight + 1},
			}.Op()
			paint.FillShape(gtx.Ops, theme.WhiteColor, line)

			drawRectangle(gtx, bounds)
			return D{Size: imageDims.Size}
		}),
	)
}

func drawRectangle(gtx C, bounds image.Rectangle) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y)))
	path.LineTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Min.Y)))
	path.LineTo(f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y)))
	path.LineTo(f32.Pt(float32(bounds.Min.X), float32(bounds.Max.Y)))
	path.Close()
	paint.FillShape(gtx.Ops, theme.WhiteColor, clip.Stroke{Path: path.End(), Width: 2}.Op())
}
