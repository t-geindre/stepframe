package widgets

import (
	"stepframe/seq"
	"stepframe/ui/theme"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayState int

const (
	PlayStateStopped PlayState = iota
	PlayStatePlaying
	PlayStateArmed
	PlayStateNone
)

type Play struct {
	*widget.Container
	icon      *widget.Graphic
	state     PlayState
	lastPulse time.Time
	id        seq.TrackId
	theme     *theme.Theme
}

func NewPlay(th *theme.Theme) *Play {
	p := &Play{
		Container: widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(th.PlayTheme.None),
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
		icon: widget.NewGraphic(
			widget.GraphicOpts.Image(th.Icons[theme.IconPlus]),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
			),
		),
		state: PlayStateNone,
		theme: th,
	}

	p.AddChild(p.icon)

	return p
}

func (p *Play) Render(screen *ebiten.Image) {
	p.Container.Render(screen)

	if p.state != PlayStatePlaying && p.state != PlayStateArmed {
		return
	}

	if p.lastPulse.IsZero() {
		return
	}

	elapsed := time.Since(p.lastPulse)
	if elapsed <= 0 || elapsed >= p.theme.PlayTheme.PulseDuration {
		return
	}

	t := 1.0 - float64(elapsed)/float64(p.theme.PlayTheme.PulseDuration)
	alpha := t * p.theme.PlayTheme.PulseStrength

	if alpha <= 0 {
		return
	}

	r := p.GetWidget().Rect
	dst := screen.SubImage(r).(*ebiten.Image)

	p.theme.PlayTheme.Pulse.Draw(dst, r.Dx(), r.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.ColorScale.Scale(float32(alpha), float32(alpha), float32(alpha), float32(alpha))
		opts.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	})
}

func (p *Play) SetState(state PlayState) {
	switch state {
	case PlayStateStopped:
		p.Container.SetBackgroundImage(p.theme.PlayTheme.Stopped)
		p.icon.Image = p.theme.Icons[theme.IconPlay]
	case PlayStatePlaying:
		p.Container.SetBackgroundImage(p.theme.PlayTheme.Playing)
		p.icon.Image = p.theme.Icons[theme.IconStop]
	case PlayStateArmed:
		p.Container.SetBackgroundImage(p.theme.PlayTheme.Armed)
		p.icon.Image = p.theme.Icons[theme.IconPause]
	case PlayStateNone:
		p.Container.SetBackgroundImage(p.theme.PlayTheme.None)
		p.icon.Image = p.theme.Icons[theme.IconPlus]
	}
	p.state = state
}

func (p *Play) GetState() PlayState {
	return p.state
}

func (p *Play) Pulse() {
	p.lastPulse = time.Now()
}

func (p *Play) SetId(id seq.TrackId) {
	p.id = id
}

func (p *Play) Id() seq.TrackId {
	return p.id
}
