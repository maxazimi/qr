package components

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
)

type AppBarEvent interface {
	AppBarEvent()
}

type AppBarNavigationClicked struct{}

func (a AppBarNavigationClicked) AppBarEvent() {}

type AppBarMoreActionClicked struct {
	Tag interface{}
}

func (a AppBarMoreActionClicked) AppBarEvent() {}

type MenuItem struct {
	Name string
	Tag  interface{}
}

type AppBarAction struct {
	*AppBarMenu
	Buttons    []*Button
	IsMainPage bool
}

type AppBar struct {
	AppBarAction
	Title           string
	navDrawer       *NavDrawer
	moreButton      *Button
	backButton      *Button
	navButton       *Button
	activeNavButton *Button
}

func NewAppBar(drawer *NavDrawer) *AppBar {
	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	moreIcon, _ := widget.NewIcon(icons.NavigationMoreVert)
	backIcon, _ := widget.NewIcon(icons.NavigationArrowBack)

	colors := theme.ButtonColors{
		TextColor:       theme.WhiteColor,
		BackgroundColor: color.NRGBA{},
	}

	navButton := NewButton(ButtonStyle{
		Icon:   navIcon,
		Colors: colors,
		Inset:  layout.UniformInset(16),
	})

	return &AppBar{
		navDrawer: drawer,
		moreButton: NewButton(ButtonStyle{
			Icon:   moreIcon,
			Colors: colors,
			Inset:  layout.UniformInset(16),
		}),
		backButton: NewButton(ButtonStyle{
			Icon:   backIcon,
			Colors: colors,
			Inset:  layout.UniformInset(16),
		}),
		navButton:       navButton,
		activeNavButton: navButton,
	}
}

func (a *AppBar) SetActions(actions AppBarAction) {
	a.AppBarAction = actions

	if actions.IsMainPage {
		a.activeNavButton = a.navButton
	} else {
		a.activeNavButton = a.backButton
	}
}

func (a *AppBar) Events(gtx C) []AppBarEvent {
	var events []AppBarEvent

	if a.activeNavButton.Clicked(gtx) {
		events = append(events, AppBarNavigationClicked{})
	}

	if a.AppBarMenu != nil {
		if a.moreButton.Clicked(gtx) {
			a.AppBarMenu.Appear()
		}
		for _, button := range a.AppBarMenu.buttons {
			if button.Clicked(gtx) {
				events = append(events, AppBarMoreActionClicked{Tag: button.Tag})
				a.AppBarMenu.Disappear()
			}
		}
	}
	return events
}

func (a *AppBar) Layout(gtx C) D {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	color1 := theme.Current().AppBarColors.BackgroundColor
	color2 := theme.Current().AppBarColors.BorderColor
	paintGradient(gtx, gtx.Constraints.Max, 0, color1, color2)

	layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return a.activeNavButton.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: 16}.Layout(gtx, func(gtx C) D {
				title := material.Body1(theme.Current().Theme, a.Title)
				title.Color = theme.Current().AppBarColors.TextColor
				title.TextSize = 18
				return title.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.E.Layout(gtx, func(gtx C) D {
				var children []layout.FlexChild
				children = toFlex(children, a.Buttons...)
				if a.AppBarMenu != nil {
					children = toFlex(children, a.moreButton)
				}
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx, children...)
			})
		}),
	)
	return D{Size: gtx.Constraints.Max}
}

type AppBarMenu struct {
	*Modal
	buttons   []*Button
	flexChild []layout.FlexChild
}

func NewAppBarMenu(items ...MenuItem) *AppBarMenu {
	m := &AppBarMenu{Modal: NewModal(layout.SW, layout.UniformInset(0), layout.UniformInset(0), 5,
		NewModalAnimationLeftDown())}

	for _, item := range items {
		button := NewButton(ButtonStyle{
			Tag:      item.Tag,
			Text:     item.Name,
			TextSize: 16,
			Radius:   5,
			Inset:    layout.UniformInset(5),
		})
		m.buttons = append(m.buttons, button)
		m.flexChild = append(m.flexChild, layout.Rigid(func(gtx C) D {
			button.Colors.TextColor = theme.Current().ModalColors.TextColor
			return button.Layout(gtx)
		}))
	}
	return m
}

func (m *AppBarMenu) Layout(gtx C) D {
	gtx.Constraints.Min = image.Pt(0, 0)

	var children []layout.FlexChild
	children = toFlex(children, m.buttons...)

	r := op.Record(gtx.Ops)
	dims := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, m.flexChild...)
	c := r.Stop()

	gtx.Constraints.Min = gtx.Constraints.Max
	pos := gtx.Constraints.Max.Sub(dims.Size)
	m.Modal.OuterInset.Bottom = unit.Dp(pos.Y)
	m.Modal.OuterInset.Left = unit.Dp(pos.X)
	m.Modal.BackgroundColor = theme.Current().ModalColors.BackgroundColor

	return m.Modal.Layout(gtx, func(gtx C) D {
		c.Add(gtx.Ops)
		return D{Size: gtx.Constraints.Max}
	})
}

func paintGradient(gtx C, max image.Point, r int, color1, color2 color.NRGBA) {
	bounds := image.Rectangle{Max: max}
	defer clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Push(gtx.Ops).Pop()

	paint.LinearGradientOp{
		Color1: color1,
		Color2: color2,
		Stop2:  f32.Pt(float32(max.X), float32(max.Y)),
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func toFlex(children []layout.FlexChild, buttons ...*Button) []layout.FlexChild {
	for _, b := range buttons {
		b := b
		children = append(children, layout.Rigid(func(gtx C) D {
			return b.Layout(gtx)
		}))
	}
	return children
}
