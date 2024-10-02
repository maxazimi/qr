// from https://github.com/g45t345rt/g45w/blob/master/components/drag_items.go

package components

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/maxazimi/v2ray-gio/ui/anim"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"image"
	"log"
	"time"
)

type draggable struct {
	Index int
	W     layout.Widget
	Dims  D
}

type DragItems struct {
	items              []draggable
	dragItem           draggable
	dragIndex          int
	drag               gesture.Drag
	dragEvent          *pointer.Event
	startPosY          float32
	dragPosY           float32
	itemMoved          bool
	canStartDrag       bool
	holdStartAnimation *anim.Animation
	holdPress          *HoldPress

	lastIndex int
	newIndex  int
}

func NewDragItems() *DragItems {
	holdStartAnimation := anim.New(false, gween.NewSequence(
		gween.New(1, 1.1, .1, ease.Linear),
		gween.New(1.1, 1, .1, ease.Linear),
	))

	return &DragItems{
		holdStartAnimation: holdStartAnimation,
		holdPress:          NewHoldPress(500 * time.Millisecond),
	}
}

func (l *DragItems) ItemMoved() (bool, int, int) {
	if l.itemMoved && l.lastIndex != l.newIndex {
		return true, l.lastIndex, l.newIndex
	}
	return false, -1, -1
}

func (l *DragItems) LayoutItem(gtx C, index int, w layout.Widget) {
	r := op.Record(gtx.Ops)
	dims := w(gtx)
	r.Stop()

	l.items = append(l.items, draggable{index, w, dims})
}

func (l *DragItems) Layout(gtx C, scroll *layout.Position, w layout.Widget) D {
	l.items = make([]draggable, 0)
	m := op.Record(gtx.Ops)
	dims := l.holdPress.Layout(gtx, w)
	c := m.Stop()

	scrollOffset := 0
	itemOffset := 0
	if scroll != nil {
		scrollOffset = scroll.Offset
		itemOffset = scroll.First
	}

	if l.holdPress.Triggered {
		l.canStartDrag = true
		l.holdStartAnimation.Reset().Start()
	}

	l.itemMoved = false
	//for _, e := range l.drag.Events(gtx.Metric, gtx.Queue, gesture.Both) {
	for {
		event, ok := gtx.Event(
			key.Filter{
				Focus: &l.drag,
			},
		)
		if !ok {
			break
		}
		e, ok := event.(pointer.Event)
		if !ok {
			continue
		}
		log.Println("\n\n\nHkdjhfgjdhfghjfdgjg\n\n\n")
		switch e.Kind {
		case pointer.Drag:
			if l.canStartDrag {
				l.dragEvent = &e
			}
		case pointer.Press:
			l.dragEvent = nil
			l.canStartDrag = false

			l.startPosY = e.Position.Y
			l.dragIndex = -1
			minY := 0 - scrollOffset
			maxY := 0 - scrollOffset
			for i, item := range l.items {
				maxY += item.Dims.Size.Y
				if l.startPosY >= float32(minY) && l.startPosY <= float32(maxY) {
					l.dragIndex = i
					l.dragItem = item
					l.dragPosY = l.startPosY - float32(item.Dims.Size.Y/2)
					break
				}

				minY += item.Dims.Size.Y
			}
		case pointer.Release, pointer.Cancel:
			l.canStartDrag = false

			if l.dragEvent != nil && l.dragIndex > -1 {
				itemPosY := float32(0) - float32(scrollOffset)
				for i, item := range l.items {
					if itemPosY+float32(item.Dims.Size.Y/2) > l.dragPosY {
						if l.dragIndex != i {
							l.itemMoved = true
							l.lastIndex = l.dragItem.Index
							l.newIndex = i + itemOffset
						}

						break
					}
					itemPosY += float32(item.Dims.Size.Y)
				}
			}

			l.dragEvent = nil
		}
	}

	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	l.drag.Add(gtx.Ops)
	c.Add(gtx.Ops)

	if l.canStartDrag {
		offsetY := float32(0)
		for i, item := range l.items {
			if i < l.dragIndex {
				offsetY += float32(item.Dims.Size.Y)
			} else {
				break
			}
		}

		if l.dragEvent != nil {
			l.dragPosY = l.dragEvent.Position.Y - l.startPosY + offsetY - float32(scrollOffset)

			if scroll != nil {
				if l.dragPosY < 0 && (scroll.Offset > 0 || scroll.First > 0) {
					v := gtx.Dp(5)
					scroll.Offset -= v
					scroll.BeforeEnd = true
				}

				itemHeight := l.dragItem.Dims.Size.Y
				if l.dragPosY+float32(itemHeight) > float32(dims.Size.Y) {
					v := gtx.Dp(5)
					scroll.Offset += v
					scroll.BeforeEnd = true
				}
			}
		}

		x := float32(0)

		value, finished := l.holdStartAnimation.Update(gtx)
		if !finished {
			origin := f32.Pt(float32(l.dragItem.Dims.Size.X/2), float32(l.dragItem.Dims.Size.Y/2))
			x := value
			scale := f32.Affine2D{}.Scale(origin, f32.Pt(x, x))
			defer op.Affine(scale).Push(gtx.Ops).Pop()
		}

		offset := f32.Affine2D{}.Offset(f32.Pt(x, l.dragPosY))
		defer op.Affine(offset).Push(gtx.Ops).Pop()
		l.dragItem.W(gtx)
		pointer.CursorGrabbing.Add(gtx.Ops)
	}

	return dims
}
