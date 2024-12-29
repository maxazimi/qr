package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/maxazimi/qr/ui/lang"
	"github.com/maxazimi/qr/ui/theme"
	"github.com/maxazimi/qr/ui/values"
	"image"
)

type TabItem struct {
	*Button
	Name string
}

type Tab struct {
	TabItems      []TabItem
	list          layout.List
	selectedIndex int
	changed       bool
}

func NewTab(items []TabItem) *Tab {
	t := &Tab{}
	for _, item := range items {
		t.TabItems = append(t.TabItems, item)
	}
	return t
}

func (t *Tab) Layout(gtx C) D {
	t.handleEvents(gtx)
	th := theme.Current()
	var selectedTabDims D

	return t.list.Layout(gtx, len(t.TabItems), func(gtx C, i int) D {
		isSelectedTab := t.SelectedIndex() == i
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Right:  values.DP24,
					Left:   values.DP24,
					Bottom: values.DP8,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						item := t.TabItems[i]
						item.Text = lang.Str(item.Name)
						item.TextSize = values.TextSize16
						if isSelectedTab {
							item.Colors.TextColor = th.PrimaryColor
						} else {
							item.Colors.TextColor = th.GrayText1Color
						}

						c := op.Record(gtx.Ops)
						dims := item.Layout(gtx)
						m := c.Stop()

						if isSelectedTab {
							selectedTabDims = dims
						}

						m.Add(gtx.Ops)
						return dims
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
	for i, item := range t.TabItems {
		if item.Clicked(gtx) {
			if t.selectedIndex != i {
				t.changed = true
				t.TabItems[t.selectedIndex].Focused = false
				item.Focused = true
			}
			t.selectedIndex = i
		}
	}
}

func (t *Tab) SelectedIndex() int {
	return t.selectedIndex
}

func (t *Tab) SelectedTab() string {
	return t.TabItems[t.selectedIndex].Name
}

func (t *Tab) Changed() bool {
	changed := t.changed
	t.changed = false
	return changed
}

func (t *Tab) SetSelectedTab(tab string) {
	for i, item := range t.TabItems {
		if item.Name == tab {
			t.selectedIndex = i
		}
	}
}

func (t *Tab) SetSelectedIndex(index int) {
	if index < len(t.TabItems) {
		t.selectedIndex = index
	}
}
