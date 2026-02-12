package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	updater []Updater
	drawer  []Drawer
	layout  Layout
}

func NewGame(obj ...any) *Game {
	g := &Game{}

	for _, o := range obj {
		if u, ok := o.(Updater); ok {
			g.updater = append(g.updater, u)
		}
		if d, ok := o.(Drawer); ok {
			g.drawer = append(g.drawer, d)
		}
		if l, ok := o.(Layout); ok {
			g.layout = l
		}
	}

	return g
}

func (g *Game) Update() error {
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
