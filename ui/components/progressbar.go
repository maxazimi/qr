// https://github.com/crypto-power/cryptopower/blob/master/ui/cryptomaterial/progressbar.go

package components

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/qr/ui/theme"
	"github.com/maxazimi/qr/ui/values"
	"image"
	"image/color"
)

type LabelStyle struct {
	material.LabelStyle
}

type ProgressBarStyle struct {
	Radius    unit.Dp
	Height    unit.Dp
	Width     unit.Dp
	Direction layout.Direction
	material.ProgressBarStyle
}

type ProgressCircleStyle struct {
	material.ProgressCircleStyle
}

type ProgressBarItem struct {
	Value float64
	Color color.NRGBA
	Label LabelStyle
}

// MultiLayerProgressBar shows the percentage of the multiple progress layer
// against the total/expected progress.
type MultiLayerProgressBar struct {
	items                []ProgressBarItem
	Radius               unit.Dp
	Height               unit.Dp
	Width                unit.Dp
	total                float64
	ShowOverlayValue     bool
	ShowOtherWidgetFirst bool
}

func Label(txt string) LabelStyle {
	th := theme.Current().Theme
	return LabelStyle{LabelStyle: material.Label(th, values.TextSize14, txt)}
}

func ProgressBar(progress int) ProgressBarStyle {
	th := theme.Current().Theme
	return ProgressBarStyle{ProgressBarStyle: material.ProgressBar(th, float32(progress)/100)}
}

func NewMultiLayerProgressBar(total float64, items []ProgressBarItem) *MultiLayerProgressBar {
	return &MultiLayerProgressBar{
		total:  total,
		Height: values.DP8,
		items:  items,
	}
}

func (p ProgressBarStyle) Layout2(gtx C) D {
	if p.Width <= unit.Dp(0) {
		p.Width = unit.Dp(gtx.Constraints.Max.X)
	}

	return p.Direction.Layout(gtx, func(gtx C) D {
		return LinearLayout{
			Width:      gtx.Dp(p.Width),
			Height:     gtx.Dp(p.Height),
			Background: p.TrackColor,
			Border:     widget.Border{CornerRadius: p.Radius},
		}.Layout(gtx, layout.Rigid(func(gtx C) D {
			return LinearLayout{
				Width:      int(float32(p.Width) * clamp1(p.Progress)),
				Height:     gtx.Dp(p.Height),
				Background: p.Color,
				Border:     widget.Border{CornerRadius: p.Radius},
			}.Layout(gtx)
		}))
	})
}

func (p ProgressBarStyle) TextLayout(gtx C, lbl layout.Widget) D {
	if p.Width <= unit.Dp(0) {
		p.Width = unit.Dp(gtx.Constraints.Max.X)
	}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return p.Direction.Layout(gtx, func(gtx C) D {
				return LinearLayout{
					Width:      MatchParent,
					Height:     gtx.Dp(p.Height),
					Background: p.TrackColor,
					Border:     widget.Border{CornerRadius: p.Radius},
				}.Layout(gtx, layout.Rigid(func(gtx C) D {
					return LinearLayout{
						Width:      int(float32(p.Width) * clamp1(p.Progress)),
						Height:     gtx.Dp(p.Height),
						Background: p.Color,
						Border:     widget.Border{CornerRadius: p.Radius},
						Direction:  layout.Center,
					}.Layout(gtx)
				}))
			})
		}),
		layout.Expanded(func(gtx C) D {
			return layout.Center.Layout(gtx, lbl)
		}),
	)
}

func (p ProgressBarStyle) Layout(gtx C) D {
	shader := func(width int, color color.NRGBA) D {
		maxHeight := p.Height
		if p.Height <= 0 {
			maxHeight = unit.Dp(4)
		}

		d := image.Point{X: width, Y: gtx.Dp(maxHeight)}
		height := gtx.Dp(maxHeight)

		defer clip.RRect{
			Rect: image.Rectangle{Max: image.Pt(width, height)},
			NW:   gtx.Dp(p.Radius),
			NE:   gtx.Dp(p.Radius),
			SE:   gtx.Dp(p.Radius),
			SW:   gtx.Dp(p.Radius),
		}.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return D{Size: d}
	}

	if p.Width <= 0 {
		p.Width = unit.Dp(gtx.Constraints.Max.X)
	}

	progressBarWidth := int(p.Width)
	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		layout.Stacked(func(_ C) D {
			return shader(progressBarWidth, p.TrackColor)
		}),
		layout.Stacked(func(_ C) D {
			fillWidth := int(float32(progressBarWidth) * clamp1(p.Progress))
			fillColor := p.Color
			fillColor = Disabled(fillColor)
			return shader(fillWidth, fillColor)
		}),
	)
}

// TODO: Allow more than just 2 layers and make it dynamic
func (mp *MultiLayerProgressBar) progressBarLayout(gtx C) D {
	th := theme.Current()
	if mp.Width <= 0 {
		mp.Width = unit.Dp(gtx.Constraints.Max.X)
	}

	pg := func(width int, lbl LabelStyle, color color.NRGBA) D {
		return LinearLayout{
			Width:      width,
			Height:     gtx.Dp(mp.Height),
			Background: color,
		}.Layout(gtx, layout.Rigid(func(gtx C) D {
			if mp.ShowOverlayValue {
				lbl.Color = th.SurfaceColor
				return LinearLayout{
					Width:      width,
					Height:     gtx.Dp(mp.Height),
					Background: color,
					Direction:  layout.Center,
				}.Layout(gtx, layout.Rigid(lbl.Layout))
			}
			return D{}
		}))
	}

	calProgressWidth := func(progress float64) float64 {
		if mp.total != 0 {
			val := (progress / mp.total) * 100
			return (float64(mp.Width) / 100) * val
		}
		return 0
	}

	// display empty gray layout when total value passed is zero (0)
	if mp.total == 0 {
		return pg(int(mp.Width), Label(""), th.Gray2Color)
	}

	// This takes only 2 layers
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(_ C) D {
			width := calProgressWidth(mp.items[0].Value)
			if width == 0 {
				return D{}
			}
			return pg(int(width), mp.items[0].Label, mp.items[0].Color)
		}),
		layout.Rigid(func(_ C) D {
			width := calProgressWidth(mp.items[1].Value)
			if width == 0 {
				return D{}
			}
			return pg(int(width), mp.items[1].Label, mp.items[1].Color)
		}),
	)
}

func (mp *MultiLayerProgressBar) Layout(gtx C, isMobileView bool, additionalWidget layout.Widget) D {
	if additionalWidget == nil {
		// We're only displaying the progress bar, no need for flex layout to wrap it.
		// TODO: Verify if a top padding is necessary if we're only displaying the progressbar.
		return layout.Inset{Top: values.DP5}.Layout(gtx, mp.progressBarLayout)
	}

	progressBarTopPadding, otherWidget := values.DP24, additionalWidget
	if isMobileView {
		progressBarTopPadding = values.DP16
	}
	if !mp.ShowOtherWidgetFirst {
		// reduce the top padding if we're showing the progress bar before the other widget
		progressBarTopPadding = values.DP5
		otherWidget = func(gtx C) D {
			return layout.Center.Layout(gtx, additionalWidget)
		}
	}

	flexWidgets := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: progressBarTopPadding}.Layout(gtx, mp.progressBarLayout)
		}),
		layout.Rigid(otherWidget),
	}

	if mp.ShowOtherWidgetFirst {
		// Swap the label and progress bar...
		flexWidgets[0], flexWidgets[1] = flexWidgets[1], flexWidgets[0]
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, flexWidgets...)
}

func ProgressBarCircle(progress int) ProgressCircleStyle {
	th := theme.Current().Theme
	return ProgressCircleStyle{ProgressCircleStyle: material.ProgressCircle(th, float32(progress)/100)}
}

func (p ProgressCircleStyle) Layout(gtx C) D {
	return p.ProgressCircleStyle.Layout(gtx)
}

// clamp1 limits mp to range [0..1]
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	}
	return v
}