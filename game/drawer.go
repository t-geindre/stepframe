package game

import "github.com/hajimehoshi/ebiten/v2"

type Drawer interface {
	Draw(screen *ebiten.Image)
}

type drawFunc struct {
	f func(screen *ebiten.Image)
}

func NewDrawFunc(f func(screen *ebiten.Image)) Drawer {
	return &drawFunc{f: f}
}

func (d *drawFunc) Draw(screen *ebiten.Image) {
	d.f(screen)
}
