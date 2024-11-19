package theme

import (
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/assets"
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

	PrimaryColor:     DarkLightPrimary,
	PageNavTextColor: DarkPageNavText,
	TextColor:        DarkText,
	GrayText1Color:   DarkGrayText1,
	GrayText2Color:   DarkGrayText2,
	GrayText3Color:   DarkGrayText3,
	GrayText4Color:   DarkGrayText4,
	DeepBlueColor:    DarkDeepBlue,
	Gray1Color:       DarkGray1,
	Gray2Color:       DarkGray2,
	Gray3Color:       DarkGray3,
	Gray4Color:       DarkGray4,
	Gray5Color:       DarkGray5,
	SurfaceColor:     DarkSurface,
	LightGrayColor:   DarkLightGray,

	// text colors
	//TextColor:     argb(0x99FFFFFF),
	TextMuteColor: color.NRGBA{R: 255, G: 255, B: 255, A: 50},

	//IndicatorColor:       color.NRGBA{A: 255},
	//DividerColor:         color.NRGBA{R: 255, G: 255, B: 255, A: 25},
	//BgColor:              BlackColor,
	//BgGradientStartColor: color.NRGBA{R: 30, G: 30, B: 30, A: 255},
	//BgGradientEndColor:   color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	//HideBalanceBgColor:   color.NRGBA{A: 255},

	ButtonColors: ButtonColors{
		TextColor:            WhiteColor,
		BackgroundColor:      BlueColor1,
		HoverBackgroundColor: nil,
		HoverTextColor:       nil,
		BorderColor:          WhiteColor,
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

	SwitchColors: SwitchColors{
		ActiveColor:       rgb(0x57B6FF),
		InactiveColor:     rgb(0x8997a5),
		ThumbColor:        rgb(0xffffff),
		ActiveTextColor:   WhiteColor,
		InactiveTextColor: GreyColor,
	},

	ModalColors: ModalColors{
		BackgroundColor:  BlueGreyDarken4,
		BackgroundColor2: BlackColor,
		TextColor:        WhiteColor,
		BackdropColor:    &color.NRGBA{R: 20, G: 20, B: 20, A: 230},
	},
	ModalButtonColors: ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &WhiteColor,
	},

	AppBarColors: ButtonColors{
		TextColor: WhiteColor,
		//BackgroundColor: BlueGreyDarken4,
		BorderColor: BlueGrey,
	},

	BottomBarBgColor: BlackColor,
	BottomButtonColors: ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &WhiteColor,
	},
	BottomButtonSelectedColor: WhiteColor,
	BottomShadowColor:         BlackColor,

	ListTextColor:        WhiteColor,
	ListBgColor:          color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	ListItemBgColor:      DarkGray3,
	ListItemHoverBgColor: color.NRGBA{R: 25, G: 25, B: 25, A: 255},
	ListItemTagBgColor:   BlackColor,
	ListItemTagTextColor: WhiteColor,

	// Images
	Images: nil,
}

func loadDarkModeImages() {
	Dark.Images = make([]*widget.Image, ImageEnd)
	for i, img := range images {
		Dark.Images[i] = img
	}

	Dark.Images[ImageSettings] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_settings_white"]), Scale: ImageScale}
	Dark.Images[ImageHistory] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_history_white"]), Scale: ImageScale}
	Dark.Images[ImageHome] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_home_white"]), Scale: ImageScale}
	Dark.Images[ImageInfo] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_info_circle_white"]), Scale: ImageScale}
	Dark.Images[ImageClose] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_close_round_white"]), Scale: ImageScale}

	Dark.Images[ImageAddItem] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_new_item_white"]), Scale: ImageScale}
	Dark.Images[ImageImportItem] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_import_item_white"]), Scale: ImageScale}
	Dark.Images[ImageScanQR] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_qrcode_scan_white"]), Scale: ImageScale}
	Dark.Images[ImageAscendingFilter] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_ascending_filter_white"]), Scale: ImageScale}
	Dark.Images[ImageDescendingFilter] = &widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_descending_filter_white"]), Scale: ImageScale}
}
