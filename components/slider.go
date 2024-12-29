// SPDX-License-Identifier: Unlicense OR MIT

package components

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/maxazimi/qr/ui/theme"
	"github.com/maxazimi/qr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"time"
)

const (
	defaultDuration   = 500 * time.Millisecond
	defaultDragEffect = 100
)

type SwipeDirection int

const (
	SwipeLeft SwipeDirection = iota
	SwipeRight
)

type Dragged func(dragDirection SwipeDirection)

type sliderItem struct {
	widgetItem layout.Widget
	button     *Clickable
}

type Slider struct {
	icon       *Icon
	nextButton *Button
	prevButton *Button
	card       Card
	slideItems []*sliderItem

	selected         int
	isSliderItemsSet bool

	// colors of the indicator and navigation button
	ButtonBackgroundColor    color.NRGBA
	IndicatorBackgroundColor color.NRGBA
	SelectedIndicatorColor   color.NRGBA // this is a full color no opacity
	slideAction              *SlideAction
	clicker                  gesture.Click
	clicked                  bool
	disableButtonDirection   bool
	ControlInset             layout.Inset
}

type SlideAction struct {
	Duration  time.Duration
	IsReverse bool
	push      int
	next      *op.Ops
	nextCall  op.CallOp
	lastCall  op.CallOp

	t0     time.Time
	offset float32

	// animation state
	dragEffect  int
	dragStarted f32.Point
	dragOffset  int
	drag        gesture.Drag
	dragged     Dragged
	isPushing   bool
}

func NewSlideAction() *SlideAction {
	return &SlideAction{
		Duration:   defaultDuration,
		dragEffect: defaultDragEffect,
	}
}

func NewSlider() *Slider {
	th := theme.Current()
	s := &Slider{
		icon:       NewIcon(icons.ImageBrightness1),
		slideItems: make([]*sliderItem, 0),
		prevButton: NewButton(ButtonStyle{
			ImgKey:    theme.ICChevronLeft,
			Animation: NewButtonAnimationDefault(),
		}),
		nextButton: NewButton(ButtonStyle{
			ImgKey:    theme.ICChevronRight,
			Animation: NewButtonAnimationDefault(),
		}),
		ButtonBackgroundColor:    th.LightGrayColor,
		IndicatorBackgroundColor: th.LightGrayColor,
		SelectedIndicatorColor:   theme.WhiteColor,
		slideAction:              NewSlideAction(),
	}

	s.ControlInset = layout.Inset{
		Right:  values.DP16,
		Left:   values.DP16,
		Bottom: values.DP16,
	}

	s.card = NewCard()
	s.card.Radius = Radius(8)

	s.slideAction.Dragged(func(dragDirection SwipeDirection) {
		isNext := dragDirection == SwipeLeft
		s.handleActionEvent(isNext)
	})

	return s
}

// GetSelectedIndex returns the index of the current slider item.
func (s *Slider) GetSelectedIndex() int {
	return s.selected
}

func (s *Slider) IsLastSlide() bool {
	return s.selected == len(s.slideItems)-1
}

func (s *Slider) NextSlide() {
	s.handleActionEvent(true)
}

func (s *Slider) ResetSlide() {
	s.selected = 0
}

func (s *Slider) SetDisableDirectionBtn(disable bool) {
	s.disableButtonDirection = disable
}

func (s *Slider) sliderItems(items []layout.Widget) []*sliderItem {
	slideItems := make([]*sliderItem, 0)
	for _, item := range items {
		slideItems = append(slideItems, &sliderItem{
			widgetItem: item,
			button:     NewClickable(false),
		})
	}

	return slideItems
}

func (s *Slider) Layout(gtx C, items []layout.Widget) D {
	// set slider items once since layout is drawn multiple times per sec.
	if !s.isSliderItemsSet {
		s.slideItems = s.sliderItems(items)
		s.isSliderItemsSet = true
	}

	if len(s.slideItems) == 0 {
		return D{}
	}

	s.handleClickEvent(gtx)
	var dims D
	var call op.CallOp
	{
		m := op.Record(gtx.Ops)
		dims = s.slideAction.DragLayout(gtx, func(gtx C) D {
			return layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return s.slideAction.TransformLayout(gtx, s.slideItems[s.selected].widgetItem)
				}),
				layout.Stacked(func(gtx C) D {
					if len(s.slideItems) == 1 {
						return D{}
					}
					return s.ControlInset.Layout(gtx, func(gtx C) D {
						return layout.Flex{
							Axis: layout.Horizontal,
						}.Layout(gtx,
							layout.Rigid(s.selectedItemIndicatorLayout),
							layout.Flexed(1, func(gtx C) D {
								if s.disableButtonDirection {
									return D{}
								}
								return layout.E.Layout(gtx, s.buttonLayout)
							}),
						)
					})
				}),
			)
		})
		call = m.Stop()
	}

	area := clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops)
	s.clicker.Add(gtx.Ops)
	defer area.Pop()

	call.Add(gtx.Ops)
	return dims
}

func (s *Slider) buttonLayout(gtx C) D {
	s.card.Radius = Radius(10)
	s.card.CardColor = s.ButtonBackgroundColor
	return s.containerLayout(gtx, func(gtx C) D {
		return layout.Inset{
			Right: values.DP4,
			Left:  values.DP4,
		}.Layout(gtx, func(gtx C) D {
			return LinearLayout{
				Width:       WrapContent,
				Height:      WrapContent,
				Orientation: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return s.prevButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return s.nextButton.Layout(gtx)
				}),
			)
		})
	})
}

func (s *Slider) selectedItemIndicatorLayout(gtx C) D {
	s.card.Radius = Radius(10)
	s.card.CardColor = s.IndicatorBackgroundColor
	return s.containerLayout(gtx, func(gtx C) D {
		return layout.Inset{
			Right: values.DP4,
			Left:  values.DP4,
		}.Layout(gtx, func(gtx C) D {
			list := &layout.List{Axis: layout.Horizontal}
			return list.Layout(gtx, len(s.slideItems), func(gtx C, i int) D {
				s.icon.Color = color.NRGBA{R: 0, G: 0, B: 0, A: opacity2Transparency(0.2)}
				if i == s.selected {
					s.icon.Color = s.SelectedIndicatorColor
				}
				return layout.Inset{
					Top:    values.DP4,
					Bottom: values.DP4,
					Right:  values.DP4,
					Left:   values.DP4,
				}.Layout(gtx, func(gtx C) D {
					return s.slideItems[i].button.Layout(gtx, func(gtx C) D {
						return s.icon.LayoutSize(gtx, values.DP12)
					})
				})
			})
		})
	})
}

func (s *Slider) containerLayout(gtx C, content layout.Widget) D {
	return s.card.Layout(gtx, content)
}

func (s *Slider) RefreshItems() {
	s.isSliderItemsSet = false
}

func (s *Slider) Clicked() bool {
	clicked := s.clicked
	s.clicked = false
	return clicked
}

func (s *Slider) handleClickEvent(gtx C) {
	if s.nextButton.Clicked(gtx) {
		s.handleActionEvent(true)
	}

	if s.prevButton.Clicked(gtx) {
		s.handleActionEvent(false)
	}

	for {
		e, ok := s.clicker.Update(gtx.Source)
		if !ok {
			break
		}
		if e.Kind == gesture.KindClick {
			if !s.clicked {
				s.clicked = true
			}
		}
	}

	for i, item := range s.slideItems {
		if item.button.Clicked(gtx) {
			if i == s.selected {
				continue
			}
			lastSelected := s.selected
			s.selected = i
			if lastSelected < i {
				s.slideAction.PushLeft()
			} else {
				s.slideAction.PushRight()
			}
			break
		}
	}
}

func (s *Slider) handleActionEvent(isNext bool) {
	if len(s.slideItems) == 1 {
		return
	}
	l := len(s.slideItems) - 1 // index starts at 0
	if isNext {
		if s.selected == l {
			s.selected = 0
		} else {
			s.selected++
		}
		s.slideAction.PushLeft()
	} else {
		if s.selected == 0 {
			s.selected = l
		} else {
			s.selected--
		}
		s.slideAction.PushRight()
	}
}

// PushLeft pushes the existing widget to the left.
func (s *SlideAction) PushLeft() { s.push = 1 }

// PushRight pushes the existing widget to the right.
func (s *SlideAction) PushRight() { s.push = -1 }

func (s *SlideAction) SetDragEffect(offset int) { s.dragEffect = offset }

func (s *SlideAction) Dragged(drag Dragged) {
	s.dragged = drag
}

func (s *SlideAction) DragLayout(gtx C, w layout.Widget) D {
	for {
		event, ok := s.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal)
		if !ok {
			break
		}
		switch event.Kind {
		case pointer.Press:
			s.dragStarted = event.Position
			s.dragOffset = 0
		case pointer.Drag:
			newOffset := int(s.dragStarted.X - event.Position.X)
			if newOffset > s.dragEffect {
				if !s.isPushing && s.dragged != nil {
					s.isPushing = true
					s.dragged(SwipeLeft)
				}
			} else if newOffset < -s.dragEffect {
				if !s.isPushing && s.dragged != nil {
					s.isPushing = true
					s.dragged(SwipeRight)
				}
			}
			s.dragOffset = newOffset
		case pointer.Release:
			fallthrough
		case pointer.Cancel:
			s.isPushing = false
		}
	}
	var dims layout.Dimensions
	var call op.CallOp
	{
		m := op.Record(gtx.Ops)
		dims = w(gtx)
		call = m.Stop()
	}

	area := clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops)
	s.drag.Add(gtx.Ops)
	defer area.Pop()

	call.Add(gtx.Ops)
	return dims
}

// TransformLayout perform transition effects between 2 widgets
func (s *SlideAction) TransformLayout(gtx C, w layout.Widget) D {
	if s.push != 0 {
		s.next = nil
		s.lastCall = s.nextCall
		s.offset = float32(s.push)
		s.t0 = gtx.Now
		s.push = 0
	}

	var delta time.Duration
	if !s.t0.IsZero() {
		now := gtx.Now
		delta = now.Sub(s.t0)
		s.t0 = now
	}

	// Calculate the duration of transition effects
	if s.offset != 0 {
		duration := s.Duration
		if duration == 0 {
			duration = defaultDuration
		}
		movement := float32(delta.Seconds()) / float32(duration.Seconds())
		if s.offset < 0 {
			s.offset += movement
			if s.offset >= 0 {
				s.offset = 0
			}
		} else {
			s.offset -= movement
			if s.offset <= 0 {
				s.offset = 0
			}
		}

		gtx.Execute(op.InvalidateCmd{})
	}

	// Record the widget presentation
	var dims layout.Dimensions
	{
		if s.next == nil {
			s.next = new(op.Ops)
		}
		gtx := gtx
		gtx.Ops = s.next
		gtx.Ops.Reset()
		m := op.Record(gtx.Ops)
		dims = w(gtx)
		s.nextCall = m.Stop()
	}

	if s.offset == 0 {
		s.nextCall.Add(gtx.Ops)
		return dims
	}

	offset := smooth(s.offset)

	reverse := 1
	if s.IsReverse {
		reverse = -1
	}

	// Implement transition effects for widgets
	if s.offset > 0 {
		defer op.Offset(image.Point{
			X: int(float32(dims.Size.X)*(offset-1)) * reverse,
		}).Push(gtx.Ops).Pop()
		s.lastCall.Add(gtx.Ops)

		defer op.Offset(image.Point{
			X: dims.Size.X * reverse,
		}).Push(gtx.Ops).Pop()
		s.nextCall.Add(gtx.Ops)
	} else {
		defer op.Offset(image.Point{
			X: int(float32(dims.Size.X)*(offset+1)) * reverse,
		}).Push(gtx.Ops).Pop()
		s.lastCall.Add(gtx.Ops)

		defer op.Offset(image.Point{
			X: -dims.Size.X * reverse,
		}).Push(gtx.Ops).Pop()
		s.nextCall.Add(gtx.Ops)
	}
	return dims
}

// smooth handles -1 to 1 with ease-in-out cubic easing func.
func smooth(t float32) float32 {
	if t < 0 {
		return -easeInOutCubic(-t)
	}
	return easeInOutCubic(t)
}

// easeInOutCubic maps a linear value to a ease-in-out-cubic easing function.
// It is a mathematical function that describes how a value changes over time.
// It can be applied to adjusting the speed of animation
func easeInOutCubic(t float32) float32 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return (t-1)*(2*t-2)*(2*t-2) + 1
}
