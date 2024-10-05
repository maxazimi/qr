//go:build (darwin || linux || windows) && !android && !ios

package values

import (
	"gioui.org/unit"
)

const (
	WindowWidth  = unit.Dp(400)
	WindowHeight = unit.Dp(600)

	ModalRadius = unit.Dp(15)
)
