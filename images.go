package theme

import (
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/assets"
)

const (
	ImageScale = 0.3
)

var images = []*widget.Image{
	nil,
	{Src: paint.NewImageOp(assets.AppIcons["ic_settings"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_history"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_home"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_info_circle"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_close_round"]), Scale: ImageScale},

	{Src: paint.NewImageOp(assets.AppIcons["ic_about"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["chevron_expand"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_copy"]), Scale: ImageScale * 1.3},
	{Src: paint.NewImageOp(assets.AppIcons["ic_delete"]), Scale: ImageScale},

	{Src: paint.NewImageOp(assets.AppIcons["ic_new_item"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_import_item"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_qrcode_scan"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_ascending_filter"]), Scale: ImageScale},
	{Src: paint.NewImageOp(assets.AppIcons["ic_descending_filter"]), Scale: ImageScale},
}

func init() {
	loadDarkModeImages()
}

// Image enums
const (
	Nil = iota
	ImageSettings
	ImageHistory
	ImageHome
	ImageInfo
	ImageClose

	ImageAbout
	ImageChevronExpand
	ImageCopy
	ImageDelete

	ImageAddItem
	ImageImportItem
	ImageScanQR
	ImageAscendingFilter
	ImageDescendingFilter

	ImageEnd
)
