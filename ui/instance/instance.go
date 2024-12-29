package instance

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/maxazimi/qr/ui/values"
	"sync"
)

var (
	window          *app.Window
	windowMtx       sync.Mutex
	currentAppWidth unit.Dp
	isMobileView    bool
)

func init() {
	window = new(app.Window)
}

func SetWindowTitle(title string) {
	window.Option(app.Title(title))
}

func SetWindowSize(width, height unit.Dp) {
	window.Option(app.Size(width, height))
}

func Window() *app.Window {
	windowMtx.Lock()
	defer windowMtx.Unlock()
	return window
}

func SetCurrentAppWidth(appWidth int, metric unit.Metric) {
	currentAppWidth = metric.PxToDp(appWidth)
	isMobileView = currentAppWidth <= values.MaxMobileWidth
}

func CurrentAppWidth() unit.Dp {
	return currentAppWidth
}

// IsMobileView returns true if the app's window width is less than the mobile view width.
func IsMobileView() bool {
	return isMobileView
}
