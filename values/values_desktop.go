//go:build (darwin || linux || windows) && !android && !ios

package values

import (
	"gioui.org/unit"
)

const (
	AppWidth  = unit.Dp(700)
	AppHeight = unit.Dp(650)

	DP16  = unit.Dp(16)
	DP18  = unit.Dp(18)
	DP24  = unit.Dp(24)
	DP30  = unit.Dp(30)
	DP32  = unit.Dp(32)
	DP420 = unit.Dp(420)

	TextSize18 = unit.Sp(18)
	TextSize20 = unit.Sp(20)
	TextSize22 = unit.Sp(22)
	TextSize24 = unit.Sp(24)
	TextSize28 = unit.Sp(28)
	TextSize30 = unit.Sp(30)
	TextSize32 = unit.Sp(32)
	TextSize34 = unit.Sp(34)
	TextSize60 = unit.Sp(60)

	ModalRadius = unit.Dp(15)
)
