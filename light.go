package theme

import (
	"github.com/maxazimi/v2ray-gio/ui/components"
	"image/color"
)

var Light = &Theme{
	Key:            "light",
	Name:           "Light",
	IndicatorColor: color.NRGBA{R: 255, G: 255, B: 255, A: 50},

	TextColor:            BlackColor,
	TextMuteColor:        color.NRGBA{A: 200},
	DividerColor:         color.NRGBA{A: 50},
	BgColor:              WhiteColor,
	BgGradientStartColor: color.NRGBA{R: 250, G: 250, B: 250, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 210, G: 210, B: 210, A: 255},
	HideBalanceBgColor:   color.NRGBA{R: 200, G: 200, B: 200, A: 255},

	HeaderBackButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &color.NRGBA{A: 255},
	},
	HeaderTopBgColor: color.NRGBA{R: 250, G: 250, B: 250, A: 255},

	BottomBarBgColor: WhiteColor,
	BottomButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &BlackColor,
	},
	BottomButtonSelectedColor: BlackColor,

	NodeStatusBgColor:        color.NRGBA{A: 255},
	NodeStatusTextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	NodeStatusDotGreenColor:  color.NRGBA{R: 0, G: 225, B: 0, A: 255},
	NodeStatusDotYellowColor: color.NRGBA{R: 255, G: 255, B: 0, A: 255},
	NodeStatusDotRedColor:    color.NRGBA{R: 225, G: 0, B: 0, A: 255},

	InputColors: components.InputColors{
		BackgroundColor: GreyColor,
		TextColor:       BlackColor,
		BorderColor:     BlackColor,
		HintColor:       color.NRGBA{A: 200},
	},

	ButtonIconPrimaryColors: components.ButtonColors{
		TextColor: BlackColor,
	},
	ButtonPrimaryColors: components.ButtonColors{
		TextColor:       WhiteColor,
		BackgroundColor: BlackColor,
	},
	ButtonSecondaryColors: components.ButtonColors{
		TextColor:   BlackColor,
		BorderColor: BlackColor,
	},
	ButtonInvertColors: components.ButtonColors{
		TextColor:       BlackColor,
		BackgroundColor: WhiteColor,
	},
	ButtonDangerColors: components.ButtonColors{
		TextColor:       WhiteColor,
		BackgroundColor: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	},

	ModalColors: components.ModalColors{
		BackgroundColor: WhiteColor,
		BackdropColor:   &color.NRGBA{A: 100},
	},
	ModalButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &BlackColor,
	},

	NotificationSuccessColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 0, G: 225, B: 0, A: 255},
		TextColor:       WhiteColor,
	},
	NotificationErrorColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 225, G: 0, B: 0, A: 255},
		TextColor:       WhiteColor,
	},
	NotificationInfoColors: components.NotificationColors{
		BackgroundColor: WhiteColor,
		TextColor:       BlackColor,
	},

	ListTextColor:        BlackColor,
	ListBgColor:          WhiteColor,
	ListItemHoverBgColor: color.NRGBA{R: 225, G: 225, B: 225, A: 255},
	ListScrollBarBgColor: BlackColor,
	ListItemTagBgColor:   color.NRGBA{R: 225, G: 225, B: 225, A: 255},
	ListItemTagTextColor: BlackColor,

	SwitchColors: SwitchColors{
		Enabled:  BlackColor,
		Disabled: WhiteColor,
		Track:    color.NRGBA{A: 100},
	},
}
