package widgets

import (
	"stepframe/ui/container"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Button struct {
	*container.Row
	hover bool
	click func()
}

func NewButton(click func()) *Button {
	return &Button{
		Row:   container.NewHorizontalRow().WithPadding(),
		click: click,
	}
}

func NewIconButton(click func(), icon theme.Icon, size theme.IconSize) *Button {
	btn := NewButton(click)
	btn.AddChild(NewIcon(icon, size))
	return btn
}

func (b *Button) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, child := range children {
		child.GetWidget().LayoutData = widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		}
	}

	return b.Row.AddChild(children...)
}

func (b *Button) Validate() {
	// Register event handlers
	b.Row.GetWidget().CursorEnterEvent.AddHandler(b.onCursorEnter)
	b.Row.GetWidget().CursorExitEvent.AddHandler(b.onCursorExit)
	b.Row.GetWidget().MouseButtonPressedEvent.AddHandler(b.onPress)
	b.Row.GetWidget().MouseButtonReleasedEvent.AddHandler(b.onRelease)

	// Default state
	b.onCursorExit(nil)

	// Validate container
	b.Row.Validate()
}

func (b *Button) onCursorEnter(any) {
	b.hover = true
	b.Row.SetBackgroundImage(theme.Current.ButtonTheme.Image.Hover)
}

func (b *Button) onCursorExit(any) {
	b.hover = false
	b.Row.SetBackgroundImage(theme.Current.ButtonTheme.Image.Idle)
}

func (b *Button) onPress(any) {
	b.Row.SetBackgroundImage(theme.Current.ButtonTheme.Image.Pressed)
}

func (b *Button) onRelease(any) {
	if b.hover {
		b.Row.SetBackgroundImage(theme.Current.ButtonTheme.Image.Hover)
		if b.click != nil {
			b.click()
		}
	} else {
		b.Row.SetBackgroundImage(theme.Current.ButtonTheme.Image.Idle)
	}
}
