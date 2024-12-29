//go:build android || ios

package values

import (
	"gioui.org/unit"
)

const (
	AppWidth  = unit.Dp(1080)
	AppHeight = unit.Dp(2166)

	DP16  = unit.Dp(12)
	DP18  = unit.Dp(12)
	DP24  = unit.Dp(16)
	DP30  = unit.Dp(16)
	DP32  = unit.Dp(16)
	DP420 = unit.Dp(340)

	TextSize18 = unit.Sp(16)
	TextSize20 = unit.Sp(18)
	TextSize22 = unit.Sp(18)
	TextSize24 = unit.Sp(20)
	TextSize28 = unit.Sp(22)
	TextSize30 = unit.Sp(24)
	TextSize32 = unit.Sp(28)
	TextSize34 = unit.Sp(28)
	TextSize60 = unit.Sp(34)

	ModalRadius = unit.Dp(45)
)
