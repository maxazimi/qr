package theme

import (
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/assets"
)

const (
	ImageScale = 0.3
)

var (
	// index0 -> light
	// index1 -> dark
	imageMap = make(map[string][]*widget.Image)
)

func init() {
	imageMap[ICSettings] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_settings"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_settings_white"]), Scale: ImageScale},
	}
	imageMap[ICHistory] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_history"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_history_white"]), Scale: ImageScale},
	}
	imageMap[ICHome] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_home"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_home_white"]), Scale: ImageScale},
	}
	imageMap[ICInfo] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_info_circle"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_info_circle_white"]), Scale: ImageScale},
	}
	imageMap[ICClose] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_close_round"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_close_round_white"]), Scale: ImageScale},
	}
	imageMap[ICChevronExpand] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["chevron_expand"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["chevron_expand"]), Scale: ImageScale},
	}
	imageMap[ICChevronLeft] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["chevron_left"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["chevron_left"]), Scale: ImageScale},
	}
	imageMap[ICChevronRight] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["chevron_coll"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["chevron_coll"]), Scale: ImageScale},
	}
	imageMap[ICAbout] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_about"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_about"]), Scale: ImageScale},
	}
	imageMap[ICEdit] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_edit"]), Scale: 0.4},
		{Src: paint.NewImageOp(assets.AppIcons["ic_edit_white"]), Scale: 0.4},
	}
	imageMap[ICCopied] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_copied"]), Scale: 0.4},
		{Src: paint.NewImageOp(assets.AppIcons["ic_copied_white"]), Scale: 0.4},
	}
	imageMap[ICDelete] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_delete"]), Scale: 0.4},
		{Src: paint.NewImageOp(assets.AppIcons["ic_delete_white"]), Scale: 0.4},
	}
	imageMap[ICShare] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_share"]), Scale: 0.4},
		{Src: paint.NewImageOp(assets.AppIcons["ic_share_white"]), Scale: 0.4},
	}
	imageMap[ICAddItem] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_new_item"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_new_item_white"]), Scale: ImageScale},
	}
	imageMap[ICImportItem] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_import_item"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_import_item_white"]), Scale: ImageScale},
	}
	imageMap[ICScanQR] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_qrcode_scan"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_qrcode_scan_white"]), Scale: ImageScale},
	}
	imageMap[ICAscendingFilter] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_ascending_filter"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_ascending_filter_white"]), Scale: ImageScale},
	}
	imageMap[ICDescendingFilter] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_descending_filter"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_descending_filter_white"]), Scale: ImageScale},
	}
	imageMap[ICLoading] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_loading"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_loading_white"]), Scale: ImageScale},
	}
	imageMap[ICSync] = []*widget.Image{
		{Src: paint.NewImageOp(assets.AppIcons["ic_sync"]), Scale: ImageScale},
		{Src: paint.NewImageOp(assets.AppIcons["ic_sync_white"]), Scale: ImageScale},
	}
}

const (
	ICSettings         = "ic_settings"
	ICHistory          = "ic_history"
	ICHome             = "ic_home"
	ICInfo             = "ic_info_circle"
	ICClose            = "ic_close_round"
	ICChevronExpand    = "chevron_expand"
	ICChevronLeft      = "chevron_left"
	ICChevronRight     = "chevron_coll"
	ICAbout            = "ic_about"
	ICEdit             = "ic_edit"
	ICCopied           = "ic_copied"
	ICDelete           = "ic_delete"
	ICShare            = "ic_share"
	ICAddItem          = "ic_new_item"
	ICImportItem       = "ic_import_item"
	ICScanQR           = "ic_qrcode_scan"
	ICAscendingFilter  = "ic_ascending_filter"
	ICDescendingFilter = "ic_descending_filter"
	ICLoading          = "ic_loading"
	ICSync             = "ic_sync"
)
