package theme

import "image/color"

var (
	WhiteColor        = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	BlurredWhiteColor = color.NRGBA{R: 250, G: 250, B: 250, A: 70}
	BlackColor        = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	GreyColor         = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	BlueColor         = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	BlueColor1        = rgb(0x3377FF)
	BlueGreyColor     = rgb(0x496495)
	BlueGreyColor1    = rgb(0x3a517b)
	SuccessColor      = rgb(0x41bf53)
	YellowColor       = rgb(0xFEE3AA)
	DangerColor       = rgb(0xed6d47)
)

// Dark Colors
var (
	DarkLightPrimary = rgb(0x57B6FF)
	DarkPageNavText  = argb(0x99FFFFFF)
	DarkText         = argb(0x99FFFFFF)
	DarkGrayText1    = argb(0xDEFFFFFF)
	DarkGrayText2    = argb(0x99FFFFFF)
	DarkGrayText3    = argb(0x61FFFFFF)
	DarkGrayText4    = argb(0x61FFFFFF)
	DarkDeepBlue     = argb(0x99FFFFFF)
	DarkGray1        = argb(0x99FFFFFF)
	DarkGray2        = rgb(0x3D3D3D)
	DarkGray3        = rgb(0x8997a5)
	DarkGray4        = rgb(0x121212)
	DarkGray5        = rgb(0x363636)
	DarkSurface      = rgb(0x252525)
	DarkLightGray    = rgb(0x2B2B2B)
)

// Light Colors
var (
	LightPrimary          = rgb(0x2970ff)
	LightPrimaryHighlight = rgb(0x1B41B3)
	LightPageNavText      = rgb(0x091440)
	LightText             = rgb(0x091440)
	LightInvText          = rgb(0xffffff)
	LightGrayText1        = rgb(0x3d5873)
	LightGrayText2        = rgb(0x596D81)
	LightGrayText3        = rgb(0x8997a5)
	LightGrayText4        = rgb(0xc4cbd2)
	LightGreenText        = rgb(0x41BE53)
	LightBackground       = argb(0x22444444)
	LightBlack            = rgb(0x000000)
	LightBlueProgressTint = rgb(0x73d7ff)
	LightDanger           = rgb(0xed6d47)
	LightDeepBlue         = rgb(0x091440)
	LightDeepBlueOrigin   = rgb(0x091440)
	LightNavyBlue         = rgb(0x1F45B0)
	LightLightBlue        = rgb(0xe4f6ff)
	LightLightBlue2       = rgb(0x75D8FF)
	LightLightBlue3       = rgb(0xBCE8FF)
	LightLightBlue4       = rgb(0xBBDEFF)
	LightLightBlue5       = rgb(0x70CBFF)
	LightLightBlue6       = rgb(0x4B91D8)
	LightLightBlue7       = rgb(0xF0F3FF)
	LightLightBlue8       = rgb(0xE3ECFC)
	LightGray1            = rgb(0x3d5873) // darkest gray #3D5873 (icon color),
	LightGray2            = rgb(0xe6eaed) // light 0xe6eaed
	LightGray3            = rgb(0xc4cbd2) // InactiveGray #C4CBD2
	LightGray4            = rgb(0xf3f5f6) // active n light gray combined f3f5f6
	LightGray5            = rgb(0xf3f5f6)
	LightGray6            = rgb(0xD8D8D8)
	LightGray7            = rgb(0x8997a5)
	LightGray8            = rgb(0xEDEFF1)
	LightGray9            = rgb(0xE6EAED)
	LightGray10           = rgb(0xC8C8C8)
	LightLightGray        = rgb(0xFBFCFC)
	LightGreen50          = rgb(0xE8F7EA)
	LightGreen500         = rgb(0x41BE53)
	LightOrange           = rgb(0xD34A21)
	LightOrange2          = rgb(0xF8E8E7)
	LightOrange3          = rgb(0xF8CABC)
	LightOrangeRipple     = rgb(0xD32F2F)
	LightSuccess          = rgb(0x41bf53)
	LightSuccess2         = rgb(0xE1F8EF)
	LightSurface          = rgb(0xffffff)
	LightTurquoise100     = rgb(0xB6EED7)
	LightTurquoise300     = rgb(0x2DD8A3)
	LightTurquoise700     = rgb(0x00A05F)
	LightTurquoise800     = rgb(0x008F52)
	LightYellow           = rgb(0xFEE3AA)
	LightOrangeYellow     = rgb(0xd8a93e)
	LightWhite            = rgb(0xffffff)
	LightWarning          = rgb(0xff9966)
)

// https://colorswall.com/palette/75024
var (
	BlueGreyDarken4  = rgb(0x263238)
	BlueGreyDarken3  = rgb(0x37474f)
	BlueGreyDarken2  = rgb(0x455a64)
	BlueGreyDarken1  = rgb(0x546e7a)
	BlueGrey         = rgb(0x607d8b)
	BlueGreyLighten1 = rgb(0x78909c)
	BlueGreyLighten2 = rgb(0x90a4ae)
	BlueGreyLighten3 = rgb(0xb0bec5) // Ocean Drive
	BlueGreyLighten4 = rgb(0xcfd8dc)
	BlueGreyLighten5 = rgb(0xeceff1)
)
