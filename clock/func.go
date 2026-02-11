package clock

import (
	"runtime"
	"time"
)

func sleepUntil(deadline time.Time) {
	// spinThreshold Spin under a given delay (300 µS)
	const spinThreshold = 300 * time.Microsecond

	for {
		now := time.Now()
		remain := deadline.Sub(now)
		if remain <= 0 {
			return
		}
		if remain > spinThreshold {
			time.Sleep(remain - spinThreshold)
			continue
		}
		for time.Now().Before(deadline) {
			runtime.Gosched()
		}
		return
	}
}
