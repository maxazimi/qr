package components

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/maxazimi/qr/ui/theme"
	"image"
)

type Switch struct {
	*theme.Theme
	*widget.Bool
	disabled bool
}

type SwitchItem struct {
	Text   string
	button Button
}

type SwitchButtonText struct {
	*theme.Theme
	items    []SwitchItem
	selected int
	changed  bool
}

func NewSwitch() *Switch {
	return &Switch{
		Theme: theme.Current(),
		Bool:  new(widget.Bool),
	}
}

func NewSwitchButtonText(i []SwitchItem) *SwitchButtonText {
	sw := &SwitchButtonText{
		Theme: theme.Current(),
		items: make([]SwitchItem, len(i)+1),
	}

	for index := range i {
		i[index].button.Text = i[index].Text
		i[index].button.Colors = sw.Theme.ButtonColors
		i[index].button.Colors.BackgroundColor = sw.SwitchColors.InactiveColor
		i[index].button.Colors.TextColor = sw.SwitchColors.InactiveTextColor
		i[index].button.TextSize = unit.Sp(14)
		sw.items[index+1] = i[index]
	}

	if len(sw.items) > 0 {
		sw.selected = 1
	}
	return sw
}

func (s *Switch) Layout(gtx C) D {
	s.Theme = theme.Current()
	dGtx := gtx

	trackWidth := dGtx.Dp(32)
	trackHeight := dGtx.Dp(20)
	thumbSize := dGtx.Dp(18)
	trackOff := (thumbSize - trackHeight) / 2

	// Draw track
	trackCorner := trackHeight / 2
	trackRect := image.Rectangle{Max: image.Point{
		X: trackWidth,
		Y: trackHeight,
	}}

	activeColor := s.SwitchColors.ActiveColor
	InactiveColor := s.SwitchColors.InactiveColor
	thumbColor := s.SwitchColors.ThumbColor

	if s.disabled {
		dGtx = gtx.Disabled()
		activeColor, InactiveColor, thumbColor = Disabled(activeColor), Disabled(InactiveColor), Disabled(thumbColor)
	}

	col := InactiveColor
	if s.IsChecked() {
		col = activeColor
	}

	trackColor := col
	t := op.Offset(image.Point{Y: trackOff}).Push(dGtx.Ops)
	cl := clip.UniformRRect(trackRect, trackCorner).Push(dGtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(dGtx.Ops)
	paint.PaintOp{}.Add(dGtx.Ops)
	cl.Pop()
	t.Pop()

	// Compute thumb offset and color.
	if s.IsChecked() {
		off := trackWidth - thumbSize
		defer op.Offset(image.Point{X: off}).Push(dGtx.Ops).Pop()
	}

	thumbRadius := thumbSize / 2

	circle := func(x, y, r int) clip.Op {
		b := image.Rectangle{
			Min: image.Pt(x-r, y-r),
			Max: image.Pt(x+r, y+r),
		}
		return clip.Ellipse(b).Op(dGtx.Ops)
	}

	// Draw thumb shadow, a translucent disc slightly larger than the thumb itself.
	// Center shadow horizontally and slightly adjust its Y.
	paint.FillShape(dGtx.Ops, col, circle(thumbRadius, thumbRadius+dGtx.Dp(.25), thumbRadius+1))

	// Draw thumb.
	paint.FillShape(dGtx.Ops, thumbColor, circle(thumbRadius, thumbRadius, thumbRadius))

	// Set up click area.
	clickSize := dGtx.Dp(38)
	clickOff := image.Point{
		X: (trackWidth) - (clickSize),
		Y: (trackHeight) - (clickSize)/2 + trackOff,
	}
	defer op.Offset(clickOff).Push(dGtx.Ops).Pop()
	sz := image.Pt(clickSize, clickSize)
	defer clip.Ellipse(image.Rectangle{Max: sz}).Push(dGtx.Ops).Pop()
	s.Bool.Layout(dGtx, func(dGtx C) D {
		semantic.Switch.Add(dGtx.Ops)
		return D{Size: sz}
	})

	dims := image.Point{X: trackWidth, Y: thumbSize}
	return D{Size: dims}
}

func (s *Switch) Changed(gtx C) bool {
	if s.disabled {
		return false
	}
	return s.Bool.Update(gtx)
}

func (s *Switch) IsChecked() bool {
	return s.Bool.Value
}

func (s *Switch) SetChecked(value bool) {
	s.Bool.Value = value
}

func (s *Switch) SetEnabled(value bool) {
	s.disabled = !value
}

func (s *SwitchButtonText) Layout(gtx C) D {
	s.Theme = theme.Current()
	s.handleClickEvent(gtx)

	m8 := unit.Dp(8)
	m4 := unit.Dp(4)
	card := NewCard()
	card.CardColor = s.Theme.CardColor // Gray1
	card.Radius = Radius(8)
	return card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
			list := &layout.List{Axis: layout.Horizontal}
			Items := s.items[1:]
			return list.Layout(gtx, len(Items), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
					index := i + 1
					btn := s.items[index].button
					btn.Inset = layout.Inset{
						Left:   m8,
						Bottom: m4,
						Right:  m8,
						Top:    m4,
					}
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(btn.Layout),
					)
				})
			})
		})
	})
}

func (s *SwitchButtonText) handleClickEvent(gtx C) {
	for index := range s.items {
		if index != 0 {
			if s.items[index].button.Clicked(gtx) {
				if s.selected != index {
					s.changed = true
				}
				s.selected = index
			}
		}

		if s.selected == index {
			s.items[s.selected].button.Colors.BackgroundColor = s.SwitchColors.ActiveColor
			s.items[s.selected].button.Colors.TextColor = s.SwitchColors.ActiveTextColor
		} else {
			s.items[index].button.Colors.BackgroundColor = s.SwitchColors.InactiveColor
			s.items[index].button.Colors.TextColor = s.SwitchColors.InactiveTextColor
		}
	}
}

func (s *SwitchButtonText) SelectedOption() string {
	return s.items[s.selected].Text
}

func (s *SwitchButtonText) SelectedIndex() int {
	return s.selected
}

func (s *SwitchButtonText) Changed() bool {
	changed := s.changed
	s.changed = false
	return changed
}

func (s *SwitchButtonText) SetSelectedIndex(index int) {
	s.selected = index
}
