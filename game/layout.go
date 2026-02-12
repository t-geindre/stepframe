package game

type Layout interface {
	Layout(x, y int) (int, int)
}

type layoutFunc struct {
	f func(x, y int) (int, int)
}

func NewLayoutFunc(f func(x, y int) (int, int)) Layout {
	return &layoutFunc{f: f}
}

func (l *layoutFunc) Layout(x, y int) (int, int) {
	return l.f(x, y)
}
