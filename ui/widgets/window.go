package widgets

import (
	"image"
	"stepframe/ui/container"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Window struct {
	*widget.Window
	options []widget.WindowOpt

	preferredW, preferredH int
	close                  widget.RemoveWindowFunc
	attached               widget.Containerer
}

func NewWindow(preferredW, preferredH int) *Window {
	var w *Window

	w = &Window{
		preferredW: preferredW,
		preferredH: preferredH,
		options: []widget.WindowOpt{
			widget.WindowOpts.Contents(
				container.NewRow().
					SetDirection(widget.DirectionVertical).
					SetPadding(theme.Current.PanelTheme.Padding).
					SetSpacing(theme.Current.PanelTheme.Spacing).
					SetBackgroundOffsetImage(theme.Current.PanelTheme.BackgroundImage),
			),
			widget.WindowOpts.Modal(),
			widget.WindowOpts.CloseMode(widget.CLICK_OUT),
			widget.WindowOpts.ClosedHandler(func(args *widget.WindowClosedEventArgs) {
				w.Close()
			}),
		},
	}

	return w
}

func (w *Window) WithTitleBar(icon theme.Icon, title string) *Window {
	header := container.NewRow().
		SetPadding(theme.Current.PanelTheme.Padding).
		SetSpacing(theme.Current.PanelTheme.Spacing).
		SetBackgroundOffsetImage(theme.Current.PanelTheme.BackgroundImage)

	if icon != theme.IconNone {
		header.AddChild(NewIcon(icon, theme.IconSizeMedium))
	}
	header.AddChild(NewLabel(title))

	w.options = append(w.options, widget.WindowOpts.TitleBar(
		header,
		theme.Current.IconsTheme.Sizes[theme.IconSizeMedium]+ // Optimistic guess of title bar height
			theme.Current.PanelTheme.Padding.Top+
			theme.Current.PanelTheme.Padding.Bottom,
	))

	return w
}

func (w *Window) AttachedTo(c widget.Containerer) *Window {
	w.attached = c
	return w
}

func (w *Window) Open() {
	w.build()
	w.relocate()
	w.close = currentUi.AddWindow(w.Window)
}

func (w *Window) Close() {
	if w.close != nil {
		w.close()
		w.close = nil
	}
}

func (w *Window) IsOpen() bool {
	return w.close != nil
}

func (w *Window) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	w.build()
	return w.Contents.AddChild(children...)
}

func (w *Window) relocate() {
	wRect := image.Rect(0, 0, w.preferredW, w.preferredH)
	uiRect := currentUi.Container.GetWidget().Rect

	// Size down if not enough space in UI
	if wRect.Dx() > uiRect.Dx() {
		wRect.Max.X = uiRect.Dx()
	}

	if wRect.Dy() > uiRect.Dy() {
		wRect.Max.Y = uiRect.Dy()
	}

	if w.attached != nil {
		// Window attached, try to place it right below the attached widget
		attachedRect := w.attached.GetWidget().Rect
		wRect = wRect.Add(image.Pt(attachedRect.Min.X+attachedRect.Dx()/2-wRect.Dx()/2, attachedRect.Max.Y))

		// Check window still fits in UI
		if wRect.Min.X < uiRect.Min.X {
			wRect = wRect.Add(image.Pt(uiRect.Min.X-wRect.Min.X, 0))
		}

		if wRect.Max.X > uiRect.Max.X {
			wRect = wRect.Add(image.Pt(uiRect.Max.X-wRect.Max.X, 0))
		}

		if wRect.Min.Y < uiRect.Min.Y {
			wRect = wRect.Add(image.Pt(0, uiRect.Min.Y-wRect.Min.Y))
		}

		if wRect.Max.Y > uiRect.Max.Y {
			wRect = wRect.Add(image.Pt(0, uiRect.Max.Y-wRect.Max.Y))
		}

	} else {
		// Center in UI
		wRect = wRect.Add(image.Pt(
			(uiRect.Dx()-wRect.Dx())/2,
			(uiRect.Dy()-wRect.Dy())/2,
		))
	}

	if wRect.Eq(w.Window.GetContainer().GetWidget().Rect) {
		return
	}

	w.Window.GetContainer().SetLocation(wRect)
}

func (w *Window) build() {
	if w.Window == nil {
		w.Window = widget.NewWindow(w.options...)
		w.Contents.GetWidget().OnUpdate = func(widget.HasWidget) {
			if w.IsOpen() {
				w.relocate()
			}
		}
	}
}
