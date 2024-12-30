package instance

import (
	"gioui.org/app"
	"gioui.org/unit"
	"sync"
)

const (
	AppWidth  = unit.Dp(500)
	AppHeight = unit.Dp(500)
)

var (
	window          *app.Window
	windowMtx       sync.Mutex
	currentAppWidth unit.Dp
	isMobileView    bool
	maxMobileWidth  unit.Dp
)

func init() {
	window = new(app.Window)
	SetWindowSize(AppWidth, AppHeight)
}

func SetWindowTitle(title string) {
	window.Option(app.Title(title))
}

func SetWindowSize(width, height unit.Dp) {
	window.Option(app.Size(width, height))
	maxMobileWidth = width
}

func Window() *app.Window {
	windowMtx.Lock()
	defer windowMtx.Unlock()
	return window
}

func SetCurrentAppWidth(appWidth int, metric unit.Metric) {
	currentAppWidth = metric.PxToDp(appWidth)
	isMobileView = currentAppWidth <= maxMobileWidth
}

func CurrentAppWidth() unit.Dp {
	return currentAppWidth
}

// IsMobileView returns true if the app's window width is less than the mobile view width.
func IsMobileView() bool {
	return isMobileView
}

func MaxMobileWidth() unit.Dp {
	return maxMobileWidth
}
