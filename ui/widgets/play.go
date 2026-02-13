package widgets

import (
	img "image"
	"image/color"
	"stepframe/seq"
	"stepframe/ui/theme"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

// Playing #22C55E or #00C853
// Stopped #EF4444 or #D50000
// Armed #F59E0B or #FFB300

var green = color.RGBA{R: 0x22, G: 0xC5, B: 0x5E, A: 0xFF}
var red = color.RGBA{R: 0xEF, G: 0x44, B: 0x44, A: 0xFF}
var orange = color.RGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF}

type PlayState int

const (
	PlayStateStopped PlayState = iota
	PlayStatePlaying
	PlayStateArmed
	PlayStateNone
)

type Play struct {
	*widget.Widget
	state     PlayState
	lastPulse time.Time
	id        seq.TrackId
	theme     *theme.Theme
}

func NewPlay(th *theme.Theme, opts ...widget.WidgetOpt) *Play {
	return &Play{
		Widget: widget.NewWidget(opts...),
		state:  PlayStateNone,
		theme:  th,
	}
}

func (p *Play) Render(screen *ebiten.Image) {
	dst := screen.SubImage(p.Widget.Rect).(*ebiten.Image)

	switch p.state {
	case PlayStateStopped:
		p.Fill(dst, red, false)
		p.drawIcon(p.theme.Icons[theme.IconPlay], dst)
	case PlayStatePlaying:
		p.Fill(dst, green, true)
		p.drawIcon(p.theme.Icons[theme.IconStop], dst)
	case PlayStateArmed:
		p.Fill(dst, orange, true)
		p.drawIcon(p.theme.Icons[theme.IconPause], dst)
	case PlayStateNone:
		p.Fill(dst, colornames.Black, false)
		p.drawIcon(p.theme.Icons[theme.IconPlus], dst)
	}
}

func (p *Play) Fill(dst *ebiten.Image, col color.Color, pulse bool) {
	const PulseDuration = 100 * time.Millisecond
	const PulseStrength = 0.45 // 0..1

	// base fill
	dst.Fill(col)

	if !pulse || p.lastPulse.IsZero() {
		return
	}

	elapsed := time.Since(p.lastPulse)
	if elapsed <= 0 || elapsed >= PulseDuration {
		return
	}

	t := 1.0 - float64(elapsed)/float64(PulseDuration)

	intensity := t * PulseStrength

	r, g, b, a := rgba8(col)
	pr := lerp8(r, 0xFF, intensity)
	pg := lerp8(g, 0xFF, intensity)
	pb := lerp8(b, 0xFF, intensity)

	dst.Fill(color.RGBA{R: pr, G: pg, B: pb, A: a})
}

func (p *Play) drawIcon(i, dst *ebiten.Image) {
	r := p.GetWidget().Rect
	op := &ebiten.DrawImageOptions{}

	rectW := r.Dx()
	rectH := r.Dy()
	imgW, imgH := i.Bounds().Dx(), i.Bounds().Dy()

	x := r.Min.X + (rectW-imgW)/2
	y := r.Min.Y + (rectH-imgH)/2

	op.GeoM.Translate(float64(x), float64(y))

	dst.DrawImage(i, op)
}

func (p *Play) GetWidget() *widget.Widget {
	return p.Widget
}

func (p *Play) PreferredSize() (int, int) {
	return 200, 200
}

func (p *Play) SetLocation(rect img.Rectangle) {
	p.Widget.SetLocation(rect)
}

func (p *Play) Validate() {
}

func (p *Play) Update(updObj *widget.UpdateObject) {
	p.Widget.Update(updObj)
}

func (p *Play) SetState(state PlayState) {
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

// UTILS
func lerp8(from, to uint8, t float64) uint8 {
	if t <= 0 {
		return from
	}
	if t >= 1 {
		return to
	}
	return uint8(float64(from) + (float64(to)-float64(from))*t)
}

func rgba8(c color.Color) (r, g, b, a uint8) {
	R, G, B, A := c.RGBA() // 0..65535
	return uint8(R >> 8), uint8(G >> 8), uint8(B >> 8), uint8(A >> 8)
}
