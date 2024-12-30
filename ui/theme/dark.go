package theme

import (
	"image/color"
)

var Dark = &Theme{
	Theme: th,
	Name:  "dark",

	// generic colors
	BackgroundColor: color.NRGBA{R: 20, G: 20, B: 20, A: 255},
	ForegroundColor: argb(0x99FFFFFF),
	GreenColor:      color.NRGBA{R: 0, G: 200, B: 0, A: 255},
	RedColor:        color.NRGBA{R: 200, G: 0, B: 0, A: 255},

	PrimaryColor:   DarkLightPrimary,
	DeepBlueColor:  DarkDeepBlue,
	Gray1Color:     DarkGray1,
	Gray2Color:     DarkGray2,
	Gray3Color:     DarkGray3,
	Gray4Color:     DarkGray4,
	Gray5Color:     DarkGray5,
	SurfaceColor:   DarkSurface,
	LightGrayColor: DarkLightGray,

	// text colors
	TextColor:        DarkText,
	TextMuteColor:    color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	PageNavTextColor: DarkPageNavText,
	GrayText1Color:   DarkGrayText1,
	GrayText2Color:   DarkGrayText2,
	GrayText3Color:   DarkGrayText3,
	GrayText4Color:   DarkGrayText4,

	ModalColors: ModalColors{
		BackgroundColor:  BlueGreyDarken4,
		BackgroundColor2: BlackColor,
		TextColor:        WhiteColor,
		BackdropColor:    &color.NRGBA{R: 20, G: 20, B: 20, A: 230},
	},
}
