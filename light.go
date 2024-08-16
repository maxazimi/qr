package theme

import (
	"image/color"
)

var Light = &Theme{
	Theme: th,

	Key:  "light",
	Name: "Light",

	// generic colors
	BackgroundColor: rgb(0xE6EAED),
	ForegroundColor: rgb(0x091440),
	PrimaryColor:    rgb(0x2970ff),
	SurfaceColor:    rgb(0xffffff),
	DeepBlueColor:   rgb(0x091440),
	GreenColor:      color.NRGBA{R: 0, G: 225, B: 0, A: 255},
	RedColor:        color.NRGBA{R: 225, G: 0, B: 0, A: 255},

	// text colors
	TextColor:     rgb(0x091440),
	TextMuteColor: color.NRGBA{A: 200},

	IndicatorColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	DividerColor:         color.NRGBA{A: 50},
	BgColor:              WhiteColor,
	BgGradientStartColor: color.NRGBA{R: 250, G: 250, B: 250, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 210, G: 210, B: 210, A: 255},
	HideBalanceBgColor:   color.NRGBA{R: 200, G: 200, B: 200, A: 255},

	ButtonColors: ButtonColors{
		TextColor:            WhiteColor,
		BackgroundColor:      BlueColor1,
		HoverBackgroundColor: nil,
		HoverTextColor:       nil,
		BorderColor:          WhiteColor,
	},

	InputColors: InputColors{
		BackgroundColor: GreyColor,
		TextColor:       BlackColor,
		BorderColor:     BlackColor,
		HintColor:       color.NRGBA{A: 200},
	},

	ClickableColor:      rgb(0), // SurfaceHighlight,
	ClickableHoverColor: rgb(0xf3f5f6),

	CardColor:      rgb(0xffffff),
	CardHoverColor: rgb(0xf3f5f6),

	SwitchActiveColor:       rgb(0x2970ff),
	SwitchInactiveColor:     rgb(0xc4cbd2), // InactiveGray #C4CBD2
	SwitchThumbColor:        rgb(0xffffff),
	SwitchActiveTextColor:   BlackColor,
	SwitchInactiveTextColor: GreyColor,

	ModalColors: ModalColors{
		BackgroundColor:  BlueGreyLighten5,
		BackgroundColor2: BlueGreyLighten2,
		TextColor:        BlackColor,
		BackdropColor:    &color.NRGBA{A: 100},
	},
	ModalButtonColors: ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &BlackColor,
	},

	AppBarColors: ButtonColors{
		TextColor:       WhiteColor,
		BackgroundColor: BlueGreyDarken4,
		BorderColor:     BlueGrey,
	},

	ListTextColor:        BlackColor,
	ListBgColor:          WhiteColor,
	ListItemHoverBgColor: color.NRGBA{R: 225, G: 225, B: 225, A: 255},
	ListScrollBarBgColor: BlackColor,
	ListItemTagBgColor:   color.NRGBA{R: 225, G: 225, B: 225, A: 255},
	ListItemTagTextColor: BlackColor,
}
