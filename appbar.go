package components

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
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
	*Button
}

func (a AppBarMoreActionClicked) AppBarEvent() {}

type AppBarAction struct {
	*AppBarMenu
	Buttons []*Button
}

type AppBar struct {
	AppBarAction
	Title      string
	moreButton *Button
}

func NewAppBar() *AppBar {
	return &AppBar{
		moreButton: NewButton(ButtonStyle{
			Icon:   NewIcon(icons.NavigationMoreVert),
			Colors: theme.Current().AppBarColors,
			Inset:  layout.UniformInset(16),
		}),
	}
}

func (a *AppBar) CurrentActions() AppBarAction {
	return a.AppBarAction
}

func (a *AppBar) SetActions(actions AppBarAction) {
	a.AppBarAction = actions
}

func (a *AppBar) Events(gtx C) []AppBarEvent {
	if a.AppBarMenu == nil {
		return nil
	}

	if a.moreButton.Clicked(gtx) {
		a.AppBarMenu.Appear()
	}

	if !a.AppBarMenu.IsVisible() {
		return nil
	}

	var events []AppBarEvent
	for _, button := range a.AppBarMenu.buttons {
		if !button.Disabled && button.Clicked(gtx) {
			button.Clickable.Click()
			events = append(events, AppBarMoreActionClicked{Button: button})
			a.AppBarMenu.Disappear()
		}
	}
	return events
}

func (a *AppBar) Layout(gtx C) D {
	gtx.Constraints.Max.Y = gtx.Dp(56)
	a.moreButton.Colors = theme.Current().AppBarColors

	layout.Flex{Alignment: layout.Middle}.Layout(gtx,
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

func NewAppBarMenu(buttons []*Button) *AppBarMenu {
	m := &AppBarMenu{
		Modal: NewModal(layout.SW, layout.UniformInset(0), layout.UniformInset(0), 10,
			NewModalAnimationLeftDown()),
	}

	m.buttons = buttons
	for _, button := range buttons {
		m.flexChild = append(m.flexChild, layout.Rigid(func(gtx C) D {
			button.Colors.TextColor = theme.Current().ModalColors.TextColor
			return button.Layout(gtx)
		}))
	}
	return m
}

func (m *AppBarMenu) Layout(gtx C) D {
	th := theme.Current()
	gtx.Constraints.Min = image.Pt(0, 0)

	var children []layout.FlexChild
	children = toFlex(children, m.buttons...)

	r := op.Record(gtx.Ops)
	dims := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, m.flexChild...)
	c := r.Stop()

	gtx.Constraints.Min = gtx.Constraints.Max
	pos := gtx.Constraints.Max.Sub(dims.Size)

	m.Modal.OuterInset.Bottom = unit.Dp(float32(pos.Y) / gtx.Metric.PxPerDp)
	m.Modal.OuterInset.Left = unit.Dp(float32(pos.X) / gtx.Metric.PxPerDp)
	m.Modal.FixedBGColor = th.SurfaceColor

	return m.Modal.Layout(gtx, func(gtx C) D {
		c.Add(gtx.Ops)
		return D{Size: gtx.Constraints.Max}
	})
}

func paintColor(gtx C, max image.Point, r int, color color.NRGBA) {
	bounds := image.Rectangle{Max: max}
	defer clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
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
		//b := b
		children = append(children, layout.Rigid(func(gtx C) D {
			return b.Layout(gtx)
		}))
	}
	return children
}
