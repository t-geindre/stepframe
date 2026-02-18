package seq

type State int

const (
	TsStopped State = iota
	TsPlaying
	TsPaused
)

type TimeShifter struct {
	state State

	lastGlobal int64 // dernier tick global reçu
	baseGlobal int64 // tick global correspondant à local==0
	nowLocal   int64 // temps local figé quand paused/stopped

	pausedAtGlobal int64 // tick global au moment du pause
}

func NewTimeShifter() *TimeShifter { return &TimeShifter{state: TsStopped} }

// Tick doit être appelé à chaque tick global (même si paused)
func (t *TimeShifter) Tick(global int64) {
	t.lastGlobal = global
	if t.state == TsPlaying {
		t.nowLocal = global - t.baseGlobal
	}
}

func (t *TimeShifter) Now() int64        { return t.nowLocal }
func (t *TimeShifter) State() State      { return t.state }
func (t *TimeShifter) LastGlobal() int64 { return t.lastGlobal }

func (t *TimeShifter) Play() {
	switch t.state {
	case TsStopped:
		t.baseGlobal = t.lastGlobal
		t.nowLocal = 0
	case TsPaused:
		pausedDur := t.lastGlobal - t.pausedAtGlobal
		t.baseGlobal += pausedDur
	}
	t.state = TsPlaying
}

func (t *TimeShifter) Pause() {
	if t.state == TsPlaying {
		t.pausedAtGlobal = t.lastGlobal
		t.state = TsPaused
	}
}

func (t *TimeShifter) Stop(rewind bool) {
	t.state = TsStopped
	if rewind {
		t.nowLocal = 0
	}
}

func (t *TimeShifter) IsBeat(ppqn int64) bool {
	if ppqn <= 0 {
		return false
	}
	return t.state == TsPlaying && (t.nowLocal%ppqn == 0)
}
