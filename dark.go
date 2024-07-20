package theme

import (
	"image/color"
)

var Dark = &Theme{
	Theme: th,

	Key:  "dark",
	Name: "Dark",

	// generic colors
	BackgroundColor: color.NRGBA{R: 20, G: 20, B: 20, A: 255},
	ForegroundColor: argb(0x99FFFFFF),
	PrimaryColor:    rgb(0x57B6FF),
	SurfaceColor:    rgb(0x252525),
	DeepBlueColor:   argb(0x99FFFFFF),

	// text colors
	TextColor:     argb(0x99FFFFFF),
	TextMuteColor: color.NRGBA{R: 255, G: 255, B: 255, A: 50},

	IndicatorColor:       color.NRGBA{A: 255},
	DividerColor:         color.NRGBA{R: 255, G: 255, B: 255, A: 25},
	BgColor:              BlackColor,
	BgGradientStartColor: color.NRGBA{R: 30, G: 30, B: 30, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	HideBalanceBgColor:   color.NRGBA{A: 255},

	ButtonColors: ButtonColors{
		TextColor:            WhiteColor,
		BackgroundColor:      BlueColor1,
		HoverBackgroundColor: nil,
		HoverTextColor:       nil,
		BorderColor:          GreyColor,
	},

	InputColors: InputColors{
		BackgroundColor: BlackColor,
		TextColor:       WhiteColor,
		BorderColor:     WhiteColor,
		HintColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	},

	ClickableColor:      rgb(0), // SurfaceHighlight,
	ClickableHoverColor: rgb(0x363636),

	CardColor:      rgb(0x252525),
	CardHoverColor: rgb(0x121212),

	SwitchActiveColor:       rgb(0x57B6FF),
	SwitchInactiveColor:     rgb(0x8997a5),
	SwitchThumbColor:        rgb(0xffffff),
	SwitchActiveTextColor:   WhiteColor,
	SwitchInactiveTextColor: GreyColor,

	ModalColors: ModalColors{
		BackgroundColor: BlackColor,
		BackdropColor:   &color.NRGBA{R: 20, G: 20, B: 20, A: 230},
	},
	ModalButtonColors: ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &WhiteColor,
	},

	NotificationSuccessColors: NotificationColors{
		BackgroundColor: color.NRGBA{R: 0, G: 200, B: 0, A: 255},
		TextColor:       WhiteColor,
	},
	NotificationErrorColors: NotificationColors{
		BackgroundColor: color.NRGBA{R: 200, G: 0, B: 0, A: 255},
		TextColor:       WhiteColor,
	},
	NotificationInfoColors: NotificationColors{
		BackgroundColor: WhiteColor,
		TextColor:       BlackColor,
	},

	ListTextColor:        WhiteColor,
	ListBgColor:          color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	ListItemHoverBgColor: color.NRGBA{R: 25, G: 25, B: 25, A: 255},
	ListScrollBarBgColor: WhiteColor,
	ListItemTagBgColor:   BlackColor,
	ListItemTagTextColor: WhiteColor,
}
