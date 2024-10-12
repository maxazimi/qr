package components

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/assets"
	"github.com/maxazimi/v2ray-gio/ui/theme"
)

type NavItem struct {
	Tag      interface{}
	Name     string
	ImgIndex int
	Icon     *widget.Icon
	button   *Button
	Selected bool
}

func (n *NavItem) Clicked(gtx C) bool {
	return n.button.Clicked(gtx)
}

func (n *NavItem) Layout(gtx C) D {
	return n.button.Layout(gtx)
}

type NavDrawer struct {
	*Modal
	Title           string
	Subtitle        string
	navList         layout.List
	selectedItem    int
	selectedChanged bool
	navigationAxis  layout.Axis
	items           []NavItem
}

func NewNavDrawer(title, subtitle string, bottombar bool) *NavDrawer {
	if bottombar {
		return &NavDrawer{
			Modal: NewModal(layout.N, layout.Inset{Top: 25}, layout.UniformInset(10), 0,
				NewModalAnimationUp()),
		}
	}
	return &NavDrawer{
		Title:    title,
		Subtitle: subtitle,
		Modal: NewModal(layout.E, layout.UniformInset(0), layout.UniformInset(10), 5,
			NewModalAnimationRight()),
		navigationAxis: layout.Vertical,
	}
}

func (n *NavDrawer) AddNavItem(navItem NavItem) {
	defer func() {
		n.items = append(n.items, navItem)
		if len(n.items) == 1 {
			n.items[0].Selected = true
		}
	}()

	if n.navigationAxis == layout.Horizontal {
		colors := theme.Current().AppBarColors.Reverse()
		colors.BackgroundColor.A = 200
		navItem.button = NewButton(ButtonStyle{
			Icon:   navItem.Icon,
			Colors: colors,
			Radius: 10,
		})
		return
	}
	navItem.button = NewButton(ButtonStyle{
		Icon:     navItem.Icon,
		Text:     navItem.Name,
		Colors:   theme.Current().AppBarColors.Reverse(),
		Radius:   5,
		TextSize: 14,
		Inset:    layout.Inset{Top: 10, Right: 16, Bottom: 10, Left: 16},
	})
}

func (n *NavDrawer) Layout(gtx C) D {
	if n.navigationAxis == layout.Vertical {
		return n.LayoutH(gtx)
	}

	n.Modal.OuterInset.Top = unit.Dp(gtx.Constraints.Max.Y * 4 / 5)
	return n.Modal.Layout(gtx, func(gtx C) D {
		layout.Center.Layout(gtx, func(gtx C) D {
			return n.layoutNavList(gtx)
		})
		return D{Size: gtx.Constraints.Max}
	})
}

func (n *NavDrawer) LayoutH(gtx C) D {
	n.Modal.OuterInset.Right = unit.Dp(gtx.Constraints.Max.X / 3)
	return n.Modal.Layout(gtx, func(gtx C) D {
		layout.Flex{
			Spacing: layout.SpaceEnd,
			Axis:    layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: 18, Right: 18, Bottom: 18}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							img := widget.Image{Src: paint.NewImageOp(assets.AppIcons["ic_about"]), Scale: 1}
							return layout.Inset{Right: 11}.Layout(gtx, img.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Max.Y = gtx.Dp(36)
							gtx.Constraints.Min = gtx.Constraints.Max
							title := material.Label(theme.Current().Theme, 18, n.Title)
							title.Font.Weight = font.Bold
							title.Color = n.Modal.TextColor
							return layout.Center.Layout(gtx, title.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Max.Y = gtx.Dp(20)
							gtx.Constraints.Min = gtx.Constraints.Max
							return layout.Center.Layout(gtx,
								material.Label(theme.Current().Theme, 12, n.Subtitle).Layout)
						}),
					)
				})
			}),
			layout.Flexed(1, func(gtx C) D {
				return n.layoutNavList(gtx)
			}),
		)
		return D{Size: gtx.Constraints.Max}
	})
}

func (n *NavDrawer) layoutNavList(gtx C) D {
	n.selectedChanged = false
	gtx.Constraints.Min.Y = 0
	n.navList.Axis = n.navigationAxis

	return n.navList.Layout(gtx, len(n.items), func(gtx C, index int) D {
		if n.navigationAxis == layout.Vertical {
			gtx.Constraints.Max.Y = gtx.Dp(48)
			gtx.Constraints.Min = gtx.Constraints.Max
		}
		dimensions := layout.Flex{Axis: n.navigationAxis}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return n.items[index].Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Spacer{Width: 20, Height: 20}.Layout(gtx)
			}),
		)
		if n.items[index].Clicked(gtx) {
			n.changeSelected(index)
			n.Disappear()
		}
		return dimensions
	})
}

func (n *NavDrawer) changeSelected(newIndex int) {
	if newIndex == n.selectedItem && n.items[n.selectedItem].Selected {
		return
	}
	n.items[n.selectedItem].Selected = false
	n.selectedItem = newIndex
	n.items[n.selectedItem].Selected = true
	n.selectedChanged = true
}

func (n *NavDrawer) UnselectNavDestination() {
	n.items[n.selectedItem].Selected = false
	n.selectedChanged = false
}

func (n *NavDrawer) SetNavDestination(tag interface{}) {
	for i, item := range n.items {
		if item.Tag == tag {
			n.changeSelected(i)
			break
		}
	}
}

func (n *NavDrawer) CurrentNavDestination() interface{} {
	return n.items[n.selectedItem].Tag
}

func (n *NavDrawer) NavDestinationChanged() bool {
	return n.selectedChanged
}
