package theme

import (
	"github.com/maxazimi/v2ray-gio/ui/components"
	"image/color"
)

var Dark = &Theme{
	Key:            "dark",
	Name:           "Dark",
	IndicatorColor: color.NRGBA{A: 255},

	TextColor:            WhiteColor,
	TextMuteColor:        color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	DividerColor:         color.NRGBA{R: 255, G: 255, B: 255, A: 25},
	BgColor:              BlackColor,
	BgGradientStartColor: color.NRGBA{R: 30, G: 30, B: 30, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	HideBalanceBgColor:   color.NRGBA{A: 255},

	HeaderBackButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &WhiteColor,
	},
	HeaderTopBgColor: color.NRGBA{R: 30, G: 30, B: 30, A: 255},

	BottomBarBgColor: BlackColor,
	BottomButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &WhiteColor,
	},
	BottomButtonSelectedColor: WhiteColor,

	NodeStatusBgColor:        color.NRGBA{A: 255},
	NodeStatusTextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	NodeStatusDotGreenColor:  color.NRGBA{R: 0, G: 200, B: 0, A: 255},
	NodeStatusDotYellowColor: color.NRGBA{R: 255, G: 255, B: 0, A: 255},
	NodeStatusDotRedColor:    color.NRGBA{R: 200, G: 0, B: 0, A: 255},

	InputColors: components.InputColors{
		BackgroundColor: BlackColor,
		TextColor:       WhiteColor,
		BorderColor:     WhiteColor,
		HintColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	},

	ButtonIconPrimaryColors: components.ButtonColors{
		TextColor: WhiteColor,
	},
	ButtonPrimaryColors: components.ButtonColors{
		TextColor:       BlackColor,
		BackgroundColor: WhiteColor,
	},
	ButtonSecondaryColors: components.ButtonColors{
		TextColor:   WhiteColor,
		BorderColor: WhiteColor,
	},
	ButtonInvertColors: components.ButtonColors{
		TextColor:       WhiteColor,
		BackgroundColor: BlackColor,
	},
	ButtonDangerColors: components.ButtonColors{
		TextColor:       WhiteColor,
		BackgroundColor: color.NRGBA{R: 200, G: 0, B: 0, A: 255},
	},

	ModalColors: components.ModalColors{
		BackgroundColor: BlackColor,
		BackdropColor:   &color.NRGBA{R: 20, G: 20, B: 20, A: 230},
	},
	ModalButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &WhiteColor,
	},

	NotificationSuccessColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 0, G: 200, B: 0, A: 255},
		TextColor:       WhiteColor,
	},
	NotificationErrorColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 200, G: 0, B: 0, A: 255},
		TextColor:       WhiteColor,
	},
	NotificationInfoColors: components.NotificationColors{
		BackgroundColor: WhiteColor,
		TextColor:       BlackColor,
	},

	ListTextColor:        WhiteColor,
	ListBgColor:          color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	ListItemHoverBgColor: color.NRGBA{R: 25, G: 25, B: 25, A: 255},
	ListScrollBarBgColor: WhiteColor,
	ListItemTagBgColor:   BlackColor,
	ListItemTagTextColor: WhiteColor,

	SwitchColors: SwitchColors{
		Enabled:  WhiteColor,
		Disabled: color.NRGBA{R: 60, G: 60, B: 60, A: 255},
		Track:    color.NRGBA{R: 60, G: 60, B: 60, A: 255},
	},
}
