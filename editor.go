// https://github.com/crypto-power/cryptopower/blob/master/ui/cryptomaterial/editor.go

package components

import (
	"gioui.org/gesture"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/maxazimi/v2ray-gio/ui/lang"
	"github.com/maxazimi/v2ray-gio/ui/theme"
	"github.com/maxazimi/v2ray-gio/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"io"
	"strings"
)

type RestoreEditor struct {
	Edit      *Editor
	Title     string
	LineColor color.NRGBA
	height    int
}

type Editor struct {
	material.EditorStyle

	Title     string
	errorText string
	LineColor color.NRGBA

	// IsRequired if true, displays a required field text at the bottom of the editor.
	IsRequired bool
	// IsTitleLabel if true makes the title label visible.
	IsTitleLabel bool

	// Set ExtraText to show a custom text in the editor.
	ExtraText string
	// Bordered if true makes the adds a border around the editor.
	Bordered bool
	// isPassword if true, displays the show and hide button.
	isPassword bool
	// If showEditorIcon is true, displays the editor widget Icon of choice
	showEditorIcon     bool
	alignEditorIconEnd bool

	HasCustomButton bool
	CustomButton    *Button

	// isEditorButtonClickable passes a clickable icon button if true and regular icon if false
	isEditorButtonClickable bool

	requiredErrorText string

	editorIcon       *Icon
	editorIconButton *Button
	showHidePassword *Button
	copy             *Button
	paste            *Button

	EditorIconButtonEvent func()

	m2 unit.Dp
	m5 unit.Dp

	click gesture.Click

	isDisableMenu bool
	isShowMenu    bool
	focused       bool

	// add space for error label if it is true
	isSpaceError bool

	isFirstFocus bool

	submitted       bool
	changed         bool
	selected        bool
	showHintOnFocus bool
}

func NewPasswordEditor(editor *widget.Editor, hint string) *Editor {
	editor.Mask = '*'
	e := NewEditor(editor, hint)
	e.isPassword = true
	e.showEditorIcon = false
	return e
}

func NewRestoreEditor(editor *widget.Editor, hint string, title string) *RestoreEditor {
	e := NewEditor(editor, hint)
	e.Bordered = false
	e.SelectionColor = color.NRGBA{}
	return &RestoreEditor{
		Edit:      e,
		Title:     title,
		LineColor: theme.Current().Gray2Color,
		height:    31,
	}
}

// NewIconEditor creates an editor widget with icon of choice
func NewIconEditor(editor *widget.Editor, hint string, icon *Icon, clickableIcon bool) *Editor {
	e := NewEditor(editor, hint)
	e.showEditorIcon = true
	e.isEditorButtonClickable = clickableIcon
	e.editorIcon = icon
	e.editorIcon.Color = theme.Current().Gray1Color
	e.editorIconButton.Icon = icon
	return e
}

func NewSearchEditor(editor *widget.Editor, hint string, icon *Icon) *Editor {
	e := NewEditor(editor, hint)
	e.showEditorIcon = true
	e.editorIcon = icon
	e.editorIcon.Color = theme.Current().Gray1Color
	e.editorIconButton.Icon = icon
	e.alignEditorIconEnd = false
	e.IsTitleLabel = false
	return e
}

func NewEditor(editor *widget.Editor, hint string) *Editor {
	th := theme.Current()
	errorLabel := material.Caption(th.Theme, "")
	errorLabel.Color = th.RedColor

	m := material.Editor(th.Theme, editor, hint)
	m.Color = th.InputColors.TextColor
	m.HintColor = th.GrayText3Color

	newEditor := &Editor{
		EditorStyle:        m,
		IsTitleLabel:       true,
		Bordered:           true,
		alignEditorIconEnd: true,
		requiredErrorText:  "Field is required",

		m2: unit.Dp(2),
		m5: unit.Dp(5),

		// Size:   values.DP24,
		editorIconButton: NewButton(ButtonStyle{
			Colors:    theme.ButtonColors{TextColor: th.GrayText3Color},
			Animation: NewButtonAnimationDefault(),
		}),
		showHidePassword: NewButton(ButtonStyle{
			Colors:    theme.ButtonColors{TextColor: th.GrayText3Color},
			Animation: NewButtonAnimationDefault(),
		}),
		CustomButton: NewButton(ButtonStyle{
			Colors:    theme.ButtonColors{TextColor: th.GrayText3Color},
			Animation: NewButtonAnimationDefault(),
		}),
		copy: NewButton(ButtonStyle{
			Text:      lang.Str("Copy"),
			TextSize:  values.TextSize10,
			Inset:     layout.UniformInset(values.DP5),
			Animation: NewButtonAnimationDefault(),
			Colors: theme.ButtonColors{
				TextColor:       th.PrimaryColor,
				BackgroundColor: color.NRGBA{},
				BorderColor:     color.NRGBA{},
			},
		}),
		paste: NewButton(ButtonStyle{
			Text:      lang.Str("Paste"),
			TextSize:  values.TextSize10,
			Inset:     layout.UniformInset(values.DP5),
			Animation: NewButtonAnimationDefault(),
			Colors: theme.ButtonColors{
				TextColor:       th.PrimaryColor,
				BackgroundColor: color.NRGBA{},
				BorderColor:     color.NRGBA{},
			},
		}),
	}

	return newEditor
}

func (e *Editor) Pressed(gtx C) bool {
	return e.click.Pressed() || e.copy.Clicked(gtx) || e.paste.Clicked(gtx)
}

func (e *Editor) FirstPressed(gtx C) bool {
	return !gtx.Source.Focused(e.Editor) && e.click.Pressed()

}

func (e *Editor) IsFocused() bool {
	return e.focused
}

func (e *Editor) SetFocus() {
	e.isFirstFocus = true
}

func (e *Editor) UpdateFocus(focus bool) {
	e.isFirstFocus = focus
}

func (e *Editor) Changed() bool {
	changed := e.changed
	e.changed = false
	return changed
}

func (e *Editor) Submitted() bool {
	submitted := e.submitted
	e.submitted = false
	return submitted
}

func (e *Editor) Selected() bool {
	selected := e.selected
	e.selected = false
	return selected
}

func (e *Editor) AlwaysShowHint() {
	e.showHintOnFocus = true
}

func (e *Editor) Text() string {
	return e.Editor.Text()
}

func (e *Editor) SetText(s string) {
	e.Editor.SetText(s)
}

func (e *Editor) Layout(gtx C) D {
	if e.isFirstFocus {
		e.isFirstFocus = false
		gtx.Execute(key.FocusCmd{Tag: e.Editor})
	}
	e.handleEvents(gtx)
	e.update(gtx)
	return e.layout(gtx)
}

func (e *Editor) update(gtx C) {
	for {
		ev, ok := e.click.Update(gtx.Source)
		if !ok {
			break
		}
		if e.click.Pressed() {
			if ev.NumClicks > 1 || (e.focused && !e.isShowMenu) {
				e.isShowMenu = true
			} else {
				e.isShowMenu = false
			}
		}
	}

	e.focused = gtx.Source.Focused(e.Editor)

	for {
		ev, ok := e.Editor.Update(gtx)
		if !ok {
			break
		}
		switch ev.(type) {
		case widget.ChangeEvent:
			e.changed = true
		case widget.SubmitEvent:
			e.submitted = true
		case widget.SelectEvent:
			e.selected = true
		}
	}
}

func (e *Editor) layout(gtx C) D {
	th := theme.Current()
	titleLabel := material.Body2(theme.Current().Theme, "")
	errorLabel := material.Body2(theme.Current().Theme, "")

	e.LineColor = th.Gray2Color
	titleLabel.Color = th.GrayText3Color
	errorLabel.Color = th.RedColor

	if e.Editor.Len() > 0 && len(e.Hint) > 0 {
		titleLabel.Text = e.Hint
	} else if e.Hint == "" {
		e.Hint = e.Title
		e.Title = ""
	}

	focused := gtx.Source.Focused(e.Editor)
	if focused {
		// Only non-read only editors should indicate an active state on focus.
		if !e.Editor.ReadOnly {
			e.LineColor = th.PrimaryColor
		}
		titleLabel.Color = th.PrimaryColor
		titleLabel.Text = e.Hint
		if !e.showHintOnFocus {
			e.Hint = ""
		}
	}

	if e.IsRequired && !focused && e.Editor.Len() == 0 {
		e.errorText = e.requiredErrorText
		e.LineColor = th.RedColor
	}

	if e.errorText != "" {
		titleLabel.Text = e.errorText
		e.LineColor = th.RedColor
		titleLabel.Color = e.LineColor
	}

	overLay := func(_ C) D { return D{} }
	//if e.Editor.ReadOnly {
	//	overLay = func(gtx C) D {
	//		gtxCopy := gtx
	//		gtxCopy.Constraints.Max.Y = gtx.Dp(values.DP46)
	//		return DisableLayout(nil, gtxCopy, nil, nil, 20, e.t.Color.Gray3, nil)
	//	}
	//	gtx = gtx.Disabled()
	//}

	dims := layout.UniformInset(e.m2).Layout(gtx, func(gtx C) D {
		return Card{CardColor: th.SurfaceColor, Radius: Radius(8)}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(e.editorLayout),
								layout.Rigid(func(gtx C) D {
									if e.errorText != "" {
										inset := layout.Inset{
											Top:  e.m2,
											Left: e.m5,
										}
										return inset.Layout(gtx, errorLabel.Layout)
									}
									if e.isSpaceError {
										return layout.Spacer{Height: values.DP18}.Layout(gtx)
									}
									return D{}
								}),
							)
						}),
						layout.Stacked(func(gtx C) D {
							if e.IsTitleLabel {
								return layout.Inset{
									Top:  values.DPMinus10,
									Left: values.DP10,
								}.Layout(gtx, func(gtx C) D {
									return Card{CardColor: th.SurfaceColor}.Layout(gtx, titleLabel.Layout)
								})
							}
							return D{}
						}),
						layout.Expanded(func(gtx C) D {
							defer pointer.PassOp{}.Push(gtx.Ops).Pop()
							defer clip.Rect(image.Rectangle{
								Max: gtx.Constraints.Min,
							}).Push(gtx.Ops).Pop()
							e.click.Add(gtx.Ops)
							return D{}
						}),
						layout.Stacked(overLay),
					)
				}),
			)
		})
	})
	if !e.isDisableMenu {
		e.editorMenusLayout(gtx, dims.Size.Y)
	}
	return dims
}

func (e *Editor) editorLayout(gtx C) D {
	if e.Bordered {
		border := widget.Border{Color: e.LineColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
		return border.Layout(gtx, func(gtx C) D {
			inset := layout.Inset{
				Top:    values.DP7,
				Bottom: values.DP7,
				Left:   values.DP12,
				Right:  values.DP12,
			}
			return inset.Layout(gtx, e.editor)
		})
	}

	inset := layout.Inset{
		Top:    values.DP3,
		Bottom: values.DP3,
		Left:   values.DP12,
		Right:  values.DP12,
	}
	return inset.Layout(gtx, e.editor)
}

func (e *Editor) editorMenusLayout(gtx C, editorHeight int) {
	th := theme.Current()
	e.isShowMenu = e.isShowMenu && (gtx.Source.Focused(e.Editor) || e.copy.Hovered() || e.paste.Hovered())
	if e.isShowMenu {
		flexChildren := make([]layout.FlexChild, 0)
		if len(e.Editor.SelectedText()) > 0 {
			flexChildren = append(flexChildren, layout.Rigid(e.copy.Layout))
			flexChildren = append(flexChildren, layout.Rigid(NewLine(20, 1).Layout))
		}
		flexChildren = append(flexChildren, layout.Rigid(e.paste.Layout))
		gtxCopy := gtx
		macro := op.Record(gtxCopy.Ops)
		LinearLayout{
			Width:      WrapContent,
			Height:     WrapContent,
			Background: th.SurfaceColor,
			Margin:     layout.Inset{Top: gtxCopy.Metric.PxToDp(-(editorHeight - 10))},
			Padding:    layout.UniformInset(values.DP5),
			Alignment:  layout.Middle,
			Border:     widget.Border{CornerRadius: 8, Color: th.Gray2Color, Width: 0.5},
		}.Layout(gtxCopy, flexChildren...)
		op.Defer(gtxCopy.Ops, macro.Stop())
	}
}

func (e *Editor) layoutIconEditor(gtx C) D {
	inset := layout.Inset{
		Top: e.m2,
	}

	if e.alignEditorIconEnd {
		inset.Left = e.m5
	} else {
		inset.Right = e.m5
	}

	return inset.Layout(gtx, func(gtx C) D {
		if e.isEditorButtonClickable {
			return e.editorIconButton.Layout(gtx)
		}
		return e.editorIcon.LayoutSize(gtx, 25)
	})
}

func (e *Editor) editor(gtx C) D {
	th := theme.Current()
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if e.showEditorIcon && !e.alignEditorIconEnd {
				return e.layoutIconEditor(gtx)
			}
			return D{}
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Inset{Top: e.m5, Bottom: e.m5}.Layout(gtx, e.EditorStyle.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			if e.ExtraText == "" {
				return D{}
			}
			lbl := material.Label(th.Theme, values.TextSize16, e.ExtraText)
			return layout.Inset{Top: values.DP5}.Layout(gtx, lbl.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			if e.showEditorIcon && e.alignEditorIconEnd {
				return e.layoutIconEditor(gtx)
			} else if e.isPassword {
				inset := layout.Inset{
					Top:  e.m2,
					Left: e.m5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					icon := NewIcon(icons.ActionVisibilityOff)
					if e.Editor.Mask == '*' {
						icon = NewIcon(icons.ActionVisibility)
					}
					e.showHidePassword.Icon = icon
					return e.showHidePassword.Layout(gtx)
				})
			}
			return D{}
		}),
		layout.Rigid(func(gtx C) D {
			if e.HasCustomButton {
				inset := layout.Inset{
					Top:   e.m5,
					Left:  e.m5,
					Right: e.m5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					e.CustomButton.TextSize = unit.Sp(10)
					return e.CustomButton.Layout(gtx)
				})
			}
			return D{}
		}),
	)
}

func (e *Editor) handleEvents(gtx C) {
	if e.showHidePassword.Clicked(gtx) {
		if e.Editor.Mask == '*' {
			e.Editor.Mask = 0
		} else if e.Editor.Mask == 0 {
			e.Editor.Mask = '*'
		}
	}

	if e.editorIconButton.Clicked(gtx) {
		e.EditorIconButtonEvent()
	}

	if e.copy.Clicked(gtx) {
		gtx.Execute(clipboard.WriteCmd{Data: io.NopCloser(strings.NewReader(e.Editor.SelectedText()))})
		e.isShowMenu = false
	}

	if e.paste.Clicked(gtx) {
		gtx.Execute(clipboard.ReadCmd{Tag: e.Editor})
		e.isShowMenu = false
	}
}

func (re RestoreEditor) Layout(gtx C) D {
	th := theme.Current()
	titleLabel := material.Body2(theme.Current().Theme, "")
	errorLabel := material.Caption(th.Theme, "")
	titleLabel.Color = th.GrayText3Color
	errorLabel.Color = th.RedColor

	width := int(gtx.Metric.PxPerDp * 2.0)
	height := int(gtx.Metric.PxPerDp * float32(re.height))
	l := SeparatorVertical(height, width)

	if gtx.Source.Focused(re.Edit.Editor) {
		titleLabel.Color, re.LineColor, l.Color = th.PrimaryColor, th.PrimaryColor, th.PrimaryColor
	} else {
		l.Color = th.Gray2Color
	}
	border := widget.Border{Color: re.LineColor, CornerRadius: values.DP8, Width: values.DP2}
	return border.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(values.DP40)
				return layout.Center.Layout(gtx, titleLabel.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: unit.Dp(-3), Right: unit.Dp(5)}.Layout(gtx, l.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				edit := re.Edit.Layout(gtx)
				re.height = edit.Size.Y
				return edit
			}),
		)
	})
}

func (e *Editor) SetRequiredErrorText(txt string) {
	e.requiredErrorText = txt
}

func (e *Editor) SetError(text string) {
	e.errorText = text
}

func (e *Editor) HasError() bool {
	return e.errorText != ""
}

func (e *Editor) ClearError() {
	e.errorText = ""
}

func (e *Editor) IsDirty() bool {
	return e.errorText == ""
}

func (e *Editor) AllowSpaceError(allow bool) {
	e.isSpaceError = allow
}

// Line represents a rectangle widget with an initial thickness of 1
type Line struct {
	Height     int
	Width      int
	Color      color.NRGBA
	isVertical bool
}

// SeparatorVertical returns a vertical line widget instance
func SeparatorVertical(height, width int) Line {
	vLine := NewLine(height, width)
	vLine.isVertical = true
	return vLine
}

// NewLine returns a line widget instance
func NewLine(height, width int) Line {
	if height == 0 {
		height = 1
	}

	col := theme.Current().PrimaryColor
	col.A = 150
	return Line{
		Height: height,
		Width:  width,
		Color:  col,
	}
}

func NewSeparator() Line {
	l := NewLine(1, 0)
	l.Color = theme.Current().Gray3Color
	return l
}

func (l Line) Layout(gtx C) D {
	if l.Width == 0 {
		l.Width = gtx.Constraints.Max.X
	}

	if l.isVertical && l.Height == 0 {
		l.Height = gtx.Constraints.Max.Y
	}

	line := image.Rectangle{
		Max: image.Point{
			X: l.Width,
			Y: l.Height,
		},
	}
	defer clip.Rect(line).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, l.Color)

	return D{Size: line.Max}
}
