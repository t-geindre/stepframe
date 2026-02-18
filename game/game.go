package game

import (
	"context"
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rs/zerolog"
)

type Game struct {
	updater []Updater
	drawer  []Drawer
	layout  Layout
	ctx     context.Context
	logger  zerolog.Logger
}

func RunGame(logger zerolog.Logger, ctx context.Context, obj ...any) {
	logger = logger.With().Str("component", "game").Logger()

	g := &Game{ctx: ctx, logger: logger}

	for _, o := range obj {
		if u, ok := o.(Updater); ok {
			g.updater = append(g.updater, u)
		}
		if ue, ok := o.(UpdaterWithoutError); ok {
			g.updater = append(g.updater, NewUpdateFunc(func() error {
				ue.Update()
				return nil
			}))
		}
		if d, ok := o.(Drawer); ok {
			g.drawer = append(g.drawer, d)
		}
		if l, ok := o.(Layout); ok {
			g.layout = l
		}
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	logger.Info().Msg("starting")
	err := ebiten.RunGame(g)

	if err != nil && !errors.Is(err, ebiten.Termination) {
		logger.Err(err).Msg("failed to run game")
	}
}

func (g *Game) Update() error {
	if g.ctx.Err() != nil || ebiten.IsWindowBeingClosed() {
		g.logger.Info().Msg("exiting")
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	for _, u := range g.updater {
		err := u.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, d := range g.drawer {
		d.Draw(screen)
	}
}

func (g *Game) Layout(x, y int) (int, int) {
	if g.layout != nil {
		return g.layout.Layout(x, y)
	}

	return x, y
}
