// https://github.com/g45t345rt/g45w/blob/master/theme/theme.go

package theme

import (
	"github.com/maxazimi/v2ray-gio/ui/components"
	"image/color"

	"gioui.org/op/paint"
)

var (
	WhiteColor = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	BlackColor = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	GreyColor  = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	BlueColor  = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
)

type Theme struct {
	Key            string
	Name           string
	IndicatorColor color.NRGBA

	TextColor            color.NRGBA
	TextMuteColor        color.NRGBA
	DividerColor         color.NRGBA
	BgColor              color.NRGBA
	BgGradientStartColor color.NRGBA
	BgGradientEndColor   color.NRGBA
	HideBalanceBgColor   color.NRGBA

	// Header
	HeaderBackButtonColors components.ButtonColors
	HeaderTopBgColor       color.NRGBA

	// Bottom Bar
	BottomBarBgColor          color.NRGBA
	BottomButtonColors        components.ButtonColors
	BottomButtonSelectedColor color.NRGBA

	// Node Status
	NodeStatusBgColor        color.NRGBA
	NodeStatusTextColor      color.NRGBA
	NodeStatusDotGreenColor  color.NRGBA
	NodeStatusDotYellowColor color.NRGBA
	NodeStatusDotRedColor    color.NRGBA

	// Input
	InputColors components.InputColors

	// Button
	ButtonIconPrimaryColors components.ButtonColors
	ButtonPrimaryColors     components.ButtonColors
	ButtonSecondaryColors   components.ButtonColors
	ButtonInvertColors      components.ButtonColors
	ButtonDangerColors      components.ButtonColors

	// Modal
	ModalColors       components.ModalColors
	ModalButtonColors components.ButtonColors

	// Notifications
	NotificationSuccessColors components.NotificationColors
	NotificationErrorColors   components.NotificationColors
	NotificationInfoColors    components.NotificationColors

	// List
	ListTextColor        color.NRGBA
	ListBgColor          color.NRGBA
	ListItemHoverBgColor color.NRGBA
	ListScrollBarBgColor color.NRGBA
	ListItemTagBgColor   color.NRGBA
	ListItemTagTextColor color.NRGBA

	// Switch
	SwitchColors SwitchColors

	// Images
	ArrowDownArcImage paint.ImageOp
	ArrowUpArcImage   paint.ImageOp
	ManageFilesImage  paint.ImageOp
}

type SwitchColors struct {
	Enabled  color.NRGBA
	Disabled color.NRGBA
	Track    color.NRGBA
}

var (
	Current *Theme = Light
	Themes         = []*Theme{Light, Dark}
)

func Get(key string) *Theme {
	for _, theme := range Themes {
		if theme.Key == key {
			return theme
		}
	}

	return nil
}
