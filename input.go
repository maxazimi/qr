// from https://github.com/g45t345rt/g45w/blob/master/components/input.go

package components

import (
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"image"
	"image/color"
)

type Input struct {
	Colors     theme.InputColors
	FontWeight font.Weight
	TextSize   unit.Sp
	Editor     *widget.Editor
	EditorMinY int
	Border     widget.Border
	Inset      layout.Inset
	Clickable  *widget.Clickable

	keyboardClick *widget.Clickable
	submitted     bool
	submitText    string
	activeSubmit  bool

	setValue bool
	newValue string
}

func NewInput() *Input {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.InputHint = key.HintText // Cap sentence flag
	border := widget.Border{
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}

	return &Input{
		Editor:        editor,
		Border:        border,
		Clickable:     new(widget.Clickable),
		keyboardClick: new(widget.Clickable),
		Inset: layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
	}
}

func NewNumberInput() *Input {
	input := NewInput()
	input.Editor.Filter = "0123456789."
	input.Editor.InputHint = key.HintNumeric
	return input
}

func NewPasswordInput() *Input {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.InputHint = key.HintPassword
	editor.Mask = rune(42) // mask with *
	border := widget.Border{
		Color:        color.NRGBA{A: 240},
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}

	return &Input{
		Editor:        editor,
		Border:        border,
		Clickable:     new(widget.Clickable),
		keyboardClick: new(widget.Clickable),
		Inset: layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
	}
}

func (t *Input) Value() string {
	return t.Editor.Text()
}

func (t *Input) SetValue(text string) {
	t.setValue = true
	t.newValue = text
}

func (t *Input) Submitted() (bool, string) {
	t.activeSubmit = true
	if t.submitted {
		t.submitted = false
		return true, t.submitText
	}

	return false, t.submitText
}

func (t *Input) Layout(gtx C, th *material.Theme, hint string) D {
	if t.activeSubmit {
		//for _, e := range t.Editor.Events() {
		//	e, ok := e.(widget.SubmitEvent)
		//	if ok {
		//		t.Editor.SetText("")
		//		t.submitText = e.Text
		//		t.submitted = true
		//	}
		//}

		for {
			event, ok := gtx.Event(
				key.Filter{
					Focus: t.Editor,
				},
			)
			if !ok {
				break
			}
			e, ok := event.(widget.EditorEvent)
			if !ok {
				continue
			}

			if e, ok := e.(widget.SubmitEvent); ok {
				t.Editor.SetText("")
				t.submitText = e.Text
				t.submitted = true
			}
		}
	}

	// always use SetText in layout, or you can sometime get nil with shaper text
	if t.setValue {
		t.setValue = false
		t.Editor.SetText(t.newValue)
		t.Editor.SetCaret(len(t.newValue), len(t.newValue))
		t.newValue = ""
	}

	gtx.Constraints.Min.Y = t.EditorMinY

	if t.keyboardClick.Clicked(gtx) {
		// on mobile if the keyboard popups and the input lose focus it will automatically close the keyboard
		// so we have to manually force keyboard request to avoid this issue
		if !t.Editor.ReadOnly {
			//key.SoftKeyboardOp{Show: true}.Add(gtx.Ops)
			gtx.Execute(key.SoftKeyboardCmd{Show: true})
		}
	}

	return t.Clickable.Layout(gtx, func(gtx C) D {
		return t.keyboardClick.Layout(gtx, func(gtx C) D {
			t.Border.Color = t.Colors.BorderColor
			return t.Border.Layout(gtx, func(gtx C) D {
				r := op.Record(gtx.Ops)
				dims := t.Inset.Layout(gtx, func(gtx C) D {
					editorStyle := material.Editor(th, t.Editor, hint)
					editorStyle.Color = t.Colors.TextColor
					editorStyle.HintColor = t.Colors.HintColor
					editorStyle.TextSize = th.TextSize
					if t.TextSize != 0 {
						editorStyle.TextSize = t.TextSize
					}
					editorStyle.Font.Weight = t.FontWeight
					return editorStyle.Layout(gtx)
				})
				c := r.Stop()

				paint.FillShape(gtx.Ops, t.Colors.BackgroundColor, clip.UniformRRect(
					image.Rectangle{Max: dims.Size},
					int(t.Border.CornerRadius),
				).Op(gtx.Ops))

				c.Add(gtx.Ops)
				return dims
			})
		})
	})
}
