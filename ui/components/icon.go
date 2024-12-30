package components

import (
	"gioui.org/unit"
	"gioui.org/widget"
	"image/color"
	"log"
)

type Icon struct {
	*widget.Icon
	Color color.NRGBA
}

func NewIcon(data []byte) *Icon {
	return &Icon{
		Icon:  MustIcon(widget.NewIcon(data)),
		Color: color.NRGBA{A: 0xff},
	}
}

func (icon *Icon) LayoutSize(gtx C, iconSize unit.Dp) D {
	gtx.Constraints.Min.X = gtx.Dp(iconSize)
	return icon.Icon.Layout(gtx, icon.Color)
}

func (icon *Icon) Layout(gtx C, col color.NRGBA) D {
	return icon.Icon.Layout(gtx, col)
}

func MustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		log.Fatal(err)
	}
	return ic
}
