package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/maxazimi/v2ray-gio/ui/anim"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type ButtonAnimation struct {
	clickable    *widget.Clickable
	animIn       *anim.Animation
	transFunc    anim.TransFunc
	animOut      *anim.Animation
	animClick    *anim.Animation
	animLoading  *anim.Animation
	transLoading anim.TransFunc
}

func NewButtonAnimationDefault() *ButtonAnimation {
	return NewButtonAnimationScale(ratio)
}

func NewButtonAnimationScale(v float32) *ButtonAnimation {
	animIn := anim.New(false, gween.NewSequence(gween.New(1, v, .1, ease.Linear)))
	animOut := anim.New(false, gween.NewSequence(gween.New(v, 1, .1, ease.Linear)))

	animClick := anim.New(false, gween.NewSequence(
		gween.New(1, v*ratio, .1, ease.Linear),
		gween.New(v*ratio, 1, .4, ease.OutBounce),
	))

	animLoading := anim.New(false, gween.NewSequence(gween.New(0, 1, 1, ease.Linear)))
	animLoading.Sequence.SetLoop(-1)

	return &ButtonAnimation{
		clickable:    new(widget.Clickable),
		animIn:       animIn,
		transFunc:    anim.TransScale,
		animOut:      animOut,
		animClick:    animClick,
		animLoading:  animLoading,
		transLoading: anim.TransRotate,
	}
}

func (b *ButtonAnimation) Hovered() bool {
	return b.clickable.Hovered()
}

func (b *ButtonAnimation) Clicked(gtx C) bool {
	return b.clickable.Clicked(gtx)
}

func (b *ButtonAnimation) Layout(gtx C, w layout.Widget) D {
	return b.clickable.Layout(gtx, func(gtx C) D {
		if b.animIn != nil {
			value, finished := b.animIn.Update(gtx)
			if !finished {
				defer b.transFunc(gtx, value).Push(gtx.Ops).Pop()
			}
		}

		if b.animOut != nil {
			value, finished := b.animOut.Update(gtx)
			if !finished {
				defer b.transFunc(gtx, value).Push(gtx.Ops).Pop()
			}
		}

		if b.animClick != nil {
			value, finished := b.animClick.Update(gtx)
			if !finished {
				defer b.transFunc(gtx, value).Push(gtx.Ops).Pop()
			}
		}

		return w(gtx)
	})
}
