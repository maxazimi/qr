package ui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/maxazimi/qr/ui/components"
	"github.com/maxazimi/qr/ui/instance"
	"github.com/maxazimi/qr/ui/theme"
	"github.com/maxazimi/qr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type UI struct {
	*components.AppBar
	qrScanner *QRScanner
}

func New() *UI {
	instance.SetWindowTitle("QR Scanner")
	instance.SetWindowSize(values.AppWidth, values.AppHeight+200)
	return &UI{
		AppBar:    components.NewAppBar(),
		qrScanner: NewQRScanner(),
	}
}

func (ui *UI) Run() error {
	return ui.loop()
}

func (ui *UI) loop() error {
	go func() {
		result := <-ui.qrScanner.Open()
		fmt.Println("Result: ", result)
	}()

	var (
		events = make(chan event.Event)
		acks   = make(chan struct{})
		ops    op.Ops
	)

	go func() {
		for {
			e := instance.Window().Event()
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
				instance.SetCurrentAppWidth(e.Size.X, e.Metric)
				gtx := app.NewContext(&ops, e)
				ui.layout(gtx)
				e.Frame(gtx.Ops)
			default:
			} // gio events
			acks <- struct{}{}
		}
	}
}

func (ui *UI) layout(gtx C) D {
	for _, e := range ui.AppBar.Events(gtx) {
		switch e.(type) {
		case components.AppBarMoreActionClicked:
		default:
		}
	}

	theme.BackdropInst.Layout(gtx)
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(.2, func(gtx C) D {
			return ui.AppBar.Layout(gtx)
		}),
		layout.Flexed(.8, func(gtx C) D {
			if ui.qrScanner.Opened() {
				return layout.UniformInset(100).Layout(gtx, ui.qrScanner.Layout)
			}
			return D{}
		}),
	)
}
