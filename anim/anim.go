// https://github.com/g45t345rt/g45w/blob/master/animation/animation.go

package anim

import (
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/tanema/gween"
	"time"
)

type (
	C = layout.Context
)

type Animation struct {
	Sequence      *gween.Sequence
	active        bool
	stop          bool
	delay         time.Duration
	lastFrameTime time.Time
	startTime     time.Time
}

func New(start bool, sequence *gween.Sequence) *Animation {
	return &Animation{
		Sequence: sequence,
		stop:     !start,
		active:   start,
	}
}

func (a *Animation) Update(gtx C) (float32, bool) {
	now := time.Now()
	var dt time.Duration

	if a.startTime.IsZero() {
		a.startTime = now
	}

	if !a.lastFrameTime.IsZero() {
		dt = now.Sub(a.lastFrameTime)
	}

	if now.Sub(a.startTime) > a.delay && !a.stop {
		a.lastFrameTime = now
	}

	seconds := float32(dt.Seconds())
	value, _, finished := a.Sequence.Update(seconds)

	if !a.stop {
		gtx.Execute(op.InvalidateCmd{})
	}

	if finished {
		a.stop = true
	}

	return value, finished
}

func (a *Animation) Start() *Animation {
	if a.stop {
		a.Reset()
		a.stop = false
		a.active = true
	}
	return a
}

func (a *Animation) StartWithDelay(delay time.Duration) *Animation {
	if a.stop {
		a.Reset()
		a.delay = delay
		a.stop = false
		a.active = true
	}
	return a
}

func (a *Animation) Resume() *Animation {
	a.lastFrameTime = time.Time{}
	a.stop = false
	return a
}

func (a *Animation) Pause() *Animation {
	a.lastFrameTime = time.Time{}
	a.stop = true
	return a
}

func (a *Animation) Reset() *Animation {
	a.active = false
	a.delay = 0
	a.lastFrameTime = time.Time{}
	a.startTime = time.Time{}
	a.stop = true
	a.Sequence.Reset()
	return a
}

func (a *Animation) IsActive() bool {
	return a.active
}
