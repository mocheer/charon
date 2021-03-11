package tm

import "time"

// SetInterval
func SetInterval(callback func(), delay time.Duration) func() {
	ticker := time.NewTicker(time.Millisecond * delay)
	go func() {
		for range ticker.C {
			callback()
		}
	}()
	return func() {
		ticker.Stop()
	}
}
