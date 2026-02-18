package game

type Updater interface {
	Update() error
}

type UpdaterWithoutError interface {
	Update()
}

type updateFunc struct {
	f func() error
}

func NewUpdateFunc(f func() error) Updater {
	return &updateFunc{f: f}
}

func (u *updateFunc) Update() error {
	return u.f()
}
