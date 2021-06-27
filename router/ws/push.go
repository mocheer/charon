package ws

import (
	"time"

	"github.com/mocheer/pluto/jsg"
)

//
func push() {
	jsg.SetInterval(getTyphoon, time.Minute)
}
