package ui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"math"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type UI struct {
}

var (
	th = material.NewTheme()
)

func New() *UI {
	SetWindowTitle("QR Scanner")
	return &UI{}
}

func (ui *UI) Run() error {
	return ui.loop()
}

func (ui *UI) loop() error {
	var (
		events = make(chan event.Event)
		acks   = make(chan struct{})
		ops    op.Ops
	)

	th := material.NewTheme(gofont.Collection())
	var button widget.Clickable

	// Image to display (this is just a placeholder)
	img := paint.ImageOp{}
	img = paint.NewImageOp(image.NewRGBA(image.Rect(0, 0, 100, 100)))

	go func() {
		for {
			e := Window().Event()
			events <- e
			<-acks
			if _, ok := e.(app.DestroyEvent); ok {
				return
			}
		}
	}()

	for {
		select {
		case e := <-events:
			switch e := e.(type) {
			case app.DestroyEvent:
				return fmt.Errorf("program terminated")
			case app.FrameEvent:
				SetCurrentAppWidth(e.Size.X, e.Metric)
				gtx := app.NewContext(&ops, e)
				layoutUI(gtx)
				e.Frame(gtx.Ops)
			default:
			} // gio events
			acks <- struct{}{}
		}
	}
}

func layoutUI(gtx C) D {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(material.H3(th, "QR test").Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			paint.FillShape(gtx.Ops,
				color.NRGBA{R: 255, G: 255, B: 255, A: 255},
				clip.Rect{Max: f32.Pt(float32(size.X), float32(size.Y))}.Op())
			paint.ImageOp{
				Src: img,
				Rect: f32.Rectangle{
					Max: f32.Point{
						X: float32(size.X),
						Y: float32(size.Y),
					},
				},
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: size}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &button, "Click Me")
				return btn.Layout(gtx)
			})
		}),
	)
}
func layoutPreview(gtx C) D {
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
