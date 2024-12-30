package theme

import (
	"image/color"
)

var Light = &Theme{
	Theme: th,
	Name:  "light",

	// generic colors
	BackgroundColor: LightGray9,
	ForegroundColor: rgb(0x091440),
	GreenColor:      color.NRGBA{R: 0, G: 225, B: 0, A: 255},
	RedColor:        color.NRGBA{R: 225, G: 0, B: 0, A: 255},

	PrimaryColor:   LightPrimary,
	DeepBlueColor:  LightDeepBlue,
	Gray1Color:     LightGray1,
	Gray2Color:     LightGray2,
	Gray3Color:     LightGray3,
	Gray4Color:     LightGray4,
	Gray5Color:     LightGray5,
	SurfaceColor:   LightSurface,
	LightGrayColor: LightLightGray,

	// text colors
	TextColor:        LightText,
	TextMuteColor:    color.NRGBA{A: 200},
	PageNavTextColor: LightPageNavText,
	GrayText1Color:   LightGrayText1,
	GrayText2Color:   LightGrayText2,
	GrayText3Color:   LightGrayText3,
	GrayText4Color:   LightGrayText4,

	ModalColors: ModalColors{
		BackgroundColor:  BlueGreyLighten5,
		BackgroundColor2: BlueGreyLighten4,
		TextColor:        BlackColor,
		BackdropColor:    &color.NRGBA{A: 100},
	},
}
