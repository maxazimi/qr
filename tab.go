package components

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/ui/lang"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"github.com/maxazimi/v2ray-gio/ui/values"
	"image"
)

type tabItem struct {
	*Clickable
	name string
}

type Tab struct {
	list          layout.List
	tabItems      []tabItem
	selectedIndex int
	changed       bool
}

func NewTab(axis layout.Axis, items []string) *Tab {
	t := &Tab{}
	t.list.Axis = axis
	for _, item := range items {
		t.tabItems = append(t.tabItems, tabItem{
			Clickable: NewClickable(false),
			name:      item,
		})
	}
	return t
}

func (t *Tab) Layout(gtx C) D {
	t.handleEvents(gtx)
	th := theme.Current()
	var selectedTabDims D

	return t.list.Layout(gtx, len(t.tabItems), func(gtx C, i int) D {
		isSelectedTab := t.SelectedIndex() == i
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Right:  values.DP24,
					Left:   values.DP24,
					Bottom: values.DP8,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						item := t.tabItems[i]
						return item.Layout(gtx, func(gtx C) D {
							lbl := material.Label(th.Theme, values.TextSize16, lang.Str(item.name))
							lbl.Color = th.GrayText1Color

							if isSelectedTab {
								lbl.Color = th.PrimaryColor
								selectedTabDims = lbl.Layout(gtx)
							}
							return lbl.Layout(gtx)
						})
					})
				})
			}),
			layout.Stacked(func(gtx C) D {
				if !isSelectedTab {
					return D{}
				}

				tabHeight := gtx.Dp(values.DP4)
				selectedTabDimsWidth := gtx.Dp(values.DP50)
				tabRect := image.Rect(0, 0, selectedTabDims.Size.X+selectedTabDimsWidth, tabHeight)
				defer clip.RRect{Rect: tabRect, SE: 0, SW: 0, NW: 10, NE: 10}.Push(gtx.Ops).Pop()

				return layout.Inset{Bottom: values.DPMinus24}.Layout(gtx, func(gtx C) D {
					paint.FillShape(gtx.Ops, th.PrimaryColor, clip.Rect(tabRect).Op())
					return D{Size: image.Point{X: selectedTabDims.Size.X + selectedTabDimsWidth, Y: tabHeight}}
				})
			}),
		)
	})
}

func (t *Tab) handleEvents(gtx C) {
	for i, item := range t.tabItems {
		if item.Clicked(gtx) {
			if t.selectedIndex != i {
				t.changed = true
			}
			t.selectedIndex = i
		}
	}
}

func (t *Tab) SelectedIndex() int {
	return t.selectedIndex
}

func (t *Tab) SelectedTab() string {
	return t.tabItems[t.selectedIndex].name
}

func (t *Tab) Changed() bool {
	changed := t.changed
	t.changed = false
	return changed
}

func (t *Tab) SetSelectedTab(tab string) {
	for i, item := range t.tabItems {
		if item.name == tab {
			t.selectedIndex = i
		}
	}
}

func (t *Tab) SetSelectedIndex(index int) {
	t.selectedIndex = index
}
