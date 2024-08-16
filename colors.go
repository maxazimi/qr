package theme

import "image/color"

var (
	WhiteColor     = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	BlackColor     = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	GreyColor      = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	BlueColor      = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	BlueColor1     = rgb(0x3377FF)
	BlueGreyColor  = rgb(0x496495)
	BlueGreyColor1 = rgb(0x3a517b)
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
