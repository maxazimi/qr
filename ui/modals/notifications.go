package modals

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/maxazimi/qr/ui/components"
	"github.com/maxazimi/qr/ui/ev"
	"github.com/maxazimi/qr/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"time"
)

const (
	closeAfterDefault = 3 * time.Second
	modalRadius       = unit.Dp(25)
)

const (
	INFO = iota
	SUCCESS
	ERROR
)

var (
	notificationModals []*components.NotificationModal
)

func init() {
	outerInset := layout.UniformInset(10)
	innerInset := layout.Inset{Top: 10, Bottom: 10, Left: 15, Right: 15}

	icon, _ := widget.NewIcon(icons.ActionInfo)
	style := components.NotificationStyle{
		Icon:       icon,
		Direction:  layout.N,
		OuterInset: outerInset,
		InnerInset: innerInset,
		Radius:     modalRadius,
		Animation:  components.NewModalAnimationDown(),
	}

	info := components.NewNotificationModal(style)
	ev.RequestEvent(ev.AddModalEvent{Modal: info})

	icon, _ = widget.NewIcon(icons.AlertError)
	style.Icon = icon
	err := components.NewNotificationModal(style)
	ev.RequestEvent(ev.AddModalEvent{Modal: err})

	icon, _ = widget.NewIcon(icons.ActionCheckCircle)
	style.Icon = icon
	success := components.NewNotificationModal(style)
	ev.RequestEvent(ev.AddModalEvent{Modal: success})

	notificationModals = append(notificationModals, success, err, info)
}

func Notify(level int, text string) {
	switch level {
	case SUCCESS:
		notificationModals[level].TitleColor = theme.Current().GreenColor
		notificationModals[level].SetText("Success", text).SetVisible(true, closeAfterDefault)
	case ERROR:
		notificationModals[level].TitleColor = theme.Current().RedColor
		notificationModals[level].SetText("Error", text).SetVisible(true, closeAfterDefault)
	case INFO:
		notificationModals[level].TitleColor = theme.Current().TextColor
		notificationModals[level].SetText("Info", text).SetVisible(true, closeAfterDefault)
	default:
	}
}
