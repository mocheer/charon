package tm

import "time"

// SetTimeout 毫秒
func SetTimeout(callback func(), delay time.Duration) func() {
	timer := time.NewTimer(time.Millisecond * delay)
	fired := false

	go func() {
		<-timer.C
		callback()
		fired = true
	}()

	return func() {
		if fired == false {
			timer.Stop()
		}
	}
}
