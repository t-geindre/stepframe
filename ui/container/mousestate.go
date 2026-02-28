package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

type MouseState[T widget.Containerer] struct {
	*OffsetBackground[T]
	outer T

	isHover bool
	onClick func()
	cursor  string

	idle, hover, pressed                         theme.Image
	idleColorize, hoverColorize, pressedColorize theme.Color
}

func NewMouseState[T widget.Containerer](offsetBg *OffsetBackground[T], outer T) *MouseState[T] {
	m := &MouseState[T]{
		OffsetBackground: offsetBg,
		outer:            outer,
		onClick:          func() {},
	}

	m.GetWidget().CursorEnterEvent.AddHandler(m.onCursorEnter)
	m.GetWidget().CursorExitEvent.AddHandler(m.onCursorExit)
	m.GetWidget().MouseButtonPressedEvent.AddHandler(m.onPress)
	m.GetWidget().MouseButtonReleasedEvent.AddHandler(m.onRelease)
	m.GetWidget().OnUpdate = m.onUpdate

	return m
}

func (m *MouseState[T]) SetIdleBackgroundOffsetImage(img theme.Image) T {
	m.idle = img
	return m.outer
}

func (m *MouseState[T]) SetIdleBackgroundOffsetColorize(color theme.Color) T {
	m.idleColorize = color
	return m.outer
}

func (m *MouseState[T]) SetHoverBackgroundOffsetImage(img theme.Image) T {
	m.hover = img
	return m.outer
}

func (m *MouseState[T]) SetHoverBackgroundOffsetColorize(color theme.Color) T {
	m.hoverColorize = color
	return m.outer
}

func (m *MouseState[T]) SetPressedBackgroundOffsetImage(img theme.Image) T {
	m.pressed = img
	return m.outer
}

func (m *MouseState[T]) SetPressedBackgroundOffsetColorize(color theme.Color) T {
	m.pressedColorize = color
	return m.outer
}

func (m *MouseState[T]) SetHoverCursor(cursor string) T {
	m.cursor = cursor
	return m.outer
}

func (m *MouseState[T]) SetOnClick(click func()) T {
	m.onClick = click
	return m.outer
}

func (m *MouseState[T]) onCursorEnter(any) {
	m.isHover = true
	if m.hover != theme.ImageNone {
		m.SetBackgroundOffsetImage(m.hover)
	}
	if m.hoverColorize != theme.ColorNone {
		m.SetColorizedBackgroundOffsetImage(m.hoverColorize)
	}
}

func (m *MouseState[T]) onCursorExit(any) {
	m.isHover = false
	if m.hover != theme.ImageNone {
		m.SetBackgroundOffsetImage(m.idle)
	}
	if m.hoverColorize != theme.ColorNone {
		m.SetColorizedBackgroundOffsetImage(m.idleColorize)
	}
}

func (m *MouseState[T]) onPress(any) {
	if m.pressed != theme.ImageNone {
		m.SetBackgroundOffsetImage(m.pressed)
	}
	if m.pressedColorize != theme.ColorNone {
		m.SetColorizedBackgroundOffsetImage(m.pressedColorize)
	}
}

func (m *MouseState[T]) onRelease(any) {
	if m.isHover {
		if m.pressed != theme.ImageNone {
			m.SetBackgroundOffsetImage(m.hover)
		}
		if m.hoverColorize != theme.ColorNone {
			m.SetColorizedBackgroundOffsetImage(m.hoverColorize)
		}
		m.onClick()
	} else {
		if m.pressed != theme.ImageNone {
			m.SetBackgroundOffsetImage(m.idle)
		}
		if m.hoverColorize != theme.ColorNone {
			m.SetColorizedBackgroundOffsetImage(m.idleColorize)
		}
	}
}

func (m *MouseState[T]) onUpdate(widget.HasWidget) {
	if m.isHover && m.cursor != input.CURSOR_NONE && m.cursor != "" {
		input.SetCursorShape(m.cursor)
	}
}
