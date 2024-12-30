package ui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/qr/platform"
	"github.com/maxazimi/qr/ui/ev"
	"github.com/maxazimi/qr/ui/instance"
	"github.com/maxazimi/qr/ui/modals"
	"github.com/maxazimi/qr/ui/theme"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type UI struct {
	qrScanner    *QRScanner
	buttonScan   widget.Clickable
	buttonCancel widget.Clickable
}

func New() *UI {
	instance.SetWindowTitle("QR Scanner")
	return &UI{
		qrScanner: NewQRScanner(),
	}
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
				ev.RequestEvent(ev.ExitEvent{})
			case app.FrameEvent:
				instance.SetCurrentAppWidth(e.Size.X, e.Metric)
				gtx := app.NewContext(&ops, e)
				ui.layout(gtx)
				e.Frame(gtx.Ops)
			default:
				platform.HandleEvent(e)
			} // gio events
			acks <- struct{}{}

		case e := <-ev.Events():
			switch e := e.(type) {
			case ev.AppLoadEvent:
				go func() {
					err := <-platform.RequestCameraPermission()
					if err != nil {
						ev.RequestEvent(ev.NotifyEvent{Level: modals.ERROR, Text: err.Error()})
						return
					}
				}()
			case ev.WindowSizeEvent:
				instance.SetWindowSize(unit.Dp(e.Width), unit.Dp(e.Height))
			case ev.AddModalEvent:
				modals.Add(e.Modal.(modals.Modal))
				instance.Window().Invalidate()
			case ev.QREvent:
				go func() {
					result := <-ui.qrScanner.Open(640, 480)
					defer ui.qrScanner.Close()
					if result == "" {
						return
					}

					err := platform.OpenURL(result)
					if err != nil {
						ev.RequestEvent(ev.NotifyEvent{Level: modals.ERROR, Text: err.Error()})
					}
				}()
			case ev.NotifyEvent:
				modals.Notify(e.Level, e.Text)
			case ev.ExitEvent:
				return fmt.Errorf("program terminated")
			} // all events
		}
	}
}

func (ui *UI) handleEvents(gtx C) {
	if ui.buttonScan.Clicked(gtx) {
		ev.RequestEvent(ev.QREvent{})
	} else if ui.buttonCancel.Clicked(gtx) {
		ui.qrScanner.Close()
	}
}

func (ui *UI) layout(gtx C) D {
	defer func() {
		// Layout modals
		for _, v := range modals.Items() {
			if v.IsVisible() {
				v.Layout(gtx)
			}
		}
	}()

	ui.handleEvents(gtx)
	theme.BackdropInst.Layout(gtx)
	th := theme.Current()

	layoutCamera := func(gtx C) D {
		if !ui.qrScanner.Opened() {
			btn := material.Button(th.Theme, &ui.buttonScan, "Scan")
			btn.CornerRadius = 5
			btn.TextSize = 14
			btn.Inset = layout.UniformInset(10)
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(.1, btn.Layout),
				layout.Flexed(.9, func(gtx C) D {
					return D{}
				}),
			)
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(ui.qrScanner.Layout),
			layout.Rigid(layout.Spacer{Height: 10}.Layout),
			layout.Rigid(func(gtx C) D {
				btn := material.Button(th.Theme, &ui.buttonCancel, "Cancel")
				btn.CornerRadius = 5
				btn.TextSize = 14
				btn.Font.Weight = font.Bold
				btn.Inset = layout.UniformInset(10)
				return btn.Layout(gtx)
			}),
		)
	}

	return layout.Inset{Top: 50, Bottom: 10, Left: 50, Right: 50}.Layout(gtx, layoutCamera)
}
