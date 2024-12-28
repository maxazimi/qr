package ui

import (
	"gioui.org/app"
	"gioui.org/unit"
	"sync"
)

var (
	window          *app.Window
	windowMtx       sync.Mutex
	currentAppWidth unit.Dp
	isMobileView    bool
	maxMobileWidth  = unit.Dp(600)
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
	isMobileView = currentAppWidth <= maxMobileWidth
}

func CurrentAppWidth() unit.Dp {
	return currentAppWidth
}

func IsMobileView() bool {
	return isMobileView
}
